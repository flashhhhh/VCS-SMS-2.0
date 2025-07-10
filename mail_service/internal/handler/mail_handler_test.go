package handler_test

import (
	"errors"
	"mail_service/internal/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
)

// mockMailService implements service.MailService for testing
type mockMailService struct {
	mock.Mock
}

func (m *mockMailService) SendServersReportEmail(to, startTime, endTime string) error {
	args := m.Called(to, startTime, endTime)
	return args.Error(0)
}

func TestSendEmail_MissingToParameter(t *testing.T) {
	mockSvc := new(mockMailService)
	h := handler.NewMailHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/send-email?start_time=2024-01-01T00:00:00Z&end_time=2024-01-02T00:00:00Z", nil)
	rr := httptest.NewRecorder()

	h.SendEmail(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Missing 'to' parameter") {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
}

func TestSendEmail_MissingStartTimeParameter(t *testing.T) {
	mockSvc := new(mockMailService)
	h := handler.NewMailHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/send-email?to=test@example.com&end_time=2024-01-02T00:00:00Z", nil)
	rr := httptest.NewRecorder()

	h.SendEmail(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Missing 'start_time' parameter") {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
}

func TestSendEmail_MissingEndTimeParameter(t *testing.T) {
	mockSvc := new(mockMailService)
	h := handler.NewMailHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/send-email?to=test@example.com&start_time=2024-01-01T00:00:00Z", nil)
	rr := httptest.NewRecorder()

	h.SendEmail(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Missing 'end_time' parameter") {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
}

func TestSendEmail_SendServersReportEmailFails(t *testing.T) {
	mockSvc := new(mockMailService)
	h := handler.NewMailHandler(mockSvc)

	mockSvc.On("SendServersReportEmail", "test@example.com", "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z").
		Return(errors.New("smtp error"))

	req := httptest.NewRequest(http.MethodGet, "/send-email?to=test@example.com&start_time=2024-01-01T00:00:00Z&end_time=2024-01-02T00:00:00Z", nil)
	rr := httptest.NewRecorder()

	h.SendEmail(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Failed to send emails: smtp error") {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
	mockSvc.AssertExpectations(t)
}

func TestSendEmail_Success(t *testing.T) {
	mockSvc := new(mockMailService)
	h := handler.NewMailHandler(mockSvc)

	mockSvc.On("SendServersReportEmail", "test@example.com", "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z").
		Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/send-email?to=test@example.com&start_time=2024-01-01T00:00:00Z&end_time=2024-01-02T00:00:00Z", nil)
	rr := httptest.NewRecorder()

	h.SendEmail(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Emails sent successfully") {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
	mockSvc.AssertExpectations(t)
}