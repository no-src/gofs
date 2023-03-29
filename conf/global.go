package conf

import (
	"errors"
	"sync"
)

var (
	globalConfigSet = &configSet{
		m: make(map[string]*Config),
	}
)

var (
	errConfigIsNil = errors.New("the config is nil")
	errConfigExist = errors.New("the config exists")
)

type configSet struct {
	m  map[string]*Config
	mu sync.RWMutex
}

func (cs *configSet) setGlobalConfig(c *Config) error {
	if c == nil {
		return errConfigIsNil
	}
	addr := c.FileServerAddr
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if _, ok := cs.m[addr]; ok {
		return errConfigExist
	}
	cs.m[addr] = c
	return nil
}

func (cs *configSet) getGlobalConfig(addr string) *Config {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.m[addr]
}

// SetGlobalConfig set the global config once per web server
func SetGlobalConfig(c *Config) error {
	return globalConfigSet.setGlobalConfig(c)
}

// GetGlobalConfig get the global config by web server address
func GetGlobalConfig(addr string) *Config {
	return globalConfigSet.getGlobalConfig(addr)
}
