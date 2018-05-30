package redis

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/session/internal/redis/mock"
	"github.com/gorilla/sessions"
)

func testConfig() *config.Config {
	return &config.Config{
		Sessions: struct {
			StoreType          string
			StoreURI           string
			MaxConnections     int
			Name               string
			CookieCodecHashKey []byte
		}{
			StoreType:          "redis",
			StoreURI:           "redis:///",
			MaxConnections:     10,
			Name:               "Sessions-test",
			CookieCodecHashKey: []byte("z&^LcHF68aV(VXqU%iLWX!sfgDc7AokASiH8YJMK&u%VVo4hHkB&kr"),
		},
	}
}
func TestRedisConnection(t *testing.T) {
	store, _, err := NewStore(testConfig())
	if err != nil {
		t.Fatal("err must be not nil")
	}

	testURL, _ := url.Parse("https://condominio.com/auth")
	req := &http.Request{
		Method: "GET",
		URL:    testURL,
	}
	mockWriter := &mock.MockResponseWriter{}
	session := sessions.NewSession(store, "SomeName")
	session.Values["key"] = "value"
	err = store.Save(req, mockWriter, session)
	if err != nil {
		t.Fatal("store can't save session")
	}

	returnedSession, err := store.Get(req, session.Name())
	if err != nil {
		t.Fatal("store can't retrieve session")
	}

	if !reflect.DeepEqual(returnedSession, session) {
		t.Fatal("different sessions")
	}

}
