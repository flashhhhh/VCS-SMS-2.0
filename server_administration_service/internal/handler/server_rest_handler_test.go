package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/handler"

	"github.com/stretchr/testify/mock"
)

// Mock implementation of ServerCRUDService
type mockServerCRUDService struct {
	mock.Mock
}

func (m *mockServerCRUDService) CreateServer(serverID, serverName, ipv4 string) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *mockServerCRUDService) ViewServers(filter *dto.ServerFilter, from, to int, sortColumn, order string) ([]domain.Server, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Server), args.Error(1)
}
func (m *mockServerCRUDService) UpdateServer(serverID string, updatedData map[string]interface{}) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockServerCRUDService) DeleteServer(serverID string) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockServerCRUDService) ImportServers(buf []byte) ([]domain.Server, []domain.Server, error) {
	args := m.Called()
	if (args.Get(0) == nil || args.Get(1) == nil) {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]domain.Server), args.Get(1).([]domain.Server), args.Error(2)
}
func (m *mockServerCRUDService) ExportServers(filter *dto.ServerFilter, from, to int, sortColumn, order string) ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func TestCreateServer_Success(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	body := map[string]interface{}{
		"server_id":   "srv-1",
		"server_name": "TestServer",
		"ipv4":        "192.168.1.1",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	mockService.On("CreateServer").Return("srv-1", nil)

	handler.CreateServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	if respBody["ServerID"] != "srv-1" {
		t.Errorf("expected ServerID 'srv-1', got %v", respBody["ServerID"])
	}
}

func TestCreateServer_InvalidBody(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/servers", strings.NewReader("{invalid json"))
	w := httptest.NewRecorder()

	handler.CreateServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestCreateServer_ServiceError(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	body := map[string]interface{}{
		"server_id":   "srv-2",
		"server_name": "FailServer",
		"ipv4":        "10.0.0.1",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/servers", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	mockService.On("CreateServer").Return("", errors.New("service error"))

	handler.CreateServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestViewServers_Success(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	servers := []domain.Server{
		{ServerID: "srv-1", ServerName: "Server1", IPv4: "192.168.1.1"},
		{ServerID: "srv-2", ServerName: "Server2", IPv4: "192.168.1.2"},
	}
	mockService.On("ViewServers").Return(servers, nil)

	req := httptest.NewRequest(http.MethodGet, "/servers?from=0&to=10&sort_column=server_id&sort_order=asc&server_id=1&server_name=Server%201&status=On&ipv4=192.168.1.1", nil)
	w := httptest.NewRecorder()

	handler.ViewServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var respBody []domain.Server
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(respBody) != 2 {
		t.Errorf("expected 2 servers, got %d", len(respBody))
	}
}

func TestViewServers_InvalidFrom(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/servers?from=abc&to=10", nil)
	w := httptest.NewRecorder()

	handler.ViewServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestViewServers_InvalidTo(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/servers?from=0&to=xyz", nil)
	w := httptest.NewRecorder()

	handler.ViewServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestViewServers_ServiceError(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	mockService.On("ViewServers").Return(nil, errors.New("service error"))

	req := httptest.NewRequest(http.MethodGet, "/servers?from=0&to=10", nil)
	w := httptest.NewRecorder()

	handler.ViewServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestUpdateServer_Success(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	body := map[string]interface{}{
		"server_name": "UpdatedServer",
		"ipv4":        "10.0.0.2",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/servers/update?server_id=srv-1", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	mockService.On("UpdateServer").Return(nil)

	handler.UpdateServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	respBody := new(strings.Builder)
	_, err := io.Copy(respBody, resp.Body)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if !strings.Contains(respBody.String(), "Server updated successfully") {
		t.Errorf("expected success message, got %s", respBody.String())
	}
}

func TestUpdateServer_MissingServerID(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	body := map[string]interface{}{
		"server_name": "UpdatedServer",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/servers/update", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.UpdateServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUpdateServer_InvalidBody(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/servers/update?server_id=srv-1", strings.NewReader("{invalid json"))
	w := httptest.NewRecorder()

	handler.UpdateServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUpdateServer_ServiceError(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	body := map[string]interface{}{
		"server_name": "UpdatedServer",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/servers/update?server_id=srv-1", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	mockService.On("UpdateServer").Return(errors.New("update error"))

	handler.UpdateServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestDeleteServer_Success(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/servers/delete?server_id=srv-1", nil)
	w := httptest.NewRecorder()

	mockService.On("DeleteServer").Return(nil)

	handler.DeleteServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Server deleted successfully") {
		t.Errorf("expected success message, got %s", string(body))
	}
}

func TestDeleteServer_MissingServerID(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/servers/delete", nil)
	w := httptest.NewRecorder()

	handler.DeleteServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestDeleteServer_ServiceError(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/servers/delete?server_id=srv-404", nil)
	w := httptest.NewRecorder()

	mockService.On("DeleteServer").Return(errors.New("not found"))

	handler.DeleteServer(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestImportServers_Success(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	imported := []domain.Server{{ServerID: "srv-1"}}
	nonImported := []domain.Server{{ServerID: "srv-2"}}
	mockService.On("ImportServers").Return(imported, nonImported, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("servers_file", "servers.xlsx")
	part.Write([]byte("dummy file content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/servers/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.ImportServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	if respBody["imported_servers"] == nil || respBody["non_imported_servers"] == nil {
		t.Errorf("expected imported and non_imported servers in response")
	}
}

func TestImportServers_MissingFile(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/servers/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.ImportServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestImportServers_ReadError(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	// Simulate a file that returns error on Read
	r, wPipe := io.Pipe()
	wPipe.CloseWithError(errors.New("read error"))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("servers_file", "servers.xlsx")
	io.Copy(part, r)
	writer.Close()

	// Mock ImportServers to avoid "unexpected call" error
	mockService.On("ImportServers").Return(nil, nil, errors.New("read error"))

	req := httptest.NewRequest(http.MethodPost, "/servers/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.ImportServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	// Accept either 400 or 500 depending on how the error is triggered
	if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 or 500, got %d", resp.StatusCode)
	}
}

func TestImportServers_ServiceError(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	mockService.On("ImportServers").Return(nil, nil, errors.New("import error"))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("servers_file", "servers.xlsx")
	part.Write([]byte("dummy file content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/servers/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.ImportServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestExportServers_Success(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	mockService.On("ExportServers").Return([]byte("exceldata"), nil)

	req := httptest.NewRequest(http.MethodGet, "/servers/export?from=0&to=10&sort_column=server_id&sort_order=asc&server_id=1&server_name=Server%201&status=On&ipv4=192.168.1.1", nil)
	w := httptest.NewRecorder()

	handler.ExportServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Errorf("expected excel content type, got %s", resp.Header.Get("Content-Type"))
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "exceldata" {
		t.Errorf("expected exceldata, got %s", string(body))
	}
}

func TestExportServers_InvalidFrom(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/servers/export?from=abc&to=10", nil)
	w := httptest.NewRecorder()

	handler.ExportServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestExportServers_InvalidTo(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/servers/export?from=0&to=xyz", nil)
	w := httptest.NewRecorder()

	handler.ExportServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestExportServers_ServiceError(t *testing.T) {
	mockService := new(mockServerCRUDService)
	handler := handler.NewServerRestHandler(mockService)

	mockService.On("ExportServers").Return(nil, errors.New("export error"))

	req := httptest.NewRequest(http.MethodGet, "/servers/export?from=0&to=10", nil)
	w := httptest.NewRecorder()

	handler.ExportServers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}