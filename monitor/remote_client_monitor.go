package monitor

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/no-src/gofs/action"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/eventlog"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/internal/cbool"
	"github.com/no-src/gofs/internal/clist"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/tran"
	"github.com/no-src/gofs/util/jsonutil"
	"github.com/no-src/gofs/util/stringutil"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

type remoteClientMonitor struct {
	baseMonitor

	client      tran.Client
	closed      *cbool.CBool
	messages    *clist.CList
	currentUser *auth.HashUser
	authorized  bool
	authChan    chan contract.Status
	infoChan    chan contract.Message
	timeout     time.Duration
}

// NewRemoteClientMonitor create an instance of remoteClientMonitor to monitor the remote file change
func NewRemoteClientMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, host string, port int, enableTLS bool, users []*auth.User, eventWriter io.Writer, enableSyncDelay bool, syncDelayEvents int, syncDelayTime time.Duration) (Monitor, error) {
	if syncer == nil {
		err := errors.New("syncer can't be nil")
		return nil, err
	}
	m := &remoteClientMonitor{
		client:      tran.NewClient(host, port, enableTLS),
		messages:    clist.New(),
		baseMonitor: newBaseMonitor(syncer, retry, syncOnce, eventWriter, enableSyncDelay, syncDelayEvents, syncDelayTime),
		authChan:    make(chan contract.Status, 100),
		infoChan:    make(chan contract.Message, 100),
		closed:      cbool.New(false),
		timeout:     time.Minute * 3,
	}
	if len(users) > 0 {
		user := users[0]
		hashUser := user.ToHashUser()
		m.currentUser = hashUser
	}
	return m, nil
}

// auth send auth request
func (m *remoteClientMonitor) auth() error {
	// if the currentUser is nil, it means to anonymous access
	if m.currentUser == nil {
		return nil
	}
	go m.retry.Do(func() error {
		m.currentUser.RefreshExpires()
		authData := auth.GenerateAuthCommandData(m.currentUser)
		err := m.client.Write(authData)
		return err
	}, "send auth request")

	var status contract.Status
	select {
	case status = <-m.authChan:
	case <-time.After(m.timeout):
		return fmt.Errorf("auth timeout for %s", m.timeout.String())
	}
	if status.Code != contract.Success {
		return errors.New("receive auth command response error => " + status.Message)
	}

	// auth success
	m.authorized = true
	log.Info("auth success, current client is authorized => [%s] ", status.Message)
	return nil
}

func (m *remoteClientMonitor) Start() error {
	if m.client == nil {
		return errors.New("remote sync client is nil")
	}
	// connect -> auth -> info|read
	err := m.client.Connect()
	if err != nil {
		return err
	}

	w := m.receive()

	if err = m.auth(); err != nil {
		return err
	}

	// execute -sync_once flag
	if m.syncOnce {
		return m.syncAndShutdown(w)
	}

	// execute -sync_cron flag
	if err := m.startCron(m.sync); err != nil {
		return err
	}

	go m.startReceiveWriteNotify()
	go m.startSyncWrite()
	go m.startProcessMessage()

	return w.Wait()
}

// sync try to sync all the files once
func (m *remoteClientMonitor) sync() (err error) {
	go func() {
		log.ErrorIf(m.client.Write(contract.InfoCommand), "write info command error")
	}()
	var info contract.FileServerInfo
	var infoMsg contract.Message
	select {
	case infoMsg = <-m.infoChan:
	case <-time.After(m.timeout):
		return fmt.Errorf("sync timeout for %s", m.timeout.String())
	}
	err = jsonutil.Unmarshal(infoMsg.Data, &info)
	if err != nil {
		return err
	}

	if info.Code != contract.Success {
		return errors.New("receive info command response error => " + info.Message)
	}

	return m.syncer.SyncOnce(info.ServerAddr + info.SourcePath)
}

// syncAndShutdown execute sync and then try to shut down
func (m *remoteClientMonitor) syncAndShutdown(w wait.Wait) (err error) {
	if err = m.sync(); err != nil {
		return err
	}
	if err = m.Shutdown(); err != nil {
		return err
	}
	return w.Wait()
}

// receive start receiving messages and parse the message, send to consumers according to the api type.
// if receive a shutdown notify, then stop reading the message.
func (m *remoteClientMonitor) receive() wait.Wait {
	wd := wait.NewWaitDone()
	shutdown := cbool.New(false)
	go m.waitShutdown(shutdown, wd)
	go m.readMessage(shutdown, wd)
	return wd
}

// waitShutdown wait for the shutdown notify then mark the work done
func (m *remoteClientMonitor) waitShutdown(st *cbool.CBool, wd wait.WaitDone) {
	select {
	case <-st.SetC(<-m.shutdown):
		{
			if st.Get() {
				log.ErrorIf(m.Close(), "close remote client monitor error")
				wd.Done()
			}
		}
	}
}

