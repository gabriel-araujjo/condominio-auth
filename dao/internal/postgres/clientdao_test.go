package postgres

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/gabriel-araujjo/condominio-auth/dao/internal/postgres/mock"
)

func TestClientDaoPG_Get(t *testing.T) {
	conf := mock.FakeDBConfig()

	db, err := sql.Open(conf.Dao.Driver, conf.Dao.URI)
	if err != nil {
		t.Fatalf("can't prepare db for test %e", err)
	}

	cleanDB(t, db)

	_, clientDao, _, err := NewDao(conf)
	if err != nil {
		t.Fatalf("can't prepare dao for test %e", err)
	}

	t.Run("ValidClient", func(t *testing.T) {
		c, err := clientDao.Get(conf.Clients[0].PublicId)
		if err != nil {
			t.Errorf("unexpected error %q", err)
		}
		if !reflect.DeepEqual(c, conf.Clients[0]) {
			t.Errorf("expecting %#v, but %#v was returned instead", conf.Clients[0], c)
		}
	})

	t.Run("InvalidClient", func(t *testing.T) {
		c, err := clientDao.Get("-1")
		if err == nil {
			t.Error("expecting error, but nil was returned")
		}
		if c != nil {
			t.Errorf("invalid client must be nil, instead of %#v", c)
		}
	})

}
