package main

import (
	"github.com/gin-gonic/gin"
	"projcet/dbconn"
	"projcet/handler"
	"projcet/repository"
	"projcet/router"
	"projcet/service"
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
