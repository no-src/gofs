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
	if err = c.Verify(); err != nil {
		return nil, err
	}
	return c, err
}

func (loader *cacheLoader) LoadContent(conf string) (content string, err error) {
	contentKey := fmt.Sprintf("%s:%s", loader.confKey, conf)
	err = loader.cache.Get(contentKey, &content)
	return content, err
}
