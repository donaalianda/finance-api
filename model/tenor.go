package model

import (
	"finance-api/model/entity"

	"github.com/jinzhu/gorm"
)

type (
	TenorModelInterface interface {
		Create(obj entity.TenorEntity) error
		GetByNIK(nik string) (*entity.TenorCustomerEntity, error)
		GetLimitByNIK(nik string) (*entity.TenorLimitCustomerEntity, error)
	}

	tenorModel struct {
		DB *gorm.DB
	}
)

func NewTenorModel(db *gorm.DB) TenorModelInterface {
	return &tenorModel{db}
}

func (t *tenorModel) Create(obj entity.TenorEntity) error {

	tx := t.DB.Begin()

	if err := tx.Save(&obj).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil

}

func (t *tenorModel) GetByNIK(nik string) (*entity.TenorCustomerEntity, error) {

	obj := entity.TenorCustomerEntity{}

	if err := t.DB.Table(`tenors`).Select(`customers.nik, customers.legal_name, tenors.*`).
		Joins(`JOIN customers ON customers.id = tenors.customer_id`).
		Where(`customers.nik = ? `, nik).First(&obj).Error; err != nil {
		return &obj, err
	}

	return &obj, nil

}

func (t *tenorModel) GetLimitByNIK(nik string) (*entity.TenorLimitCustomerEntity, error) {

	obj := entity.TenorLimitCustomerEntity{}

	if err := t.DB.Table(`tenors`).
		Select(`tenors.customer_id, customers.nik, customers.legal_name, one_month+two_month+three_month+four_month AS 'limit'`).
		Joins(`JOIN customers ON customers.id = tenors.customer_id`).
		Where(`customers.nik = ? `, nik).First(&obj).Error; err != nil {
		return &obj, err
	}

	return &obj, nil

}
