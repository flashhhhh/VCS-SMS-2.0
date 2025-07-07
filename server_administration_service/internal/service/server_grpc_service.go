package service

import (
	"server_administration_service/internal/dto"
	"server_administration_service/internal/repository"
)

type ServerGRPCService interface {
	GetServerAddresses() ([]dto.ServerAddress, error)
	UpdateStatus(server_id, status string) (error)
}

type serverGRPCService struct {
	serverGRPCRepository repository.ServerGRPCRepository
}

func NewServerGRPCService(serverGRPCRepository repository.ServerGRPCRepository) ServerGRPCService {
	return &serverGRPCService{
		serverGRPCRepository: serverGRPCRepository,
	}
}

func (s *serverGRPCService) GetServerAddresses() ([]dto.ServerAddress, error) {
	return s.serverGRPCRepository.GetServerAddresses()
}

func (s *serverGRPCService) UpdateStatus(server_id, status string) (error) {
	return s.serverGRPCRepository.UpdateStatus(server_id, status)
}