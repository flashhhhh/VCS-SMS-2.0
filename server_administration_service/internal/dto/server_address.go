package dto

type ServerAddress struct {
	ServerID string `json:"server_id"`
	IPv4 string `json:"ipv4"`
	Status string `json:"status"`
}