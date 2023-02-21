package sync

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/driver/minio"
)

type minIOPullClientSync struct {
	driverPullClientSync

	endpoint    string
	bucketName  string
	secure      bool
	currentUser *auth.User
}

// NewMinIOPullClientSync create an instance of the minIOPullClientSync
func NewMinIOPullClientSync(opt Option) (Sync, error) {
	// the fields of option
	source := opt.Source
	users := opt.Users
	chunkSize := opt.ChunkSize
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

	s := &minIOPullClientSync{
		driverPullClientSync: driverPullClientSync{
			diskSync: *ds,
		},
		endpoint:    source.Addr(),
		bucketName:  source.RemotePath(),
		secure:      source.Secure(),
		currentUser: users[0],
	}

	s.driver = minio.NewMinIODriver(s.endpoint, s.bucketName, s.secure, s.currentUser.UserName(), s.currentUser.Password(), true, r)

	err = s.start()
	if err != nil {
		return nil, err
	}

	// reset the sourceAbsPath to current path
	s.diskSync.sourceAbsPath = "."

	// reset some functions for MinIO
	s.diskSync.isDirFn = s.IsDir
	s.diskSync.statFn = s.driver.Stat
	s.diskSync.getFileTimeFn = s.driver.GetFileTime

	return s, nil
}
