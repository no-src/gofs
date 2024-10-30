package apiclient

import (
	"github.com/no-src/gofs/api/info"
	"github.com/no-src/gofs/api/monitor"
	"github.com/no-src/gofs/api/task"
)

// Client the gofs api client
type Client interface {
	// Start running the api client
	Start() error
	// Stop stop the client
	Stop() error
	// GetInfo get the file server info
	GetInfo() (*info.FileServerInfo, error)
	// Monitor monitor the remote server
	Monitor() (monitor.MonitorService_MonitorClient, error)
	// IsClosed is connection closed of the current client
	IsClosed(err error) bool
	// IsUnauthenticated check whether the error is unauthorized
	IsUnauthenticated(err error) bool
	// SubscribeTask register a task client to the task server and wait to receive task
	SubscribeTask(clientInfo *task.ClientInfo) (task.TaskService_SubscribeTaskClient, error)
	// Login login to the server
	Login() (err error)
}
