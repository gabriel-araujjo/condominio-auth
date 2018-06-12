package daos

import (
	"github.com/gabriel-araujjo/condominio-auth/domain"
	jp "github.com/gabriel-araujjo/json-patcher"
)

// ClientDao manage all queries related to clients
type ClientDao interface {
	Create(u *domain.Client) error
	Delete(u *domain.Client) error
	Update(u *domain.Client) error
	Get(publicID string) (*domain.Client, error)
	Auth(publicID string, secret string) (pubID string, err error)
	GetAuthorizedScopesByUser(publicID string, userID int64) domain.Scope
}

// UserDao manage all queries related to users
type UserDao interface {
	Create(u *domain.User) error
	Delete(id int64) error
	Update(id int64, patch jp.Patch) error
	Get(id int64) (*domain.User, error)
	Authenticate(credential string, password string) (int64, error)
	AuthorizeClient(userID int64, clientPublicID string, scope domain.Scope) error
	//GetAuthorizedScopeForClient(clientPublicID string) []domain.Permission
}

// TokenDao manage all queries related to users
type TokenDao interface {
	Verify(token *domain.Claims) error
	Revoke(token *domain.Claims) error
}
