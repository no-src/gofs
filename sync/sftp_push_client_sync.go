package sync

import (
	"errors"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/driver/sftp"
	"github.com/no-src/gofs/retry"
)

type sftpPushClientSync struct {
	driverPushClientSync

	remoteAddr  string
	remotePath  string
	currentUser *auth.User
}

// NewSftpPushClientSync create an instance of the sftpPushClientSync
func NewSftpPushClientSync(source, dest core.VFS, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool, r retry.Retry) (Sync, error) {
	if chunkSize <= 0 {
		return nil, errors.New("chunk size must greater than zero")
	}

	if len(users) == 0 {
		return nil, errors.New("user account is required")
	}

	ds, err := newDiskSync(source, dest, enableLogicallyDelete, chunkSize, checkpointCount, forceChecksum)
	if err != nil {
		return nil, err
	}

	s := &sftpPushClientSync{
		driverPushClientSync: driverPushClientSync{
			diskSync: *ds,
			basePath: dest.RemotePath(),
		},
		remoteAddr:  dest.Addr(),
		remotePath:  dest.RemotePath(),
		currentUser: users[0],
	}

	s.client = sftp.NewSFTPClient(s.remoteAddr, s.currentUser.UserName(), s.currentUser.Password(), true, r)

	err = s.start()
	if err != nil {
		return nil, err
	}
	return s, nil
}
