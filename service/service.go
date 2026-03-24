package service

import (
	"errors"
	"fmt"
	"project/dbconn"
	"project/repository"
	"project/types"
	"strings"
)

type Services struct {
	Repository *repository.Repository
}

func NewServices(rep *repository.Repository) *Services {
	return &Services{Repository: rep}
}

func (s *Services) GenerateCard(input types.Account) (types.Card, error) {
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
		CardNumberHash: HashData(cardNum, dbconn.Secret()),
		Holder:         "",
		ExpMonth:       expMonth,
		ExpYear:        expYear,
		CVVHash:        hashCVV,
	}
	card, err = s.FillingCard(input, card)
	if err != nil {
		return types.Card{}, err
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
func (s *Services) AccountRegistration(input types.Account) (types.Account, error) {
	age, err := GetAge(input.DateOfBirth)
	if err != nil {
		return types.Account{}, err
	}
	if age < types.MinAgeToCreateAccount {
		return types.Account{}, errors.New("you are less than 16 ")
	}
	input.PhoneNumber = ParcePhoneNumber(input.PhoneNumber)
	if err := ParseName(input.FirstName, input.LastName); err != nil {
		return types.Account{}, err
	}
	if err := ParseMail(input.Email); err != nil {
		return types.Account{}, err
	}
	input.Password = HashData(input.Password, dbconn.Secret())

	err = s.SaveAccountDB(input)

	return input, err
}
func (s *Services) SaveAccountDB(account types.Account) error {
	return s.Repository.AddAccount(account)
}

//todo сохранить эту новую карту

func (s *Services) SaveCardDB(card types.Card) error {
	return s.Repository.AddCard(card)
}
func (s *Services) HasCardByID(id int) bool {
	var check bool
	card, err := s.Repository.HasCardByID(id)
	if err != nil {
		return false
	}
	if card == (types.Card{}) {
		check = true
	}

	return check
}
