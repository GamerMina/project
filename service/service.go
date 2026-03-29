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

	ok, err := CompareHash(account.Password, input.Password)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("password incorrect")
	}

	if account.AccountStatus != "active" {
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

	ok, err := CompareHash(account.Password, input.Password)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("password incorrect")
	}

	if account.AccountStatus != "active" {
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
	if err := ValidateAccountToAccountTransfer(
		transfer.PhoneNumberSender,
		transfer.PhoneNumberReceiver,
		transfer.Amount,
	); err != nil {
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

	ok, err := CompareHash(sender.Password, transfer.PasswordSender)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("password incorrect")
	}

	if sender.AccountStatus != "active" {
		return errors.New("sender account is not active")
	}
	if receiver.AccountStatus != "active" {
		return errors.New("receiver account is not active")
	}
	if sender.Currency != receiver.Currency {
		return errors.New("currencies do not match")
	}

	newSenderBalance, newReceiverBalance, err := TransferMoney(sender.Balance, receiver.Balance, transfer.Amount)
	if err != nil {
		return err
	}
	//сохраняем все это дело
	err = s.Repository.TransferMoneyTransaction(
		sender.ID,
		receiver.ID,
		newSenderBalance,
		newReceiverBalance,
	)
	if err != nil {
		return err
	}

	return nil
}
func (s *Services) MoneyTransferAccountToCard(transfer types.MoneyTransferAccountToCard) error {
	var err error

	transfer.PhoneNumberSender = ParcePhoneNumber(transfer.PhoneNumberSender)

	if transfer.PhoneNumberSender == "" {
		return errors.New("invalid sender phone")
	}
	if transfer.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	sender, err := s.Repository.GetAccountByPhone(transfer.PhoneNumberSender)
	if err != nil {
		return err
	}

	ok, err := CompareHash(sender.Password, transfer.PasswordSender)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("password incorrect")
	}
	// прверящаем наш инпут в хеш чтобы найты данные в бд
	cardHash := HashData(transfer.CardReceiver, dbconn.Secret())

	receiver, err := s.Repository.GetCardByNumberHash(cardHash)
	if err != nil {
		return err
	}

	if sender.AccountStatus != "active" {
		return errors.New("sender account is not active")
	}
	if receiver.Status != "active" {
		return errors.New("receiver card is not active")
	}
	if sender.Currency != receiver.Currency {
		return errors.New("currencies do not match")
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

	if transfer.PhoneNumberReceiver == "" {
		return errors.New("invalid receiver phone")
	}
	if transfer.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	cardHash := HashData(transfer.CardSender, dbconn.Secret())

	sender, err := s.Repository.GetCardByNumberHash(cardHash)
	if err != nil {
		return err
	}
	// тут еще проверяем его Cvv
	ok, err := CompareHash(sender.CVVHash, transfer.CVVSender)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("cvv incorrect")
	}

	receiver, err := s.Repository.GetAccountByPhone(transfer.PhoneNumberReceiver)
	if err != nil {
		return err
	}

	if sender.Status != "active" {
		return errors.New("sender card is not active")
	}
	if receiver.AccountStatus != "active" {
		return errors.New("receiver account is not active")
	}
	if sender.Currency != receiver.Currency {
		return errors.New("currencies do not match")
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

	if transfer.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if transfer.CardSender == transfer.CardReceiver {
		return errors.New("cannot transfer to the same card")
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

	ok, err := CompareHash(sender.CVVHash, transfer.CVVSender)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("cvv incorrect")
	}

	if sender.Status != "active" {
		return errors.New("sender card is not active")
	}
	if receiver.Status != "active" {
		return errors.New("receiver card is not active")
	}
	if sender.Currency != receiver.Currency {
		return errors.New("currencies do not match")
	}

	newSenderBalance, newReceiverBalance, err := TransferMoney(sender.Balance, receiver.Balance, transfer.Amount)
	if err != nil {
		return err
	}

	return s.Repository.TransferCardToCard(
		sender.ID,
		receiver.ID,
		newSenderBalance,
		newReceiverBalance,
	)
}
