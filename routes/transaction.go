package routes

import (
	"finance-api/http/api"

	"github.com/labstack/echo"
)

func TransactionRoutes(group *echo.Group, injector api.InjectFinanceAPIHandler) {

	group.POST("/transactions", injector.NewTransactionHandler)

}