// readMessage loop read the messages, if receive a message, parse the message then send to consumers according to the api type.
// if receive a shutdown notify, then stop reading the message.
func (m *remoteClientMonitor) readMessage(st *cbool.CBool, wd wait.WaitDone) {
	for {
		if m.closed.Get() {
			wd.DoneWithError(errors.New("remote monitor is closed"))
			break
		}
		data, err := m.client.ReadAll()
		if err != nil {
			if st.Get() {
				break
			}
			log.Error(err, "remote client monitor read data error")
			if m.client.IsClosed() {
				m.authorized = false
				log.Debug("try reconnect to server %s:%d", m.client.Host(), m.client.Port())
				m.retry.Do(func() error {
					if m.client.IsClosed() {
						innerErr := m.client.Connect()
						if innerErr != nil {
							return innerErr
						}
					}
					return nil
				}, fmt.Sprintf("client reconnect to %s:%d", m.client.Host(), m.client.Port()))

				if !m.authorized {
					go m.auth()
				}
			}
		} else {
			m.parseMessage(data)
		}
	}
}

// parseMessage arse the message then send to consumers according to the api type
func (m *remoteClientMonitor) parseMessage(data []byte) error {
	var status contract.Status
	err := jsonutil.Unmarshal(data, &status)
	if err != nil {
		log.Error(err, "remote client monitor unmarshal data error")
		return err
	}

	switch status.ApiType {
	case contract.AuthApi:
		m.authChan <- status
	case contract.InfoApi:
		m.infoChan <- contract.NewMessage(data)
	case contract.SyncMessageApi:
		m.messages.PushBack(contract.NewMessage(data))
	default:
		log.Warn("remote client monitor received a unknown data => %s", string(data))
	}
	return nil
}

// startProcessMessage start loop to process the file change messages
func (m *remoteClientMonitor) startProcessMessage() {
	for {
		m.waitSyncDelay(m.messages.Len)

		element := m.messages.Front()
		if element == nil || element.Value == nil {
			if element != nil {
				m.messages.Remove(element)
			}
			m.resetSyncDelay()
			<-time.After(time.Second)
			continue
		}
		message := element.Value.(contract.Message)
		log.Info("client read request => %s", message.String())
		var msg sync.Message
		err := jsonutil.Unmarshal(message.Data, &msg)
		if err != nil {
			log.Error(err, "client unmarshal data error")
		} else if msg.Code != contract.Success {
			log.Error(errors.New(msg.Message), "remote monitor received the error message")
		} else if ignore.MatchPath(msg.Path, "remote client monitor", msg.Action.String()) {
			// ignore match
		} else {
			m.execSync(msg)
		}
		m.messages.Remove(element)
	}
}

// execSync execute the file change message to sync
func (m *remoteClientMonitor) execSync(msg sync.Message) (err error) {
	values := url.Values{}
	values.Add(contract.FsDir, msg.IsDir.String())
	values.Add(contract.FsSize, stringutil.String(msg.Size))
	values.Add(contract.FsHash, msg.Hash)
	values.Add(contract.FsCtime, stringutil.String(msg.CTime))
	values.Add(contract.FsAtime, stringutil.String(msg.ATime))
	values.Add(contract.FsMtime, stringutil.String(msg.MTime))
	if len(msg.HashValues) > 0 {
		values.Add(contract.FsHashValues, stringutil.String(msg.HashValues))
	}

	// replace question marks with "%3F" to avoid parse the path is breaking when it contains some question marks
	path := msg.BaseUrl + strings.ReplaceAll(msg.Path, "?", "%3F") + fmt.Sprintf("?%s", values.Encode())

	switch msg.Action {
	case action.CreateAction:
		err = m.syncer.Create(path)
	case action.WriteAction:
		err = m.syncer.Create(path)
		// ignore is not exist error
		if err != nil && os.IsNotExist(err) {
			err = nil
		}
		m.addWrite(path)
	case action.RemoveAction:
		m.removeWrite(path)
		err = m.syncer.Remove(path)
	case action.RenameAction:
		err = m.syncer.Rename(path)
	case action.ChmodAction:
		err = m.syncer.Chmod(path)
	}

	m.el.Write(eventlog.NewEvent(path, msg.Action.String()))

	if err != nil {
		log.Error(err, "%s action execute error => [%s]", msg.Action.String(), path)
	}
	return err
}

// Close mark the monitor is closed, then close the connection
func (m *remoteClientMonitor) Close() error {
	m.closed.Set(true)
	if m.client != nil {
		return m.client.Close()
	}
	return nil
}
