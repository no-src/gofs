package server

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
)

var (
	errInvalidSession      = errors.New("invalid session connection")
	errUnsupportedSession  = errors.New("unsupported session connection")
	errInvalidRedisDB      = errors.New("invalid redis db")
	errInvalidRedisMaxIdle = errors.New("invalid redis max idle")
)

// NewSessionStore create a session store, stored in memory or redis
func NewSessionStore(sessionConnection string) (sessions.Store, error) {
	secret := make([]byte, 64)
	_, err := rand.Reader.Read(secret)
	if err != nil {
		return nil, err
	}
	connUrl, err := url.Parse(sessionConnection)
	if err != nil {
		return nil, fmt.Errorf("%w => %s", errors.Join(errInvalidSession, err), sessionConnection)
	}
	switch strings.ToLower(connUrl.Scheme) {
	case "memory":
		return memstore.NewStore(secret), nil
	case "redis":
		return redisSessionStore(connUrl, secret)
	default:
		return nil, fmt.Errorf("%w => %s", errUnsupportedSession, sessionConnection)
	}
}

func redisSessionStore(redisUrl *url.URL, secret []byte) (sessions.Store, error) {
	maxIdle, network, address, password, db, redisSecret, err := parseRedisConnection(redisUrl)
	if err != nil {
		return nil, err
	}
	if len(redisSecret) > 0 {
		secret = redisSecret
	}
	// get the existing secret in the redis, if not exist, set the new secret
	// TODO
	return redis.NewStoreWithDB(maxIdle, network, address, password, strconv.Itoa(db), secret)
}

// parseRedisConnection parse the redis connection string
// for example => redis://127.0.0.1:6379?password=redis_password&db=10&max_idle=10&secret=redis_secret
func parseRedisConnection(u *url.URL) (maxIdle int, network, address, password string, db int, secret []byte, err error) {
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
			err = fmt.Errorf("%w => %d", errInvalidRedisMaxIdle, maxIdle)
			return
		} else if maxIdle <= 0 {
			err = fmt.Errorf("%w => %d, max idle must be greater than zero", errInvalidRedisMaxIdle, maxIdle)
			return
		}
	}

	// address
	address = u.Host

	// password
	password = u.Query().Get("password")

	// secret
	secret = []byte(u.Query().Get("secret"))

	// db
	dbValue := u.Query().Get("db")
	if len(dbValue) == 0 {
		return
	}
	if dbInt, dbErr := strconv.Atoi(dbValue); dbErr != nil {
		err = fmt.Errorf("%w => %s", errInvalidRedisDB, dbValue)
	} else if dbInt < 0 || dbInt > 15 {
		err = fmt.Errorf("%w => %s, db must be between 0 and 15", errInvalidRedisDB, dbValue)
	} else {
		db = dbInt
	}
	return
}
