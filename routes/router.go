package routes

import (
	"net/http"

	"github.com/gabriel-araujjo/condominio-auth/dao"
	"github.com/gabriel-araujjo/condominio-auth/sessions"
)

//TODO: Make Dao an interface

func NewServeAuth(dao *dao.Dao, s sessions.Store) http.Handler {

	routes := http.NewServeMux()
	// context = &context{dao, sessions}
	// userRoutes := newUserRoutes(context)
	// routes.Handle()
	// r.HandleFunc
	// user := &userRouter{dao}
	// oidc := &oidcRouter{}

	// router.GET("/oidc/auth", oidc.auth)
	// router.POST("/oidc/auth", oidc.auth)
	// router.POST("/user/login", user.login)
	// router.GET("/user/:id", user.get)
	// router.POST("/user", user.create)
	// router.DELETE("/user/:id", user.delete)

	//client := &ClientRouter{dao}
	//router.GET("/client/token", )
	return routes
}
