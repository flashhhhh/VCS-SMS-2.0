package repository

import (
	"context"
	"encoding/json"
	"server_administration_service/infrastructure/elasticsearch"
	"server_administration_service/internal/domain"
	"time"

	"github.com/flashhhhh/pkg/env"
	"gorm.io/gorm"
)

type ServerKafkaRepository interface {
	UpdateStatus(server_id, status string) (error)
}

type serverKafkaRepository struct {
	db  *gorm.DB
	esc elasticsearch.ElasticsearchClient
}

func NewServerKafkaRepository(db *gorm.DB, esc elasticsearch.ElasticsearchClient) ServerKafkaRepository {
	return &serverKafkaRepository{
		db:  db,
		esc: esc,
	}
}

func (r *serverKafkaRepository) UpdateStatus(server_id, status string) (error) {
	if err := r.db.Model(&domain.Server{}).Where("server_id = ?", server_id).Update("status", status).Error; err != nil {
		return err
	}

	docs := map[string]any {
		"ID": server_id,
		"Status": status,
		"Timestamp": time.Now(),
	}

	data, err := json.Marshal(docs)

	err = r.esc.Index(context.Background(), env.GetEnv("ES_NAME", "ping_status"), data)
	return err
}