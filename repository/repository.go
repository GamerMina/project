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

func (r *Repository) GetAccount(id int) (types.Account, error) {
	var account types.Account
	err := r.Connection.Table("bank_accounts").First(&account, id).Error
	if err != nil {
		return types.Account{}, err
	}
	return account, nil
}
func (r *Repository) HasCardByID(idAccount int) (types.Card, error) {
	var card types.Card
	if err := r.Connection.Table("visa_cards").Where("account_id = ?", idAccount).First(&card).Error; err != nil {
		return types.Card{}, err
	}
	return types.Card{}, nil
}
