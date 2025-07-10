package service_test

import (
	"fmt"
	"mail_service/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for mailsending.MailSending
type mockMailSending struct {
	mock.Mock
}

func (m *mockMailSending) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

// Mock for repository.MailGRPCClientRepository
type mockMailGRPCClientRepository struct {
	mock.Mock
}

func (m *mockMailGRPCClientRepository) GetServersInformation(startTime, endTime string) (int, int, int, float64, error) {
	args := m.Called(startTime, endTime)
	return args.Int(0), args.Int(1), args.Int(2), args.Get(3).(float64), args.Error(4)
}

func TestSendServersReportEmail_Success(t *testing.T) {
	mockMail := new(mockMailSending)
	mockRepo := new(mockMailGRPCClientRepository)

	to := "admin@example.com"
	startTime := "2024-06-01"
	endTime := "2024-06-02"
	numServers := 10
	numOnServers := 7
	numOffServers := 3
	meanUpTimeRatio := 87.5

	mockRepo.On("GetServersInformation", startTime, endTime).
		Return(numServers, numOnServers, numOffServers, meanUpTimeRatio, nil)

	expectedSubject := "Daily Server Status Report for " + time.Now().Format("2006-01-02")
	expectedBody := fmt.Sprintf(
		"Dear server administrator,\n\nThe server status is as follows:\n\nTotal servers: %d\nServers on: %d\nServers off: %d\nMean uptime rate: %.2f%%\n\nBest regards,\nYour Server Monitoring System",
		numServers, numOnServers, numOffServers, meanUpTimeRatio,
	)
	mockMail.On("SendEmail", to, expectedSubject, expectedBody).Return(nil)

	svc := service.NewMailService(mockMail, mockRepo)
	err := svc.SendServersReportEmail(to, startTime, endTime)
	assert.NoError(t, err)
	mockMail.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestSendServersReportEmail_GetServersInformationError(t *testing.T) {
	mockMail := new(mockMailSending)
	mockRepo := new(mockMailGRPCClientRepository)

	to := "admin@example.com"
	startTime := "2024-06-01"
	endTime := "2024-06-02"
	expectedErr := assert.AnError

	mockRepo.On("GetServersInformation", startTime, endTime).
		Return(0, 0, 0, 0.0, expectedErr)

	svc := service.NewMailService(mockMail, mockRepo)
	err := svc.SendServersReportEmail(to, startTime, endTime)
	assert.EqualError(t, err, expectedErr.Error())
	mockRepo.AssertExpectations(t)
}

func TestSendServersReportEmail_SendEmailError(t *testing.T) {
	mockMail := new(mockMailSending)
	mockRepo := new(mockMailGRPCClientRepository)

	to := "admin@example.com"
	startTime := "2024-06-01"
	endTime := "2024-06-02"
	numServers := 5
	numOnServers := 3
	numOffServers := 2
	meanUpTimeRatio := 75.0
	expectedErr := assert.AnError

	mockRepo.On("GetServersInformation", startTime, endTime).
		Return(numServers, numOnServers, numOffServers, meanUpTimeRatio, nil)

	expectedSubject := "Daily Server Status Report for " + time.Now().Format("2006-01-02")
	expectedBody := fmt.Sprintf(
		"Dear server administrator,\n\nThe server status is as follows:\n\nTotal servers: %d\nServers on: %d\nServers off: %d\nMean uptime rate: %.2f%%\n\nBest regards,\nYour Server Monitoring System",
		numServers, numOnServers, numOffServers, meanUpTimeRatio,
	)
	mockMail.On("SendEmail", to, expectedSubject, expectedBody).Return(expectedErr)

	svc := service.NewMailService(mockMail, mockRepo)
	err := svc.SendServersReportEmail(to, startTime, endTime)
	assert.EqualError(t, err, expectedErr.Error())
	mockMail.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}