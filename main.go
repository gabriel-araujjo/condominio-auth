package app

import (
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/dao"
	"github.com/gabriel-araujjo/condominio-auth/routes"
	"github.com/gabriel-araujjo/condominio-auth/session"
	"github.com/gin-gonic/gin"
)

func sessions(config *config.Config) *session.Session {
	s, err := session.NewFromConfig(config)
	if err != nil {
		panic(err)
	}
	return s
}

func database(config *config.Config) *dao.Dao {
	dao, err := dao.NewFromConfig(config)
	if err != nil {
		panic(err)
	}
	return dao
}

func main() {
	conf := config.DefaultConfig()
	db := database(conf)
	defer db.Close()
	sess := sessions(conf)
	defer sess.Close()

	engine := gin.Default()
	engine.Use(sess.Middleware())

	v1 := engine.Group("/api/v1")
	routes.ConfigureEngine(v1, db)

	engine.Run()
}
