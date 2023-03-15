package command

import "os"

type rm struct {
	Source string `yaml:"source"`
}

func (c rm) Exec() error {
	return os.RemoveAll(c.Source)
}

func (c rm) Name() string {
	return "rm"
}

func init() {
	registerCommand[rm]()
}
