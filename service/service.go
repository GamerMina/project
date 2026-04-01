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

	input.Password = HashPassword(input.Password, dbconn.Secret())
	return input, err
}
func (s *Services) SaveAccountDB(account types.Account) error {
	return s.Repository.AddAccount(account)
}
func (s *Services) CreateCard(account types.Account) (types.CardResponse, error) {
	card, panPlain, cvvPlain, err := s.GenerateCard(account)
	if err != nil {
		return types.CardResponse{}, err
	}

	if err := s.Repository.AddCard(card); err != nil {
		return types.CardResponse{}, err
	}

	resp := types.CardResponse{
		IDAccount:  card.IDAccount,
		CardNumber: panPlain,
		Holder:     card.Holder,
		ExpMonth:   card.ExpMonth,
		ExpYear:    card.ExpYear,
		CVV:        cvvPlain,
		Balance:    card.Balance,
		Currency:   card.Currency,
		Status:     card.Status,
	}

	return resp, nil
}
func (s *Services) GenerateCard(input types.Account) (types.Card, string, string, error) {
	var card types.Card

	cardNum, err := s.generateCardNumber()
	if err != nil {
		return types.Card{}, "", "", err
	}

	cvvPlain := GenerateCVV()

	hashCVV, err := HashCVV(cvvPlain)
	if err != nil {
		return types.Card{}, "", "", err
	}

	expYear, expMonth := AddYearsMonths(5, 0)

	card = types.Card{
		CardNumberHash: HashData(cardNum, dbconn.Secret()),
		Holder:         "",
		ExpMonth:       expMonth,
		ExpYear:        expYear,
		CVVHash:        hashCVV,
	}

	card, err = s.FillingCard(input, card)
	if err != nil {
		return types.Card{}, "", "", err
	}

	return card, cardNum, cvvPlain, nil
}

