package tran

type Client interface {
	Connect() error
	Write(data []byte) error
	ReadAll() (data []byte, err error)
	Host() string
	Port() int
	Close() error
	IsClosed() bool
}
