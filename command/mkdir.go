package command

import "os"

type mkdir struct {
	Source string `yaml:"source"`
}

func (c mkdir) Exec() error {
	return os.MkdirAll(c.Source, defaultDirPerm)
}
