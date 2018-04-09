package auth

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"time"
	"errors"
)


type Auth struct {
	method jwt.SigningMethod
	privateKey interface{}
	publicKey interface{}
}

func (a *Auth) Sign(token jwt.MapClaims) string {
	token["exp"] = time.Now().Add(30 * 24 * time.Hour).Unix()
	token["nbf"] = time.Now().Unix()
	res, _ := jwt.NewWithClaims(a.method, token).SignedString(a.privateKey)
	return res
}

func (a *Auth) Verify(tokenString string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return a.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	nbtHead, hasNbt := token.Header["nbt"]

	if !hasNbt {
		return nil, errors.New("expired")
	}

	nbt, isInt64 := nbtHead.(*int64)

	if !isInt64 || time.Now().Before(time.Unix(*nbt, 0)) {
		return nil, errors.New("expired")
	}

	expHead, hasExp := token.Header["exp"]

	if !hasExp {
		return nil, errors.New("expired")
	}

	exp, isInt64 := expHead.(*int64)

	if !isInt64 || time.Now().After(time.Unix(*exp, 0)) {
		return nil, errors.New("expired")
	}

	return token, nil
}