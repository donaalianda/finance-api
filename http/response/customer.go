package response

type CustomerResponse struct {
	ID            uint   `json:"id"`
	NIK           string `json:"nik"`
	FullName      string `json:"full_name"`
	LegalName     string `json:"legal_name"`
	BirthPlace    string `json:"birth_place"`
	BirthDate     string `json:"birth_date"`
	Salary        int    `json:"salary"`
	PhotoCustomer string `json:"photo_customer"`
	PhotoSelfie   string `json:"photo_selfie"`
}
