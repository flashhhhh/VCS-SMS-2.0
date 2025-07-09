package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"server_administration_service/internal/domain"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/flashhhhh/pkg/env"
	"gorm.io/gorm"
)

type ServerKafkaRepository interface {
	UpdateStatus(server_id, status string) (error)
}

type serverKafkaRepository struct {
	db *gorm.DB
	es *elasticsearch.Client
}

func NewServerKafkaRepository(db *gorm.DB, es *elasticsearch.Client) ServerKafkaRepository {
	return &serverKafkaRepository{
		db: db,
		es: es,
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
	if err != nil {
		return errors.New("can't convert document to JSON")
	}

	req := esapi.IndexRequest{
		Index:   env.GetEnv("ES_NAME", "ping_status"),
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), r.es)
	if err != nil {
		return errors.New("can't send request to ES")
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("Error response from ES: " + res.String())
	}

	return nil
}