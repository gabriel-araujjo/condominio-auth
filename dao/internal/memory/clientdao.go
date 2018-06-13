package memory

import (
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

type clientDaoMemory struct {
}

func (d *clientDaoMemory) Create(u *domain.Client) error {
	panic("implement me")
}

func (d *clientDaoMemory) Delete(u *domain.Client) error {
	panic("implement me")
}

func (d *clientDaoMemory) Update(u *domain.Client) error {
	panic("implement me")
}

func (d *clientDaoMemory) Get(publicID string) (*domain.Client, error) {
	panic("implement me")
}

func (d *clientDaoMemory) Auth(publicID string, secret string) (string, error) {
	panic("implement me")
}

func (d *clientDaoMemory) GetAuthorizedScopesByUser(publicID string, userID int64) domain.Scope {
	//TODO:
	return nil
}
