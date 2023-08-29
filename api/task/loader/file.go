package loader

import (
	"io/fs"
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
	if err = c.verify(); err != nil {
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

func (loader *fileLoader) SaveConfig(c *TaskConfig) error {
	if c == nil {
		return errNilTaskConfig
	}
	if err := c.verify(); err != nil {
		return err
	}
	data, err := conf.ToString(filepath.Ext(loader.path), c)
	if err != nil {
		return err
	}
	return loader.write(loader.path, data)
}

func (loader *fileLoader) SaveContent(conf string, content string) error {
	if !filepath.IsAbs(conf) {
		conf = filepath.Join(filepath.Dir(loader.path), conf)
	}
	return loader.write(conf, content)
}

func (loader *fileLoader) Close() error {
	return nil
}

func (loader *fileLoader) write(path string, content string) error {
	return os.WriteFile(path, []byte(content), fs.ModePerm)
}
