package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"server_administration_service/internal/handler"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockServerKafkaService struct {
	mock.Mock
}

func (m *mockServerKafkaService) UpdateStatus(serverID, status string) error {
	args := m.Called(serverID, status)
	return args.Error(0)
}

type mockConsumerGroupSession struct {
	mock.Mock
}

func (m *mockConsumerGroupSession) Claims() map[string][]int32                      { return nil }
func (m *mockConsumerGroupSession) MemberID() string                                { return "" }
func (m *mockConsumerGroupSession) GenerationID() int32                             { return 0 }
func (m *mockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}
func (m *mockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}
func (m *mockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.Called(msg, metadata)
}
func (m *mockConsumerGroupSession) Context() context.Context { return nil }

// Add missing Commit method to satisfy sarama.ConsumerGroupSession interface
func (m *mockConsumerGroupSession) Commit() {}

type mockConsumerGroupClaim struct {
	mock.Mock
	messages chan *sarama.ConsumerMessage
}

func (m *mockConsumerGroupClaim) Topic() string               { return "test-topic" }
func (m *mockConsumerGroupClaim) Partition() int32            { return 0 }
func (m *mockConsumerGroupClaim) InitialOffset() int64        { return 0 }
func (m *mockConsumerGroupClaim) HighWaterMarkOffset() int64  { return 0 }
func (m *mockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	return m.messages
}

func TestSetup(t *testing.T) {
	mockService := new(mockServerKafkaService)
	handler := handler.NewServerConsumerHandler(mockService)

	mockSession := new(mockConsumerGroupSession)
	err := handler.Setup(mockSession)

	assert.NoError(t, err)
}

func TestCleanup(t *testing.T) {
	mockService := new(mockServerKafkaService)
	handler := handler.NewServerConsumerHandler(mockService)

	mockSession := new(mockConsumerGroupSession)
	err := handler.Cleanup(mockSession)

	assert.NoError(t, err)
}

func TestConsumeClaim(t *testing.T) {
	mockService := new(mockServerKafkaService)
	handler := handler.NewServerConsumerHandler(mockService)

	mockSession := new(mockConsumerGroupSession)
	mockClaim := &mockConsumerGroupClaim{
		messages: make(chan *sarama.ConsumerMessage, 1),
	}

	serverMessage := map[string]string{
		"server_id": "srv123",
		"status":    "running",
	}
	msgValue, _ := json.Marshal(serverMessage)

	message := &sarama.ConsumerMessage{
		Value: msgValue,
	}

	mockSession.On("MarkMessage", message, "").Return()
	mockService.On("UpdateStatus", "srv123", "running").Return(nil)

	// Send the message into the channel and close it after a short delay
	mockClaim.messages <- message
	close(mockClaim.messages)

	err := handler.ConsumeClaim(mockSession, mockClaim)
	assert.NoError(t, err)

	// Allow goroutine to complete
	time.Sleep(100 * time.Millisecond)

	mockSession.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// Additional test: error in unmarshaling
func TestConsumeClaim_InvalidJSON(t *testing.T) {
	mockService := new(mockServerKafkaService)
	handler := handler.NewServerConsumerHandler(mockService)

	mockSession := new(mockConsumerGroupSession)
	mockClaim := &mockConsumerGroupClaim{
		messages: make(chan *sarama.ConsumerMessage, 1),
	}

	message := &sarama.ConsumerMessage{
		Value: []byte("invalid-json"),
	}

	mockSession.On("MarkMessage", message, "").Return()

	mockClaim.messages <- message
	close(mockClaim.messages)

	err := handler.ConsumeClaim(mockSession, mockClaim)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockSession.AssertExpectations(t)
	mockService.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
}

// Additional test: service returns error
func TestConsumeClaim_UpdateStatusFails(t *testing.T) {
	mockService := new(mockServerKafkaService)
	handler := handler.NewServerConsumerHandler(mockService)

	mockSession := new(mockConsumerGroupSession)
	mockClaim := &mockConsumerGroupClaim{
		messages: make(chan *sarama.ConsumerMessage, 1),
	}

	serverMessage := map[string]string{
		"server_id": "srv123",
		"status":    "down",
	}
	msgValue, _ := json.Marshal(serverMessage)

	message := &sarama.ConsumerMessage{
		Value: msgValue,
	}

	mockSession.On("MarkMessage", message, "").Return()
	mockService.On("UpdateStatus", "srv123", "down").Return(errors.New("DB error"))

	mockClaim.messages <- message
	close(mockClaim.messages)

	err := handler.ConsumeClaim(mockSession, mockClaim)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockSession.AssertExpectations(t)
	mockService.AssertExpectations(t)
}
