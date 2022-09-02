package sync

import (
	"errors"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/driver/minio"
	"github.com/no-src/gofs/retry"
)

type minIOPullClientSync struct {
	driverPullClientSync

	endpoint    string
	bucketName  string
	secure      bool
	currentUser *auth.User
}

// NewMinIOPullClientSync create an instance of the minIOPullClientSync
func NewMinIOPullClientSync(source, dest core.VFS, users []*auth.User, enableLogicallyDelete bool, chunkSize int64, checkpointCount int, forceChecksum bool, r retry.Retry) (Sync, error) {
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

	s := &minIOPullClientSync{
		driverPullClientSync: driverPullClientSync{
			diskSync: *ds,
		},
		endpoint:    source.Addr(),
		bucketName:  source.RemotePath(),
		secure:      source.Secure(),
		currentUser: users[0],
	}

	s.client = minio.NewMinIOClient(s.endpoint, s.bucketName, s.secure, s.currentUser.UserName(), s.currentUser.Password(), true, r)

	err = s.start()
	if err != nil {
		return nil, err
	}

	// reset the sourceAbsPath to current path
	s.diskSync.sourceAbsPath = "."

	// reset some functions for MinIO
	s.diskSync.isDirFn = s.IsDir
	s.diskSync.statFn = s.client.Stat
	s.diskSync.getFileTimeFn = s.client.GetFileTime

	return s, nil
}
