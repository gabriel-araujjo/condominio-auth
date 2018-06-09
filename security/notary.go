package security

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

type Blacklist interface {
	Contains(*jwt.Token) bool
}

type Notary struct {
	method     jwt.SigningMethod
	blacklist  Blacklist
	privateKey interface{}
	publicKey  interface{}
}

func (a *Notary) SignClaims(claims *domain.Claims) string {
	claims.ExpiresAt = time.Now().Add(30 * 24 * time.Hour).Unix()
	claims.NotBefore = time.Now().Unix()
	res, _ := jwt.NewWithClaims(a.method, claims).SignedString(a.privateKey)
	return res
}

func (a *Notary) VerifyClaims(tokenString string) (*domain.Claims, error) {

	var claims domain.Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != a.method.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}

		if a.blacklist.Contains(token) {
			return nil, fmt.Errorf("revoked token")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return a.publicKey, nil
	})
	return &claims, err
}
