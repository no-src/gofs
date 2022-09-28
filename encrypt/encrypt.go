package encrypt

import (
	"io"
	"strings"

	"github.com/no-src/gofs/fs"
	"github.com/no-src/log"
)

// Encrypt the encryption component
type Encrypt struct {
	opt        Option
	parentPath string
}

// NewEncrypt create an encryption component
func NewEncrypt(opt Option, parentPath string) Encrypt {
	enc := Encrypt{
		opt:        opt,
		parentPath: parentPath,
	}
	if enc.opt.Decrypt {
		isSub, err := fs.IsSub(parentPath, opt.EncryptPath)
		if err != nil || !isSub {
			log.Warn("disable encrypt because the encrypt path is not a subdirectory of the source path, source=%s encrypt=%s", parentPath, opt.EncryptPath)
			enc.opt.Decrypt = false
		}
	}
	return enc
}

// NewWriter create an encryption writer
func (e Encrypt) NewWriter(w io.Writer, source string, name string) (io.WriteCloser, error) {
	if e.opt.Encrypt {
		isSub, err := fs.IsSub(e.opt.EncryptPath, source)
		if err == nil && isSub {
			return newEncryptWriter(w, name, e.opt.EncryptSecret, e.opt.EncryptSuffix)
		}
	}
	return newBufferWriter(w), nil
}

// BuildEncryptName returns the encryption name of the destination file
func (e Encrypt) BuildEncryptName(source, dest string) string {
	if e.opt.Encrypt {
		suffix := e.opt.EncryptSuffix
		suffix = strings.TrimSpace(suffix)
		isSub, err := fs.IsSub(e.opt.EncryptPath, source)
		if err == nil && isSub && len(suffix) > 0 {
			dest = dest + suffix
		}
	}
	return dest
}
