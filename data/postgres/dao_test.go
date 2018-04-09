package postgres

import (
	"github.com/gabriel-araujjo/condominio-auth/data"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"github.com/gabriel-araujjo/condominio-auth/data/postgres/mock"
	"github.com/gabriel-araujjo/versioned-database"
	"database/sql"
)

func TestNewDao(t *testing.T) {

	t.Run("SchemeOK", func(t *testing.T) {
		db, m, _ := sqlmock.New()

		// version.PersistScheme(db, scheme) always make a transaction
		expectCommitedTx(&m)
		dao, err := newDaoInternal(db, mock.SchemeOK())

		want := data.NewDao(&userDaoPG{db:db}, &clientDaoPG{db:db}, db)
		if !reflect.DeepEqual(dao, want) {
			t.Errorf("newdao: got=%v, want=%v", dao, want)
		}

		if err != nil {
			t.Errorf("newdao: err should be nil instead of %q", err)
		}
	})

	t.Run("CrashedScheme", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		expectRollbackTx(&m)
		m.ExpectClose()
		dao, err := newDaoInternal(db, mock.CrashedScheme())

		if dao != nil {
			t.Error("newdao: dao should be nil")
		}

		if err == nil {
			t.Error("newdao: err should be nil")
		}

		if err := m.ExpectationsWereMet(); err != nil {
			t.Error("newdao: db should be closed")
		}
	})

	t.Run("Integration", func(t *testing.T) {
		tests := []struct{
			name string
			mocker func() (version.Scheme, func() bool)
		}{{
			name: "SchemeCreation",
			mocker: mock.CreationScheme,
		},{
			name: "SchemeUpdate",
			mocker: mock.UpdatingScheme,
		}}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				db, m, _ := sqlmock.New()
				expectCommitedTx(&m)
				scheme, check := tt.mocker()
				dao, err := newDaoInternal(db, scheme)

				if dao == nil {
					t.Error("newdao: dao should not be nil")
				}

				if err != nil {
					t.Errorf("newdao: err should be nil instead of %q", err)
				}

				if !check() {
					t.Error("newdao: calls not checked")
				}
			})
		}
	})

	t.Run("FakeDatabase", func(t *testing.T) {
		dao, err := NewDao(mock.FakeDBConfig())
		if err != nil {
			t.Errorf("newdao: err should be nil instead of %q", err)
		}

		if dao == nil {
			t.Fatal("newdao: dao should not be nil")
		}


		cleanDB(t, dao.User.(*userDaoPG).db)
		dao.Close()
	})

	t.Run("InvalidDriver", func(t *testing.T) {
		dao, err := NewDao(mock.InvalidDBConfig())
		if dao != nil {
			t.Fatal("newdao: dao should be nil")
		}

		if err == nil {
			t.Error("newdao: err should not be nil")
		}
	})
}

func expectCommitedTx(m *sqlmock.Sqlmock)  {
	(*m).ExpectBegin()
	(*m).ExpectCommit()
}

func expectRollbackTx(m *sqlmock.Sqlmock)  {
	(*m).ExpectBegin()
	(*m).ExpectRollback()
}

func cleanDB(t *testing.T, db *sql.DB) {
	if db == nil {
		return
	}
	_, err := db.Exec("COMMENT ON DATABASE test IS NULL")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
DROP SCHEMA "public" CASCADE;
CREATE SCHEMA "public";
GRANT ALL ON SCHEMA "public" TO postgres;
GRANT ALL ON SCHEMA "public" TO public;
`)

	if err != nil {
		t.Fatal(err)
	}
}