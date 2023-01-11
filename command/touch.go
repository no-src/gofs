package command

import (
	"github.com/no-src/gofs/fs"
)

type touch struct {
	Source string `yaml:"source"`
}

func (c touch) Exec() error {
	f, err := fs.CreateFile(c.Source)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func (c touch) Name() string {
	return "touch"
}

func init() {
	registerCommand("touch", func(a Action) (Command, error) {
		return parse[touch](a)
	})
}
