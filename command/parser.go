package command

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	errInfiniteRecursion = errors.New("find an infinite recursion")
)

const (
	maxRecursionDepth = 10000
)

type parser struct {
	c int // check infinite recursion
}

func newParser() *parser {
	return &parser{}
}

func (p *parser) ParseConfigFile(path string) (commands *Commands, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return
	}
	if len(c.Root) == 0 {
		c.Root = filepath.Dir(path)
	}
	return p.ParseConfig(c)
}

func (p *parser) ParseConfig(conf Config) (commands *Commands, err error) {
	commands, err = p.parseConfig(conf)
	if err != nil {
		return nil, err
	}

	var (
		init    []Command
		actions []Command
		clear   []Command
	)
	for _, f := range conf.Include {
		baseCommands, err := p.ParseConfigFile(filepath.Join(conf.Root, f))
		if err != nil {
			return nil, err
		}
		init = append(init, baseCommands.Init...)
		actions = append(actions, baseCommands.Actions...)
		clear = append(clear, baseCommands.Clear...)
	}

	commands.Init = append(init, commands.Init...)
	commands.Actions = append(actions, commands.Actions...)
	commands.Clear = append(clear, commands.Clear...)
	return commands, nil
}

func (p *parser) parseConfig(conf Config) (commands *Commands, err error) {
	p.c++
	if p.c > maxRecursionDepth {
		return nil, errInfiniteRecursion
	}
	commands = &Commands{
		Name: conf.Name,
	}
	commands.Init, err = p.parseCommands(conf.Init)
	if err != nil {
		return nil, err
	}
	commands.Actions, err = p.parseCommands(conf.Actions)
	if err != nil {
		return nil, err
	}
	commands.Clear, err = p.parseCommands(conf.Clear)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func (p *parser) parseCommands(actions []Action) (commands []Command, err error) {
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

// ParseConfigFile parse the config file to a command list
func ParseConfigFile(path string) (commands *Commands, err error) {
	return newParser().ParseConfigFile(path)
}

// ParseConfig parse the config to a command list
func ParseConfig(conf Config) (commands *Commands, err error) {
	return newParser().ParseConfig(conf)
}
