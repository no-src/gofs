package sync

import (
	"github.com/no-src/gofs/driver/sftp"
)

type sftpPushClientSync struct {
	driverPushClientSync

	remoteAddr string
}

// NewSftpPushClientSync create an instance of the sftpPushClientSync
func NewSftpPushClientSync(opt Option) (Sync, error) {
	// the fields of option
	dest := opt.Dest
	users := opt.Users
	chunkSize := opt.ChunkSize
	maxTranRate := opt.MaxTranRate
	r := opt.Retry

	if chunkSize <= 0 {
		return nil, errInvalidChunkSize
	}

	if len(users) == 0 {
		return nil, errUserIsRequired
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
		remoteAddr: dest.Addr(),
	}

	currentUser := users[0]
	s.driver = sftp.NewSFTPDriver(s.remoteAddr, currentUser.UserName(), currentUser.Password(), opt.SSHKey, true, r, maxTranRate)

	err = s.start()
	if err != nil {
		return nil, err
	}
	return s, nil
}
