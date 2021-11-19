package core

import (
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

// VFS virtual file system
type VFS struct {
	path     string
	fsType   VFSType
	host     string
	port     int
	server   bool
	fsServer string
	msgQueue int
}

const (
	paramPath            = "path"
	paramMode            = "mode"
	paramFsServer        = "fs_server"
	paramMsgQueue        = "msg_queue"
	valueModeServer      = "server"
	valueDefaultMsgQueue = 500
	RemoteServerSchema   = "rs://"
)

// Path file path
func (vfs *VFS) Path() string {
	return vfs.path
}

// Type file system type
func (vfs *VFS) Type() VFSType {
	return vfs.fsType
}

// Host server or client host
func (vfs *VFS) Host() string {
	return vfs.host
}

// Port server or client port
func (vfs *VFS) Port() int {
	return vfs.port
}

// IsDisk is local file system
func (vfs *VFS) IsDisk() bool {
	return vfs.Is(Disk)
}

// Is current VFS is type of t
func (vfs *VFS) Is(t VFSType) bool {
	return vfs.fsType == t
}

// Server is server mode
func (vfs *VFS) Server() bool {
	return vfs.server
}

// FsServer file server access addr
func (vfs *VFS) FsServer() string {
	return vfs.fsServer
}

// MessageQueue receive message queue size
func (vfs *VFS) MessageQueue() int {
	return vfs.msgQueue
}

func NewDiskVFS(path string) VFS {
	vfs := VFS{
		fsType: Disk,
		path:   filepath.Clean(path),
	}
	return vfs
}

func NewEmptyVFS() VFS {
	vfs := VFS{
		fsType: Unknown,
	}
	return vfs
}

func NewVFS(path string) VFS {
	vfs := NewDiskVFS(path)
	lowerPath := strings.ToLower(path)
	var err error
	if strings.HasPrefix(lowerPath, RemoteServerSchema) {
		// example of rs protocol to see README.md
		vfs.fsType = RemoteDisk
		_, vfs.host, vfs.port, vfs.path, vfs.server, vfs.fsServer, vfs.msgQueue, err = parse(path)
	}
	if err != nil {
		return NewEmptyVFS()
	}
	return vfs
}

func parse(path string) (scheme string, host string, port int, localPath string, isServer bool, fsServer string, msgQueue int, err error) {
	parseUrl, err := url.Parse(path)
	if err != nil {
		return
	}
	scheme = parseUrl.Scheme
	host = parseUrl.Hostname()
	port, err = strconv.Atoi(parseUrl.Port())
	if err != nil {
		return
	}
	localPath = filepath.Clean(parseUrl.Query().Get(paramPath))
	mode := parseUrl.Query().Get(paramMode)
	if strings.ToLower(mode) == valueModeServer {
		isServer = true
	}
	fsServer = parseUrl.Query().Get(paramFsServer)
	if len(fsServer) > 0 {
		fsServerLower := strings.ToLower(fsServer)
		if !strings.HasPrefix(fsServerLower, "http://") && !strings.HasPrefix(fsServerLower, "https://") {
			fsServer = "http://" + fsServer
		}
	}

	defaultMsgQueue := valueDefaultMsgQueue
	msgQueueStr := parseUrl.Query().Get(paramMsgQueue)
	if len(msgQueueStr) > 0 {
		msgQueue, err = strconv.Atoi(msgQueueStr)
		if err != nil || msgQueue <= 0 {
			// default is 500 of message queue size
			msgQueue = defaultMsgQueue
		}
	} else {
		msgQueue = defaultMsgQueue
	}
	return
}
