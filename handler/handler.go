package handler

import (
	"github.com/gin-gonic/gin"
	"log"
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
	if check = h.Service.HasCardByID(input.ID); check == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "card already exist",
		})
		return
	}
	log.Println("hascardbyID", err)

	if card, err = h.Service.GenerateCard(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
		log.Println("generate card", err)
	}

	if err := h.Service.SaveDB(card); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("saveBD", err)

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
