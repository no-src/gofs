package command

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	errUnsupportedCommand = errors.New("unsupported command")
	errNotExpected        = errors.New("the check result did not match the expectation")
)

var (
	defaultDirPerm  os.FileMode = 0700
	defaultFilePerm os.FileMode = 0666
)

// Command defined some common commands abstraction
type Command interface {
	// Exec execute the command
	Exec() error

	// Name return the command name
	Name() string
}

// Exec parse the config file to a command list and execute them in turn
func Exec(conf string) error {
	commands, err := ParseConfigFile(conf)
	if err != nil {
		return err
	}
	return commands.Exec()
}

func parse[T Command](a Action) (c T, err error) {
	out, err := yaml.Marshal(a)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(out, &c)
	return
}

func newNotExpectedError(c Command, actual any) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	expect := strings.Join(strings.Split(strings.TrimSpace(string(b)), "\n"), ", ")
	return fmt.Errorf("%w, get result => [%v], expect => [%s]", errNotExpected, actual, expect)
}
