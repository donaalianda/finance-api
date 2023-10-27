package model

import (
	"finance-api/model/entity"

	"github.com/jinzhu/gorm"
)

type (
	CustomerModelInterface interface {
		Create(obj entity.CustomerEntity) error
		GetByNIK(nik string) (*entity.CustomerEntity, error)
		Update(obj entity.CustomerEntity) error
	}

	customerModel struct {
		DB *gorm.DB
	}
)

func NewCustomerModel(db *gorm.DB) CustomerModelInterface {
	return &customerModel{db}
}

func (t *customerModel) Create(obj entity.CustomerEntity) error {

	tx := t.DB.Begin()

	if err := tx.Save(&obj).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil

}

func (t *customerModel) GetByNIK(nik string) (*entity.CustomerEntity, error) {

	obj := entity.CustomerEntity{}

	if err := t.DB.Where(`nik = ? `, nik).First(&obj).Error; err != nil {
		return &obj, err
	}

	return &obj, nil

}

func (t *customerModel) Update(obj entity.CustomerEntity) error {

	tx := t.DB.Begin()

	updatedField := map[string]interface{}{
		"updated_at":     obj.UpdatedAt,
		"updated_by":     obj.UpdatedBy,
		"nik":            obj.NIK,
		"full_name":      obj.FullName,
		"legal_name":     obj.LegalName,
		"birth_place":    obj.BirthPlace,
		"birth_date":     obj.BirthDate,
		"photo_customer": obj.PhotoCustomer,
		"photo_selfie":   obj.PhotoSelfie,
	}

	err := tx.Model(&obj).Where(`id = ?`, obj.ID).Updates(updatedField).Error
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
