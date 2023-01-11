package command

import "os"

type mv struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest"`
}

func (c mv) Exec() error {
	return os.Rename(c.Source, c.Dest)
}

func (c mv) Name() string {
	return "mv"
}

func init() {
	registerCommand("mv", func(a Action) (Command, error) {
		return parse[mv](a)
	})
}
