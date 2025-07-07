package repository

import (
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	eslib "server_administration_service/infrastructure/elasticsearch"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/flashhhhh/pkg/env"
	"gorm.io/gorm"
)

type ServerGRPCRepository interface {
	GetServerAddresses() ([]dto.ServerAddress, error)
	UpdateStatus(server_id, status string) (error)
}

type serverGRPCRepository struct {
	db *gorm.DB
	es *elasticsearch.Client
}

func NewServerGRPCRepository(db *gorm.DB, es *elasticsearch.Client) ServerGRPCRepository {
	return &serverGRPCRepository{
		db: db,
		es: es,
	}
}

func (r *serverGRPCRepository) GetServerAddresses() ([]dto.ServerAddress, error) {
	var serverAddresses []dto.ServerAddress
	if err := r.db.Model(&domain.Server{}).
		Select("server_id", "ipv4", "status").
		Find(&serverAddresses).Error; err != nil {
			return nil, err
		}
	
	return serverAddresses, nil
}

func (r *serverGRPCRepository) UpdateStatus(server_id, status string) (error) {
	if err := r.db.Model(&domain.Server{}).Where("server_id = ?", server_id).Update("status", status).Error; err != nil {
		return err
	}

	docs := map[string]any {
		"ID": server_id,
		"Status": status,
		"Timestamp": time.Now(),
	}
	eslib.CreateDocument(r.es, env.GetEnv("ES_NAME", "ping_status"), docs)

	return nil
}