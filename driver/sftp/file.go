package sftp

import (
	"io/fs"
	"net/http"

	"github.com/pkg/sftp"
)

type file struct {
	*sftp.File

	driver *sftpDriver
	name   string
}

func newFile(f *sftp.File, driver *sftpDriver, name string) http.File {
	return &file{
		File:   f,
		driver: driver,
		name:   name,
	}
}

func (f *file) Readdir(count int) (fis []fs.FileInfo, err error) {
	fis, err = f.driver.ReadDir(f.name)
	if err == nil && count > 0 && len(fis) > count {
		fis = fis[:count]
	}
	return fis, err
}
