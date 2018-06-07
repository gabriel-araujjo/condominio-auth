package app

import (
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/dao"
	"github.com/gabriel-araujjo/condominio-auth/routes"
	"github.com/gabriel-araujjo/condominio-auth/sessions"
	"github.com/gin-gonic/gin"
)

func database(config *config.Config) *dao.Dao {
	dao, err := dao.NewFromConfig(config)
	if err != nil {
		panic(err)
	}
	return dao
}

func sessionsStore(config *config.Config) sessions.Store {
	store, err := sessions.NewStoreFromConfig(config)
	if err != nil {
		panic(err)
	}
	return store
}

func main() {
	conf := config.DefaultConfig()

	db := database(conf)
	session := sessionsStore(conf)
	defer db.Close()
	defer session.Close()

	engine := gin.Default()

	v1 := engine.Group("/api/v1")
	routes.ConfigureEngine(v1, db)

	engine.Run()
}
