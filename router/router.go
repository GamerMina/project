package router

import (
	"github.com/gin-gonic/gin"
	"project/handler"
)

func InitRoutes(r *gin.Engine, H *handler.Handler) {
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
	r.POST("/register_card", H.CardRegistration)
	r.GET("/check", H.HasCardByID)
	r.POST("/transfer")
}
