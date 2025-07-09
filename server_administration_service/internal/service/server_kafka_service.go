package service

import (
	"server_administration_service/internal/repository"
)

type ServerKafkaService interface {
	UpdateStatus(server_id, status string) (error)
}

type serverKafkaService struct {
	serverKafkaRepository repository.ServerKafkaRepository
}

func NewServerKafaService(serverKafkaRepository repository.ServerKafkaRepository) ServerKafkaService {
	return &serverKafkaService{
		serverKafkaRepository: serverKafkaRepository,
	}
}

func (s *serverKafkaService) UpdateStatus(server_id, status string) (error) {
	return s.serverKafkaRepository.UpdateStatus(server_id, status)
}