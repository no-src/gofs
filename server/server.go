package server

import (
	"github.com/no-src/log"
	"net"
)

var serverAddr *net.TCPAddr

const (
	SrcRoutePrefix    = "/src/"
	TargetRoutePrefix = "/target/"
	QueryRoute        = "/query"
)

func initServerAddr(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err == nil {
		serverAddr = tcpAddr
	} else {
		log.Error(err, "invalid server addr => %s", addr)
	}
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
