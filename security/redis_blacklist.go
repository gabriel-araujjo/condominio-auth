package security

import (
	"errors"
	"io"

	"github.com/gabriel-araujjo/condominio-auth/domain"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gomodule/redigo/redis"
)

type redisTokenStore struct {
	pool *redis.Pool
}

func (b *redisTokenStore) Add(token string, expiresAt int64, userID int64, scope domain.Scope) error {
	conn := b.pool.Get()
	defer conn.Close()
	args := make([]interface{}, len(scope)+2)
	args[0] = token
	args[1] = userID
	for i := range scope {
		args[i+1] = scope[i]
	}
	conn.Send("SET", args...)
	conn.Send("EXPIREAT", token, expiresAt)
	return conn.Flush()
}

func (b *redisTokenStore) Get(token string) (userID int64, scopes domain.Scope, err error) {
	conn := b.pool.Get()
	defer conn.Close()
	result, err := conn.Do("SMEMBERS", token)
	if err != nil {
		return
	}
	setMembers, isArray := result.([]interface{})
	if !isArray {
		err = errors.New("unexpected type")
		return
	}

	scopes = make([]string, len(setMembers)-1)
	var isString bool
	userID = setMembers[0].(int64)
	for i := range setMembers[1:] {
		scopes[i], isString = setMembers[i].(string)
		if !isString {
			err = errors.New("unexpected type")
			return
		}
	}
	return
}

func (b *redisTokenStore) Contains(token string) (bool, error) {
	conn := b.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("EXISTS", token)
	exists, _ := reply.(string)
	return exists == "1", err
}

func (b *redisTokenStore) Remove(token string) error {
	conn := b.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", token)
	return err
}

func newRedisTokenStore(config *config.Config) (TokenStore, io.Closer, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(config.Notary.TokenStoreURI)
		},
		MaxIdle: 10,
	}
	return &redisTokenStore{pool: pool}, pool, nil
}
