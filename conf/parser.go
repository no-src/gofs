package conf

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/no-src/nsgo/jsonutil"
	"github.com/no-src/nsgo/yamlutil"
)

var (
	errEmptyConfigPath         = errors.New("the config file is empty")
	errUnSupportedConfigFormat = errors.New("unsupported config format")
)

// Parse read and parse the config file, support json and yaml format currently
func Parse[T any](path string, config *T) error {
	if len(path) == 0 {
		return errEmptyConfigPath
	}
	confBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ext := filepath.Ext(path)
	return ParseContent(confBytes, ext, config)
}

// ParseContent parse the config content, support json and yaml format currently
func ParseContent[T any](content []byte, ext string, config *T) (err error) {
	if JsonFormat.MatchExt(ext) {
		err = jsonutil.Unmarshal(content, &config)
	} else if YamlFormat.MatchExt(ext) {
		err = yamlutil.Unmarshal(content, &config)
	} else {
		err = errUnSupportedConfigFormat
	}
	return err
}

// ToString convert the config object to string, support json and yaml format currently
func ToString(ext string, config any) (s string, err error) {
	var data []byte
	if JsonFormat.MatchExt(ext) {
		data, err = jsonutil.Marshal(config)
	} else if YamlFormat.MatchExt(ext) {
		data, err = yamlutil.Marshal(config)
	} else {
		err = errUnSupportedConfigFormat
	}
	return string(data), err
}
