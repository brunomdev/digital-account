package presenter

type ErrorsResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Status int    `json:"status,omitempty"`
	Source string `json:"source,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}
