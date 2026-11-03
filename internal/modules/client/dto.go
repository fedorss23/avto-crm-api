package client

type ClientWithTotal struct {
	Clients []Client `json:"clients"`
	Total int64 `json:"total"`
}

type UpdateClientRequest struct {
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

type OneClientResponse struct {
	Client Client `json:"client"`
}