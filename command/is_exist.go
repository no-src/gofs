package command

import (
	"fmt"

	"github.com/no-src/gofs/fs"
)

var errIsExistNotExpected = fmt.Errorf("[is-exist] %w", errNotExpected)

type isExist struct {
	Source string `yaml:"source"`
	Expect bool   `yaml:"expect"`
}

func (c isExist) Exec() error {
	exist, err := fs.FileExist(c.Source)
	if err != nil {
		return err
	}
	if exist != c.Expect {
		err = newNotExpectedError(errIsExistNotExpected, c.Expect, exist)
	}
	return err
}
