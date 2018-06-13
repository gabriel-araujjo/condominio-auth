package routes

import (
	"fmt"
	"net/http"
	"net/url"
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
	req.ParseForm()
	redirectUri, _ := url.Parse(req.Form.Get("redirect_uri"))
	responseType := req.Form.Get("response_type")
	scope := strings.Fields(req.Form.Get("scope"))
	clientID := req.Form.Get("cleint_id")
	state := req.Form.Get("state")

	client, err := o.context.dao.Client.Get(clientID)
	var scopeIDs []int64
	var code string

	if err != nil || strings.EqualFold(responseType, "code") || redirectUri == nil {
		errors.WriteErrorWithCode(w, http.StatusNotFound, "not found")
		return
	}

	userID, err := o.context.CurrentUserID(req)
	if err != nil {
		redirectUri.Query().Set("error", "login_required")
		goto respond
	}

	err = o.context.dao.User.AuthorizeClient(userID, clientID, scope)
	if err != nil {
		redirectUri.Query().Set("error", "invalid_request_uri")
		goto respond
	}

	scopeIDs, err = o.context.dao.Permission.ScopeIntoPermissionIDs(scope)

	if len(scopeIDs) == 0 {
		redirectUri.Query().Set("error", "invalid_request_uri")
		goto respond
	}

	code, err = o.notary.NewClientCode(client.ID, scopeIDs, userID)

	if err != nil {
		errors.WriteErrorWithCode(w, http.StatusInternalServerError, "can't generate code")
		return
	}

	redirectUri.Query().Set("code", code)

respond:
	if len(state) != 0 {
		redirectUri.Query().Set("state", state)
	}
	w.Header().Set("Location", redirectUri.String())
	w.WriteHeader(http.StatusFound)
}

// Token exchange - get an access_token sending an id_token
//https://tools.ietf.org/html/draft-ietf-oauth-token-exchange-12

func (o *oAuth2) token(w http.ResponseWriter, req *http.Request) {
	// 	req.ParseForm()
	// 	grantType := req.Form.Get("grant_type")
	// 	code := req.Form.Get("code")
	// 	clientPublicID := req.Form.Get("cleint_id")
	// 	clientSecret := req.Form.Get("client_secret")

	// 	if !strings.EqualFold(grantType, "authorization_code") ||
	// 		code == "" ||
	// 		clientSecret == "" ||
	// 		clientPublicID == "" {
	// 		errors.WriteErrorWithCode(http.StatusBadRequest, "invalid_request")
	// 		return
	// 	}

	// 	client, err := o.context.dao.Client.Get(clientPublicID)
	// 	if client == nil {
	// 		errors.WriteErrorWithCode(w, http.StatusNotFound, "not found")
	// 		return
	// 	}

	// 	clientID, scope, uID, err := o.notary.DecipherCode(code)

	// 	if clientID != client.ID || clientSecret != clientSecret {
	// 		errors.WriteErrorWithCode(w, http.StatusNotFound, "not found")
	// 		return
	// 	}

	// 	expiresAt time.Now() + 60 * time.Minute
	// 	accessToken, err := o.notary.NewAccessToken(60 *time.Minute, uID, scope...)
	// 	var idToken string

	// 	for s := range scope {
	// 		if strings.EqualFold(s, "openid") {
	// 			idToken := o.notary.NewIDTokenWithClaims()
	// 		}
	// 	}

	// respond:
	// 	if len(state) != 0 {
	// 		redirectUri.Query().Set("state", state)
	// 	}
	// 	w.Header().Set("Location", redirectUri.String())
	// 	w.WriteHeader(http.StatusFound)
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
