package command

import (
	"github.com/no-src/gofs/util/hashutil"
)

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
	hf, err := hashutil.HashFromFileName(c.Source, h)
	if err != nil {
		return err
	}
	if hf != c.Expect {
		err = newNotExpectedError(c, hf)
	}
	return err
}

func (c hash) Name() string {
	return "hash"
}

func init() {
	registerCommand[hash]()
}
