package repository

import (
	"fmt"
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"

	"github.com/flashhhhh/pkg/logging"
	"gorm.io/gorm"
)

type ServerCRUDRepository interface {
	CreateServer(server *domain.Server) (int, error)
	CreateServers(servers []domain.Server) ([]domain.Server, []domain.Server, error)
	ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error)
	UpdateServer(server_id string, updatedData map[string]interface{}) error
	DeleteServer(serverID string) error
}

type serverCRUDRepository struct {
	db *gorm.DB
}

func NewServerCRUDRepository(db *gorm.DB) ServerCRUDRepository {
	return &serverCRUDRepository{
		db: db,
	}
}

func (r *serverCRUDRepository) CreateServer(server *domain.Server) (int, error) {
	err := r.db.Create(server).Error
	if (err != nil) {
		return 0, err
	}

	// Write to Elasticsearch
	// doc := map[string]any{
	// 	"ID": server.ID,
	// 	"Status": "On",
	// 	"Timestamp": time.Now(),
	// }

	// eslib.CreateDocument(r.es, env.GetEnv("ES_NAME", "ping_status"), doc)

	// doc = map[string]any {
	// 	"ID": server.ID,
	// 	"Status": "Off",
	// 	"Timestamp": time.Now(),
	// }

	// eslib.CreateDocument(r.es, env.GetEnv("ES_NAME", "ping_status"), doc)
	
	return server.ID, nil
}

func (r *serverCRUDRepository) CreateServers(servers []domain.Server) ([]domain.Server, []domain.Server, error) {
	query := `
		INSERT INTO servers (server_id, server_name, status, ipv4) VALUES 
	`

	for i, server := range servers {
		query += fmt.Sprintf("('%s', '%s', '%s', '%s')",
			server.ServerID, server.ServerName, server.Status, server.IPv4)
		
		if i < len(servers)-1 {
			query += ", "
		}
	}

	query += " ON CONFLICT DO NOTHING RETURNING *"

	var result []domain.Server
	err := r.db.Raw(query).Scan(&result).Error
	if err != nil {
		logging.LogMessage("server_administration_service", "Error inserting servers: " + err.Error(), "ERROR")
		return nil, nil, err
	}

	// Determine non-inserted records
	insertedMap := make(map[string]bool)
	var insertedServer, nonInsertedServer []domain.Server

	for _, server := range result {
		insertedMap[server.ServerID] = true
		insertedServer = append(insertedServer, server)
	}

	for _, server := range servers {
		if !insertedMap[server.ServerID] {
			nonInsertedServer = append(nonInsertedServer, server)
		}
	}

	return insertedServer, nonInsertedServer, nil
}

func (r *serverCRUDRepository) ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error) {
	var servers []domain.Server
	query := r.db.Model(&domain.Server{})

	if serverFilter.ServerID != "" {
		query = query.Where("server_id = ?", serverFilter.ServerID)
	}

	if serverFilter.ServerName != "" {	
		query = query.Where("server_name LIKE ?", "%"+serverFilter.ServerName+"%")
	}

	if serverFilter.Status != "" {
		query = query.Where("status = ?", serverFilter.Status)
	}

	if serverFilter.IPv4 != "" {
		query = query.Where("ipv4 = ?", serverFilter.IPv4)
	}

	// sortedColumn is mandatory
	err := query.Order(sortedColumn + " " + order).Offset(from).Limit(to - from).Find(&servers).Error
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *serverCRUDRepository) UpdateServer(serverID string, updatedData map[string]interface{}) error {
	if err := r.db.Model(&domain.Server{}).
		Where("server_id = ?", serverID).
		Updates(updatedData).Error; err != nil {
			return err
		}

	// Get the server's id
	var server domain.Server
	if err := r.db.Where("server_id = ?", serverID).First(&server).Error; err != nil {
		return err
	}

	return nil
}

func (r *serverCRUDRepository) DeleteServer(serverID string) error {
	// Get the server's ID before deleting
	var server domain.Server
	if err := r.db.Where("server_id = ?", serverID).First(&server).Error; err != nil {
		return err
	}

	// Delete the server
	if err := r.db.Where("server_id = ?", serverID).Delete(&domain.Server{}).Error; err != nil {
		return err
	}

	return nil
}