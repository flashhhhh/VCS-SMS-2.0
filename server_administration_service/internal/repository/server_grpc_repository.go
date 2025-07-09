package repository

import (
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"

	"gorm.io/gorm"
)

type ServerGRPCRepository interface {
	GetServerAddresses() ([]dto.ServerAddress, error)
}

type serverGRPCRepository struct {
	db *gorm.DB
}

func NewServerGRPCRepository(db *gorm.DB) ServerGRPCRepository {
	return &serverGRPCRepository{
		db: db,
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