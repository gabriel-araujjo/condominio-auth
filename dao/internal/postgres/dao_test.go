package postgres

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"./mock"
	"github.com/gabriel-araujjo/versioned-database"
	"database/sql"
)

func TestNewDao(t *testing.T) {

	// It should return a userDaoPG a clientDaoPG.
	// The passed database must be returned as a closer
	// Error must be nil
	t.Run("SchemeOK", func(t *testing.T) {
		db, m, _ := sqlmock.New()

		// version.PersistScheme(db, scheme) always make a transaction
		expectCommitedTx(&m)
		userDao, clientDao, closer, err := newDaoInternal(db, mock.SchemeOK())

		wantedUserDao := &userDaoPG{db:db}
		wantedClientDao := &clientDaoPG{db:db}

		if err != nil {
			t.Errorf("newdao: err should be nil instead of %q", err)
		}

		if !reflect.DeepEqual(userDao, wantedUserDao) {
			t.Errorf("newdao: wrong user dao got=%v, want=%v", userDao, wantedUserDao)
		}

		if !reflect.DeepEqual(clientDao, wantedClientDao) {
			t.Errorf("newdao: wrong client dao got=%v, want=%v", clientDao, wantedClientDao)
		}

		if db != closer {
			t.Error("newdao: passed db is not closer")
		}
	})

	// It should return an error and close db
	t.Run("CrashedScheme", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		expectRollbackTx(&m)
		m.ExpectClose()
		_, _, _, err := newDaoInternal(db, mock.CrashedScheme())

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
				_, _, _, err := newDaoInternal(db, scheme)

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
		userDao, clientDao, db, err := NewDao(mock.FakeDBConfig())
		if err != nil {
			t.Errorf("newdao: err should be nil instead of %q", err)
		}

		if userDao == nil || clientDao == nil || db == nil {
			t.Fatal("newdao: dao should not be nil")
		}

		cleanDB(t, db.(*sql.DB))
	})

	t.Run("InvalidDriver", func(t *testing.T) {
		_, _, _, err := NewDao(mock.InvalidDBConfig())

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