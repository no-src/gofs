package command

import (
	"fmt"
	"os"
)

var (
	errIsEqualTextNotExpected = fmt.Errorf("[is-equal-text] %w", errNotExpected)
)

type isEqualText struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest"`
	Expect bool   `yaml:"expect"`
}

func (c isEqualText) Exec() error {
	srcStat, err := os.Stat(c.Source)
	if err != nil {
		return err
	}
	actual := srcStat.Size() == int64(len(c.Dest))
	if c.Expect && !actual {
		return newNotExpectedError(errIsEqualTextNotExpected, c.Expect, actual)
	}
	if !c.Expect && !actual {
		return nil
	}
	srcData, err := os.ReadFile(c.Source)
	if err == nil {
		actual = string(srcData) == c.Dest
		if actual != c.Expect {
			err = newNotExpectedError(errIsEqualTextNotExpected, c.Expect, actual)
		}
	}
	return err
}

func (c isEqualText) Name() string {
	return "is-equal-text"
}

func init() {
	registerCommand("is-equal-text", func(a Action) (Command, error) {
		return parse[isEqualText](a)
	})
}
