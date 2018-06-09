package security

import (
	"io"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gomodule/redigo/redis"
)

type redisBlacklist struct {
	conn redis.Conn
}

func (b *redisBlacklist) Add(tokenSignature string, expiresAt int64) error {
	_, err := b.conn.Do("SET", tokenSignature, true)
	if err != nil {
		return err
	}
	_, err = b.conn.Do("EXPIREAT", tokenSignature, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (b *redisBlacklist) Contains(tokenSignature string) (bool, error) {
	reply, err := b.conn.Do("GET", tokenSignature)
	return reply == nil, err
}

func newRedisBlackist(config *config.Config) (Blacklist, io.Closer, error) {
	conn, err := redis.DialURL(config.Notary.BlacklistURI)
	if err != nil {
		return nil, nil, err
	}
	return &redisBlacklist{conn}, conn, nil
}
