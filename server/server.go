package server

import (
	"github.com/no-src/log"
	"net"
)

var serverAddr *net.TCPAddr

func initServerAddr(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err == nil {
		serverAddr = tcpAddr
	} else {
		log.Error(err, "invalid server addr => %s", addr)
	}
}

func ServerAddr() *net.TCPAddr {
	return serverAddr
}

func ServerPort() int {
	if serverAddr != nil {
		return serverAddr.Port
	}
	return 0
}
