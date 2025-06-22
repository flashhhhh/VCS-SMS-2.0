package dto

type ServerAddress struct {
	ID int `json:"id"`
	IPv4 string `json:"ipv4"`
	Status string `json:"status"`
}