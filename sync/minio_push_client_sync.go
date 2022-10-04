package sync

import (
	"errors"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/driver/minio"
	"github.com/no-src/gofs/encrypt"
	"github.com/no-src/gofs/retry"
)

type minIOPushClientSync struct {
	driverPushClientSync

	endpoint    string
	bucketName  string
	secure      bool
	currentUser *auth.User
}

// NewMinIOPushClientSync create an instance of the minIOPushClientSync
func NewMinIOPushClientSync(source, dest core.VFS, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool, r retry.Retry, encOpt encrypt.Option) (Sync, error) {
	if chunkSize <= 0 {
		return nil, errors.New("chunk size must greater than zero")
	}

	if len(users) == 0 {
		return nil, errors.New("user account is required")
	}

	ds, err := newDiskSync(source, dest, enableLogicallyDelete, chunkSize, checkpointCount, forceChecksum, encOpt)
	if err != nil {
		return nil, err
	}

	s := &minIOPushClientSync{
		driverPushClientSync: driverPushClientSync{
			diskSync: *ds,
			basePath: "",
		},
		endpoint:    dest.Addr(),
		bucketName:  dest.RemotePath(),
		secure:      dest.Secure(),
		currentUser: users[0],
	}

	s.client = minio.NewMinIOClient(s.endpoint, s.bucketName, s.secure, s.currentUser.UserName(), s.currentUser.Password(), true, r)

	err = s.start()
	if err != nil {
		return nil, err
	}
	return s, nil
}
