package service_test

import (
	"bytes"
	"errors"
	"testing"

	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

// Mock for ServerCRUDRepository
type mockServerCRUDRepository struct {
	mock.Mock
}

func (m *mockServerCRUDRepository) CreateServer(server *domain.Server) (string, error) {
	args := m.Called(server)
	return args.String(0), args.Error(1)
}

func (m *mockServerCRUDRepository) ViewServers(filter *dto.ServerFilter, from, to int, sortedColumn, order string) ([]domain.Server, error) {
	args := m.Called(filter, from, to, sortedColumn, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Server), args.Error(1)
}

func (m *mockServerCRUDRepository) UpdateServer(server_id string, updatedData map[string]interface{}) error {
	args := m.Called(server_id, updatedData)
	return args.Error(0)
}

func (m *mockServerCRUDRepository) DeleteServer(server_id string) error {
	args := m.Called(server_id)
	return args.Error(0)
}

func (m *mockServerCRUDRepository) CreateServers(servers []domain.Server) ([]domain.Server, []domain.Server, error) {
	args := m.Called(servers)
	if args.Get(0) == nil || args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]domain.Server), args.Get(1).([]domain.Server), args.Error(2)
}

func TestCreateServer_Success(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	server := &domain.Server{
		ServerID:   "srv1",
		ServerName: "Server One",
		Status:     "Off",
		IPv4:       "192.168.1.1",
	}
	mockRepo.On("CreateServer", server).Return("srv1", nil)

	id, err := service.CreateServer("srv1", "Server One", "192.168.1.1")
	assert.NoError(t, err)
	assert.Equal(t, "srv1", id)
	mockRepo.AssertExpectations(t)
}

func TestCreateServer_Error(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	server := &domain.Server{
		ServerID:   "srv2",
		ServerName: "Server Two",
		Status:     "Off",
		IPv4:       "10.0.0.2",
	}
	mockRepo.On("CreateServer", server).Return("", errors.New("db error"))

	id, err := service.CreateServer("srv2", "Server Two", "10.0.0.2")
	assert.Error(t, err)
	assert.Empty(t, id)
	mockRepo.AssertExpectations(t)
}

func TestViewServers_Success(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	filter := &dto.ServerFilter{}
	expected := []domain.Server{
		{ServerID: "srv1", ServerName: "Server One", Status: "Off", IPv4: "192.168.1.1"},
	}
	mockRepo.On("ViewServers", filter, 0, 10, "ServerID", "asc").Return(expected, nil)

	servers, err := service.ViewServers(filter, 0, 10, "ServerID", "asc")
	assert.NoError(t, err)
	assert.Equal(t, expected, servers)
	mockRepo.AssertExpectations(t)
}

func TestViewServers_Error(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	filter := &dto.ServerFilter{}
	mockRepo.On("ViewServers", filter, 0, 10, "ServerID", "asc").Return(nil, errors.New("db error"))

	servers, err := service.ViewServers(filter, 0, 10, "ServerID", "asc")
	assert.Error(t, err)
	assert.Nil(t, servers)
	mockRepo.AssertExpectations(t)
}

