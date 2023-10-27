package api

import (
	"finance-api/http/request"
	"finance-api/model/entity"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func (_h *InjectFinanceAPIHandler) NewTenorHandler(c echo.Context) error {
	var (
		err   error
		input request.NewTenorRequest
	)

	err = c.Bind(&input)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	if err = _h.Helper.Validate.Struct(input); err != nil {
		return _h.Helper.SendValidationError(c, err.(validator.ValidationErrors))
	}

	customer, err := _h.CustomerModel.GetByNIK(input.NIK)
	if err != nil {
		return _h.Helper.SendBadRequest(c, "Customer not found >>> "+err.Error(), _h.Helper.EmptyJsonMap())
	}

	newTenor := entity.TenorEntity{
		CustomerID: customer.ID,
		OneMonth:   input.OneMonth,
		TwoMonth:   input.TwoMonth,
		ThreeMonth: input.ThreeMonth,
		FourMonth:  input.FourMonth,
	}

	err = _h.TenorModel.Create(newTenor)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	return _h.Helper.SendSuccess(c, "Success", _h.Helper.EmptyJsonMap())
}

func (_h *InjectFinanceAPIHandler) GetTenorCustomerByNIKHandler(c echo.Context) error {
	var (
		err error
		nik = c.QueryParam("nik")
	)

	obj, err := _h.TenorModel.GetByNIK(nik)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	tenor := map[string]interface{}{
		"tenor_id":    obj.ID,
		"one_month":   obj.OneMonth,
		"two_month":   obj.TwoMonth,
		"three_month": obj.ThreeMonth,
		"four_month":  obj.FourMonth,
	}

	tenorCustomer := map[string]interface{}{
		"customer_id": obj.CustomerID,
		"nik":         obj.NIK,
		"legal_name":  obj.LegalName,
		"tenor":       tenor,
	}

	return _h.Helper.SendSuccess(c, "Success", tenorCustomer)
}

func (_h *InjectFinanceAPIHandler) GetTenorLimitCustomerByNIKHandler(c echo.Context) error {
	var (
		err error
		nik = c.QueryParam("nik")
	)

	obj, err := _h.TenorModel.GetLimitByNIK(nik)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	tenorCustomer := map[string]interface{}{
		"nik":        obj.NIK,
		"legal_name": obj.LegalName,
		"limit":      obj.Limit,
	}

	return _h.Helper.SendSuccess(c, "Success", tenorCustomer)
}
