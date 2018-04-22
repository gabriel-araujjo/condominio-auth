package dao

import (
	"io"
	"github.com/gabriel-araujjo/condominio-auth/dao/daos"
	"github.com/gabriel-araujjo/condominio-auth/dao/internal/postgres"
	"github.com/gabriel-araujjo/condominio-auth/dao/internal/memory"
	"github.com/gabriel-araujjo/condominio-auth/config"
	"errors"
)

// Dao contains all Domain related daos
type Dao struct {
	closer io.Closer
	User   daos.UserDao
	Client daos.ClientDao
}

// Close the dao connections, this method must be
// called when the Dao is not necessary anymore.
func (d *Dao) Close() error {
	if d.closer != nil {
		return d.closer.Close()
	}
	return nil
}

// NewFromConfig create a new Dao following a Config specification
func NewFromConfig(config *config.Config) (*Dao, error) {
	d := Dao{}
	var err error
	switch config.Dao.Driver {
	case "postgres":
		d.User, d.Client, d.closer, err = postgres.NewDao(config)
	case "memory":
		d.User, d.Client, err = memory.NewDao(config)
	default:
		err = errors.New("dao: Undefined DB driver")
	}

	if err != nil {
		return nil, err
	}

	return &d, nil
}
