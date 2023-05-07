package monitor

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/no-src/gofs/api/apiclient"
	"github.com/no-src/gofs/api/task"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/internal/cbool"
	"github.com/no-src/gofs/result"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/util/randutil"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

type taskClientMonitor struct {
	shutdown chan struct{}
	retry    retry.Retry
	client   apiclient.Client
	closed   *cbool.CBool
	runFn    func(content string, ext string) result.Result
	clientId string
	labels   []string
	tasks    sync.Map
}

// NewTaskClientMonitor create an instance of taskClientMonitor to receive the task from remote server
func NewTaskClientMonitor(opt Option, run runFn) (Monitor, error) {
	source := opt.Syncer.Source()
	host := source.Host()
	port := source.Port()
	enableTLS := opt.EnableTLS
	certFile := opt.TLSCertFile
	users := opt.Users
	labels := opt.TaskClientLabels
	retry := opt.Retry

	var user *auth.User
	if len(users) > 0 {
		user = users[0]
	}
	m := &taskClientMonitor{
		shutdown: make(chan struct{}, 1),
		retry:    retry,
		client:   apiclient.New(host, port, enableTLS, certFile, user),
		closed:   cbool.New(false),
		runFn:    run,
		clientId: randutil.RandomString(10),
		labels:   labels,
	}
	return m, nil
}

func (m *taskClientMonitor) Start() (wait.Wait, error) {
	if m.client == nil {
		return nil, errors.New("remote task client is nil")
	}
	err := m.client.Start()
	if err != nil {
		return nil, err
	}
	return m.receive(), nil
}

// receive start receiving messages and parse the message, send to consumers according to the api type.
// if receive a shutdown notify, then stop reading the message.
func (m *taskClientMonitor) receive() wait.Wait {
	wd := wait.NewWaitDone()
	shutdown := cbool.New(false)
	go m.waitShutdown(shutdown, wd)
	go m.readMessage(shutdown, wd)
	return wd
}

// waitShutdown wait for the shutdown notify then mark the work done
func (m *taskClientMonitor) waitShutdown(st *cbool.CBool, wd wait.Done) {
	select {
	case <-m.shutdown:
		{
			st.Set(true)
			log.ErrorIf(m.Close(), "close remote client monitor error")
			wd.Done()
		}
	}
}

// readMessage loop read the messages, if receive a message, parse the message then send to consumers according to the api type.
// if receive a shutdown notify, then stop reading the message.
func (m *taskClientMonitor) readMessage(st *cbool.CBool, wd wait.Done) {
	clientInfo := &task.ClientInfo{
		ClientId: m.clientId,
		Labels:   m.labels,
	}
	rc, err := m.client.SubscribeTask(clientInfo)
	if err != nil {
		return
	}
	for {
		if m.closed.Get() {
			wd.DoneWithError(errors.New("remote task client is closed"))
			break
		}
		t, err := rc.Recv()
		if err != nil {
			if st.Get() {
				break
			}
			log.Error(err, "subscribe task message error")
			if m.client.IsClosed(err) {
				m.retry.Do(func() error {
					nrc, err := m.client.SubscribeTask(clientInfo)
					if err == nil {
						rc = nrc
					}
					return err
				}, "subscribe the task server")
			} else {
				wd.DoneWithError(fmt.Errorf("remote task server is return error %w", err))
				break
			}
		} else {
			m.run(*t)
		}
	}
}

// Close mark the monitor is closed, then close the connection
func (m *taskClientMonitor) Close() error {
	m.closed.Set(true)
	if m.client != nil {
		return m.client.Stop()
	}
	return nil
}

func (m *taskClientMonitor) SyncCron(spec string) error {
	spec = strings.TrimSpace(spec)
	if len(spec) == 0 {
		return nil
	}
	return errors.New("the usage of the -sync_cron flag is incompatible with enabling the -task_client flag")
}

func (m *taskClientMonitor) Shutdown() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	close(m.shutdown)
	return err
}

func (m *taskClientMonitor) run(t task.TaskInfo) {
	if _, loaded := m.tasks.LoadOrStore(t.Name, t); loaded {
		return
	}
	go func() {
		log.Info("running gofs task [%s]", t.Name)
		r := m.runFn(t.Content, t.Ext)
		done := make(chan struct{}, 1)
		go func() {
			log.ErrorIf(r.Wait(), "running gofs task error [%s]", t.Name)
		}()
		select {
		case <-m.shutdown:
			log.ErrorIf(r.Shutdown(), "shutdown gofs task error [%s]", t.Name)
		case <-done:
		}
		m.tasks.Delete(t.Name)
		log.Info("running gofs task finished [%s]", t.Name)
	}()
}
