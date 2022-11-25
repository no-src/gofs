package command

import (
	"fmt"
	"os"

	"github.com/no-src/gofs/util/hashutil"
)

var (
	errIsEqualNotExpected  = fmt.Errorf("[is-equal] %w", errNotExpected)
	errIsEqualMustNonEmpty = fmt.Errorf("[is-equal] %w, must be non-empty", errNotExpected)
)

type isEqual struct {
	Source       string `yaml:"source"`
	Dest         string `yaml:"dest"`
	Expect       bool   `yaml:"expect"`
	MustNonEmpty bool   `yaml:"must-non-empty"`
	Algorithm    string `yaml:"algorithm"`
}

func (c isEqual) Exec() error {
	if len(c.Algorithm) == 0 {
		c.Algorithm = hashutil.MD5Hash
	}
	h, err := hashutil.NewHash(c.Algorithm)
	if err != nil {
		return err
	}
	srcStat, err := os.Stat(c.Source)
	if err != nil {
		return err
	}
	if c.MustNonEmpty && srcStat.Size() == 0 {
		return errIsEqualMustNonEmpty
	}
	dstStat, err := os.Stat(c.Dest)
	if err != nil {
		return err
	}
	if c.MustNonEmpty && dstStat.Size() == 0 {
		return errIsEqualMustNonEmpty
	}
	actual := srcStat.Size() == dstStat.Size()
	if c.Expect && !actual {
		return newNotExpectedError(errIsEqualNotExpected, c.Expect, actual)
	}
	if !c.Expect && !actual {
		return nil
	}
	srcHash, err := hashutil.HashFromFileName(c.Source, h)
	if err == nil {
		h.Reset()
		var dstHash string
		dstHash, err = hashutil.HashFromFileName(c.Dest, h)
		if err == nil {
			actual = srcHash == dstHash
			if actual != c.Expect {
				err = newNotExpectedError(errIsEqualNotExpected, c.Expect, actual)
			}
		}
	}
	return err
}
