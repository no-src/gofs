package tran

import "net"

// Server a network communication server
type Server interface {
	// Listen listen the specified port to wait client connect
	Listen() error
	// Accept accept the client connection
	Accept(process func(client net.Conn, data []byte)) error
	// ClientCount the client count of the connected
	ClientCount() int
	// Send send the data to the client
	Send([]byte) error
	// Host return the server host
	Host() string
	// Port return the server bind port
	Port() int
	// Close close the server
	Close() error
}
