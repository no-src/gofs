package tran

import "net"

type Server interface {
	Listen() error
	Accept(process func(client net.Conn, data []byte)) error
	ClientCount() int
	Send([]byte) error
	Host() string
	Port() int
	Close() error
}
