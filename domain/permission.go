package domain

import (
	"strings"
)

// Permission describe a permission which an user can grant to a clients
type Permission struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ClientAuthorizationRequest stores data of authorization requests
type ClientAuthorizationRequest struct {
	clientPublicID string
	Scopes         []string
}

// NewClientAuthorizationRequest create an authorization
func NewClientAuthorizationRequest(clientPublicID string, scopes string) *ClientAuthorizationRequest {
	return &ClientAuthorizationRequest{clientPublicID, strings.Split(scopes, " ")}
}
