package apiserver

import "github.com/no-src/gofs/api/monitor"

// Server the gofs api server
type Server interface {
	// Start running the server
	Start() error
	// Stop stop the server
	Stop()
	// SendMonitorMessage send monitor message
	SendMonitorMessage(message *monitor.MonitorMessage) error
}
