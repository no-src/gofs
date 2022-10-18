package sync

import (
	"errors"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/driver/sftp"
)

type sftpPushClientSync struct {
	driverPushClientSync

	remoteAddr  string
	remotePath  string
	currentUser *auth.User
}

// NewSftpPushClientSync create an instance of the sftpPushClientSync
func NewSftpPushClientSync(opt Option) (Sync, error) {
	// the fields of option
	dest := opt.Dest
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
