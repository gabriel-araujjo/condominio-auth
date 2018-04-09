package app

import ( // Standard library packages
	// Third party packages
	"github.com/gabriel-araujjo/condominio-auth/data/factory"
	"github.com/gabriel-araujjo/condominio-auth/routes"
	"github.com/gin-gonic/gin"
	"github.com/gabriel-araujjo/condominio-auth/config"
)

func main() {

	conf := config.DefaultConfig()
	dao, err := factory.NewDao(conf)

	if err != nil {
		panic(err)
	}
	defer dao.Close()

	engine := gin.Default()

	routes.ConfigRoutes(engine, dao)

	engine.Run()
}
