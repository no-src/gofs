package command

import (
	"io"
	"os"
)

type cp struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest"`
}

func (c cp) Exec() error {
	src, err := os.Open(c.Source)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(c.Dest)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func (c cp) Name() string {
	return "cp"
}

func init() {
	registerCommand[cp]()
}
