package routes

import (
	"finance-api/http/api"

	"github.com/labstack/echo"
)

func CustomerRoutes(group *echo.Group, injector api.InjectFinanceAPIHandler) {

	group.POST("/customers", injector.NewCustomerHandler)
	group.GET("/customers", injector.GetCustomerByNIKHandler)
	group.PUT("/customers/:id", injector.UpdateCustomerIDHandler)

}
