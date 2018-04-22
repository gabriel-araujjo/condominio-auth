package routes

import (
	"github.com/gabriel-araujjo/condominio-auth/dao"
	"github.com/gin-gonic/gin"
	"github.com/gabriel-araujjo/condominio-auth/errors"
	"net/http"
	"github.com/gabriel-araujjo/condominio-auth/auth"
	"strings"
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

type ClientRouter struct {
	dao *dao.Dao
	jwt *auth.Auth
}

func (e *ClientRouter) Auth(c *gin.Context)  {
	pubID := c.Query("pub_id")
	secret := c.Query("secret")

	pubId, err := e.dao.Client.Auth(pubID, secret)
	if err!= nil {
		c.Error(httperrors.Forbidden("unauthorized client"))
		return
	}

	c.String(http.StatusOK, e.jwt.Sign(&domain.Claims{
		Audience: pubId,
	}))
}

func (e *ClientRouter) AuthJwt(c *gin.Context) {
	authHead := c.GetHeader("Authorization")
	tokenString := strings.Trim(authHead, "Bearer ")
	_/*token*/, err := e.jwt.Verify(tokenString)
	if err != nil {
		c.Error(httperrors.Forbidden("unauthorized client"))
	}
	// TODO: Register access of client
	c.Next()
}
