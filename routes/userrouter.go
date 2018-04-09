package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gabriel-araujjo/condominio-auth/data"
	"github.com/gabriel-araujjo/base62"
	"github.com/gabriel-araujjo/condominio-auth/errors"
	"fmt"
)

type UserRouter struct {
	dao *data.Dao
}

func (e *UserRouter) Auth(c *gin.Context) {
	var body struct{
		Credential string `json: "cred" binding:"required"`
		Password string `json:"passwd" binding:"required"`
	}

	if err := c.BindJSON(body); err != nil {
		c.Error(httperrors.PreconditionFailed("cred and passws fields are required"))
		return
	}
	var id int64
	if id, err :=e.dao.User.Auth(body.Credential, body.Password); err != nil {

	}
}

func (e *UserRouter) Get(c *gin.Context) {
	id, err := base62.ParseInt(c.Param("id"))
	if err != nil {
		c.Error(httperrors.NotFound(fmt.Sprintf("not found")))
		return
	}
	u, err := e.dao.User.Get(id)
	if err != nil {
		c.Error(httperrors.NotFound(fmt.Sprintf("not found")))
		return
	}

	c.JSON(200, u)
}

func (e *UserRouter) Create(c *gin.Context) {

}

func (e *UserRouter) Delete(c *gin.Context) {

}
