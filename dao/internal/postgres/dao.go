package postgres

import (
	"database/sql"
	"io"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/dao/daos"
	// Versioning strategy
	_ "github.com/gabriel-araujjo/psql-versioning"
	"github.com/gabriel-araujjo/versioned-database"
	// Postgres driver
	_ "github.com/lib/pq"
)

// NewDao creates a dao following the passed config
func NewDao(conf *config.Config) (daos.UserDao, daos.ClientDao, daos.PermissionDao, io.Closer, error) {
	db, err := sql.Open(conf.Dao.Driver, conf.Dao.URI)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return newDaoInternal(db, &scheme{conf: conf})
}

func newDaoInternal(db *sql.DB, scheme version.Scheme) (daos.UserDao, daos.ClientDao, daos.PermissionDao, io.Closer, error) {
	err := version.PersistScheme(db, scheme)
	if err != nil {
		db.Close()
		return nil, nil, nil, nil, err
	}

	return &userDaoPG{db: db}, &clientDaoPG{db: db}, newPGPermissionDao(db), db, nil
}
