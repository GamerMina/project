package repository

import "gorm.io/gorm"

type Repository struct {
	Connection *gorm.DB
}

func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{Connection: conn}
}