func (s *Services) FillingCard(input types.Account, card types.Card) (types.Card, error) {
	id := input.ID // Acount id
	holderAccount, err := s.Repository.GetAccount(id)
	// превращаю его в тип ИМЯ ФАМИЛИЯ
	holder := fmt.Sprintf("%s %s", holderAccount.FirstName, holderAccount.LastName)
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

func (s *Services) BlockCardByPhone(input types.BlockCardByPhone) error {
	var err error

	input.PhoneNumber = ParcePhoneNumber(input.PhoneNumber)
	if input.PhoneNumber == "" {
		return errors.New("invalid phone number")
	}

	account, err := s.Repository.GetAccountByPhone(input.PhoneNumber)
	if err != nil {
		return err
	}

	if ComparePassword(account.Password, input.Password, dbconn.Secret()) == false {
		return errors.New("password incorrect")
	}

	if account.Status != "active" {
		return errors.New("account is not active")
	}

	err = s.Repository.BlockCardByAccountID(account.ID)
	if err != nil {
		return err
	}

	return nil
}
func (s *Services) ActivateCardByPhone(input types.BlockCardByPhone) error {
	var err error

	input.PhoneNumber = ParcePhoneNumber(input.PhoneNumber)
	if input.PhoneNumber == "" {
		return errors.New("invalid phone number")
	}

	account, err := s.Repository.GetAccountByPhone(input.PhoneNumber)
	if err != nil {
		return err
	}

	if ComparePassword(account.Password, input.Password, dbconn.Secret()) == false {
		return errors.New("password incorrect")
	}
	if account.Status != "active" {
		return errors.New("account is not active")
	}

	err = s.Repository.ActivateCardByAccountID(account.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Services) MoneyTransferAccountToAccount(transfer types.MoneyTransferAccountToAccount) error {
	var err error
	//парсим номер
	transfer.PhoneNumberSender = ParcePhoneNumber(transfer.PhoneNumberSender)
	transfer.PhoneNumberReceiver = ParcePhoneNumber(transfer.PhoneNumberReceiver)
	// проверяем на наличие ошибок
	if err := ValidateTransfer(transfer.PhoneNumberSender, transfer.PhoneNumberReceiver, transfer.Amount); err != nil {
		return err
	}
	//получаем данные отправителя
	sender, err := s.Repository.GetAccountByPhone(transfer.PhoneNumberSender)
	if err != nil {
		return err
	}
	//получаем данные получателя
	receiver, err := s.Repository.GetAccountByPhone(transfer.PhoneNumberReceiver)
	if err != nil {
		return err
	}

	ok := ComparePassword(sender.Password, transfer.PasswordSender, dbconn.Secret())
	if ok == false {
		return errors.New("password incorrect")
	}

	if err := ValidateTransactionAfterGettingData(sender.Status, receiver.Status, sender.Currency, receiver.Currency); err != nil {
		return err
	}
	newSenderBalance, newReceiverBalance, err := TransferMoney(sender.Balance, receiver.Balance, transfer.Amount)
	if err != nil {
		return err
	}
	//сохраняем все это дело
	err = s.Repository.TransferMoneyTransaction(sender.ID, receiver.ID, newSenderBalance, newReceiverBalance)
	if err != nil {
		return err
	}

	return nil
}
func (s *Services) MoneyTransferAccountToCard(transfer types.MoneyTransferAccountToCard) error {
	var err error

	transfer.PhoneNumberSender = ParcePhoneNumber(transfer.PhoneNumberSender)

	if err := ValidateTransfer(transfer.PhoneNumberSender, transfer.CardReceiver, transfer.Amount); err != nil {
		return err
	}
	sender, err := s.Repository.GetAccountByPhone(transfer.PhoneNumberSender)
	if err != nil {
		return err
	}

	if ComparePassword(sender.Password, transfer.PasswordSender, dbconn.Secret()) == false {
		return errors.New("password incorrect")
	}

	// прверящаем наш инпут в хеш чтобы найты данные в бд
	cardHash := HashData(transfer.CardReceiver, dbconn.Secret())

	receiver, err := s.Repository.GetCardByNumberHash(cardHash)
	if err != nil {
		return err
	}

	if err := ValidateTransactionAfterGettingData(sender.Status, receiver.Status, sender.Currency, receiver.Currency); err != nil {
		return err
	}
	newSenderBalance, newReceiverBalance, err := TransferMoney(sender.Balance, receiver.Balance, transfer.Amount)
	if err != nil {
		return err
	}
	return s.Repository.TransferAccountToCard(
		sender.ID,
		receiver.ID,
		newSenderBalance,
		newReceiverBalance,
	)
}
func (s *Services) MoneyTransferCardToAccount(transfer types.MoneyTransferCardToAccount) error {
	var err error

	transfer.PhoneNumberReceiver = ParcePhoneNumber(transfer.PhoneNumberReceiver)

	if err := ValidateTransfer(transfer.CardSender, transfer.PhoneNumberReceiver, transfer.Amount); err != nil {
		return err
	}
	cardHash := HashData(transfer.CardSender, dbconn.Secret())

	sender, err := s.Repository.GetCardByNumberHash(cardHash)
	if err != nil {
		return err
	}
	// тут еще проверяем его Cvv
	if CompareCVV(sender.CVVHash, transfer.CVVSender) == false {
		return errors.New("cvv incorrect")
	}

	receiver, err := s.Repository.GetAccountByPhone(transfer.PhoneNumberReceiver)
	if err != nil {
		return err
	}

	if err := ValidateTransactionAfterGettingData(sender.Status, receiver.Status, sender.Currency, receiver.Currency); err != nil {
		return err
	}
	newSenderBalance, newReceiverBalance, err := TransferMoney(sender.Balance, receiver.Balance, transfer.Amount)
	if err != nil {
		return err
	}

	return s.Repository.TransferCardToAccount(
		sender.ID,
		receiver.ID,
		newSenderBalance,
		newReceiverBalance,
	)
}
func (s *Services) MoneyTransferCardToCard(transfer types.MoneyTransferCardToCard) error {
	var err error

	if err := ValidateTransfer(transfer.CardSender, transfer.CardReceiver, transfer.Amount); err != nil {
		return err
	}
	senderHash := HashData(transfer.CardSender, dbconn.Secret())
	receiverHash := HashData(transfer.CardReceiver, dbconn.Secret())

	sender, err := s.Repository.GetCardByNumberHash(senderHash)
	if err != nil {
		return err
	}

	receiver, err := s.Repository.GetCardByNumberHash(receiverHash)
	if err != nil {
		return err
	}

	if CompareCVV(sender.CVVHash, transfer.CVVSender) == false {
		return errors.New("cvv incorrect")
	}
	if err := ValidateTransactionAfterGettingData(sender.Status, receiver.Status, sender.Currency, receiver.Currency); err != nil {
		return err
	}

	newSenderBalance, newReceiverBalance, err := TransferMoney(sender.Balance, receiver.Balance, transfer.Amount)
	if err != nil {
		return err
	}

	err = s.Repository.TransferCardToCard(sender.ID, receiver.ID, newSenderBalance, newReceiverBalance)

	return err
}
