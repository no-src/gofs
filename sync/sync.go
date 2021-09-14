package sync

type Sync interface {
	Create(path string) error
	Write(path string) error
	Remove(path string) error
	Rename(path string) error
	Chmod(path string) error
	IsDir(path string) (bool, error)
	SyncOnce() error
}
