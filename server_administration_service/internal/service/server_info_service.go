package service

import (
	"server_administration_service/internal/repository"
)

type ServerInfoService interface {
	GetNumServers() (int, error)
	GetNumOnServers() (int, error)
	GetNumOffServers() (int, error)
	GetServerMeanUpTimeRatio(startTime, endTime string) (float64, error)
}

type serverInfoService struct {
	serverInfoRepository repository.ServerInfoRepository
}

func NewServerInfoService(serverInfoRepository repository.ServerInfoRepository) ServerInfoService {
	return &serverInfoService{
		serverInfoRepository: serverInfoRepository,
	}
}

func (s *serverInfoService) GetNumServers() (int, error) {
	return s.serverInfoRepository.GetNumServers()
}

func (s *serverInfoService) GetNumOnServers() (int, error) {
	return s.serverInfoRepository.GetNumOnServers()
}

func (s *serverInfoService) GetNumOffServers() (int, error) {
	return s.serverInfoRepository.GetNumOffServers()
}

func (s *serverInfoService) GetServerMeanUpTimeRatio(startTime, endTime string) (float64, error) {
	sumUpTimeRatio, err := s.serverInfoRepository.GetServerSumUpTimeRatio(startTime, endTime)
	if err != nil {
		return 0, err
	}
	numServers, err := s.GetNumServers()
	if err != nil {
		return 0, err
	}
	if numServers == 0 {
		return 0, nil
	}
	return sumUpTimeRatio / float64(numServers), nil
}