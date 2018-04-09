package postgres

import (
	"database/sql"
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/data"
	_ "github.com/gabriel-araujjo/psql-versioning"
	"github.com/gabriel-araujjo/versioned-database"
	_ "github.com/lib/pq"
)

func NewDao(conf *config.Config) (*data.Dao, error) {
	db, err := sql.Open(conf.Dao.Driver, conf.Dao.DNS)
	if err != nil {
		return nil, err
	}
	return newDaoInternal(db, &scheme{conf: conf})
}

func newDaoInternal(db *sql.DB, scheme version.Scheme) (*data.Dao, error) {
	err := version.PersistScheme(db, scheme)
	if err != nil {
		db.Close()
		return nil, err
	}

	return data.NewDao(&userDaoPG{db: db}, &clientDaoPG{db: db}, db), nil
}
