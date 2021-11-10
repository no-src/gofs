package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/tran"
	"github.com/no-src/log"
	"strings"
)

type remoteMonitor struct {
	syncer   sync.Sync
	retry    retry.Retry
	client   tran.Client
	closed   bool
	messages chan message
}

type message struct {
	data []byte
}

func NewRemoteMonitor(syncer sync.Sync, retry retry.Retry, host string, port int, messageQueue int) (m Monitor, err error) {
	if syncer == nil {
		err = errors.New("syncer can't be nil")
		return nil, err
	}
	m = &remoteMonitor{
		syncer:   syncer,
		retry:    retry,
		client:   tran.NewClient(host, port),
		messages: make(chan message, messageQueue),
	}
	return m, nil
}

func (m *remoteMonitor) Start() error {
	if m.client == nil {
		return errors.New("remote sync client is nil")
	}
	err := m.client.Connect()
	if err != nil {
		return err
	}
	go m.processingMessage()
	for {
		if m.closed {
			return errors.New("remote monitor is closed")
		}
		data, err := m.client.ReadAll()
		if err != nil {
			log.Error(err, "client read data error")
			if m.client.IsClosed() {
				log.Debug("try reconnect to server %s:%d", m.client.Host(), m.client.Port())
				m.retry.Do(func() error {
					return m.client.Connect()
				}, fmt.Sprintf("client reconnect to %s:%d", m.client.Host(), m.client.Port()))
			}
		} else {
			m.messages <- message{
				data: data,
			}
		}
	}
	return nil
}

func (m *remoteMonitor) processingMessage() {
	for {
		message := <-m.messages
		log.Info("client read request => %s", string(message.data))
		var req sync.Request
		err := json.Unmarshal(message.data, &req)
		if err != nil {
			log.Error(err, "client unmarshal data error")
		} else {
			// append is dir, 1 or 0,-1 mean unknown
			// replace question marks with "%3F" to avoid parse the path is breaking when it contains some question marks
			path := req.BaseUrl + strings.ReplaceAll(req.Path, "?", "%3F") + fmt.Sprintf("?dir=%d", req.IsDir)
			// append file size, bytes
			path += fmt.Sprintf("&size=%d", req.Size)
			// append file hash
			path += fmt.Sprintf("&hash=%s", req.Hash)

			switch req.Action {
			case sync.CreateAction:
				m.syncer.Create(path)
				break
			case sync.WriteAction:
				m.syncer.Write(path)
				break
			case sync.RemoveAction:
				m.syncer.Remove(path)
				break
			case sync.RenameAction:
				m.syncer.Rename(path)
				break
			case sync.ChmodAction:
				m.syncer.Chmod(path)
				break
			}
		}
	}
}

func (m *remoteMonitor) Close() error {
	m.closed = true
	return m.client.Close()
}
