package command

import "os"

type mkdir struct {
	Source string `yaml:"source"`
}

func (c mkdir) Exec() error {
	return os.MkdirAll(c.Source, defaultDirPerm)
}

func (c mkdir) Name() string {
	return "mkdir"
}

func init() {
	registerCommand("mkdir", func(a Action) (Command, error) {
		return parse[mkdir](a)
	})
}
