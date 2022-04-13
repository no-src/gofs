package conf

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/no-src/gofs/util/jsonutil"
	"gopkg.in/yaml.v3"
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
		err = jsonutil.Unmarshal(confBytes, &config)
	} else if YamlFormat.MatchExt(ext) {
		err = yaml.Unmarshal(confBytes, &config)
	} else {
		err = errUnSupportedConfigFormat
	}
	return err
}
