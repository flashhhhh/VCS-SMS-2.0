package service_test

import (
	"errors"
	"reflect"
	"testing"

	"server_administration_service/internal/dto"
	"server_administration_service/internal/service"

	"github.com/stretchr/testify/mock"
)

// Mock implementation of ServerGRPCRepository
type mockServerGRPCRepository struct {
	mock.Mock
}

func (m *mockServerGRPCRepository) GetServerAddresses() ([]dto.ServerAddress, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]dto.ServerAddress), args.Error(1)
}

func (m *mockServerGRPCRepository) UpdateStatus(server_id, status string) error {
	args := m.Called(server_id, status)
	return args.Error(0)
}

func TestServerGRPCService_GetServerAddresses_Success(t *testing.T) {
	mockRepo := new(mockServerGRPCRepository)
	expected := []dto.ServerAddress{
		{ServerID: "1", IPv4: "127.0.0.1", Status: "On"},
		{ServerID: "2", IPv4: "192.168.1.1", Status: "Off"},
	}
	mockRepo.On("GetServerAddresses").Return(expected, nil)

	svc := service.NewServerGRPCService(mockRepo)
	result, err := svc.GetServerAddresses()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
	mockRepo.AssertExpectations(t)
}

func TestServerGRPCService_GetServerAddresses_Error(t *testing.T) {
	mockRepo := new(mockServerGRPCRepository)
	mockErr := errors.New("db error")
	mockRepo.On("GetServerAddresses").Return(nil, mockErr)

	svc := service.NewServerGRPCService(mockRepo)
	result, err := svc.GetServerAddresses()

	if err != mockErr {
		t.Errorf("expected error %v, got %v", mockErr, err)
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
	mockRepo.AssertExpectations(t)
}