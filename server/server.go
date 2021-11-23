package server

import (
	"github.com/no-src/log"
	"net"
)

var serverAddr *net.TCPAddr
var enableTLS bool

const (
	SrcRoutePrefix    = "/src/"
	TargetRoutePrefix = "/target/"
	QueryRoute        = "/query"
)

func initServerInfo(addr string, tls bool) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err == nil {
		serverAddr = tcpAddr
	} else {
		log.Error(err, "invalid server addr => %s", addr)
	}
	enableTLS = tls
}

// ServerAddr the addr of file server running
func ServerAddr() *net.TCPAddr {
	return serverAddr
}

// ServerPort the port of file server running
func ServerPort() int {
	if serverAddr != nil {
		return serverAddr.Port
	}
	return 0
}

// EnableTLS is using https on the file server
func EnableTLS() bool {
	return enableTLS
}
