package repository

import (
	"gorm.io/gorm"
	"project/types"
)

type Repository struct {
	Connection *gorm.DB
}

func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{Connection: conn}
}

func (r *Repository) AddCard(card types.Card) error {
	err := r.Connection.Table("visa_cards").Create(&card).Error
	return err
}
func (r *Repository) AddAccount(account types.Account) error {
	err := r.Connection.Table("bank_accounts").Create(&account).Error
	return err
}
func (r *Repository) GetCard(id int) (types.Card, error) {
	var card types.Card
	err := r.Connection.Table("visa_cards").First(&card, id).Error
	if err != nil {
		return types.Card{}, err
	}
	return card, nil
}
func (r *Repository) GetAccount(id int) (types.Account, error) {
	var account types.Account
	err := r.Connection.Table("bank_accounts").First(&account, id).Error
	if err != nil {
		return types.Account{}, err
	}
	return account, nil
}
func (r *Repository) GetAccountBalance(phone string) (types.Balance, error) {
	var account types.Account
	err := r.Connection.Table("bank_accounts").
		Where("phone_number = ?", phone).
		First(&account).Error
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}
func (r *Repository) HasCardByID(idAccount int) (types.Card, error) {
	var card types.Card
	if err := r.Connection.Table("visa_cards").Where("account_id = ?", idAccount).First(&card).Error; err != nil {
		return types.Card{}, err
	}
	return types.Card{}, nil
}
func (r *Repository) BlockCardByAccountID(accountID int) error {
	return r.Connection.Table("visa_cards").
		Where("account_id = ?", accountID).
		Update("card_status", types.CardStatusBlocked).Error
}
func (r *Repository) ActivateCardByAccountID(accountID int) error {
	return r.Connection.Table("visa_cards").
		Where("account_id = ?", accountID).
		Update("card_status", types.CardStatusActive).Error
}
func (r *Repository) TransferMoneyTransaction(senderID int, receiverID int, senderBalance types.Balance, receiverBalance types.Balance) error {
	return r.Connection.Transaction(func(tx *gorm.DB) error {

		if err := tx.Table("bank_accounts").
			Where("id = ?", senderID).
			Update("balance", senderBalance).Error; err != nil {
			return err
		}

		if err := tx.Table("bank_accounts").
			Where("id = ?", receiverID).
			Update("balance", receiverBalance).Error; err != nil {
			return err
		}

		return nil
	})
}
func (r *Repository) GetAccountByPhone(phone string) (types.Account, error) {
	var account types.Account
	err := r.Connection.Table("bank_accounts").
		Where("phone_number = ?", phone).
		First(&account).Error
	if err != nil {
		return types.Account{}, err
	}
	return account, nil
}
func (r *Repository) GetCardByNumberHash(hash string) (types.Card, error) {
	var card types.Card
	err := r.Connection.Table("visa_cards").
		Where("pan_hash = ?", hash).
		First(&card).Error
	if err != nil {
		return types.Card{}, err
	}
	return card, nil
}
func (r *Repository) TransferAccountToCard(senderID int, receiverID int, senderBalance types.Balance, receiverBalance types.Balance) error {
	return r.Connection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("bank_accounts").
			Where("id = ?", senderID).
			Update("balance", senderBalance).Error; err != nil {
			return err
		}

		if err := tx.Table("visa_cards").
			Where("id = ?", receiverID).
			Update("balance", receiverBalance).Error; err != nil {
			return err
		}

		return nil
	})
}
func (r *Repository) TransferCardToAccount(senderID int, receiverID int, senderBalance types.Balance, receiverBalance types.Balance) error {
	return r.Connection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("visa_cards").
			Where("id = ?", senderID).
			Update("balance", senderBalance).Error; err != nil {
			return err
		}

		if err := tx.Table("bank_accounts").
			Where("id = ?", receiverID).
			Update("balance", receiverBalance).Error; err != nil {
			return err
		}

		return nil
	})
}
func (r *Repository) TransferCardToCard(senderID int, receiverID int, senderBalance types.Balance, receiverBalance types.Balance) error {
	return r.Connection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("visa_cards").
			Where("id = ?", senderID).
			Update("balance", senderBalance).Error; err != nil {
			return err
		}

		if err := tx.Table("visa_cards").
			Where("id = ?", receiverID).
			Update("balance", receiverBalance).Error; err != nil {
			return err
		}

		return nil
	})
}
