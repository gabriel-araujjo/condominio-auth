package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
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
	Get(token string) (userID int64, scope domain.Scope, err error)
	Add(token string, expiresAt int64, userID int64, scope domain.Scope) error
	Remove(token string) error
}

// Notary controls the bureaucracy of access tokens
type Notary struct {
	method     jwt.SigningMethod
	tokenStore TokenStore
	privateKey interface{}
	publicKey  interface{}
	codeKey    *rsa.PrivateKey
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

// VerifyAccessToken verifies if the access token is for userID and whether the scope iscovered
func (a *Notary) VerifyAccessToken(accessToken string, userID int64, scope ...string) error {
	uID, allowedScope, err := a.tokenStore.Get(accessToken)
	if err != nil {
		return err
	}

	if uID != userID {
		return errors.New("invalid access_toke")
	}

	if !allowedScope.HasSubscope(scope) {
		return fmt.Errorf("invalid scope %q", strings.Join(scope, " "))
	}

	return nil
}

// NewAccessToken generate an access or a refresh token
func (a *Notary) NewAccessToken(duration time.Duration, userID int64, scope ...string) (string, error) {
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
	return token, a.tokenStore.Add(token, time.Now().Add(duration).Unix(), userID, scope)
}

// RevokeAccessToken revokes an access token
func (a *Notary) RevokeAccessToken(accessToken string) error {
	return a.tokenStore.Remove(accessToken)
}

func (a *Notary) NewClientCode(clientID int64, scope []int64) (string, error) {
	if len(scope) > 21 {
		return "", errors.New("max of 21 scopes per code")
	}
	var message [96]byte
	binary.BigEndian.PutUint32(message[:4], uint32(clientID))
	message[4] = uint8(len(scope))
	i := 5
	for _, s := range scope {
		binary.BigEndian.PutUint32(message[i:i+4], uint32(s))
		i += 4
	}
	rng := rand.Reader
	rand.Read(message[i:])
	cipher, err := rsa.EncryptPKCS1v15(rng, &a.codeKey.PublicKey, message[:])

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(cipher), nil
}

func (a *Notary) DecipherCode(code string) (clientID int64, scope []int64, err error) {
	cipher, err := base64.URLEncoding.DecodeString(code)
	if err != nil {
		return
	}
	rng := rand.Reader
	message, err := rsa.DecryptPKCS1v15(rng, a.codeKey, cipher)
	if err != nil {
		return
	}

	scopeLen := int32(message[4])

	if scopeLen > 21 {
		return
	}

	clientID = int64(binary.BigEndian.Uint32(message[:4]))
	i := 5
	scope = make([]int64, scopeLen)
	for s := range scope {
		scope[s] = int64(binary.BigEndian.Uint32(message[i : i+4]))
		i += 4
	}

	return
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

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	return &Notary{
		method:     jwt.GetSigningMethod(config.Notary.JWTAlgorithm),
		tokenStore: tokenStore,
		privateKey: config.Notary.JWTSigningKey,
		publicKey:  config.Notary.JWTVerifyingKey,
		codeKey:    privateKey,
		closer:     closer,
	}, nil
}
