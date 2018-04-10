package routes

import (
	"github.com/gabriel-araujjo/condominio-auth/data"
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(router gin.IRouter, dao *data.Dao) {
	user := &UserRouter{dao}

	router.POST("/auth", user.Auth)
	router.GET("/user/:id", user.Get)
	router.POST("/user", user.Create)
	router.DELETE("/user/:id", user.Delete)

	//client := &ClientRouter{dao}
	//router.GET("/client/token", )

}