func TestUpdateServer_Success(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	updatedData := map[string]interface{}{"ServerName": "Updated"}
	mockRepo.On("UpdateServer", "srv1", updatedData).Return(nil)

	err := service.UpdateServer("srv1", updatedData)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateServer_Error(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	updatedData := map[string]interface{}{"ServerName": "Updated"}
	mockRepo.On("UpdateServer", "srv1", updatedData).Return(errors.New("update error"))

	err := service.UpdateServer("srv1", updatedData)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteServer_Success(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	mockRepo.On("DeleteServer", "srv1").Return(nil)

	err := service.DeleteServer("srv1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteServer_Error(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	mockRepo.On("DeleteServer", "srv1").Return(errors.New("delete error"))

	err := service.DeleteServer("srv1")
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestImportServers_Success(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	// Prepare Excel file in memory
	buf := new(bytes.Buffer)
	f := createTestExcelFile()
	_ = f.Write(buf)

	expectedServers := []domain.Server{
		{ServerID: "srv1", ServerName: "Server One", Status: "Off", IPv4: "192.168.1.1"},
	}
	mockRepo.On("CreateServers", mock.Anything).Return(expectedServers, []domain.Server{}, nil)

	inserted, nonInserted, err := service.ImportServers(buf.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, expectedServers, inserted)
	assert.Empty(t, nonInserted)
	mockRepo.AssertExpectations(t)
}

func TestImportServers_InvalidFile(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	invalidBuf := []byte("not an excel file")
	inserted, nonInserted, err := service.ImportServers(invalidBuf)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Nil(t, nonInserted)
}

func TestImportServers_MissingSheet(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	// Excel file with missing columns
	buf := new(bytes.Buffer)
	f := createTestExcelFileMissingSheet()
	_ = f.Write(buf)

	inserted, nonInserted, err := service.ImportServers(buf.Bytes())
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Nil(t, nonInserted)
}

func TestImportServers_MissingColumns(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	// Excel file with missing columns
	buf := new(bytes.Buffer)
	f := createTestExcelFileMissingColumns()
	_ = f.Write(buf)

	inserted, nonInserted, err := service.ImportServers(buf.Bytes())
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Nil(t, nonInserted)
}

func TestImportServers_MissingRows(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	// Excel file with missing rows
	buf := new(bytes.Buffer)
	f := createTestExcelFileMissingRows()
	_ = f.Write(buf)

	inserted, nonInserted, err := service.ImportServers(buf.Bytes())
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Nil(t, nonInserted)
}

func TestImportServers_Error(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	// Excel file with missing rows
	buf := new(bytes.Buffer)
	f := createTestExcelFile()
	_ = f.Write(buf)

	mockRepo.On("CreateServers", mock.Anything).Return(nil, nil, errors.New("import error"))

	inserted, nonInserted, err := service.ImportServers(buf.Bytes())
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Nil(t, nonInserted)
}

func TestExportServers_Success(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	filter := &dto.ServerFilter{}
	servers := []domain.Server{
		{ServerID: "srv1", ServerName: "Server One", Status: "Off", IPv4: "192.168.1.1"},
	}
	mockRepo.On("ViewServers", filter, 0, 10, "ServerID", "asc").Return(servers, nil)

	data, err := service.ExportServers(filter, 0, 10, "ServerID", "asc")
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	mockRepo.AssertExpectations(t)
}

func TestExportServers_Error(t *testing.T) {
	mockRepo := new(mockServerCRUDRepository)
	service := service.NewServerCRUDService(mockRepo)

	filter := &dto.ServerFilter{}
	mockRepo.On("ViewServers", filter, 0, 10, "ServerID", "asc").Return(nil, errors.New("db error"))

	data, err := service.ExportServers(filter, 0, 10, "ServerID", "asc")
	assert.Error(t, err)
	assert.Nil(t, data)
	mockRepo.AssertExpectations(t)
}

// Helpers

func createTestExcelFile() *excelize.File {
	f := excelize.NewFile()
	sheet := "Servers"
	f.SetSheetName("Sheet1", sheet)
	f.SetCellValue(sheet, "A1", "Server ID")
	f.SetCellValue(sheet, "B1", "Server Name")
	f.SetCellValue(sheet, "C1", "IPv4")
	f.SetCellValue(sheet, "A2", "srv1")
	f.SetCellValue(sheet, "B2", "Server One")
	f.SetCellValue(sheet, "C2", "192.168.1.1")
	return f
}

func createTestExcelFileMissingSheet() *excelize.File {
	f := excelize.NewFile()
	sheet := "Server"
	f.SetSheetName("Sheet1", sheet)
	f.SetCellValue(sheet, "A1", "Server ID")
	f.SetCellValue(sheet, "B1", "IPv4")
	f.SetCellValue(sheet, "A2", "srv1")
	f.SetCellValue(sheet, "B2", "192.168.1.1")
	return f
}

func createTestExcelFileMissingColumns() *excelize.File {
	f := excelize.NewFile()
	sheet := "Servers"
	f.SetSheetName("Sheet1", sheet)
	f.SetCellValue(sheet, "A1", "Server ID")
	f.SetCellValue(sheet, "B1", "IPv4")
	f.SetCellValue(sheet, "A2", "srv1")
	f.SetCellValue(sheet, "B2", "192.168.1.1")
	return f
}

func createTestExcelFileMissingRows() *excelize.File {
	f := excelize.NewFile()
	sheet := "Servers"
	f.SetSheetName("Sheet1", sheet)
	f.SetCellValue(sheet, "A1", "Server ID")
	f.SetCellValue(sheet, "B1", "Server Name")
	f.SetCellValue(sheet, "C1", "IPv4")
	return f
}