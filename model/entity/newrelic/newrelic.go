package newrelic

type NewRelicLogs struct {
	RequestUrl     string `json:"request_url"`
	RequestBody    string `json:"request_body"`
	ResponseStatus string `json:"response_status"`
	ResponseBody   string `json:"response_body"`
}
