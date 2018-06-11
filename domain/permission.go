package domain

import (
	"strings"
)

type Scope []string

// ClientAuthorizationRequest stores data of authorization requests
type ClientAuthorizationRequest struct {
	clientPublicID string
	Scopes         Scope
}

// NewClientAuthorizationRequest create an authorization
func NewClientAuthorizationRequest(clientPublicID string, scopes string) *ClientAuthorizationRequest {
	return &ClientAuthorizationRequest{clientPublicID, strings.Split(scopes, " ")}
}
