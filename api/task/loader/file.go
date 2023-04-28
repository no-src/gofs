package loader

import (
	"os"
	"path/filepath"

	"github.com/no-src/gofs/conf"
)

type fileLoader struct {
	path string
}

func newFileLoader(path string) Loader {
	return &fileLoader{
		path: path,
	}
}

func (loader *fileLoader) LoadConfig() (c *TaskConfig, err error) {
	if err = conf.Parse(loader.path, &c); err != nil {
		return nil, err
	}
	if err = c.Verify(); err != nil {
		return nil, err
	}
	return c, err
}

func (loader *fileLoader) LoadContent(conf string) (content string, err error) {
	if !filepath.IsAbs(conf) {
		conf = filepath.Join(filepath.Dir(loader.path), conf)
	}
	bytes, err := os.ReadFile(conf)
	return string(bytes), err
}
