package tran

// Client a network communication client
type Client interface {
	// Connect connect to the server
	Connect() error
	// Write write data to the server
	Write(data []byte) error
	// ReadAll read all the data of the server
	ReadAll() (data []byte, err error)
	// Host return the client host
	Host() string
	// Port return the client bind port
	Port() int
	// Close close the client connection
	Close() error
	// IsClosed is connection closed of the current client
	IsClosed() bool
}
