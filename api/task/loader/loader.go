package loader

import (
	"errors"
	"strings"
)

const (
	filePrefix = "file://"
	confKey    = "nosrc-gofs-task-conf"
)

var (
	errNilTaskConfig = errors.New("task config is nil")
	errDuplicateTask = errors.New("duplicate task")
)

// Loader a task config loader interface
type Loader interface {
	// LoadConfig load the task config
	LoadConfig() (*TaskConfig, error)
	// LoadContent load the specified task item content
	LoadContent(conf string) (string, error)
	// SaveConfig save the task config
	SaveConfig(c *TaskConfig) error
	// SaveContent save the specified task item content
	SaveContent(conf string, content string) error
	// Close release the dependent resource
	Close() error
}

// NewLoader return a task config loader instance and currently support file, memory, redis, buntdb memory, buntdb, etcd.
// Examples like the following:
// file://task.yaml
// memory:
// redis://127.0.0.1:6379
// buntdb://:memory: or buntdb://buntdb.db
// etcd://127.0.0.1:2379?dial_timeout=5s
func NewLoader(path string) (Loader, error) {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return newEmptyLoader(), nil
	}
	if strings.HasPrefix(path, filePrefix) {
		return newFileLoader(strings.TrimPrefix(path, filePrefix)), nil
	}
	return newCacheLoader(path, confKey)
}
