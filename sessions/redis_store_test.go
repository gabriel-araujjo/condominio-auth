package sessions

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/sessions/test"
)

func testConfig() *config.Config {
	return &config.Config{
		Session: config.Session{
			StoreType:  "redis",
			StoreURI:   "redis:",
			PoolSize:   10,
			CookieName: "Sessions-test",
			HashKey:    []byte("z&^LcHF68aV(VXqU%iLWX!sfgDc7AokASiH8YJMK&u%VVo4hHkB&kr"),
		},
	}
}
func TestRedisConnection(t *testing.T) {
	store, err := NewRedisStore(testConfig())
	if err != nil {
		t.Fatalf("store can't be created due to %q", err.Error())
	}
	defer store.Close()

	testURL, _ := url.Parse("https://condominio.com/auth")
	req := &http.Request{
		Method: "GET",
		URL:    testURL,
	}
	mockWriter := test.NewMoockWriter()
	session, _ := store.Get(req, "some-name")
	session.Set("key", "value2")
	err = store.Save(req, mockWriter, session)
	if err != nil {
		t.Fatalf("store can't save session due to %q", err.Error())
	}

	returnedSession, err := store.Get(req, "some-name")
	if err != nil {
		t.Fatalf("store can't retrieve session due to %q", err.Error())
	}

	if !reflect.DeepEqual(returnedSession, session) {
		t.Fatalf("different sessions: %#v and %#v", returnedSession, session)
	}

}
