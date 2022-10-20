package encrypt

import (
	"archive/zip"
	"bufio"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"

	"github.com/no-src/gofs/fs"
	"github.com/no-src/log"
)

type decryptReader struct {
	zrc    *zip.ReadCloser
	secret []byte
}

func (r *decryptReader) WriteTo(path string) (err error) {
	for _, file := range r.zrc.File {
		outPath := filepath.Join(path, file.Name)

		// check zip slip
		var isSub bool
		isSub, err = fs.IsSub(path, outPath)
		if err != nil {
			return err
		}
		if !isSub {
			return fmt.Errorf("illegal file path => %s", file.Name)
		}

		// path is directory
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(outPath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// path is a file
		var f iofs.File
		f, err = r.zrc.Open(file.Name)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return err
		}

		var out *os.File
		out, err = os.Create(outPath)
		if err != nil {
			return err
		}

		br := bufio.NewReader(f)
		_, err = br.WriteTo(newDecryptWriter(out, r.secret))
		if err != nil {
			out.Close()
			return err
		}
		err = out.Close()
		if err != nil {
			return err
		}
		log.Info("save decryption file success => %s", outPath)
	}
	return err
}

// NewDecryptReader create a decryption reader
func NewDecryptReader(path string, secret []byte) (*decryptReader, error) {
	zrc, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return &decryptReader{
		zrc:    zrc,
		secret: secret,
	}, nil
}
