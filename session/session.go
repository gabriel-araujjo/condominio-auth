package session

import (
	"fmt"
	"io"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/session/internal/redis"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Session struct {
	name   string
	store  sessions.Store
	closer io.Closer
}

func (s *Session) Middleware() gin.HandlerFunc {
	return sessions.Sessions(s.name, s.store)
}

func (s *Session) Close() error {
	if s.closer == nil {
		return nil
	}
	return s.closer.Close()
}

func NewFromConfig(config *config.Config) (*Session, error) {
	var (
		store  sessions.Store
		closer io.Closer
		err    error
	)
	switch config.Sessions.StoreType {
	case "redis":
		store, closer, err = redis.NewStore(config)
	default:
		err = fmt.Errorf("invalid sessions sore type: %q", config.Sessions.StoreType)
	}

	if err != nil {
		return nil, err
	}
	return &Session{config.Sessions.Name, store, closer}, nil
}
