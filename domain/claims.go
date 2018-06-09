package domain

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Claims is a jwt as defined in https://tools.ietf.org/html/rfc7519
// with some useful claims
type Claims struct {
	jwt.StandardClaims
	// AuthTime is the time when the End-User authentication occurred
	AuthTime int64 `json:"auth_time"`
	// Nonce is sent by client
	Nonce string `json:"nonce"`
	// Picture of the user
	Picture string `json:"picture"`
	// Email of the user
	Email string `json:"email"`
	// EmailVerified is whether the email is verified
	EmailVerified bool `json:"email_verified"`
	// Gender of the user
	Gender string `json:"gender"`
	// ZoneInfo is the user's time zone
	ZoneInfo string `json:"zone_info"`
	// Locale is the user's
	Locale string `json:"locale"`
	// Roles has the roles allowed
	Scope []string `json:"scope"`
}

// ContainScope checks whether this claim cover the scope passed
func (c *Claims) ContainScope(scope ...string) bool {
	for _, s := range scope {
		if binarySearch(c.Scope, s) == -1 {
			return false
		}
	}
	return true
}

func binarySearch(target_map []string, value string) int {

	left := 0
	right := len(target_map)

	for left < right {
		half := (left + right) / 2
		comp := strings.Compare(value, target_map[half])

		if comp == 0 {
			return half
		} else if comp < 0 {
			right = half - 1
		} else {
			left = half + 1
		}

	}

	return -1
}
