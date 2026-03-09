package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"projcet/service"
	"projcet/types"
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
func (h *Handler) Registration(c *gin.Context) {
	var input types.Card

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card := service.FillingCard(input, h.Service.GenerateCard())
	//TODO: сохранение в БД
	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}
