package client

type ClientWithTotal struct {
	Clients []Client `json:"clients"`
	Total int64 `json:"total"`
}