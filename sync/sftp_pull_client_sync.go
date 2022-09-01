package sync

import (
	"errors"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/driver/sftp"
	"github.com/no-src/gofs/retry"
)

type sftpPullClientSync struct {
	driverPullClientSync

	remoteAddr  string
	remotePath  string
	currentUser *auth.User
}

// NewSftpPullClientSync create an instance of the sftpPullClientSync
func NewSftpPullClientSync(source, dest core.VFS, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool, r retry.Retry) (Sync, error) {
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

	s := &sftpPullClientSync{
		driverPullClientSync: driverPullClientSync{
			diskSync: *ds,
		},
		remoteAddr:  source.Addr(),
		remotePath:  source.RemotePath(),
		currentUser: users[0],
	}

	s.client = sftp.NewSFTPClient(s.remoteAddr, s.currentUser.UserName(), s.currentUser.Password(), true, r)

	err = s.start()
	if err != nil {
		return nil, err
	}

	// reset the sourceAbsPath because the source.Path() or source.RemotePath() is absolute representation of path and the source.Path() or source.RemotePath() may be cross-platform
	// source.Path() and source.RemotePath() are equivalent here, and source.RemotePath() has higher priority
	s.diskSync.sourceAbsPath = source.RemotePath()
	if len(s.diskSync.sourceAbsPath) == 0 {
		s.diskSync.sourceAbsPath = source.Path()
	}
	// reset some functions for sftp
	s.diskSync.isDirFn = s.IsDir
	s.diskSync.statFn = s.client.Stat
	s.diskSync.getFileTimeFn = s.client.GetFileTime

	return s, nil
}
