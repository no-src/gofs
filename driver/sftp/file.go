package sftp

import (
	"io/fs"
	"net/http"

	"github.com/pkg/sftp"
)

type file struct {
	*sftp.File

	client *sftpClient
	name   string
}

func newFile(f *sftp.File, client *sftpClient, name string) http.File {
	return &file{
		File:   f,
		client: client,
		name:   name,
	}
}

func (f *file) Readdir(count int) (fis []fs.FileInfo, err error) {
	fis, err = f.client.ReadDir(f.name)
	if err == nil && count > 0 && len(fis) > count {
		fis = fis[:count]
	}
	return fis, err
}
