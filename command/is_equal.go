package command

import (
	stdhash "hash"
	"os"

	"github.com/no-src/gofs/util/hashutil"
)

type isEqual struct {
	Source       string `yaml:"source"`
	Dest         string `yaml:"dest"`
	Expect       bool   `yaml:"expect"`
	MustNonEmpty bool   `yaml:"must-non-empty"`
	Algorithm    string `yaml:"algorithm"`
}

func (c isEqual) Exec() error {
	h, err := c.newHash()
	if err != nil {
		return err
	}
	srcStat, err := os.Stat(c.Source)
	if err != nil {
		return err
	}
	if c.MustNonEmpty && srcStat.Size() == 0 {
		return newNotExpectedError(c, "empty source file")
	}
	dstStat, err := os.Stat(c.Dest)
	if err != nil {
		return err
	}
	if c.MustNonEmpty && dstStat.Size() == 0 {
		return newNotExpectedError(c, "empty dest file")
	}
	actual := srcStat.Size() == dstStat.Size()
	if c.Expect && !actual {
		return newNotExpectedError(c, actual)
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
				err = newNotExpectedError(c, actual)
			}
		}
	}
	return err
}

func (c isEqual) Name() string {
	return "is-equal"
}

func (c isEqual) newHash() (stdhash.Hash, error) {
	algorithm := c.Algorithm
	if len(algorithm) == 0 {
		algorithm = hashutil.MD5Hash
	}
	return hashutil.NewHash(algorithm)
}

func init() {
	registerCommand[isEqual]()
}
