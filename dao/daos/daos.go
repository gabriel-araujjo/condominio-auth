package daos

import (
	"github.com/gabriel-araujjo/condominio-auth/domain"
	jp "github.com/gabriel-araujjo/json-patcher"
)

type ClientDao interface {
	Create(u *domain.Client) error
	Delete(u *domain.Client) error
	Update(u *domain.Client) error
	Get(publicId string) (*domain.Client, error)
	Auth(publicId string, secret string) (pubId string, err error)
}

type UserDao interface {
	Create(u *domain.User) error
	Delete(id int64) error
	Update(id int64, patch jp.Patch) error
	Get(id int64) (*domain.User, error)
	Auth(credential string, password string) (int64, error)
}

type TokenDao interface {
	Verify(token *domain.Claims) error
	Revoke(token *domain.Claims) error
}
