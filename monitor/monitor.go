package monitor

type Monitor interface {
	Monitor(dir string) error
	Start() error
	Close() error
}
