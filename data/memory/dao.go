package memory

import (
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/data"
)

// NewDao create a dao in memory
func NewDao(conf *config.Config) (*data.Dao, error) {
	return data.NewDao(
			&userDaoMemory{},
			&clientDaoMemory{},
			nil),
		nil
}
