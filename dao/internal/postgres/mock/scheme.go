package mock

import (
	"database/sql"
	"github.com/gabriel-araujjo/versioned-database"
)

type Scheme struct {
	Vers   func() int
	Strat  func() string
	Create func(*sql.DB) error
	Update func(*sql.DB, int) error
}

func (m *Scheme) Version() int {
	return m.Vers()
}

func (m *Scheme) VersionStrategy() string {
	return m.Strat()
}

func (m *Scheme) OnCreate(db *sql.DB) error {
	return m.Create(db)
}

func (m *Scheme) OnUpdate(db *sql.DB, oldVersion int) error {
	return m.Update(db, oldVersion)
}

func SchemeOK() version.Scheme {
	return &Scheme{
		Vers: func() int { return 1 },
		Strat: func() string { return "version1" },
	}
}

func CrashedScheme() version.Scheme {
	return &Scheme{
		Vers: func() int { return 1 },
		Strat: func() string { return "version0" },
		Create: func(db *sql.DB) error { return someErr },
	}
}

func CreationScheme() (version.Scheme, func() bool) {
	called := false
	return &Scheme{
		Vers: func() int { return 1 },
		Strat: func() string { return "version0" },
		Create: func(*sql.DB) error {
			called = true
			return nil
		},
	}, func() bool {
		return called
	}
}

func UpdatingScheme() (version.Scheme, func() bool) {
	called := false
	return &Scheme{
		Vers: func() int { return 2 },
		Strat: func() string { return "version1" },
		Update: func(*sql.DB, int) error {
			called = true
			return nil
		},
	}, func() bool {
		return called
	}
}



