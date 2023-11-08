package encrypt

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/no-src/gofs/logger"
)

type decryptReader struct {
	zrc    *zip.ReadCloser
	secret []byte
	logger *logger.Logger
}

func (r *decryptReader) WriteTo(path string) (err error) {
	for _, file := range r.zrc.File {
		// check zip slip
		isValid := fs.ValidPath(file.Name)
		if !isValid {
			return fmt.Errorf("%w => %s", errIllegalPath, file.Name)
		}

		outPath := filepath.Join(path, file.Name)

		// path is directory
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(outPath, fs.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// path is a file
		var f fs.File
		f, err = r.zrc.Open(file.Name)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(outPath), fs.ModePerm)
		if err != nil {
			return err
		}

		var out *os.File
		out, err = os.Create(outPath)
		if err != nil {
			return err
		}

		var dw io.Writer
		dw, err = newDecryptWriter(out, r.secret, aesIV)
		if err != nil {
			out.Close()
			return err
		}

		br := bufio.NewReader(f)
		_, err = br.WriteTo(dw)
		if err != nil {
			out.Close()
			return err
		}
		err = out.Close()
		if err != nil {
			return err
		}
		r.logger.Info("save decryption file success => %s", outPath)
	}
	return err
}

// newDecryptReader create a decryption reader
func newDecryptReader(path string, secret []byte, logger *logger.Logger) (*decryptReader, error) {
	zrc, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return &decryptReader{
		zrc:    zrc,
		secret: secret,
		logger: logger,
	}, nil
}
