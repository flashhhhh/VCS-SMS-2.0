package repository

import (
	"context"
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	eslib "server_administration_service/infrastructure/elasticsearch"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/flashhhhh/pkg/env"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServerGRPCRepository interface {
	GetServerAddresses() ([]dto.ServerAddress, error)
	UpdateStatus(id int, status string) (error)
}

type serverGRPCRepository struct {
	db *gorm.DB
	redis *redis.Client
	es *elasticsearch.Client
}

func NewServerGRPCRepository(db *gorm.DB, redis *redis.Client, es *elasticsearch.Client) ServerGRPCRepository {
	return &serverGRPCRepository{
		db: db,
		redis: redis,
		es: es,
	}
}

func (r *serverGRPCRepository) GetServerAddresses() ([]dto.ServerAddress, error) {
	var serverAddresses []dto.ServerAddress
	if err := r.db.Model(&domain.Server{}).
		Select("id", "ipv4", "status").
		Find(&serverAddresses).Error; err != nil {
			return nil, err
		}
	
	return serverAddresses, nil
}

func (r *serverGRPCRepository) UpdateStatus(id int, status string) (error) {
	statusInt := 0
	if status == "On" {
		statusInt = 1
	}

	r.redis.SetBit(context.Background(), env.GetEnv("REDIS_BITMAP", "server_status"), int64(id), statusInt)

	if err := r.db.Model(&domain.Server{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	docs := map[string]any {
		"ID": id,
		"Status": status,
		"Timestamp": time.Now(),
	}
	eslib.CreateDocument(r.es, env.GetEnv("ES_NAME", "ping_status"), docs)

	return nil
}