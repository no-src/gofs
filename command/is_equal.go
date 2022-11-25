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
}

func (c isEqual) Exec() error {
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
	srcHash, err := hashutil.HashFromFileName(c.Source)
	if err != nil {
		return err
	}
	dstHash, err := hashutil.HashFromFileName(c.Dest)
	if err != nil {
		return err
	}
	actual = srcHash == dstHash
	if actual != c.Expect {
		err = newNotExpectedError(errIsEqualNotExpected, c.Expect, actual)
	}
	return err
}
