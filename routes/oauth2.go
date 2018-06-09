package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gabriel-araujjo/condominio-auth/security"

	"github.com/gabriel-araujjo/condominio-auth/errors"
)

type oAuth2 struct {
	notary security.Notary
}

func (o *oAuth2) verifyTokenScope(accessToken string, scope ...string) bool {
	claims, err := o.notary.VerifyAccessToken(accessToken)
	return err == nil && claims.ContainScope(scope...)
}

func (o *oAuth2) coverScope(scopes ...string) *Middleware {
	return newMiddleware(func(w http.ResponseWriter, req *http.Request) bool {
		fields := strings.Fields(req.Header.Get("Authorization"))
		if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") || !o.verifyTokenScope(fields[1], scopes...) {
			w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, strings.Join(scopes, " ")))
			errors.WriteErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
			return true // shortcut
		}
		return false // continue
	})
}

func (o *oAuth2) revokeAccess() *Middleware {
	return newMiddleware(func(w http.ResponseWriter, req *http.Request) bool {
		fields := strings.Fields(req.Header.Get("Authorization"))
		if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
			if err := o.notary.RevokeAccessToken(fields[1]); err != nil {
				errors.WriteErrorWithCode(w, http.StatusInternalServerError, "Unexpected error")
				return true
			}
		}
		return false
	})
}
