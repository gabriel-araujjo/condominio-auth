package mock

import (
	"database/sql"
	"github.com/gabriel-araujjo/versioned-database"
	"errors"
)

var someErr = errors.New("some error")

type versionStub struct {
	v       int
	getVErr error
	setVErr error
}

func (s *versionStub) Version(db *sql.DB) (int, error) {
	return s.v, s.getVErr
}

func (s *versionStub) SetVersion(db *sql.DB, v int) error {
	return s.setVErr
}

func init() {
	version.Register("version0", &versionStub{v: 0})
	version.Register("version1", &versionStub{v: 1})
	version.Register("crashed",  &versionStub{v: -1, getVErr: someErr, setVErr: someErr})
	version.Register("crashedGet",  &versionStub{v: -1, getVErr: someErr})
	version.Register("crashedSet",  &versionStub{v: 1, setVErr: someErr})
}