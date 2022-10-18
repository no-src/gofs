package sync

import (
	"errors"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/driver/sftp"
)

type sftpPullClientSync struct {
	driverPullClientSync

	remoteAddr  string
	remotePath  string
	currentUser *auth.User
}

// NewSftpPullClientSync create an instance of the sftpPullClientSync
func NewSftpPullClientSync(opt Option) (Sync, error) {
	// the fields of option
	source := opt.Source
	users := opt.Users
	chunkSize := opt.ChunkSize
	r := opt.Retry

	if chunkSize <= 0 {
		return nil, errors.New("chunk size must greater than zero")
	}

	if len(users) == 0 {
		return nil, errors.New("user account is required")
	}

	ds, err := newDiskSync(opt)
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
