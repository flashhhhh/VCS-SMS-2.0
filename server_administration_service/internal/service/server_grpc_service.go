package service

import (
	"server_administration_service/internal/dto"
	"server_administration_service/internal/repository"
)

type ServerGRPCService interface {
	GetServerAddresses() ([]dto.ServerAddress, error)
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