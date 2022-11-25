package command

import (
	"fmt"

	"github.com/no-src/gofs/util/hashutil"
)

var errHashNotExpected = fmt.Errorf("[hash] %w", errNotExpected)

type hash struct {
	Algorithm string `yaml:"algorithm"`
	Source    string `yaml:"source"`
	Expect    string `yaml:"expect"`
}

func (c hash) Exec() error {
	h, err := hashutil.NewHash(c.Algorithm)
	if err != nil {
		return err
	}
	hash, err := hashutil.HashFromFileName(c.Source, h)
	if err != nil {
		return err
	}
	if hash != c.Expect {
		err = newNotExpectedError(errHashNotExpected, c.Expect, hash)
	}
	return err
}
