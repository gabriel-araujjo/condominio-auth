package postgres

import (
	"testing"
	"github.com/gabriel-araujjo/condominio-auth/dao/internal/postgres/postgres/mock"
	"database/sql"
	"fmt"
	"github.com/gabriel-araujjo/condominio-auth/config"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"errors"
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

var validCpfs = []int64{
	95115469103,
}
var invalidCpfs = []int64{
	0,
	11111111111,
	22222222222,
	33333333333,
	44444444444,
	55555555555,
	66666666666,
	77777777777,
	88888888888,
	99999999999,
}
var tableNames = []string{
	"user",
	"user_email",
	"user_phone",
	"client",
}

var someErr = errors.New("some err")

func TestScheme(t *testing.T) {
	t.Run("OnCreate", func(t *testing.T) {
		conf := mock.FakeDBConfig()
		db, err := sql.Open(conf.Dao.Driver, conf.Dao.DNS)

		cleanDB(t, db)

		if err != nil {
			t.Fatalf("can't initialize Fake DB. err = %q", err)
		}
		s := scheme{conf:conf}

		err = s.OnCreate(db)

		if err != nil {
			t.Errorf("err should be nil instead of %q", err)
		}

		for _, cpf := range validCpfs {
			t.Run(fmt.Sprintf("ValidCpf:%d", cpf), func(t *testing.T) {
				var valid bool
				err := db.QueryRow(fmt.Sprintf(`SELECT check_cpf(%d)`, cpf)).Scan(&valid)
				if err != nil {
					t.Errorf("err should be nil instead of %e", err)
				}
				if !valid {
					t.Error("cpf should be valid", err)
				}
			})
		}

		for _, cpf := range invalidCpfs {
			t.Run(fmt.Sprintf("InvalidCpf:%d", cpf), func(t *testing.T) {
				var valid bool
				err := db.QueryRow(fmt.Sprintf(`SELECT check_cpf(%d)`, cpf)).Scan(&valid)
				if err != nil {
					t.Errorf("err should be nil instead of %e", err)
				}
				if valid {
					t.Error("cpf should be invalid", err)
				}
			})
		}

		for _, table := range tableNames {
			t.Run(fmt.Sprintf("Table:%q", table), func(t *testing.T) {
				var exists bool
				err := db.QueryRow(`SELECT 1 from information_schema.tables WHERE table_name = $1`, table).Scan(&exists)
				if err != nil {
					t.Errorf("err should be nil instead of %e", err)
				}
				if !exists {
					t.Errorf("table %q should exist", table)
				}
			})
		}

		t.Run("PreRegisteredClients", func(t *testing.T) {
			for _, client := range conf.Clients {
				rows, err := db.Query(`SELECT * FROM "client" WHERE public_id = $1`, client.PublicId)
				if err != nil {
					t.Error(err)
				}
				if rows == nil {
					t.Fatalf("no client was returned for public_id = %q", client.PublicId)
				}
				if !rows.Next() {
					t.Errorf("No entry for client %q", client.Name)
				}
				if rows.Next() {
					t.Errorf("2 entries for client %q", client.Name)
				}
				rows.Close()
			}
		})

		cleanDB(t, db)
		db.Close()

	})

	t.Run("ExecutionError", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		m.ExpectExec(".*").WillReturnError(someErr)

		s := &scheme{}
		if err:= s.OnCreate(db); err != someErr {
			t.Errorf("should return 'some error' instead of %q", err)
		}
	})

	t.Run("ClientInsertionError", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		r := sqlmock.NewResult(0, 0)
		m.ExpectExec(".*").WillReturnResult(r)
		m.ExpectQuery("INSERT.*").WillReturnError(someErr)

		conf := &config.Config{
			Clients:[]*domain.Client {{
				Name:     "fake1",
				PublicId: "1",
				Secret:   "1",
			}},
		}
		s := &scheme{conf:conf}
		if err:= s.OnCreate(db); err != someErr {
			t.Errorf("should return 'some error' instead of %q", err)
		}
	})

	t.Run("OnUpdate", func(t *testing.T) {
		s := scheme{conf:nil}
		err := s.OnUpdate(nil, 0)

		if err != nil {
			t.Errorf("err should be nil instead of %q", err)
		}
	})

	t.Run("Version", func(t *testing.T) {
		s := scheme{conf:nil}
		if s.Version() <= 0 {
			t.Errorf("scheme.Version() should return a integer greater then zero instead of %d", s.Version())
		}
		err := s.OnUpdate(nil, 0)

		if err != nil {
			t.Errorf("err should be nil instead of %q", err)
		}
	})

	t.Run("VersionStrategy", func(t *testing.T) {
		conf := &config.Config{Dao: struct {
			Driver          string
			DNS             string
			VersionStrategy string
		}{VersionStrategy: "strategy"}}
		s := scheme{conf: conf}

		if s.VersionStrategy() != conf.Dao.VersionStrategy {
			t.Errorf("strategy must be %q instead of %q", conf.Dao.VersionStrategy, s.VersionStrategy())
		}
	})

}
