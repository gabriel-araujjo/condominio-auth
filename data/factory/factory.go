package factory

import (
	"errors"
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/data"
	"github.com/gabriel-araujjo/condominio-auth/data/postgres"
	"github.com/gabriel-araujjo/condominio-auth/data/memory"
)

func NewDao(config *config.Config) (*data.Dao, error) {
	switch config.Dao.Driver {
	case "postgres":
		return postgres.NewDao(config)
	case "memory":
		return memory.NewDao(config)
	default:
		return nil, errors.New("dao_factory: Undefined DB driver")
	}
}
