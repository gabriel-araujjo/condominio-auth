package domain

// Claims is a jwt as defined in https://tools.ietf.org/html/rfc7519
// with some useful claims
type Claims struct {
	// ID is the id of the token
	ID int64 `json:"jti"`
	// Issuer is the url of this service
	Issuer string `json:"iss"`
	// Audience is the id of the client that require the token
	Audience string `json:"aud"`
	// Subject identifier, aka user id
	Subject string `json:"sub"`
	// ExpirationTime is the end of the token's period of validity
	ExpirationTime int64 `json:"exp"`
	// NotBefore is the date from which the token is valid in unix time
	NotBefore int64 `json:"nbf"`
	// IssuedAt is the creation data of the token
	IssuedAt int64 `json:"ait"`
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
	Scope map[string]interface{} `json:"scope"`
}

func (c *Claims) Valid() error {
	panic("implement me")
}



