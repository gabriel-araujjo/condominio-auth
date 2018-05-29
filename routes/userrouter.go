package routes

import (
	"fmt"

	"github.com/gabriel-araujjo/base62"
	"github.com/gabriel-araujjo/condominio-auth/dao"
	"github.com/gabriel-araujjo/condominio-auth/errors"
	"github.com/gin-gonic/gin"
)

type userRouter struct {
	dao *dao.Dao
}

func (e *userRouter) login(c *gin.Context) {
	var params struct {
		Credential string `json:"cred" form:"cred" binding:"required"`
		Password   string `json:"passwd" binding:"required"`
	}

	if err := c.Bind(params); err != nil {
		c.Error(httperrors.PreconditionFailed("cred and passws fields are required"))
		return
	}

	if _, err := e.dao.User.Auth(params.Credential, params.Password); err != nil {
		c.Error(httperrors.Forbidden("forbidden"))
	}
	// TODO: send JWT back
}

func (e *userRouter) get(c *gin.Context) {
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

func (e *userRouter) create(c *gin.Context) {

}

func (e *userRouter) delete(c *gin.Context) {

}
