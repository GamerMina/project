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
	input := types.Account{}
	card := types.Card{}
	var check bool
	var err error
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if check, err = h.Service.HasCardByID(input.ID); err != nil {

	}
	if check == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "card already exists",
		})
		return
	}

	if card, err = h.Service.GenerateCard(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if card, err = h.Service.FillingCard(input, card); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.SaveDB(card); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func (h *Handler) HasCardByID(c *gin.Context) {
	input := types.Account{}
	var err error
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
