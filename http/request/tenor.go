package request

type NewTenorRequest struct {
	NIK        string `json:"nik" validate:"max=16"`
	OneMonth   uint   `json:"one_month"  validate:"required"`
	TwoMonth   uint   `json:"two_month" validate:"required"`
	ThreeMonth uint   `json:"three_month" validate:"required"`
	FourMonth  uint   `json:"four_month" validate:"required"`
}
