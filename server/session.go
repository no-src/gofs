package server

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
)

const (
	// MemorySession represent a memory session store
	MemorySession = iota + 1
	// RedisSession represent a redis session store
	RedisSession
)

// NewSessionStore create a session store, stored in memory or redis, default is memory
func NewSessionStore(sessionMode int, sessionConnection string) (sessions.Store, error) {
	secret := make([]byte, 32)
	_, err := rand.Reader.Read(secret)
	if err != nil {
		return nil, err
	}
	switch sessionMode {
	case MemorySession:
		return memstore.NewStore(secret), nil
	case RedisSession:
		return redisSessionStore(sessionConnection, secret)
	default:
		return memstore.NewStore(secret), nil
	}
}

func redisSessionStore(redisUrl string, secret []byte) (sessions.Store, error) {
	maxIdle, network, address, password, db, err := parseRedisConnection(redisUrl)
	if err != nil {
		return nil, err
	}
	// get the existing secret in the redis, if not exist, set the new secret
	// TODO
	return redis.NewStoreWithDB(maxIdle, network, address, password, db, secret)
}

// parseRedisConnection parse the redis connection string
// for example => redis://127.0.0.1:6379?password=secret&db=10&max_idle=10
func parseRedisConnection(redisUrl string) (maxIdle int, network, address, password string, db string, err error) {
	u, err := url.Parse(redisUrl)
	if err != nil {
		return
	}
	// network
	network = "tcp"

	// maxIdle
	defaultMaxIdle := 10
	maxIdleStr := u.Query().Get("max_idle")
	if len(maxIdleStr) == 0 {
		maxIdle = defaultMaxIdle
	} else {
		maxIdle, err = strconv.Atoi(maxIdleStr)
		if err != nil {
			err = fmt.Errorf("invalid redis max idle => %d", maxIdle)
			return
		} else if maxIdle <= 0 {
			err = fmt.Errorf("invalid redis max idle => %d, max idle must be greater than zero", maxIdle)
			return
		}
	}

	// address
	address = u.Host

	// password
	password = u.Query().Get("password")

	// db
	defaultDB := "0"
	db = u.Query().Get("db")
	if len(db) == 0 {
		db = defaultDB
		return
	}
	if dbInt, dbErr := strconv.Atoi(db); dbErr != nil {
		err = fmt.Errorf("invalid redis db => %s", db)
	} else if dbInt < 0 || dbInt > 15 {
		err = fmt.Errorf("invalid redis db => %s, db must be between 0 and 15", db)
	}
	return
}
