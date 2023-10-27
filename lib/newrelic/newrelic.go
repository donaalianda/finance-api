package newrelic

import (
	"finance-api/config/env"
	"finance-api/config/helper/str"
	nrEntity "finance-api/model/entity/newrelic"

	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var SuccessResponseCodes = []string{"200", "201"}

func SendOutgoingLogToNewRelic(config env.Config, url, thirdParty string, reqBody interface{}, restyResp *resty.Response, app *newrelic.Application) {
	if config == nil || app == nil || reqBody == nil || restyResp == nil {
		return
	}

	type RequestBodyLogs struct {
		PayloadRequest string `json:"payload_request"`
		Purpose        string `json:"purpose"`
		ThirdParty     string `json:"third_party"`
	}

	if config.GetBool(`new_relic.enable`) && config.GetBool(`new_relic.enable_outgoing_logging`) {

		responseCode := strconv.Itoa(restyResp.StatusCode())

		// ignore send logs to new relic if response success and enable_logging_on_success_outgoing is false
		if !config.GetBool(`new_relic.enable_logging_on_success_outgoing`) && str.StringContains(SuccessResponseCodes, responseCode) {
			return
		}

		severity := "error"
		if str.StringContains(SuccessResponseCodes, responseCode) {
			severity = "info"
		}

		reqBodyByte, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println(err)
			return
		}

		payloadRequest := RequestBodyLogs{
			PayloadRequest: "###" + string(reqBodyByte) + "###",
			Purpose:        "outgoing-logging",
			ThirdParty:     thirdParty,
		}

		requestBody, err := json.Marshal(payloadRequest)
		if err != nil {
			fmt.Println("unable to marshal newrelic outgoing request body", err)
			return
		}

		newRelicLogs := nrEntity.NewRelicLogs{
			RequestUrl:     url,
			RequestBody:    string(requestBody),
			ResponseStatus: responseCode,
			ResponseBody:   string(restyResp.Body()),
		}

		newRelicLogsJson, err := json.Marshal(newRelicLogs)
		if err != nil {
			fmt.Println("marshal error newrelic outgoing logs >>> " + err.Error())
			return
		}

		fmt.Println("send logs to newrelic outgoing logs")
		loc, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(loc)
		app.RecordLog(newrelic.LogData{
			Timestamp: now.UnixNano() / int64(time.Millisecond),
			Message:   string(newRelicLogsJson),
			Severity:  severity,
		})
	}
}

func SendIncomingLogToNewRelic(config env.Config, url, reqBody, resBody, responseCode string, app *newrelic.Application) {
	if config == nil || app == nil {
		return
	}

	type RequestBodyLogs struct {
		PayloadRequest string `json:"payload_request"`
		Purpose        string `json:"purpose"`
	}

	if config.GetBool(`new_relic.enable`) && config.GetBool(`new_relic.enable_incoming_logging`) {

		// ignore send logs to new relic if response success and enable_logging_on_success_incoming is false
		if !config.GetBool(`new_relic.enable_logging_on_success_incoming`) && str.StringContains(SuccessResponseCodes, responseCode) {
			return
		}

		severity := "error"
		if str.StringContains(SuccessResponseCodes, responseCode) {
			severity = "info"
		}

		payloadRequest := RequestBodyLogs{
			PayloadRequest: "###" + reqBody + "###",
			Purpose:        "incoming-logging",
		}

		requestBody, err := json.Marshal(payloadRequest)
		if err != nil {
			fmt.Println("unable to marshal newrelic incoming request body", err)
			return
		}

		newRelicLogs := nrEntity.NewRelicLogs{
			RequestUrl:     url,
			RequestBody:    string(requestBody),
			ResponseStatus: responseCode,
			ResponseBody:   resBody,
		}

		newRelicLogsJson, err := json.Marshal(newRelicLogs)
		if err != nil {
			fmt.Println("marshal error newrelic incoming logs >>> " + err.Error())
			return
		}

		fmt.Println("send logs to newrelic incoming logs")
		loc, _ := time.LoadLocation("Asia/Jakarta")
		now := time.Now().In(loc)
		app.RecordLog(newrelic.LogData{
			Timestamp: now.UnixNano() / int64(time.Millisecond),
			Message:   string(newRelicLogsJson),
			Severity:  severity,
		})
	}
}
