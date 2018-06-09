package routes

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gabriel-araujjo/condominio-auth/security"

	"github.com/gabriel-araujjo/condominio-auth/dao"
	"github.com/gabriel-araujjo/condominio-auth/domain"
	"github.com/gabriel-araujjo/condominio-auth/errors"
)

type ClientRouter struct {
	dao *dao.Dao
	jwt *security.Notary
}

func (e *ClientRouter) Auth(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	pubID := req.Form.Get("pub_id") // c.Query("pub_id")
	secret := req.Form.Get("secret")

	pubID, err := e.dao.Client.Auth(pubID, secret)
	if err != nil {
		errors.WriteErrorWithCode(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	w.Write([]byte(e.jwt.NewAccessTokenWithClaims(&domain.Claims{
		StandardClaims: jwt.StandardClaims{
			Audience: pubID,
		},
	})))
}

func (e *ClientRouter) AuthJwt(w http.ResponseWriter, req *http.Request) {
	authHead := req.Header.Get("Authorization")
	tokenString := strings.Trim(authHead, "Bearer ")
	_ /*token*/, err := e.jwt.VerifyAccessToken(tokenString)
	if err != nil {
		errors.WriteErrorWithCode(w, http.StatusUnauthorized, "Unauthorized")
	}
	// TODO: Register access of client
}
