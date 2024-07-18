package core

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/no-src/gofs/logger"
)

// VFS virtual file system
type VFS struct {
	original          string
	path              Path
	remotePath        Path
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
	paramSSHConfig          = "ssh_config"
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
func (vfs *VFS) Path() Path {
	return vfs.path
}

// RemotePath the remote file path
func (vfs *VFS) RemotePath() Path {
	return vfs.remotePath
}

// Abs returns an absolute representation of Path
func (vfs *VFS) Abs() (string, error) {
	return filepath.Abs(vfs.Path().Base())
}

// IsEmpty whether the local file path is empty
func (vfs *VFS) IsEmpty() bool {
	return len(vfs.Path().String()) == 0
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
		path:     newPath(filepath.Clean(path), Disk),
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
		_, vfs.host, vfs.port, vfs.path, vfs.remotePath, vfs.server, vfs.fsServer, vfs.localSyncDisabled, vfs.secure, _, err = parse(path, vfs.fsType)
	} else if strings.HasPrefix(lowerPath, sftpServerScheme+schemeDelimiter) {
		vfs.fsType = SFTP
		_, vfs.host, vfs.port, vfs.path, vfs.remotePath, vfs.server, vfs.fsServer, vfs.localSyncDisabled, vfs.secure, vfs.sshConf, err = parse(path, vfs.fsType)
	} else if strings.HasPrefix(lowerPath, minIOServerScheme+schemeDelimiter) {
		vfs.fsType = MinIO
		_, vfs.host, vfs.port, vfs.path, vfs.remotePath, vfs.server, vfs.fsServer, vfs.localSyncDisabled, vfs.secure, _, err = parse(path, vfs.fsType)
	}
	if err != nil {
		return NewEmptyVFS()
	}
	return vfs
}

func parse(path string, fsType VFSType) (scheme string, host string, port int, localPath Path, remotePath Path, isServer bool, fsServer string, localSyncDisabled bool, secure bool, sshConf SSHConfig, err error) {
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
			logger.InnerLogger().Info("no remote server source port is specified, use default port => %d", port)
		} else if scheme == sftpServerScheme {
			port = sftpServerDefaultPort
			err = nil
			logger.InnerLogger().Info("no sftp server destination port is specified, use default port => %d", port)
		} else if scheme == minIOServerScheme {
			port = minIOServerDefaultPort
			err = nil
			logger.InnerLogger().Info("no MinIO server destination port is specified, use default port => %d", port)
		}
	}

	localPath = newPath(parseUrl.Query().Get(paramPath), Disk)
	remotePath = newPath(parseUrl.Query().Get(paramRemotePath), fsType)

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

	// process ssh_config file
	if fsType == SFTP && strings.TrimSpace(parseUrl.Query().Get(paramSSHConfig)) == valueTrue {
		sshConfigHost := host
		sshConfigHostName := ssh_config.Get(sshConfigHost, "HostName")
		if len(sshConfigHostName) > 0 {
			host = sshConfigHostName
			if sshConfigPort, portErr := strconv.Atoi(ssh_config.Get(sshConfigHost, "Port")); portErr == nil {
				if port != sshConfigPort {
					logger.InnerLogger().Info("the port of the SFTP server has been changed from %d to %d.", port, sshConfigPort)
					port = sshConfigPort
				}
			}
			sshConf.Username = ssh_config.Get(sshConfigHost, "User")
			sshConf.Key = findSSHKey(sshConfigHost)
		}
	}

	// parse SSH config
	usernameFromUrl := strings.TrimSpace(parseUrl.Query().Get(paramSSHUsername))
	if len(usernameFromUrl) > 0 {
		sshConf.Username = usernameFromUrl
	}
	sshConf.Password = strings.TrimSpace(parseUrl.Query().Get(paramSSHPassword))
	keyFromUrl := strings.TrimSpace(parseUrl.Query().Get(paramSSHKey))
	if len(keyFromUrl) > 0 {
		sshConf.Key = keyFromUrl
	}
	sshConf.KeyPass = strings.TrimSpace(parseUrl.Query().Get(paramSSHKeyPassphrase))
	sshConf.HostKey = strings.TrimSpace(parseUrl.Query().Get(paramSSHHostKey))

	if fsType == SFTP {
		logger.InnerLogger().Info("ssh config info: Host: %s Port: %d Username: %s Key: %s", host, port, sshConf.Username, sshConf.Key)
	}
	return
}

// findSSHKey find SSH IdentityFile based on the ssh config and usual locations
func findSSHKey(sshConfigHost string) string {
	possibleKeys := ssh_config.GetAll(sshConfigHost, "IdentityFile")
	firstExistingKey := ""
	for _, keyFile := range possibleKeys {
		keyFile = parsePath(keyFile)
		if _, err := os.Stat(keyFile); err == nil {
			firstExistingKey = keyFile
			break
		} else {
			logger.InnerLogger().Warn("SSH IdentityFile %s does not exist", keyFile)
		}
	}
	if len(firstExistingKey) > 0 {
		return firstExistingKey
	}

	if len(possibleKeys) == 1 && possibleKeys[0] == ssh_config.Default("IdentityFile") {
		home, _ := os.UserHomeDir()
		fallbackDefault := filepath.Join(home, ".ssh", "id_rsa")
		if _, err := os.Stat(fallbackDefault); err == nil {
			return fallbackDefault
		}
	}
	return ""
}

func parsePath(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~/"))
}
