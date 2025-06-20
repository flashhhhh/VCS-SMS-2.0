package domain

import "time"

type Server struct {
	ID int `json:"id" gorm:"autoIncrement"`
	ServerID string `json:"server_id" gorm:"primary_key;unique"`
	ServerName string `json:"server_name" gorm:"unique;not null"`
	Status string `json:"status" gorm:"not null"`
	CreatedTime time.Time `json:"created_time" gorm:"autoCreateTime"`
	LastUpdated time.Time `json:"last_updated" gorm:"autoUpdateTime"`
	IPv4 string `json:"ipv4" gorm:"not null;unique"`
}
