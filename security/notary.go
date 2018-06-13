package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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
	codeCipher cipher.Block
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

// rawClientCode stores the client authorization code
//
// XX XX = client id            (4 bytes)
// N     = permission count     (1 byte )
// PP    = permission           (2 bytes)
// HH HH = hash                 (4 bytes)
// R*    = random section
//
//         01 23 4 5 67 89 AB CD EF
// 0x00    XX XX|N|R|PP|PP|PP|RR RR
// 0x10    PP|PP|P P|PP|PP|PP|RR RR
// 0x20    PP|PP|P P|PP|PP|PP|RR RR
// 0x30    PP|PP|P P|PP|PP|PP|RR RR
// 0x40    PP|PP|P P|PP|UU|UU|RR RR
// 0x50    UU|UU|H H HH HH HH|RR RR
type rawClientCode [0x60]byte

func (code rawClientCode) clientID() int64 {
	return int64(binary.BigEndian.Uint32(code[0:4]))
}

func (code rawClientCode) scopeIDs() []int64 {
	toRead := code[4]
	s := make([]int64, 0, toRead)
	i := 6
	for i < 0x60 && toRead > 0 {
		// Skip random bytes
		if i%0x10 == 0xC {
			i += 4
		}
		s = append(s, int64(binary.BigEndian.Uint16(code[i:i+2])))
		i += 2
		toRead--
	}
	return s
}

func (code rawClientCode) hash() (h []byte) {
	h = make([]byte, 8)
	copy(h, code[0x54:0x5C])
	return
}

func (code rawClientCode) userID() (uID int64) {
	uID = int64((uint64(binary.BigEndian.Uint32(code[0x48:0x4C])) << 32) | uint64(binary.BigEndian.Uint32(code[0x50:0x54])))
	return
}

func (code rawClientCode) strip() (stripped []byte) {
	toRead := code[4]
	stripped = make([]byte, 0, 4+toRead*2+8)
	stripped = append(stripped, code[0:4]...)
	i := 6
	for i < 0x60 && toRead > 0 {
		if i%0x10 == 0xC {
			i += 4
		}
		stripped = append(stripped, code[i:i+2]...)
		i += 2
		toRead--
	}
	stripped = append(stripped, code[0x48:0x4C]...)
	stripped = append(stripped, code[0x50:0x54]...)
	return
}

// NewClientCode generate a new code to be used on authorization end point
func (a *Notary) NewClientCode(clientID int64, scope []int64, userID int64) (string, error) {
	if len(scope) > 25 {
		return "", errors.New("max of 25 scopes per code")
	}
	rng := rand.Reader

	var message rawClientCode

	// write client id
	binary.BigEndian.PutUint32(message[:4], uint32(clientID))

	// write scope length
	message[4] = uint8(len(scope))

	// fifth byte is random
	rng.Read(message[5:6])

	// start permission writing at sixth byte
	i := 6
	for _, s := range scope {
		// stuff random bytes at end of each block
		if i%0x10 == 0xC {
			rng.Read(message[i : i+4])
			i += 4
		}

		// write permission
		binary.BigEndian.PutUint16(message[i:i+2], uint16(s))
		i += 2
	}

	// stuff scope space with random bytes
	rng.Read(message[i:0x48])

	// uid first 4 bytes
	binary.BigEndian.PutUint32(message[0x48:0x4C], uint32(userID>>32))
	// random bytes
	rng.Read(message[0x4C:0x50])
	// uid last 4 bytes
	binary.BigEndian.PutUint32(message[0x50:0x54], uint32(userID))

	// write hash
	hash := sha256.Sum256(message.strip())
	copy(message[0x54:0x5C], hash[0:8])

	// stuff final code random bytes
	rng.Read(message[0x5C:0x60])

	blockSize := a.codeCipher.BlockSize()
	for i = 0; i < 0x60; i += blockSize {
		a.codeCipher.Encrypt(message[i:i+blockSize], message[i:i+blockSize])
	}

	return base64.URLEncoding.EncodeToString(message[:]), nil
}

// DecipherCode get the client and the scope of a code
func (a *Notary) DecipherCode(code string) (clientID int64, scope []int64, uID int64, err error) {
	var message rawClientCode
	_, err = base64.URLEncoding.Decode(message[:], []byte(code))
	if err != nil {
		return
	}

	blockSize := a.codeCipher.BlockSize()
	for i := 0; i < 0x60; i += blockSize {
		a.codeCipher.Decrypt(message[i:i+blockSize], message[i:i+blockSize])
	}

	hash := sha256.Sum256(message.strip())

	if !bytes.Equal(hash[0:8], message.hash()) {
		err = errors.New("invalid code")
		return
	}

	clientID = message.clientID()
	scope = message.scopeIDs()
	uID = message.userID()
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

	privateKey, _ := aes.NewCipher(config.Notary.CodeCipherSecret)

	return &Notary{
		method:     jwt.GetSigningMethod(config.Notary.JWTAlgorithm),
		tokenStore: tokenStore,
		privateKey: config.Notary.JWTSigningKey,
		publicKey:  config.Notary.JWTVerifyingKey,
		codeCipher: privateKey,
		closer:     closer,
	}, nil
}
