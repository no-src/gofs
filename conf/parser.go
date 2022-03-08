package conf

import (
	"errors"
	"github.com/no-src/gofs/util"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var (
	errEmptyConfigPath         = errors.New("the config file is empty")
	errUnSupportedConfigFormat = errors.New("unsupported config format")
)

// Parse read and parse the config file, support json and yaml format currently
func Parse(path string, config *Config) error {
	if len(path) == 0 {
		return errEmptyConfigPath
	}
	confBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	if JsonFormat.MatchExt(ext) {
		err = util.Unmarshal(confBytes, &config)
	} else if YamlFormat.MatchExt(ext) {
		err = yaml.Unmarshal(confBytes, &config)
	} else {
		err = errUnSupportedConfigFormat
	}
	return err
}
