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
	chunkSize := opt.ChunkSize
	maxTranRate := opt.MaxTranRate
	r := opt.Retry
	logger := opt.Logger
	syncOnce := opt.SyncOnce
	syncCron := opt.SyncCron

	if chunkSize <= 0 {
		return nil, errInvalidChunkSize
	}

	ds, err := newDiskSync(opt)
	if err != nil {
		return nil, err
	}

	s := &sftpPushClientSync{
		driverPushClientSync: newDriverPushClientSync(*ds, dest.RemotePath().Base()),
		remoteAddr:           dest.Addr(),
	}

	s.driver = sftp.NewSFTPDriver(s.remoteAddr, dest.SSHConfig(), true, r, maxTranRate, logger)

	isSync := syncOnce || len(syncCron) > 0
	err = s.start(isSync)
	if err != nil {
		return nil, err
	}
	return s, nil
}
