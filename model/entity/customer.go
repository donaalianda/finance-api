package entity

import (
	"time"
)

// CustomerEntity ...
type CustomerEntity struct {
	ID            uint      `gorm:"primary_key;column:id" json:"id"`
	NIK           string    `gorm:"column:nik;unique_index" json:"nik"`
	FullName      string    `gorm:"column:full_name" json:"full_name"`
	LegalName     string    `gorm:"column:legal_name" json:"legal_name"`
	BirthPlace    string    `gorm:"column:birth_place" json:"birth_place"`
	BirthDate     string    `gorm:"column:birth_date" json:"birth_date"`
	Salary        uint      `gorm:"column:salary" json:"salary"`
	PhotoCustomer string    `gorm:"column:photo_customer" json:"photo_customer"`
	PhotoSelfie   string    `gorm:"column:photo_selfie" json:"photo_selfie"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy     string    `gorm:"column:created_by" json:"created_by"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy     string    `gorm:"column:updated_by" json:"updated_by"`
}

// Customer Table Name ...
func (CustomerEntity) TableName() string {
	return "customers"
}
