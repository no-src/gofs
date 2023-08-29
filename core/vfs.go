package core

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/no-src/log"
)

// VFS virtual file system
type VFS struct {
	original          string
	path              string
	remotePath        string
	fsType            VFSType
	host              string
	port              int
	server            bool
	fsServer          string
	localSyncDisabled bool
	secure            bool
	sshConf           SSHConfig
}

const (
	paramPath               = "path"
	paramRemotePath         = "remote_path"
	paramMode               = "mode"
	paramFsServer           = "fs_server"
	paramLocalSyncDisabled  = "local_sync_disabled"
	paramSecure             = "secure"
	paramSSHUsername        = "ssh_user"
	paramSSHPassword        = "ssh_pass"
	paramSSHKey             = "ssh_key"
	paramSSHKeyPassphrase   = "ssh_key_pass"
	paramSSHHostKey         = "ssh_host_key"
	valueModeServer         = "server"
	valueTrue               = "true"
	schemeDelimiter         = "://"
	remoteServerScheme      = "rs"
	remoteServerDefaultPort = 8105
	sftpServerScheme        = "sftp"
	sftpServerDefaultPort   = 22
	minIOServerScheme       = "minio"
	minIOServerDefaultPort  = 9000
)

// Path the local file path
func (vfs *VFS) Path() string {
	return vfs.path
}

// RemotePath the remote file path
func (vfs *VFS) RemotePath() string {
	return vfs.remotePath
}

// Abs returns an absolute representation of Path
func (vfs *VFS) Abs() (string, error) {
	return filepath.Abs(vfs.Path())
}

// IsEmpty whether the local file path is empty
func (vfs *VFS) IsEmpty() bool {
	return len(vfs.Path()) == 0
}

// Type file system type
func (vfs *VFS) Type() VFSType {
	return vfs.fsType
}

// Host returns the server host
func (vfs *VFS) Host() string {
	return vfs.host
}

// Port returns the server port
func (vfs *VFS) Port() int {
	return vfs.port
}

// Addr returns the server address
func (vfs *VFS) Addr() string {
	return fmt.Sprintf("%s:%d", vfs.Host(), vfs.Port())
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

// Secure use secure connection
func (vfs *VFS) Secure() bool {
	return vfs.secure
}

// SSHConfig returns the SSH config
func (vfs *VFS) SSHConfig() SSHConfig {
	return vfs.sshConf
}

// NewDiskVFS create an instance of VFS for the local disk file system
func NewDiskVFS(path string) VFS {
	vfs := VFS{
		fsType:   Disk,
		path:     filepath.Clean(path),
		original: path,
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
	if strings.HasPrefix(lowerPath, remoteServerScheme+schemeDelimiter) {
		// example of rs protocol to see README.md
		vfs.fsType = RemoteDisk
		_, vfs.host, vfs.port, vfs.path, vfs.remotePath, vfs.server, vfs.fsServer, vfs.localSyncDisabled, vfs.secure, _, err = parse(path)
	} else if strings.HasPrefix(lowerPath, sftpServerScheme+schemeDelimiter) {
		vfs.fsType = SFTP
		_, vfs.host, vfs.port, vfs.path, vfs.remotePath, vfs.server, vfs.fsServer, vfs.localSyncDisabled, vfs.secure, vfs.sshConf, err = parse(path)
	} else if strings.HasPrefix(lowerPath, minIOServerScheme+schemeDelimiter) {
		vfs.fsType = MinIO
		_, vfs.host, vfs.port, vfs.path, vfs.remotePath, vfs.server, vfs.fsServer, vfs.localSyncDisabled, vfs.secure, _, err = parse(path)
	}
	if err != nil {
		return NewEmptyVFS()
	}
	return vfs
}

func parse(path string) (scheme string, host string, port int, localPath string, remotePath string, isServer bool, fsServer string, localSyncDisabled bool, secure bool, sshConf SSHConfig, err error) {
	parseUrl, err := url.Parse(path)
	if err != nil {
		return
	}
	scheme = parseUrl.Scheme
	host = parseUrl.Hostname()
	port, err = strconv.Atoi(parseUrl.Port())
	if err != nil {
		if scheme == remoteServerScheme {
			port = remoteServerDefaultPort
			err = nil
			log.Info("no remote server source port is specified, use default port => %d", port)
		} else if scheme == sftpServerScheme {
			port = sftpServerDefaultPort
			err = nil
			log.Info("no sftp server destination port is specified, use default port => %d", port)
		} else if scheme == minIOServerScheme {
			port = minIOServerDefaultPort
			err = nil
			log.Info("no MinIO server destination port is specified, use default port => %d", port)
		}
	}

	localPath = filepath.Clean(parseUrl.Query().Get(paramPath))
	// maybe the remote os is different from the current os, force convert remote path to slash
	remotePath = filepath.ToSlash(filepath.Clean(parseUrl.Query().Get(paramRemotePath)))

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
	if strings.ToLower(localSyncDisabledValue) == valueTrue {
		localSyncDisabled = true
	}

	isSecure := parseUrl.Query().Get(paramSecure)
	if strings.ToLower(isSecure) == valueTrue {
		secure = true
	}

	// parse SSH config
	sshConf.Username = strings.TrimSpace(parseUrl.Query().Get(paramSSHUsername))
	sshConf.Password = strings.TrimSpace(parseUrl.Query().Get(paramSSHPassword))
	sshConf.Key = strings.TrimSpace(parseUrl.Query().Get(paramSSHKey))
	sshConf.KeyPass = strings.TrimSpace(parseUrl.Query().Get(paramSSHKeyPassphrase))
	sshConf.HostKey = strings.TrimSpace(parseUrl.Query().Get(paramSSHHostKey))
	return
}
