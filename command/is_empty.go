package command

import (
	"fmt"
	"os"
)

var errIsEmptyNotExpected = fmt.Errorf("[is-empty] %w", errNotExpected)

type isEmpty struct {
	Source string `yaml:"source"`
	Expect bool   `yaml:"expect"`
}

func (c isEmpty) Exec() error {
	stat, err := os.Stat(c.Source)
	if err != nil {
		return err
	}
	isEmpty := stat.Size() == 0
	if isEmpty != c.Expect {
		err = newNotExpectedError(errIsEmptyNotExpected, c.Expect, isEmpty)
	}
	return err
}
