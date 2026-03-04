package router

import (
	"github.com/gin-gonic/gin"
	"projcet/handler"
)

func InitRoutes(r *gin.Engine, H *handler.Handler) {
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})

}
