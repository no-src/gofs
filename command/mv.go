package command

import "os"

type mv struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest"`
}

func (c mv) Exec() error {
	return os.Rename(c.Source, c.Dest)
}
