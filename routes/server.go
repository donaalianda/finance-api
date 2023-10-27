package routes

import (
	"finance-api/config/env"
	"finance-api/config/helper/str"
	"finance-api/http/api"
	httpHelper "finance-api/http/helper"
	nr "finance-api/lib/newrelic"
	"finance-api/model"

	"fmt"
	"log"
	"strconv"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v3"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gopkg.in/go-playground/validator.v9"
)

// HTTPHandler ...
type HTTPHandler struct {
	E               *echo.Echo
	Config          env.Config
	Helper          httpHelper.HTTPHelper
	ValidatorDriver *validator.Validate
	Translator      ut.Translator
}

// RegisterAPIHandler ...
func (h *HTTPHandler) RegisterAPIHandler() *HTTPHandler {
	h.Helper = httpHelper.HTTPHelper{
		Validate:   h.ValidatorDriver,
		Translator: h.Translator,
	}

	// model initialize
	db := model.Info{Config: h.Config}
	customerModel := model.NewCustomerModel(db.Connect())
	tenorModel := model.NewTenorModel(db.Connect())
	transactionModel := model.NewTransactionModel(db.Connect())

	financeAPIHandler := api.InjectFinanceAPIHandler{
		Config:           h.Config,
		Helper:           h.Helper,
		CustomerModel:    customerModel,
		TenorModel:       tenorModel,
		TransactionModel: transactionModel,
	}

	group := h.E.Group(`api/v1`)
	group.GET("/ping", financeAPIHandler.PingHandler)

	//customer routes
	CustomerRoutes(group, financeAPIHandler)

	//tenor routes
	TenorRoutes(group, financeAPIHandler)

	//transaction routes
	TransactionRoutes(group, financeAPIHandler)

	return h
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "finance-api/1.0")
		return next(c)
	}
}

// RegisterMiddleware ...
func (h *HTTPHandler) RegisterMiddleware() {
	h.E.Use(serverHeader)
	h.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	h.E.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	if h.Config.GetBool(`app.debug`) == true {
		h.E.Use(middleware.Logger())
		h.E.HideBanner = true
		h.E.Debug = true
	} else {
		h.E.HideBanner = true
		h.E.Debug = false
		h.E.Use(middleware.Recover())
	}
}

func (h *HTTPHandler) RegisterRequestReturnLogger() {
	h.E.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		req := c.Request()
		res := c.Response()
		if req.URL.String() != "/" && h.Config.GetBool(`app.debug`) == true {
			log.Println("Request URL => " + req.URL.String())
			log.Println("Request Body => " + string(reqBody))
			log.Println("Response Status => " + strconv.Itoa(res.Status))
			log.Println("Response Body => " + string(resBody))
		}
	}))
}

func (h *HTTPHandler) RegisterNewRelic() *newrelic.Application {
	//newRelic
	if h.Config.GetBool(`new_relic.enable`) == true {
		app, errNewRelic := newrelic.NewApplication(
			newrelic.ConfigAppName(h.Config.GetString(`new_relic.name`)),
			newrelic.ConfigLicense(h.Config.GetString(`new_relic.key`)),
			newrelic.ConfigAppLogForwardingEnabled(true),
			//newrelic.ConfigDebugLogger(os.Stdout),
		)
		if errNewRelic == nil {
			fmt.Println("use new relic >>> true")
			// The New Relic Middleware should be the first middleware registered
			h.E.Use(nrecho.Middleware(app))

			h.E.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
				req := c.Request()
				res := c.Response()

				disabledMonitoringEndpoint := h.Config.GetStringSlice("new_relic.disabled_monitoring_endpoint")
				isDisableSendToNewrelic := str.StringContainsPrefix(disabledMonitoringEndpoint, req.URL.String())

				if isDisableSendToNewrelic {
					return
				}

				if req.URL.String() != "/" && !strings.Contains(req.URL.String(), "echo.go") {
					nr.SendIncomingLogToNewRelic(h.Config, req.URL.String(), string(reqBody), string(resBody), strconv.Itoa(res.Status), app)
				}
			}))

			return app
		}
	}

	return nil
}
