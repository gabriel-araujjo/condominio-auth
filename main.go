package app

import (
	"github.com/gabriel-araujjo/condominio-auth/routes"
	"github.com/gin-gonic/gin"
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/dao"
)

func main() {
	conf := config.DefaultConfig()
	dao, err := dao.NewFromConfig(conf)

	if err != nil {
		panic(err)
	}

	defer dao.Close()

	engine := gin.Default()

	routes.ConfigRoutes(engine, dao)

	engine.Run()
}
