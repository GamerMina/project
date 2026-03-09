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


func (s *Services) generateCard(card types.Card) (types.Card ,error) { //Todo: хеширование или шифрование данных дописать
	card.CardNumber := s.generateCardNumber()
	card.CardNumberHash := HashCardData(card.CardNumber,dbconn.Secret())
	card.CardNumber := HidePAN (card.CardNumber)

	return card, nil
}
