package command

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestParseConfigFile_ConfigFileNotExist(t *testing.T) {
	_, err := ParseConfigFile("./example/notexist.yaml")
	if !os.IsNotExist(err) {
		t.Errorf("ParseConfigFile expect to get a not exist error, but get %v", err)
	}
}

func TestParseConfigFile_InvalidConfigFile(t *testing.T) {
	_, err := ParseConfigFile("./command_test.go")
	if err == nil {
		t.Errorf("ParseConfigFile expect get an error, but get nil")
	}
}

func TestParseConfig_UnsupportedCommand(t *testing.T) {
	testCases := []struct {
		name string
		conf Config
	}{
		{"unsupported command in init", Config{Init: []Action{{"unsupported-command": ""}}}},
		{"unsupported command in actions", Config{Actions: []Action{{"unsupported-command": ""}}}},
		{"unsupported command in clear", Config{Clear: []Action{{"unsupported-command": ""}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseConfig(tc.conf)
			if !errors.Is(err, errUnsupportedCommand) {
				t.Errorf("ParseConfig expect get error => %v, but get %v", errUnsupportedCommand, err)
			}
		})
	}
}

func TestParseConfig_WithIllegalField(t *testing.T) {
	conf := Config{
		Name: "invalid command",
	}
	action := make(Action)
	action["cp"] = ""
	action["source"] = errMarshaler{}
	conf.Actions = append(conf.Actions, action)
	_, err := ParseConfig(conf)
	if !errors.Is(err, errMarshalYamlMock) {
		t.Errorf("ParseConfig expect get error => %v, but get %v", errMarshalYamlMock, err)
	}
}

func TestParseConfigFile_InfiniteRecursion(t *testing.T) {
	_, err := ParseConfigFile("./example/infinite_recursion.yaml")
	if !errors.Is(err, errInfiniteRecursion) {
		t.Errorf("ParseConfigFile expect to get %v, but get %v", errInfiniteRecursion, err)
	}
}

func TestParseConfigFile_WithInclude(t *testing.T) {
	commands, err := ParseConfigFile("./example/error.yaml")
	if err != nil {
		t.Errorf("ParseConfigFile expect to get an nil error, but get %v", err)
		return
	}
	if len(commands.Init) != 1 || len(commands.Actions) != 1 || len(commands.Clear) != 1 {
		t.Errorf("parse config file failed with include section")
	}
}

func ExampleParseConfigFile_WithInclude() {
	commands, err := ParseConfigFile("./example/include/include.yaml")
	if err != nil {
		panic(fmt.Sprintf("ParseConfigFile expect to get an nil error, but get %v", err))
	}
	err = commands.Exec()
	if err != nil {
		panic(fmt.Sprintf("Exec expect to get an nil error, but get %v", err))
	}

	// Output:
	//call init from init_source.yaml
	//call init from init_dest.yaml
	//call init from actions_step1.yaml
	//call init from actions_step2.yaml
	//call init from actions.yaml
	//call init from clear_step1.yaml
	//call init from clear_step2.yaml
	//call init from clear.yaml
	//call init from include.yaml
	//call actions from actions_step1.yaml
	//call actions from actions_step2.yaml
	//call actions from actions.yaml
	//call actions from clear_step1.yaml
	//call actions from clear_step2.yaml
	//call actions from clear.yaml
	//call actions from include.yaml
	//call clear from actions_step1.yaml
	//call clear from actions_step2.yaml
	//call clear from actions.yaml
	//call clear from clear_step1.yaml
	//call clear from clear_step2.yaml
	//call clear from clear.yaml
	//call clear from include.yaml
}
