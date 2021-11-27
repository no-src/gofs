package server

import (
	"crypto/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
)

func DefaultSessionStore() (sessions.Store, error) {
	secret := make([]byte, 32)
	_, err := rand.Reader.Read(secret)
	if err != nil {
		return nil, err
	}
	store := memstore.NewStore(secret)
	return store, nil
}
