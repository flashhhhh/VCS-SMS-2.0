package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

// Mock implementation of ServerInfoRepository
type mockServerInfoRepository struct {
	mock.Mock
}

func (m *mockServerInfoRepository) GetNumServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *mockServerInfoRepository) GetNumOnServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *mockServerInfoRepository) GetNumOffServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *mockServerInfoRepository) GetServerSumUpTimeRatio(startTime, endTime string) (float64, error) {
	args := m.Called(startTime, endTime)
	return args.Get(0).(float64), args.Error(1)
}

func TestGetNumServers(t *testing.T) {
	mockRepo := new(mockServerInfoRepository)
	mockRepo.On("GetNumServers").Return(5, nil)

	service := NewServerInfoService(mockRepo)
	num, err := service.GetNumServers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if num != 5 {
		t.Errorf("expected 5, got %d", num)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetNumOnServers(t *testing.T) {
	mockRepo := new(mockServerInfoRepository)
	mockRepo.On("GetNumOnServers").Return(3, nil)

	service := NewServerInfoService(mockRepo)
	num, err := service.GetNumOnServers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if num != 3 {
		t.Errorf("expected 3, got %d", num)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetNumOffServers(t *testing.T) {
	mockRepo := new(mockServerInfoRepository)
	mockRepo.On("GetNumOffServers").Return(2, nil)

	service := NewServerInfoService(mockRepo)
	num, err := service.GetNumOffServers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if num != 2 {
		t.Errorf("expected 2, got %d", num)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetServerMeanUpTimeRatio_Success(t *testing.T) {
	mockRepo := new(mockServerInfoRepository)
	mockRepo.On("GetServerSumUpTimeRatio", "2024-01-01", "2024-01-31").Return(10.0, nil)
	mockRepo.On("GetNumServers").Return(5, nil)

	service := NewServerInfoService(mockRepo)
	ratio, err := service.GetServerMeanUpTimeRatio("2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := 2.0
	if ratio != expected {
		t.Errorf("expected %v, got %v", expected, ratio)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetServerMeanUpTimeRatio_RepoError(t *testing.T) {
	mockRepo := new(mockServerInfoRepository)
	mockRepo.On("GetServerSumUpTimeRatio", "2024-01-01", "2024-01-31").Return(0.0, errors.New("repo error"))

	service := NewServerInfoService(mockRepo)
	_, err := service.GetServerMeanUpTimeRatio("2024-01-01", "2024-01-31")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	mockRepo.AssertExpectations(t)
}

func TestGetServerMeanUpTimeRatio_NumServersError(t *testing.T) {
	mockRepo := new(mockServerInfoRepository)
	mockRepo.On("GetServerSumUpTimeRatio", "2024-01-01", "2024-01-31").Return(10.0, nil)
	mockRepo.On("GetNumServers").Return(0, errors.New("num servers error"))

	service := NewServerInfoService(mockRepo)
	_, err := service.GetServerMeanUpTimeRatio("2024-01-01", "2024-01-31")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	mockRepo.AssertExpectations(t)
}

func TestGetServerMeanUpTimeRatio_ZeroServers(t *testing.T) {
	mockRepo := new(mockServerInfoRepository)
	mockRepo.On("GetServerSumUpTimeRatio", "2024-01-01", "2024-01-31").Return(10.0, nil)
	mockRepo.On("GetNumServers").Return(0, nil)

	service := NewServerInfoService(mockRepo)
	ratio, err := service.GetServerMeanUpTimeRatio("2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ratio != 0 {
		t.Errorf("expected 0, got %v", ratio)
	}
	mockRepo.AssertExpectations(t)
}