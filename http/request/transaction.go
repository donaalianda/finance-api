package request

type NewTransactionRequest struct {
	NIK               string  `json:"nik" validate:"max=16"`
	ContractNumber    string  `json:"contract_number"  validate:"required"`
	OTR               uint    `json:"otr" validate:"required"`
	InstallmentAmount uint    `json:"installment_amount" validate:"gt=1,lt=60"`
	AdminFee          uint    `json:"admin_fee" validate:"required"`
	AmountOfInterest  float64 `json:"amount_of_interest" validate:"required"`
	AssetName         string  `json:"asset_name" validate:"required"`
}
