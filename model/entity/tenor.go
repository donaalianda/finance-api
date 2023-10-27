package entity

// TenorEntity ...
type TenorEntity struct {
	ID         uint `gorm:"primary_key;column:id" json:"id"`
	CustomerID uint `gorm:"column:customer_id" json:"customer_id"`
	OneMonth   uint `gorm:"column:one_month" json:"one_month"`
	TwoMonth   uint `gorm:"column:two_month" json:"two_month"`
	ThreeMonth uint `gorm:"column:three_month" json:"trhee_month"`
	FourMonth  uint `gorm:"column:four_month" json:"four_month"`
}

// Tenor Table Name ...
func (TenorEntity) TableName() string {
	return "tenors"
}

// TenorCustomerEntity ...
type TenorCustomerEntity struct {
	ID         uint   `gorm:"primary_key;column:id" json:"id"`
	CustomerID uint   `gorm:"column:customer_id" json:"customer_id"`
	NIK        string `gorm:"column:nik" json:"nik"`
	LegalName  string `gorm:"column:legal_name" json:"legal_name"`
	OneMonth   uint   `gorm:"column:one_month" json:"one_month"`
	TwoMonth   uint   `gorm:"column:two_month" json:"two_month"`
	ThreeMonth uint   `gorm:"column:three_month" json:"trhee_month"`
	FourMonth  uint   `gorm:"column:four_month" json:"four_month"`
}

// TenorLimitCustomerEntity ...
type TenorLimitCustomerEntity struct {
	CustomerID uint   `gorm:"column:customer_id" json:"customer_id"`
	NIK        string `gorm:"column:nik" json:"nik"`
	LegalName  string `gorm:"column:legal_name" json:"legal_name"`
	Limit      uint   `gorm:"column:limit" json:"limit"`
}
