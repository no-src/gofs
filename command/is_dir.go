package command

import (
	"fmt"

	"github.com/no-src/gofs/fs"
)

var errIsDirNotExpected = fmt.Errorf("[is-dir] %w", errNotExpected)

type isDir struct {
	Source string `yaml:"source"`
	Expect bool   `yaml:"expect"`
}

func (c isDir) Exec() error {
	isDir, err := fs.IsDir(c.Source)
	if err != nil {
		return err
	}
	if isDir != c.Expect {
		err = newNotExpectedError(errIsDirNotExpected, c.Expect, isDir)
	}
	return err
}
