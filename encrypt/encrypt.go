package encrypt

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/no-src/gofs/logger"
	"github.com/no-src/nsgo/fsutil"
)

// Encrypt the encryption component
type Encrypt struct {
	opt        Option
	parentPath string
	logger     *logger.Logger
}

// NewEncrypt create an encryption component
func NewEncrypt(opt Option, parentPath string) (*Encrypt, error) {
	enc := &Encrypt{
		opt:        opt,
		parentPath: parentPath,
		logger:     opt.Logger,
	}
	if enc.opt.Encrypt {
		isSub, err := fsutil.IsSub(parentPath, opt.EncryptPath)
		if err != nil {
			return nil, err
		}
		if !isSub {
			return nil, fmt.Errorf("%w, source=%s encrypt=%s", errNotSubDir, parentPath, opt.EncryptPath)
		}
		if err = checkAESKey(opt.EncryptSecret); err != nil {
			return nil, err
		}
	}
	return enc, nil
}

// NewWriter create an encryption writer
func (e *Encrypt) NewWriter(w io.Writer, source string, name string) (io.WriteCloser, error) {
	if e.NeedEncrypt(source) {
		return newEncryptWriter(w, name, e.opt.EncryptSecret, aesIV)
	}
	return newBufferWriter(w), nil
}

// NeedEncrypt encryption is enabled and path is matched
func (e *Encrypt) NeedEncrypt(path string) bool {
	if e.opt.Encrypt {
		isSub, err := fsutil.IsSub(e.opt.EncryptPath, path)
		if err == nil && isSub {
			return true
		}
	}
	return false
}

// CreateEncryptTemp create an encryption temporary file if enable encrypt and the path is matched
func (e *Encrypt) CreateEncryptTemp(path string) (tempPath string, removeTemp func() error, err error) {
	removeTemp = func() error {
		return nil
	}
	if !e.NeedEncrypt(path) {
		return path, removeTemp, nil
	}
	sourceFile, err := os.Open(path)
	if err != nil {
		return tempPath, removeTemp, err
	}
	defer func() {
		e.logger.ErrorIf(sourceFile.Close(), "[encrypt temp] close the source file error")
	}()
	sourceStat, err := sourceFile.Stat()
	if err != nil {
		return tempPath, removeTemp, err
	}

	fileName := sourceStat.Name()
	reader := bufio.NewReader(sourceFile)

	tempFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return tempPath, removeTemp, err
	}

	defer func() {
		e.logger.ErrorIf(tempFile.Close(), "[encrypt temp] close the temporary file error")
	}()

	removeTemp = func() error {
		return e.logger.ErrorIf(os.Remove(tempFile.Name()), "[encrypt temp] remove the temporary file error")
	}

	w, err := e.NewWriter(tempFile, path, fileName)
	if err != nil {
		return tempPath, removeTemp, err
	}
	defer func() {
		e.logger.ErrorIf(w.Close(), "[encrypt temp] close the encrypt writer error")
	}()
	_, err = reader.WriteTo(w)
	if err != nil {
		return tempPath, removeTemp, err
	}
	return tempFile.Name(), removeTemp, nil
}
