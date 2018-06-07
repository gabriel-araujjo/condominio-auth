package sessions

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gabriel-araujjo/condominio-auth/config"
)

// Session interface allows access and editions of values
type Session interface {
	// Set creates or updates a value.
	// To persists alterations call Store#Save() with this session as a parameter
	Set(key string, value interface{})
	// Get provides the value stored at key or nil
	Get(key string) interface{}
}

// Store is a gorilla sessions.Store that is a closer
type Store interface {
	io.Closer
	// Get should return a cached session.
	Get(r *http.Request, name string) (Session, error)

	// New should create and return a new session.
	//
	// Note that New should never return a nil session, even in the case of
	// an error if using the Registry infrastructure to cache the session.
	New(r *http.Request, name string) (Session, error)

	// Save should persist session to the underlying store implementation.
	Save(r *http.Request, w http.ResponseWriter, s Session) error
}

// NewStoreFromConfig should creates a new Store from config
func NewStoreFromConfig(config *config.Config) (Store, error) {
	switch config.Session.StoreType {
	case "redis":
		return NewRedisStore(config)
	default:
		return nil, fmt.Errorf("invalid sessions sore type: %q", config.Session.StoreType)
	}
}
