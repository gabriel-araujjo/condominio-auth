package data

import (
	"github.com/gabriel-araujjo/condominio-auth/domain"
	"io"
	jp "github.com/gabriel-araujjo/json-patcher"
)

type Dao struct {
	closer io.Closer
	User UserDao
	Client ClientDao
}

func NewDao(user UserDao, client ClientDao, closer io.Closer) *Dao {
	return &Dao{closer:closer, User: user, Client: client}
}

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

func (d *Dao) Close() error {
	if d.closer != nil {
		return d.closer.Close()
	}
	return nil
}
