package sessions

import (
	"errors"
	"io"
	"net/http"

	gsessions "github.com/gorilla/sessions"
)

// AdaptGorillaStore wrapes a gorilla sessions.Store into a sessions.Store
func AdaptGorillaStore(store gsessions.Store, closer io.Closer) Store {
	return &gorillaStoreAdapter{store, closer}
}

type gorillaStoreAdapter struct {
	store  gsessions.Store
	closer io.Closer
}

func (s *gorillaStoreAdapter) Close() error {
	if s.closer == nil {
		return nil
	}
	return s.closer.Close()
}

func (s *gorillaStoreAdapter) Get(r *http.Request, name string) (Session, error) {
	gs, err := s.store.Get(r, name)
	if err != nil {
		return nil, err
	}
	return (*gorillaSession)(gs), nil
}

func (s *gorillaStoreAdapter) New(r *http.Request, name string) (Session, error) {
	gs, err := s.store.New(r, name)
	if err != nil {
		return nil, err
	}
	return (*gorillaSession)(gs), nil
}

func (s *gorillaStoreAdapter) Save(r *http.Request, w http.ResponseWriter, session Session) error {
	gs, isGorillaSession := session.(*gorillaSession)
	if !isGorillaSession {
		return errors.New("session type not supported, only <github.com/gorilla/sessions.go#Session>s are supported")
	}
	return s.store.Save(r, w, (*gsessions.Session)(gs))
}

type gorillaSession gsessions.Session

func (gs *gorillaSession) Set(key string, value interface{}) {
	if gs == nil {
		return
	}
	gs.Values[key] = value
}

func (gs *gorillaSession) Get(key string) interface{} {
	if gs == nil {
		return nil
	}
	return gs.Values[key]
}
