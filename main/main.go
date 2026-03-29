package main

import (
	"github.com/gin-gonic/gin"
	"project/dbconn"
	"project/handler"
	"project/repository"
	"project/router"
	"project/service"
)

func main() {
	r := gin.Default()
	appConfig := dbconn.LoadAppCfg()
	Db := dbconn.DbCon(appConfig.Db)

	newRepository := repository.NewRepository(Db)
	newService := service.NewServices(newRepository)
	newHandler := handler.NewHandler(r, newService)

	router.InitRoutes(r, newHandler)

	r.Run("localhost:8080")

}
for save