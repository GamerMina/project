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

func (s *Services) GenerateCard() (types.Card, error) {
	var card types.Card
	cardNum, err := s.generateCardNumber()
	if err != nil {
		return card, err
	}

	hiddenCardNum := HidePAN(cardNum)

	hashCVV, err := HashCVV(GenerateCVV())
	if err != nil {
		return card, err
	}

	expYear, expMonth := AddYearsMonths(5, 0) //5 лет это пример

	card = types.Card{
		CardNumber:     hiddenCardNum,
		CardNumberHash: HashCardData(cardNum, dbconn.Secret()),
		Holder:         "",
		ExpMonth:       expMonth,
		ExpYear:        expYear,
		CVVHash:        hashCVV,
	}

	return card, nil
}
func (s *Services) FillingCard(input types.Account, card types.Card) (types.Card, error) {
	id := input.ID // Acount id
	Holder, err := s.Repository.GetAccount(id)
	// превращаю его в тип ИМЯ ФАМИЛИЯ
	holder := fmt.Sprintf("%s %s", Holder.FirstName, Holder.LastName)
	holder = strings.ToUpper(holder)

	filler := types.Card{
		IDAccount:      input.ID,
		CardNumber:     card.CardNumber,
		CardNumberHash: card.CardNumberHash,
		Holder:         holder,
		ExpMonth:       card.ExpMonth,
		ExpYear:        card.ExpYear,
		CVV:            card.CVV,
		CVVHash:        card.CVVHash,
		Balance:        0,                      // default
		Currency:       types.CurrencyTJS,      // input.Currency, используем если карты будут не только TJS
		Status:         types.CardStatusActive, // default
	}
	return filler, err
}
func (s *Services) SaveDB(card types.Card) error {

	err := s.Repository.AddCard(card)
	return err
}
func (s *Services) HasCardByID(id int) (bool, error) {
	var check bool
	card, err := s.Repository.HasCardByID(id)
	if err != nil {
		return false, err
	}
	if card == (types.Card{}) {
		check = true
	}

	return check, err
}
