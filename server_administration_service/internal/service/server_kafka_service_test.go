package service_test

import (
	"errors"
	"testing"

	"server_administration_service/internal/service"

	"github.com/stretchr/testify/mock"
)

// Mock implementation of ServerKafkaRepository
type mockServerKafkaRepository struct {
	mock.Mock
}

func (m *mockServerKafkaRepository) UpdateStatus(server_id, status string) error {
	args := m.Called(server_id, status)
	return args.Error(0)
}

func TestServerKafkaService_UpdateStatus_Success(t *testing.T) {
	mockRepo := new(mockServerKafkaRepository)
	service := service.NewServerKafaService(mockRepo)

	serverID := "server123"
	status := "active"

	mockRepo.On("UpdateStatus", serverID, status).Return(nil)

	err := service.UpdateStatus(serverID, status)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestServerKafkaService_UpdateStatus_Error(t *testing.T) {
	mockRepo := new(mockServerKafkaRepository)
	service := service.NewServerKafaService(mockRepo)

	serverID := "server123"
	status := "inactive"
	expectedErr := errors.New("update failed")

	mockRepo.On("UpdateStatus", serverID, status).Return(expectedErr)

	err := service.UpdateStatus(serverID, status)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	mockRepo.AssertExpectations(t)
}