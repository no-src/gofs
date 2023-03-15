package command

import (
	"os"
)

type isEmpty struct {
	Source string `yaml:"source"`
	Expect bool   `yaml:"expect"`
}

func (c isEmpty) Exec() error {
	stat, err := os.Stat(c.Source)
	if err != nil {
		return err
	}
	empty := stat.Size() == 0
	if empty != c.Expect {
		err = newNotExpectedError(c, empty)
	}
	return err
}

func (c isEmpty) Name() string {
	return "is-empty"
}

func init() {
	registerCommand[isEmpty]()
}
