package service

import (
	"projcet/dbconn"
	"projcet/repository"
	"projcet/types"
)

type Services struct {
	Repository *repository.Repository
}

func NewServices(rep *repository.Repository) *Services {
	return &Services{Repository: rep}
}

func (s *Services) GenerateCard() (types.Card, ) { //Todo: хеширование или шифрование данных дописать
	var card types.Card
	card.CardNumber, _ = s.generateCardNumber()
	card.CardNumberHash = HashCardData(card.CardNumber, dbconn.Secret())
	card.CardNumber = HidePAN(card.CardNumber)
	card.CVVHash, _ = HashCVV(GenerateCVV())
	_, card.ExpMonth = AddYearsMonths(5, 0) //5 лет это пример
	card.ExpYear, _ = AddYearsMonths(5, 0)  // 5 лет это пример

	return card
}
