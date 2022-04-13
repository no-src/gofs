package server

import (
	"crypto/rand"
	
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
)

// DefaultSessionStore create a default session store, current is stored in memory
func DefaultSessionStore() (sessions.Store, error) {
	secret := make([]byte, 32)
	_, err := rand.Reader.Read(secret)
	if err != nil {
		return nil, err
	}
	store := memstore.NewStore(secret)
	return store, nil
}
