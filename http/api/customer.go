package api

import (
	"finance-api/http/request"
	"finance-api/http/response"
	"finance-api/model/entity"
	"strconv"

	"strings"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func (_h *InjectFinanceAPIHandler) NewCustomerHandler(c echo.Context) error {
	var (
		err   error
		input request.NewCustomerRequest
	)

	err = c.Bind(&input)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	if err = _h.Helper.Validate.Struct(input); err != nil {
		return _h.Helper.SendValidationError(c, err.(validator.ValidationErrors))
	}

	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return _h.Helper.SendBadRequest(c, "Invalid birth date (YYYY-MM-DD) "+err.Error(), _h.Helper.EmptyJsonMap())
	}

	newCustomer := entity.CustomerEntity{
		NIK:           input.NIK,
		FullName:      input.FullName,
		LegalName:     input.LegalName,
		BirthPlace:    strings.ToUpper(input.BirthPlace),
		BirthDate:     birthDate.Format("2006-01-02"),
		Salary:        input.Salary,
		PhotoCustomer: input.PhotoCustomer,
		PhotoSelfie:   input.PhotoSelfie,
		CreatedBy:     "SYSTEM",
		CreatedAt:     time.Now(),
	}

	err = _h.CustomerModel.Create(newCustomer)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	return _h.Helper.SendSuccess(c, "Success", _h.Helper.EmptyJsonMap())
}

func (_h *InjectFinanceAPIHandler) GetCustomerByNIKHandler(c echo.Context) error {
	var (
		err    error
		nik    = c.QueryParam("nik")
		result = response.CustomerResponse{}
	)

	obj, err := _h.CustomerModel.GetByNIK(nik)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	result = response.CustomerResponse{
		ID:            obj.ID,
		NIK:           obj.NIK,
		FullName:      obj.FullName,
		LegalName:     obj.LegalName,
		BirthPlace:    obj.BirthPlace,
		BirthDate:     obj.BirthDate[0:10],
		Salary:        int(obj.Salary),
		PhotoCustomer: obj.PhotoCustomer,
		PhotoSelfie:   obj.PhotoSelfie,
	}

	return _h.Helper.SendSuccess(c, "Success", result)
}

func (_h *InjectFinanceAPIHandler) UpdateCustomerIDHandler(c echo.Context) error {
	var (
		err   error
		id    = c.Param("id")
		input request.NewCustomerRequest
	)

	if err = c.Bind(&input); err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	if err = _h.Helper.Validate.Struct(input); err != nil {
		return _h.Helper.SendValidationError(c, err.(validator.ValidationErrors))
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")

	idInt, _ := strconv.Atoi(id)

	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return _h.Helper.SendBadRequest(c, "Invalid birth date (YYYY-MM-DD) "+err.Error(), _h.Helper.EmptyJsonMap())
	}

	obj := entity.CustomerEntity{
		ID:            uint(idInt),
		NIK:           input.NIK,
		FullName:      input.FullName,
		LegalName:     input.LegalName,
		BirthPlace:    strings.ToUpper(input.BirthPlace),
		BirthDate:     birthDate.Format("2006-01-02"),
		Salary:        input.Salary,
		PhotoCustomer: input.PhotoCustomer,
		PhotoSelfie:   input.PhotoSelfie,
		CreatedBy:     "SYSTEM",
		CreatedAt:     time.Now().In(loc),
		UpdatedBy:     "SYSTEM",
		UpdatedAt:     time.Now().In(loc),
	}

	err = _h.CustomerModel.Update(obj)
	if err != nil {
		return _h.Helper.SendBadRequest(c, err.Error(), _h.Helper.EmptyJsonMap())
	}

	return _h.Helper.SendSuccess(c, "Success", _h.Helper.EmptyJsonMap())
}
