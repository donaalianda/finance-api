package api

import (
	"finance-api/http/request"
	"finance-api/model/entity"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func (_h *InjectFinanceAPIHandler) NewTransactionHandler(c echo.Context) error {
	var (
		err   error
		input request.NewTransactionRequest
	)

	err = c.Bind(&input)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	if err = _h.Helper.Validate.Struct(input); err != nil {
		return _h.Helper.SendValidationError(c, err.(validator.ValidationErrors))
	}

	customer, err := _h.TenorModel.GetLimitByNIK(input.NIK)
	if err != nil {
		return _h.Helper.SendBadRequest(c, "Customer not found >>> "+err.Error(), _h.Helper.EmptyJsonMap())
	}

	if customer.Limit > 0 {
		newTrx := entity.TransactionEntity{
			CustomerID:        customer.CustomerID,
			ContractNumber:    input.ContractNumber,
			OTR:               input.OTR,
			AdminFee:          input.AdminFee,
			AmountOfInterest:  input.AmountOfInterest,
			InstallmentAmount: input.InstallmentAmount,
			AssetName:         input.AssetName,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err = _h.TransactionModel.Create(newTrx)
		if err != nil {
			return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
		}
	} else {
		return _h.Helper.SendBadRequest(c, "Customer has no limit...", _h.Helper.EmptyJsonMap())
	}

	return _h.Helper.SendSuccess(c, "Success", _h.Helper.EmptyJsonMap())
}
