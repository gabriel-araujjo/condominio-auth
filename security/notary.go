package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

//TODO: Add refresh token support

// TokenStore is a map with a token and its scopes
type TokenStore interface {
	// tokenSignature param is the third part of the token
	// the base64 encoded signature bytes
	Contains(token string) (bool, error)
	Get(token string) ([]string, error)
	Add(token string, expiresAt int64, scope ...string) error
	Remove(token string) error
}

// Notary controls the bureaucracy of access tokens
type Notary struct {
	method     jwt.SigningMethod
	tokenStore TokenStore
	privateKey interface{}
	publicKey  interface{}
	closer     io.Closer
}

// NewIDTokenWithClaims creates a new access token with the especified claims
func (a *Notary) NewIDTokenWithClaims(claims *domain.Claims) string {
	claims.ExpiresAt = time.Now().Add(30 * 24 * time.Hour).Unix()
	claims.NotBefore = time.Now().Unix()
	res, _ := jwt.NewWithClaims(a.method, claims).SignedString(a.privateKey)
	return res
}

// VerifyIDToken checks the access token signature and whether the token is revoked
func (a *Notary) VerifyIDToken(tokenString string) (*domain.Claims, error) {

	var claims domain.Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != a.method.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return a.publicKey, nil
	})
	return &claims, err
}

// NewAccessToken generate an access or a refresh token
func (a *Notary) NewAccessToken(duration time.Duration, scope ...string) (string, error) {
	var (
		tokenBytes [33]byte
		token      string
	)
	for {
		binary.BigEndian.PutUint64(tokenBytes[:8], uint64(time.Now().Unix()))
		rand.Read(tokenBytes[8:])
		hash := sha256.Sum256(tokenBytes[:])
		copy(tokenBytes[:32], hash[:])
		token = base64.StdEncoding.EncodeToString(tokenBytes[:])
		contains, err := a.tokenStore.Contains(token)
		if err != nil {
			return "", err
		}
		if !contains {
			break
		}
	}
	return token, a.tokenStore.Add(token, time.Now().Add(duration).Unix(), scope...)
}

// RevokeAccessToken revokes an access token
func (a *Notary) RevokeAccessToken(accessToken string) error {
	return a.tokenStore.Remove(accessToken)
}

// Close closes any remain connection
func (a *Notary) Close() error {
	return a.closer.Close()
}

// NewNotary creates a notary following config specs
func NewNotary(config *config.Config) (*Notary, error) {
	var (
		tokenStore TokenStore
		closer     io.Closer
		err        error
	)
	switch config.Notary.TokenStoreType {
	case "redis":
		tokenStore, closer, err = newRedisTokenStore(config)
	default:
		return nil, errors.New("invalid TokenStoreType")
	}
	if err != nil {
		return nil, err
	}

	return &Notary{
		method:     jwt.GetSigningMethod(config.Notary.JWTAlgorithm),
		tokenStore: tokenStore,
		privateKey: config.Notary.JWTSigningKey,
		publicKey:  config.Notary.JWTVerifyingKey,
		closer:     closer,
	}, nil
}
