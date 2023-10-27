package api

import (
	"finance-api/config/env"
	httpHelper "finance-api/http/helper"
	"finance-api/http/helper/logger"
	"finance-api/model"

	"github.com/labstack/echo"
)

// InjectFinanceAPIHandler ...
type InjectFinanceAPIHandler struct {
	Config           env.Config
	Helper           httpHelper.HTTPHelper
	CustomerModel    model.CustomerModelInterface
	TenorModel       model.TenorModelInterface
	TransactionModel model.TransactionModelInterface
}

// PingHandler ...
func (_h *InjectFinanceAPIHandler) PingHandler(c echo.Context) error {
	var err error
	/*********************** Jaeger Log ************************/
	closer, span, starttime := logger.JaegerStart(_h.Config, "PricetimeService", "PingHandler", "PingHandler Process started")
	if _h.Config.GetBool(`logger.jaeger.jaeger_enabled`) {
		defer func() {
			if err != nil {
				logger.JaegerEnd(_h.Config, starttime, span, "Error", "PingHandler Process Failed : "+err.Error())
			} else {
				logger.JaegerEnd(_h.Config, starttime, span, "Success", "PingHandler Process has been completed")
			}
			span.Finish()
			closer.Close()
		}()
	}
	/*********************** End of Jaeger Log ************************/
	return _h.Helper.SendSuccess(c, "Success", map[string]interface{}{
		"ping": "pong",
	})
}
