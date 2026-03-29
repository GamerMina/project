package router

import (
	"github.com/gin-gonic/gin"
	"project/handler"
)

func InitRoutes(r *gin.Engine, H *handler.Handler) {
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
	r.GET("/check", H.HasCardByID)
	register := r.Group("/register")
	{
		register.POST("/card", H.CardRegistration)
		register.POST("/account", H.AccountRegistration)
	}

	card := r.Group("/card")
	{
		card.POST("/block", H.BlockCardByPhone)
		card.POST("/activate", H.ActivateCardByPhone)

	}
	transfer := r.Group("/transfer")
	{
		transfer.POST("/account_to_account", H.MoneyTransferAccountToAccount)
		transfer.POST("/account_to_card", H.MoneyTransferAccountToCard)
		transfer.POST("/card_to_account", H.MoneyTransferCardToAccount)
		transfer.POST("/card_to_card", H.MoneyTransferCardToCard)
	}

}
