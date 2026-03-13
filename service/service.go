package service

import (
	"fmt"
	"projcet/dbconn"
	"projcet/repository"
	"projcet/types"
	"strings"
)

type Services struct {
	Repository *repository.Repository
}

func NewServices(rep *repository.Repository) *Services {
	return &Services{Repository: rep}
}

func (s *Services) GenerateCard() types.Card { //Todo: хеширование или шифрование данных дописать//done
	var card types.Card
	card.CardNumber, _ = s.generateCardNumber()
	card.CardNumberHash = HashCardData(card.CardNumber, dbconn.Secret())
	card.CardNumber = HidePAN(card.CardNumber)
	card.CVVHash, _ = HashCVV(GenerateCVV())
	_, card.ExpMonth = AddYearsMonths(5, 0) //5 лет это пример
	card.ExpYear, _ = AddYearsMonths(5, 0)  // 5 лет это пример
	return card
}
func (s *Services) FillingCard(input types.Card, card types.Card) (types.Card, error) {
	id := input.IDAccount
	Holder, err := s.Repository.GetAccount(id)
	// превращаю его в тип ИМЯ ФАМИЛИЯ
	holder := fmt.Sprintf("%s %s", Holder.FirstName, Holder.LastName)
	holder = strings.ToUpper(holder)

	filler := types.Card{
		IDAccount:      input.IDAccount,
		CardNumber:     card.CardNumber,
		CardNumberHash: card.CardNumberHash,
		Holder:         holder,
		ExpMonth:       card.ExpMonth,
		ExpYear:        card.ExpYear,
		CVV:            card.CVV,
		CVVHash:        card.CVVHash,
		Balance:        0,        // default
		Currency:       "TJS",    // input.Currency, используем если карты будут не только TJS
		Status:         "active", // default
	}
	return filler, err
}
func (s *Services) SaveDB(card types.Card) error {
	err := s.Repository.AddCard(card)
	return err
}
