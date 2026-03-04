package handler

import (
	"github.com/gin-gonic/gin"
	"projcet/service"
)


type Handler struct {
	Engine  *gin.Engine
	Service *service.Services
	//	Authorization *service.Authentication
}

func NewHandler(engine *gin.Engine, services *service.Services) *Handler {
	return &Handler{
		Engine:  engine,
		Service: services,
	}
}
