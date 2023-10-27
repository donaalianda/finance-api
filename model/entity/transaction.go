package entity

import (
	"time"
)

// TransactionEntity ...
type TransactionEntity struct {
	ID                uint      `gorm:"primary_key;column:id" json:"id"`
	CustomerID        uint      `gorm:"column:customer_id" json:"customer_id"`
	ContractNumber    string    `gorm:"column:contract_number" json:"contract_number"`
	OTR               uint      `gorm:"column:otr" json:"otr"`
	AdminFee          uint      `gorm:"column:admin_fee" json:"admin_fee"`
	InstallmentAmount uint      `gorm:"column:installment_amount" json:"installment_amount"`
	AmountOfInterest  float64   `gorm:"column:amount_of_interest" json:"amount_of_interest"`
	AssetName         string    `gorm:"column:asset_name" json:"asset_name"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// Transaction Table Name ...
func (TransactionEntity) TableName() string {
	return "transactions"
}
