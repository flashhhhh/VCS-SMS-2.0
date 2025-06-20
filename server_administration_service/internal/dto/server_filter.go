package dto

type ServerFilter struct {
	ServerID string `json:"server_id"`
	ServerName string `json:"server_name"`
	Status	 string `json:"status"`
	IPv4	  string `json:"ipv4"`
}