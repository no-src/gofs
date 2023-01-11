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
	if err == nil && exist != c.Expect {
		err = newNotExpectedError(errIsExistNotExpected, c.Expect, exist)
	}
	return err
}

func (c isExist) Name() string {
	return "is-exist"
}

func init() {
	registerCommand("is-exist", func(a Action) (Command, error) {
		return parse[isExist](a)
	})
}
