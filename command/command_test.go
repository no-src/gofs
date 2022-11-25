package command

import (
	"errors"
	"os"
	"testing"
)

func TestExec(t *testing.T) {
	err := Exec("./example/command.yaml")
	if err != nil {
		t.Errorf("execute commands error, err=%v", err)
	}
}

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
	conf := Config{
		Name: "unsupported command",
	}
	action := make(Action)
	action["unsupported-command"] = ""
	conf.Actions = append(conf.Actions, action)
	_, err := ParseConfig(conf)
	if !errors.Is(err, errUnsupportedCommand) {
		t.Errorf("ParseConfig expect get error => %v, but get %v", errUnsupportedCommand, err)
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

var errMarshalYamlMock = errors.New("marshal yaml error mock")

type errMarshaler struct {
}

func (m errMarshaler) MarshalYAML() (interface{}, error) {
	return nil, errMarshalYamlMock
}
