package routes

import (
	"github.com/gabriel-araujjo/condominio-auth/dao"
	"github.com/gin-gonic/gin"
)

func ConfigureEngine(router gin.IRouter, dao *dao.Dao) {
	user := &userRouter{dao}
	oidc := &oidcRouter{}

	router.GET("/oidc/auth", oidc.auth)
	router.POST("/oidc/auth", oidc.auth)
	router.POST("/user/login", user.login)
	router.GET("/user/:id", user.get)
	router.POST("/user", user.create)
	router.DELETE("/user/:id", user.delete)

	//client := &ClientRouter{dao}
	//router.GET("/client/token", )

}
