package loader

import (
	"fmt"

	"github.com/no-src/nscache"
	_ "github.com/no-src/nscache/boltdb"
	_ "github.com/no-src/nscache/buntdb"
	_ "github.com/no-src/nscache/etcd"
	_ "github.com/no-src/nscache/memory"
	_ "github.com/no-src/nscache/redis"
)

type cacheLoader struct {
	cache   nscache.NSCache
	confKey string
}

func newCacheLoader(conn string, confKey string) (Loader, error) {
	cache, err := nscache.NewCache(conn)
	if err != nil {
		return nil, err
	}
	return &cacheLoader{
		cache:   cache,
		confKey: confKey,
	}, nil
}

func (loader *cacheLoader) LoadConfig() (c *TaskConfig, err error) {
	if err = loader.cache.Get(loader.confKey, &c); err != nil {
		return nil, err
	}
	if err = c.verify(); err != nil {
		return nil, err
	}
	return c, err
}

func (loader *cacheLoader) LoadContent(conf string) (content string, err error) {
	err = loader.cache.Get(loader.contentKey(conf), &content)
	return content, err
}

func (loader *cacheLoader) SaveConfig(c *TaskConfig) error {
	if c == nil {
		return errNilTaskConfig
	}
	if err := c.verify(); err != nil {
		return err
	}
	return loader.cache.Set(loader.confKey, c, 0)
}

func (loader *cacheLoader) SaveContent(conf string, content string) error {
	return loader.cache.Set(loader.contentKey(conf), content, 0)
}

func (loader *cacheLoader) Close() error {
	return loader.cache.Close()
}

func (loader *cacheLoader) contentKey(conf string) string {
	return fmt.Sprintf("%s:%s", loader.confKey, conf)
}
