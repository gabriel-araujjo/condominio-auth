package routes

import (
	"net/http"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/dao"
	"github.com/gabriel-araujjo/condominio-auth/sessions"
)

const userKey = "user"

type context struct {
	sessionName   string
	dao           *dao.Dao
	sessionsStore sessions.Store
}

func newContext(conf *config.Config, dao *dao.Dao, sessionsStore sessions.Store) *context {
	return &context{conf.Session.CookieName, dao, sessionsStore}
}

func (c *context) Session(req *http.Request) (sessions.Session, error) {
	return c.sessionsStore.Get(req, c.sessionName)
}

func (c *context) PersistSession(req *http.Request, w http.ResponseWriter) error {
	s, err := c.Session(req)
	if err != nil {
		return err
	}
	return c.sessionsStore.Save(req, w, s)
}

func (c *context) SetCurrentUserID(req *http.Request, userID int64) error {
	session, err := c.Session(req)
	if err != nil {
		return err
	}
	session.Set(userKey, userID)
	return nil
}

func (c *context) CurrentUserID(req *http.Request) (int64, error) {
	session, err := c.Session(req)
	if err != nil {
		return 0, err
	}
	userID := session.Get(userKey)
	if userID == nil {
		return 0, nil
	}
	return userID.(int64), nil
}
