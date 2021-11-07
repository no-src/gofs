package core

import (
	"net/url"
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
}

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

func NewDiskVFS(path string) VFS {
	vfs := VFS{
		fsType: Disk,
		path:   path,
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
	if strings.HasPrefix(lowerPath, "rs://") {
		// rs://127.0.0.1:9016?mode=server&path=/var/source&fs_server=https://fs-server-domain.com
		vfs.fsType = RemoteDisk
		_, vfs.host, vfs.port, vfs.path, vfs.server, vfs.fsServer, err = parse(path)
	}
	if err != nil {
		return NewEmptyVFS()
	}
	return vfs
}

func parse(path string) (scheme string, host string, port int, localPath string, isServer bool, fsServer string, err error) {
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
	localPath = parseUrl.Query().Get("path")
	mode := parseUrl.Query().Get("mode")
	if strings.ToLower(mode) == "server" {
		isServer = true
	}
	fsServer = parseUrl.Query().Get("fs_server")
	if len(fsServer) > 0 {
		fsServerLower := strings.ToLower(fsServer)
		if !strings.HasPrefix(fsServerLower, "http://") && !strings.HasPrefix(fsServerLower, "https://") {
			fsServer = "http://" + fsServer
		}
	}
	return
}
