package response

type PaginationResponse struct {
	CurrentPage      int           `json:"current_page"`
	PerPage          int           `json:"per_page"`
	TotalPages       int           `json:"total_pages"`
	TotalRecords     int           `json:"total_records"`
	LinkParameter    LinksResponse `json:"link_parameter"`
	Links            LinksResponse `json:"links"`
	CurrentParameter string        `json:"current_parameter"`
}

type LinksResponse struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}
