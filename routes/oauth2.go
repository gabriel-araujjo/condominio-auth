package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gabriel-araujjo/condominio-auth/security"

	"github.com/gabriel-araujjo/condominio-auth/errors"
)

type oAuth2 struct {
	*context
	notary security.Notary
}

func (o *oAuth2) verifyTokenScope(req *http.Request, scope ...string) bool {
	userID, err := o.CurrentUserID(req)
	fields := strings.Fields(req.Header.Get("Authorization"))

	return err != nil &&
		len(fields) == 2 &&
		strings.EqualFold(fields[0], "Bearer") &&
		o.notary.VerifyAccessToken(fields[1], userID, scope...) != nil
}

func (o *oAuth2) requireScope(scopes ...string) *Middleware {
	return newMiddleware(func(w http.ResponseWriter, req *http.Request) bool {
		if !o.verifyTokenScope(req, scopes...) {
			w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, strings.Join(scopes, " ")))
			errors.WriteErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
			return true // shortcut
		}
		return false // continue
	})
}

func (o *oAuth2) authorize(w http.ResponseWriter, req *http.Request) {
	// req.ParseForm()
	// redirectUri := req.Form.Get("redirect_uri")
	// responseType := req.Form.Get("response_type")
	// scope := req.Form.Get("scope")
	// clientID := req.Form.Get("cleint_id")

	// client, err := o.context.dao.Client.Get(clientID)

	// userId, err := o.context.CurrentUserID(req)
	// if err != nil {

	// }
}

// Token exchange - get an access_token sending an id_token
//https://tools.ietf.org/html/draft-ietf-oauth-token-exchange-12

func (o *oAuth2) token(w http.ResponseWriter, req *http.Request) {
	// return access token and id token
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
		w.WriteHeader(http.StatusNoContent)
		return true
	})
}
