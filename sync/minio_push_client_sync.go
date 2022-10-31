package sync

import (
	"errors"

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

	s.driver = minio.NewMinIODriver(s.endpoint, s.bucketName, s.secure, s.currentUser.UserName(), s.currentUser.Password(), true, r)

	err = s.start()
	if err != nil {
		return nil, err
	}
	return s, nil
}
