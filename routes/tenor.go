package routes

import (
	"finance-api/http/api"

	"github.com/labstack/echo"
)

func TenorRoutes(group *echo.Group, injector api.InjectFinanceAPIHandler) {

	group.POST("/tenors", injector.NewTenorHandler)
	group.GET("/tenors", injector.GetTenorCustomerByNIKHandler)
	group.GET("/tenors/limit", injector.GetTenorLimitCustomerByNIKHandler)

}
