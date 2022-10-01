package encrypt

import (
	"fmt"
	"io"

	"github.com/no-src/gofs/fs"
)

// Encrypt the encryption component
type Encrypt struct {
	opt        Option
	parentPath string
}

// NewEncrypt create an encryption component
func NewEncrypt(opt Option, parentPath string) (*Encrypt, error) {
	enc := &Encrypt{
		opt:        opt,
		parentPath: parentPath,
	}
	if enc.opt.Encrypt {
		isSub, err := fs.IsSub(parentPath, opt.EncryptPath)
		if err != nil {
			return nil, err
		}
		if !isSub {
			return nil, fmt.Errorf("the encrypt path is not a subdirectory of the source path, source=%s encrypt=%s", parentPath, opt.EncryptPath)
		}
	}
	return enc, nil
}

// NewWriter create an encryption writer
func (e *Encrypt) NewWriter(w io.Writer, source string, name string) (io.WriteCloser, error) {
	if e.NeedEncrypt(source) {
		return newEncryptWriter(w, name, e.opt.EncryptSecret)
	}
	return newBufferWriter(w), nil
}

// NeedEncrypt encryption is enabled and path is matched
func (e *Encrypt) NeedEncrypt(path string) bool {
	if e.opt.Encrypt {
		isSub, err := fs.IsSub(e.opt.EncryptPath, path)
		if err == nil && isSub {
			return true
		}
	}
	return false
}
