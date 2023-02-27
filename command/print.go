package command

import (
	"fmt"
)

type print struct {
	Input string `yaml:"input"`
}

func (c print) Exec() error {
	_, err := fmt.Println(c.Input)
	return err
}

func (c print) Name() string {
	return "print"
}

func init() {
	registerCommand("print", func(a Action) (Command, error) {
		return parse[print](a)
	})
}
