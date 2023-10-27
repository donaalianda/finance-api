package model

import (
	"finance-api/model/entity"

	"github.com/jinzhu/gorm"
)

type (
	TransactionModelInterface interface {
		Create(obj entity.TransactionEntity) error
	}

	transactionModel struct {
		DB *gorm.DB
	}
)

func NewTransactionModel(db *gorm.DB) TransactionModelInterface {
	return &transactionModel{db}
}

func (t *transactionModel) Create(obj entity.TransactionEntity) error {

	tx := t.DB.Begin()

	if err := tx.Save(&obj).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil

}
