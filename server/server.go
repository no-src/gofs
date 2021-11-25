package server

import (
	"fmt"
	"github.com/no-src/log"
	"net"
	"strings"
)

var serverAddr *net.TCPAddr
var enableTLS bool

const (
	SrcRoutePrefix    = "/src/"
	TargetRoutePrefix = "/target/"
	QueryRoute        = "/query"
)

const (
	DefaultAddrHttps = ":443"
	DefaultAddrHttp  = ":80"
	ProtocolHttp     = "http"
	ProtocolHttps    = "https"
	DefaultPortHttp  = 80
	DefaultPortHttps = 443
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

// GenerateAddr generate http or https address
func GenerateAddr(protocol, host string, port int) string {
	addr := ""
	protocol = strings.ToLower(protocol)
	if protocol == ProtocolHttp && port == DefaultPortHttp {
		addr = fmt.Sprintf("%s://%s", ProtocolHttp, host)
	} else if protocol == ProtocolHttps && port == DefaultPortHttps {
		addr = fmt.Sprintf("%s://%s", ProtocolHttps, host)
	} else {
		addr = fmt.Sprintf("%s://%s:%d", protocol, host, port)
	}
	return addr
}
