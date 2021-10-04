package monitor

import "github.com/no-src/gofs/core"

type Monitor interface {
	Monitor(vfs core.VFS) error
	Start() error
	Close() error
}
