package conf

import (
	"os"
	"testing"
)

const (
	jsonConfigPath = "./example/gofs-remote-client.json"
	yamlConfigPath = "./example/gofs-remote-server.yaml"
)

func TestParseJsonConfiguration(t *testing.T) {
	c := Config{}
	path := jsonConfigPath
	err := Parse(path, &c)
	if err != nil {
		t.Errorf("parse json configuration error => %s", path)
	}
}

func TestParseYamlConfiguration(t *testing.T) {
	c := Config{}
	path := yamlConfigPath
	err := Parse(path, &c)
	if err != nil {
		t.Errorf("parse yaml configuration error => %s", path)
	}
}

func TestParseEmptyPath(t *testing.T) {
	c := Config{}
	path := ""
	err := Parse(path, &c)
	if err != errEmptyConfigPath {
		t.Errorf("parse configuration error, expert err:%s, actual err:%s", errEmptyConfigPath, err)
	}
}

func TestParseNotExistPath(t *testing.T) {
	c := Config{}
	path := "./not-exist-file.json"
	err := Parse(path, &c)
	if !os.IsNotExist(err) {
		t.Errorf("parse configuration error, expert error not exist, actual err:%s", err)
	}
}

func TestParseUnSupportedConfigFormat(t *testing.T) {
	c := Config{}
	path := "./parser_test.go"
	err := Parse(path, &c)
	if err != errUnSupportedConfigFormat {
		t.Errorf("parse configuration error, expert err:%s, actual err:%s", errUnSupportedConfigFormat, err)
	}
}
