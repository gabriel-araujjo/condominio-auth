package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gabriel-araujjo/condominio-auth/errors"
	"github.com/gabriel-araujjo/go-jws"
)

type oAuth2 struct {
	sign jws.Algorithm
}

func (o *oAuth2) verify(token string, scopes ...string) bool {
	// verify token scope
	// verify required scope against the database
	return true
}

func (o *oAuth2) authorize(scopes ...string) *Middleware {
	return newMiddleware(func(w http.ResponseWriter, req *http.Request) bool {
		token := strings.Split(req.Header.Get("Authorization"), " ")
		if len(token) != 2 || !strings.EqualFold(token[0], "Bearer") || !o.verify(token[1], scopes...) {
			w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, strings.Join(scopes, " ")))
			errors.WriteErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
			return true // shortcut
		}
		return false // continue
	})
}
