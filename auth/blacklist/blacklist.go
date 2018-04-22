package blacklist

import (
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

// BlackList must revoke and verify whether a token is revoked
type BlackList interface {

	// Verify whether a token is revoked returning
	// an error in case the token is revoked
	Verify(token *domain.Claims) error

	// Revoke the token. In success nil must be returned
	Revoke(token *domain.Claims) error


}