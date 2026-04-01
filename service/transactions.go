package service

import (
	"errors"
	"project/dbconn"
	"project/types"
)

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
