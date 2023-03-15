package command

import (
	"os"
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
		return newNotExpectedError(c, actual)
	}
	if !c.Expect && !actual {
		return nil
	}
	srcData, err := os.ReadFile(c.Source)
	if err == nil {
		actual = string(srcData) == c.Dest
		if actual != c.Expect {
			err = newNotExpectedError(c, actual)
		}
	}
	return err
}

func (c isEqualText) Name() string {
	return "is-equal-text"
}

func init() {
	registerCommand[isEqualText]()
}
