package monitor

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/contract"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/tran"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"net/url"
	"os"
	"strings"
	"time"
)

type remoteClientMonitor struct {
	baseMonitor
	client      tran.Client
	closed      bool
	messages    *list.List
	syncOnce    bool
	currentUser *auth.HashUser
	authorized  bool
}

type message struct {
	data []byte
}

// NewRemoteClientMonitor create an instance of remoteClientMonitor to monitor the remote file change
func NewRemoteClientMonitor(syncer sync.Sync, retry retry.Retry, syncOnce bool, host string, port int, enableTLS bool, users []*auth.User) (Monitor, error) {
	if syncer == nil {
		err := errors.New("syncer can't be nil")
		return nil, err
	}
	m := &remoteClientMonitor{
		client:      tran.NewClient(host, port, enableTLS),
		messages:    list.New(),
		syncOnce:    syncOnce,
		baseMonitor: newBaseMonitor(syncer, retry),
	}
	if len(users) > 0 {
		user := users[0]
		hashUser, err := user.ToHashUser()
		if err != nil {
			return nil, err
		}
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
		authData := auth.GenerateAuthCommandData(m.currentUser.UserNameHash, m.currentUser.PasswordHash)
		err := m.client.Write(authData)
		return err
	}, "send auth request")

	var status contract.Status
	m.retry.Do(func() error {
		result, err := m.client.ReadAll()
		if err != nil {
			return err
		}
		err = util.Unmarshal(result, &status)
		if err != nil {
			return err
		}

		if status.ApiType != contract.AuthApi {
			err = errors.New("auth command is received other message")
			log.Error(err, "retry to get auth response error")
			return err
		}

		return nil
	}, "receive auth command response").Wait()

	if status.Code != contract.Success {
		return errors.New("receive auth command response error => " + status.Message)
	}

	if status.ApiType != contract.AuthApi {
		return errors.New("auth command is received other message")
	}
	// auth success
	m.authorized = true
	log.Debug("auth success, current client is authorized")
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

	if err = m.auth(); err != nil {
		return err
	}

	// check sync once command
	if m.syncOnce {
		go func() {
			err = m.client.Write(contract.InfoCommand)
			if err != nil {
				log.Error(err, "write info command error")
			}
		}()

		var info contract.FileServerInfo
		m.retry.Do(func() error {
			infoResult, err := m.client.ReadAll()
			if err != nil {
				return err
			}
			err = util.Unmarshal(infoResult, &info)
			if err != nil {
				return err
			}

			if info.ApiType != contract.InfoApi {
				err = errors.New("info command is received other message")
				log.Error(err, "retry to get info response error")
				return err
			}

			return nil
		}, "receive info command response").Wait()

		if info.Code != contract.Success {
			return errors.New("receive info command response error => " + info.Message)
		}

		if info.ApiType != contract.InfoApi {
			return errors.New("info command is received other message")
		}
		return m.syncer.SyncOnce(info.ServerAddr + info.SrcPath)
	}
	go m.processWrite()
	go m.startSyncWrite()
	go m.processingMessage()
	for {
		if m.closed {
			return errors.New("remote monitor is closed")
		}
		data, err := m.client.ReadAll()
		if err != nil {
			log.Error(err, "client read data error")
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
					if !m.authorized {
						return m.auth()
					}
					return nil
				}, fmt.Sprintf("client reconnect to %s:%d", m.client.Host(), m.client.Port()))
			}
		} else {
			m.messages.PushBack(message{
				data: data,
			})
		}
	}
	return nil
}

func (m *remoteClientMonitor) processingMessage() {
	for {
		element := m.messages.Front()
		if element == nil || element.Value == nil {
			if element != nil {
				m.messages.Remove(element)
			}
			<-time.After(time.Second)
			continue
		}
		message := element.Value.(message)
		log.Info("client read request => %s", string(message.data))
		var msg sync.Message
		err := util.Unmarshal(message.data, &msg)
		if err != nil {
			log.Error(err, "client unmarshal data error")
		} else if msg.Code != contract.Success {
			log.Error(errors.New(msg.Message), "remote monitor received the error message")
		} else if msg.ApiType != contract.SyncMessageApi {
			log.Debug("received other message, ignore it => %s", string(message.data))
		} else {
			values := url.Values{}
			values.Add(contract.FsDir, msg.IsDir.String())
			values.Add(contract.FsSize, util.String(msg.Size))
			values.Add(contract.FsHash, msg.Hash)
			values.Add(contract.FsCtime, util.String(msg.CTime))
			values.Add(contract.FsAtime, util.String(msg.ATime))
			values.Add(contract.FsMtime, util.String(msg.MTime))

			// replace question marks with "%3F" to avoid parse the path is breaking when it contains some question marks
			path := msg.BaseUrl + strings.ReplaceAll(msg.Path, "?", "%3F") + fmt.Sprintf("?%s", values.Encode())

			switch msg.Action {
			case sync.CreateAction:
				err = m.syncer.Create(path)
				break
			case sync.WriteAction:
				err = m.syncer.Create(path)
				// ignore is not exist error
				if err != nil && os.IsNotExist(err) {
					err = nil
				}
				m.addWrite(path)
				break
			case sync.RemoveAction:
				m.removeWrite(path)
				err = m.syncer.Remove(path)
				break
			case sync.RenameAction:
				err = m.syncer.Rename(path)
				break
			case sync.ChmodAction:
				err = m.syncer.Chmod(path)
				break
			}
			if err != nil {
				log.Error(err, "%s action execute error => [%s]", msg.Action.String(), path)
			}
		}
		m.messages.Remove(element)
	}
}

func (m *remoteClientMonitor) Close() error {
	m.closed = true
	if m.client != nil {
		return m.client.Close()
	}
	return nil
}
