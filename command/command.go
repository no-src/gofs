package command

import (
	"errors"
	"fmt"
	"os"

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

// ParseConfigFile parse the config file to a command list
func ParseConfigFile(path string) (commands *Commands, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return
	}
	return ParseConfig(c)
}

// ParseConfig parse the config to a command list
func ParseConfig(conf Config) (commands *Commands, err error) {
	commands = &Commands{
		Name: conf.Name,
	}
	commands.Init, err = parseCommands(conf.Init)
	if err != nil {
		return nil, err
	}
	commands.Actions, err = parseCommands(conf.Actions)
	if err != nil {
		return nil, err
	}
	commands.Clear, err = parseCommands(conf.Clear)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func parseCommands(actions []Action) (commands []Command, err error) {
	for _, action := range actions {
		var c Command
		for name, fn := range allCommands {
			if _, ok := action[name]; ok {
				c, err = fn(action)
				break
			}
		}
		if err != nil {
			return nil, err
		}
		if c == nil {
			return nil, errUnsupportedCommand
		}
		commands = append(commands, c)
	}
	return commands, nil
}

func parse[T Command](a Action) (c T, err error) {
	out, err := yaml.Marshal(a)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(out, &c)
	return
}

func newNotExpectedError(err error, expect any, actual any) error {
	return fmt.Errorf("%w, expect to get %v, but get %v", err, expect, actual)
}
