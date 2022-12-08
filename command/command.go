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
	var c Command
	for _, action := range actions {
		if _, ok := action["cp"]; ok {
			c, err = parse[cp](action)
		} else if _, ok = action["mv"]; ok {
			c, err = parse[mv](action)
		} else if _, ok = action["rm"]; ok {
			c, err = parse[rm](action)
		} else if _, ok = action["touch"]; ok {
			c, err = parse[touch](action)
		} else if _, ok = action["echo"]; ok {
			c, err = parse[echo](action)
		} else if _, ok = action["mkdir"]; ok {
			c, err = parse[mkdir](action)
		} else if _, ok = action["run"]; ok {
			c, err = parse[run](action)
		} else if _, ok = action["sleep"]; ok {
			c, err = parse[sleep](action)
		} else if _, ok = action["is-equal"]; ok {
			c, err = parse[isEqual](action)
		} else if _, ok = action["is-equal-text"]; ok {
			c, err = parse[isEqualText](action)
		} else if _, ok = action["is-empty"]; ok {
			c, err = parse[isEmpty](action)
		} else if _, ok = action["is-exist"]; ok {
			c, err = parse[isExist](action)
		} else if _, ok = action["is-dir"]; ok {
			c, err = parse[isDir](action)
		} else if _, ok = action["hash"]; ok {
			c, err = parse[hash](action)
		} else {
			err = errUnsupportedCommand
		}

		if err != nil {
			return nil, err
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
