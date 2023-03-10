package sync

import (
	"github.com/no-src/gofs/driver/sftp"
)

type sftpPullClientSync struct {
	driverPullClientSync

	remoteAddr string
}

// NewSftpPullClientSync create an instance of the sftpPullClientSync
func NewSftpPullClientSync(opt Option) (Sync, error) {
	// the fields of option
	source := opt.Source
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

	s := &sftpPullClientSync{
		driverPullClientSync: driverPullClientSync{
			diskSync: *ds,
		},
		remoteAddr: source.Addr(),
	}
	currentUser := users[0]
	s.driver = sftp.NewSFTPDriver(s.remoteAddr, currentUser.UserName(), currentUser.Password(), opt.SSHKey, true, r, maxTranRate)

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
	s.diskSync.statFn = s.driver.Stat
	s.diskSync.getFileTimeFn = s.driver.GetFileTime

	return s, nil
}
