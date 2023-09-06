package monitor

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/no-src/gofs/api/apiclient"
	"github.com/no-src/gofs/api/task"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/internal/clist"
	"github.com/no-src/gofs/internal/logger"
	"github.com/no-src/gofs/result"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/util/randutil"
	"github.com/no-src/gofs/wait"
)

type taskClientMonitor struct {
	shutdown chan struct{}
	retry    retry.Retry
	client   apiclient.Client
	closed   atomic.Bool
	runFn    func(content string, ext string) result.Result
	clientId string
	labels   []string
	tasks    sync.Map
	workers  chan struct{}
	queue    *clist.CList
	logger   *logger.Logger
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
	maxWorker := opt.TaskClientMaxWorker
	logger := opt.Logger

	var user *auth.User
	if len(users) > 0 {
		user = users[0]
	}
	m := &taskClientMonitor{
		shutdown: make(chan struct{}, 1),
		retry:    retry,
		client:   apiclient.New(host, port, enableTLS, certFile, user),
		runFn:    run,
		clientId: randutil.RandomString(10),
		labels:   labels,
		queue:    clist.New(),
		workers:  make(chan struct{}, maxWorker),
		logger:   logger,
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
	shutdown := &atomic.Bool{}
	go m.waitShutdown(shutdown, wd)
	go m.readMessage(shutdown, wd)
	go m.dequeue()
	return wd
}

// waitShutdown wait for the shutdown notify then mark the work done
func (m *taskClientMonitor) waitShutdown(st *atomic.Bool, wd wait.Done) {
	<-m.shutdown
	st.Store(true)
	m.logger.ErrorIf(m.Close(), "close remote client monitor error")
	wd.Done()
}

// readMessage loop read the messages, if receive a message, parse the message then send to consumers according to the api type.
// if receive a shutdown notify, then stop reading the message.
func (m *taskClientMonitor) readMessage(st *atomic.Bool, wd wait.Done) {
	clientInfo := &task.ClientInfo{
		ClientId: m.clientId,
		Labels:   m.labels,
	}
	rc, err := m.client.SubscribeTask(clientInfo)
	if err != nil {
		return
	}
	for {
		if m.closed.Load() {
			wd.DoneWithError(errors.New("remote task client is closed"))
			break
		}
		t, err := rc.Recv()
		if err != nil {
			if st.Load() {
				break
			}
			m.logger.Error(err, "subscribe task message error")
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
			m.enqueue(t)
		}
	}
}

// Close mark the monitor is closed, then close the connection
func (m *taskClientMonitor) Close() error {
	m.closed.Store(true)
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
	// broadcast shutdown notify
	close(m.shutdown)
	return err
}

func (m *taskClientMonitor) run(t *task.TaskInfo) {
	m.logger.Info("running gofs task [%s]", t.Name)
	r := m.runFn(t.Content, t.Ext)
	go func() {
		m.logger.ErrorIf(r.Wait(), "running gofs task error [%s]", t.Name)
	}()

	<-m.shutdown
	m.logger.ErrorIf(r.Shutdown(), "shutdown gofs task error [%s]", t.Name)

	m.tasks.Delete(t.Name)
	<-m.workers
	m.logger.Info("running gofs task finished [%s]", t.Name)
}

func (m *taskClientMonitor) enqueue(t *task.TaskInfo) {
	if t != nil {
		m.queue.PushBack(t)
	}
}

func (m *taskClientMonitor) dequeue() {
	for {
		if m.closed.Load() {
			break
		}
		e := m.queue.Front()
		if e == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		t := e.Value.(*task.TaskInfo)
		m.queue.Remove(e)

		if _, loaded := m.tasks.LoadOrStore(t.Name, t); loaded {
			m.logger.Info("[ignore task] task already exists => %s", t.Name)
			continue
		}
		m.workers <- struct{}{}
		go m.run(t)
	}
}
