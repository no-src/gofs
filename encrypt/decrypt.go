package encrypt

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	nsfs "github.com/no-src/gofs/fs"
)

var (
	errDecryptOutNotDir = errors.New("the decrypt output path must be directory")
	errIllegalPath      = errors.New("illegal file path")
	errNotSubDir        = errors.New("the encrypt path is not a subdirectory of the source path")
)

// Decrypt the decryption component
type Decrypt struct {
	opt Option
}

// NewDecrypt create a decryption component
func NewDecrypt(opt Option) (*Decrypt, error) {
	if opt.Decrypt {
		if err := checkAESKey(opt.DecryptSecret); err != nil {
			return nil, err
		}
	}
	return &Decrypt{
		opt: opt,
	}, nil
}

// Decrypt uses the decryption option to decrypt the files
func (dec *Decrypt) Decrypt() error {
	if !dec.opt.Decrypt {
		return nil
	}
	isDir, err := nsfs.IsDir(dec.opt.DecryptOut)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dec.opt.DecryptOut, fs.ModePerm)
		}
		if err != nil {
			return err
		}
		isDir = true
	}
	if !isDir {
		return errDecryptOutNotDir
	}
	return filepath.WalkDir(dec.opt.DecryptPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		r, err := NewDecryptReader(path, dec.opt.DecryptSecret)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dec.opt.DecryptPath, path)
		if err != nil {
			return err
		}
		outPath := filepath.Join(dec.opt.DecryptOut, rel)
		outPath = filepath.Dir(outPath)
		return r.WriteTo(outPath)
	})
}
