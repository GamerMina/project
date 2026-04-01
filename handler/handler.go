package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/service"
	"project/types"
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

func (h *Handler) AccountRegistration(c *gin.Context) {
	input := types.Account{}
	var err error
	/* i need from input:
	FirstName        имя на кириллице
	LastName         фамилия на кириллице
	DateOfBirth      Др
	PhoneNumber      номер телефона в виде 12 чисел с +ом (если не так то само появиться	)
	Email            gmail
	Password         пароль   4 значный*/
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := h.Service.AccountRegistration(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = h.Service.SaveAccountDB(account); err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, gin.H{"message": "registration success"})

}
func (h *Handler) CardRegistration(c *gin.Context) {
	input := types.Account{}
	var err error
	var card types.CardResponse
	var check bool

	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if check = h.Service.HasCardByID(input.ID); check {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "card already exist",
		})
		return
	}

	card, err = h.Service.CreateCard(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Handler card=%+v\n", card)
	c.JSON(http.StatusOK, gin.H{
		"message": "card created",
		"card":    card,
	})
}
func (h *Handler) HasCardByID(c *gin.Context) {
	input := types.Account{}
	var err error
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
func (h *Handler) BlockCardByPhone(c *gin.Context) {
	var input types.BlockCardByPhone

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.BlockCardByPhone(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "card blocked successfully",
	})
}
func (h *Handler) ActivateCardByPhone(c *gin.Context) {
	var input types.BlockCardByPhone

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.ActivateCardByPhone(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "card activated successfully",
	})
}
func (h *Handler) MoneyTransferAccountToAccount(c *gin.Context) {
	var input types.MoneyTransferAccountToAccount

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.Service.MoneyTransferAccountToAccount(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "money transfer successful",
	})
}
func (h *Handler) MoneyTransferAccountToCard(c *gin.Context) {
	var input types.MoneyTransferAccountToCard

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.MoneyTransferAccountToCard(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "money transfer successful"})
}
func (h *Handler) MoneyTransferCardToAccount(c *gin.Context) {
	var input types.MoneyTransferCardToAccount

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.MoneyTransferCardToAccount(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "money transfer successful"})
}
func (h *Handler) MoneyTransferCardToCard(c *gin.Context) {
	var input types.MoneyTransferCardToCard

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.MoneyTransferCardToCard(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "money transfer successful"})
}
