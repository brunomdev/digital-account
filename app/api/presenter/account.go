package presenter

type AccountResponse struct {
	ID                   int     `json:"account_id"`
	DocumentNumber       string  `json:"document_number"`
	AvailableCreditLimit float64 `json:"available_credit_limit"`
}
