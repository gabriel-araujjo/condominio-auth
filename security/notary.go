package security

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

// Blacklist is the set of signatures of revoked tokens
type Blacklist interface {
	// tokenSignature param is the third part of the token
	// the base64 encoded signature bytes
	Contains(tokenSignature string) (bool, error)
	Add(tokenSignature string, expiresAt int64) error
}

// Notary controls the bureaucracy of access tokens
type Notary struct {
	method     jwt.SigningMethod
	blacklist  Blacklist
	privateKey interface{}
	publicKey  interface{}
	closer     io.Closer
}

// NewAccessTokenWithClaims creates a new access token with the especified claims
func (a *Notary) NewAccessTokenWithClaims(claims *domain.Claims) string {
	claims.ExpiresAt = time.Now().Add(30 * 24 * time.Hour).Unix()
	claims.NotBefore = time.Now().Unix()
	res, _ := jwt.NewWithClaims(a.method, claims).SignedString(a.privateKey)
	return res
}

// VerifyAccessToken checks the access token signature and whether the token is revoked
func (a *Notary) VerifyAccessToken(tokenString string) (*domain.Claims, error) {

	var claims domain.Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != a.method.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		contains, err := a.blacklist.Contains(token.Signature)
		if contains || err != nil {
			return nil, fmt.Errorf("revoked token")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return a.publicKey, nil
	})
	return &claims, err
}

// RevokeAccessToken revokes an access token and notify other microservices
func (a *Notary) RevokeAccessToken(accessToken string) error {

	claims, err := a.VerifyAccessToken(accessToken)
	if err != nil {
		return err
	}
	parts := strings.Split(accessToken, ".")

	// TODO: notify the revocation
	return a.blacklist.Add(parts[2], claims.ExpiresAt)
}

// Close closes any remain connection
func (a *Notary) Close() error {
	return a.closer.Close()
}

// NewNotary creates a notary following config specs
func NewNotary(config *config.Config) (*Notary, error) {
	var (
		blacklist Blacklist
		closer    io.Closer
		err       error
	)
	switch config.Notary.BlacklistType {
	case "redis":
		blacklist, closer, err = newRedisBlackist(config)
	}
	if err != nil {
		return nil, err
	}

	return &Notary{
		method:     jwt.GetSigningMethod(config.Notary.JWTAlgorithm),
		blacklist:  blacklist,
		privateKey: config.Notary.JWTSigningKey,
		publicKey:  config.Notary.JWTVerifyingKey,
		closer:     closer,
	}, nil
}
