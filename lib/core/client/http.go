package client

import (
	"encoding/json"
	"errors"

	"finance-api/config/env"
	nr "finance-api/lib/newrelic"

	"github.com/go-resty/resty"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Info struct {
	XPlayer     string
	Method      string
	Url         string
	Auth        string
	PaxelKey    string
	Payload     map[string]interface{}
	Config      env.Config
	NewRelicApp *newrelic.Application
}

func (i *Info) prepare() *resty.Request {
	req := resty.DefaultClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Player", i.XPlayer).
		SetHeader("Paxel-Key", i.PaxelKey)

	if i.Auth != "" {
		req.SetAuthToken(i.Auth)
	}

	return req
}

func (i *Info) Dispatch(target interface{}) error {
	var (
		err  error
		resp *resty.Response
	)
	req := i.prepare()
	if i.Method == "" {
		return errors.New("Missing HTTP method for dispatch request to core.")
	}

	if len(i.Payload) == 0 || i.Payload == nil {
		resp, err = req.Execute(i.Method, i.Url)
	} else {
		resp, err = req.SetBody(i.Payload).Execute(i.Method, i.Url)
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(resp.String()), target)
	if err != nil {
		return err
	}

	if i.Config != nil && i.NewRelicApp != nil {
		nr.SendOutgoingLogToNewRelic(i.Config, i.Url, i.Url, i.Payload, resp, i.NewRelicApp)
	}

	return nil
}
