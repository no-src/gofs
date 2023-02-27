package command

import (
	"errors"
	"os"
	"testing"
)

func TestExec(t *testing.T) {
	testCases := []struct {
		conf string
	}{
		{"./example/command.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.conf, func(t *testing.T) {
			err := Exec(tc.conf)
			if err != nil {
				t.Errorf("execute commands error, err=%v", err)
			}
		})
	}
}

func TestExec_ReturnError(t *testing.T) {
	testCases := []struct {
		conf string
	}{
		{"./example/error_init.yaml"},
		{"./example/error_actions.yaml"},
		{"./example/error_clear.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.conf, func(t *testing.T) {
			err := Exec(tc.conf)
			if err == nil {
				t.Errorf("execute commands expect to get en error but get nil")
			}
		})
	}
}

func TestExec_ConfigFileNotExist(t *testing.T) {
	err := Exec("./example/notexist.yaml")
	if !os.IsNotExist(err) {
		t.Errorf("Exec expect to get a not exist error, but get %v", err)
	}
}

var (
	errMarshalYamlMock = errors.New("marshal yaml error mock")

	testExecReturnErrorFailedMessage = "execute command expect to get an error but get nil"
)

type errMarshaler struct {
}

func (m errMarshaler) MarshalYAML() (interface{}, error) {
	return nil, errMarshalYamlMock
}

type commandCase struct {
	name string
	cmd  Command
}
