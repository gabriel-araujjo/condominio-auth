package memory

import (
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/dao/daos"
)

// NewDao create a dao in memory
func NewDao(conf *config.Config) (daos.UserDao, daos.ClientDao, daos.PermissionDao, error) {
	return &userDaoMemory{}, &clientDaoMemory{}, &permissionDaoMemory{}, nil
}
