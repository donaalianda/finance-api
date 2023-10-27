package request

type NewCustomerRequest struct {
	NIK           string `json:"nik" validate:"max=16"`
	FullName      string `json:"full_name"  validate:"required"`
	LegalName     string `json:"legal_name" validate:"required"`
	BirthPlace    string `json:"birth_place" validate:"required"`
	BirthDate     string `json:"birth_date" validate:"required"`
	Salary        uint   `json:"salary" validate:"required"`
	PhotoCustomer string `json:"photo_customer" validate:"required"`
	PhotoSelfie   string `json:"photo_selfie" validate:"required"`
}
