package command

import (
	"os"
)

type echo struct {
	Source    string `yaml:"source"`
	Input     string `yaml:"input"`
	Append    bool   `yaml:"append"`
	NoNewline bool   `yaml:"no-newline"`
}

func (c echo) Exec() (err error) {
	var f *os.File
	if c.Append {
		f, err = os.OpenFile(c.Source, os.O_RDWR|os.O_CREATE|os.O_APPEND, defaultFilePerm)
	} else {
		f, err = os.OpenFile(c.Source, os.O_RDWR|os.O_CREATE|os.O_TRUNC, defaultFilePerm)
	}
	if err != nil {
		return err
	}
	defer f.Close()
	if !c.NoNewline {
		c.Input += "\n"
	}
	_, err = f.WriteString(c.Input)
	return err
}
