package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/driver/minio"
)

type minIOPushClientSync struct {
	driverPushClientSync

	endpoint    string
	bucketName  string
	secure      bool
	currentUser *auth.User
}

// NewMinIOPushClientSync create an instance of the minIOPushClientSync
func NewMinIOPushClientSync(opt Option) (Sync, error) {
	// the fields of option
	dest := opt.Dest
	users := opt.Users
	chunkSize := opt.ChunkSize
	maxTranRate := opt.MaxTranRate
	r := opt.Retry
	logger := opt.Logger

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

	s := &minIOPushClientSync{
		driverPushClientSync: newDriverPushClientSync(*ds, ""),
		endpoint:             dest.Addr(),
		bucketName:           dest.RemotePath().Bucket(),
		secure:               dest.Secure(),
		currentUser:          users[0],
	}

	s.driver = minio.NewMinIODriver(s.endpoint, s.bucketName, s.secure, s.currentUser.UserName(), s.currentUser.Password(), true, r, maxTranRate, logger)

	err = s.start()
	if err != nil {
		return nil, err
	}
	return s, nil
}
