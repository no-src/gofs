package core

import (
	"github.com/no-src/log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

// VFS virtual file system
type VFS struct {
	path              string
	fsType            VFSType
	host              string
	port              int
	server            bool
	fsServer          string
	localSyncDisabled bool
}

const (
	paramPath                = "path"
	paramMode                = "mode"
	paramFsServer            = "fs_server"
	paramLocalSyncDisabled   = "local_sync_disabled"
	valueModeServer          = "server"
	valueLocalSyncIsDisabled = "true"
	remoteServerScheme       = "rs://"
	remoteServerDefaultPort  = 8105
)

// Path file path
func (vfs *VFS) Path() string {
	return vfs.path
}

// Abs returns an absolute representation of path
func (vfs *VFS) Abs() (string, error) {
	return filepath.Abs(vfs.Path())
}

// IsEmpty whether the file path is empty
func (vfs *VFS) IsEmpty() bool {
	return len(vfs.Path()) == 0
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

// LocalSyncDisabled is local disk sync disabled
func (vfs *VFS) LocalSyncDisabled() bool {
	return vfs.localSyncDisabled
}

// NewDiskVFS create an instance of VFS for the local disk file system
func NewDiskVFS(path string) VFS {
	vfs := VFS{
		fsType: Disk,
		path:   filepath.Clean(path),
	}
	return vfs
}

// NewEmptyVFS create an instance of VFS for the unknown file system
func NewEmptyVFS() VFS {
	vfs := VFS{
		fsType: Unknown,
	}
	return vfs
}

// NewVFS auto recognition the file system and create an instance of VFS according to the path
func NewVFS(path string) VFS {
	vfs := NewDiskVFS(path)
	lowerPath := strings.ToLower(path)
	var err error
	if strings.HasPrefix(lowerPath, remoteServerScheme) {
		// example of rs protocol to see README.md
		vfs.fsType = RemoteDisk
		_, vfs.host, vfs.port, vfs.path, vfs.server, vfs.fsServer, vfs.localSyncDisabled, err = parse(path)
	}
	if err != nil {
		return NewEmptyVFS()
	}
	return vfs
}

func parse(path string) (scheme string, host string, port int, localPath string, isServer bool, fsServer string, localSyncDisabled bool, err error) {
	parseUrl, err := url.Parse(path)
	if err != nil {
		return
	}
	scheme = parseUrl.Scheme
	host = parseUrl.Hostname()
	port, err = strconv.Atoi(parseUrl.Port())
	if err != nil {
		port = remoteServerDefaultPort
		err = nil
		log.Info("no remote server source port is specified, use default port => %d", port)
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
			fsServer = "https://" + fsServer
		}
	}

	localSyncDisabledValue := parseUrl.Query().Get(paramLocalSyncDisabled)
	if strings.ToLower(localSyncDisabledValue) == valueLocalSyncIsDisabled {
		localSyncDisabled = true
	}
	return
}
