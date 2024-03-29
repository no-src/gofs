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
	chunkSize := opt.ChunkSize
	maxTranRate := opt.MaxTranRate
	r := opt.Retry
	logger := opt.Logger

	if chunkSize <= 0 {
		return nil, errInvalidChunkSize
	}

	ds, err := newDiskSync(opt)
	if err != nil {
		return nil, err
	}

	s := &sftpPullClientSync{
		driverPullClientSync: newDriverPullClientSync(*ds),
		remoteAddr:           source.Addr(),
	}
	s.driver = sftp.NewSFTPDriver(s.remoteAddr, source.SSHConfig(), true, r, maxTranRate, logger)

	err = s.start()
	if err != nil {
		return nil, err
	}

	// reset the sourceAbsPath because the source.RemotePath() is absolute representation of path and the source.RemotePath() may be cross-platform
	s.diskSync.sourceAbsPath = source.RemotePath().Base()

	// reset some functions for sftp
	s.diskSync.isDirFn = s.IsDir
	s.diskSync.statFn = s.driver.Stat
	s.diskSync.getFileTimeFn = s.driver.GetFileTime

	return s, nil
}
