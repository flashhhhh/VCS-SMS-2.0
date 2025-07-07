package service

import (
	"bytes"
	"errors"
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/repository"
	"strconv"
	"strings"

	"github.com/flashhhhh/pkg/logging"
	"github.com/xuri/excelize/v2"
)

type ServerCRUDService interface {
	CreateServer(server_id, server_name, ipv4 string) (string, error)
	ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error)
	UpdateServer(server_id string, updatedData map[string]interface{}) error
	DeleteServer(server_id string) error
	ImportServers(buf []byte) ([]domain.Server, []domain.Server, error)
	ExportServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]byte, error)
}

type serverCRUDService struct {
	serverCRUDRepository repository.ServerCRUDRepository
}

func NewServerCRUDService(serverCRUDRepository repository.ServerCRUDRepository) ServerCRUDService {
	return &serverCRUDService{
		serverCRUDRepository: serverCRUDRepository,
	}
}

func (s *serverCRUDService) CreateServer(server_id, server_name, ipv4 string) (string, error) {
	server := &domain.Server{
		ServerID:   server_id,
		ServerName: server_name,
		Status: "Off",
		IPv4:  ipv4,
	}

	id, err := s.serverCRUDRepository.CreateServer(server)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *serverCRUDService) ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error) {
	servers, err := s.serverCRUDRepository.ViewServers(serverFilter, from, to, sortedColumn, order)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (s *serverCRUDService) UpdateServer(server_id string, updatedData map[string]interface{}) error {
	err := s.serverCRUDRepository.UpdateServer(server_id, updatedData)
	return err
}

func (s *serverCRUDService) DeleteServer(server_id string) error {
	err := s.serverCRUDRepository.DeleteServer(server_id)
	return err
}

func (s *serverCRUDService) ImportServers(buf []byte) ([]domain.Server, []domain.Server, error) {
	f, err := excelize.OpenReader(strings.NewReader(string(buf)))
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to open Excel file: "+err.Error(), "ERROR")
		return nil, nil, err
	}

	rows, err := f.GetRows("Servers")
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to get rows from Excel file: "+err.Error(), "ERROR")
		return nil, nil, err
	}

	if len(rows) < 2 {
		logging.LogMessage("server_administration_service", "Excel files doesn't have any row data", "ERROR")
		return nil, nil, errors.New("Excel files doesn't have any row data")
	}

	servers := make([]domain.Server, 0)

	serverID_col := -1
	serverName_col := -1
	ipv4_col := -1

	for id, val := range rows[0] {
		if val == "Server ID" {
			serverID_col = id
		} else if val == "Server Name" {
			serverName_col = id
		} else if val == "IPv4" {
			ipv4_col = id
		}
	}

	if serverID_col == -1 || serverName_col == -1 || ipv4_col == -1 {
		logging.LogMessage("server_administration_service", "Servers file doesn't contain enough information for importing", "INFO")
		return nil, nil, errors.New("Failed to import servers: Not enough information")
	}

	for _, row := range rows[1:] {
		serverID := row[serverID_col]
		serverName := row[serverName_col]
		ipv4 := row[ipv4_col]

		server := domain.Server{
			ServerID:   serverID,
			ServerName: serverName,
			Status: "Off",
			IPv4:       ipv4,
		}

		servers = append(servers, server)
	}

	insertedServers, nonInsertedServers, err := s.serverCRUDRepository.CreateServers(servers)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to import servers: "+err.Error(), "ERROR")
		return nil, nil, err
	}
	
	logging.LogMessage("server_administration_service", "Servers imported successfully", "INFO")
	return insertedServers, nonInsertedServers, nil	
}

func (s *serverCRUDService) ExportServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]byte, error) {
	servers, err := s.serverCRUDRepository.ViewServers(serverFilter, from, to, sortedColumn, order)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to export servers: "+err.Error(), "ERROR")
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Servers"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Server ID", "Server Name", "Status", "IPv4"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	for i, server := range servers {
		f.SetCellValue(sheet, "A"+strconv.Itoa(i+2), server.ServerID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(i+2), server.ServerName)
		f.SetCellValue(sheet, "C"+strconv.Itoa(i+2), server.Status)
		f.SetCellValue(sheet, "D"+strconv.Itoa(i+2), server.IPv4)
	}

	var buf bytes.Buffer
	_ = f.Write(&buf)

	logging.LogMessage("server_administration_service", "Servers exported successfully", "INFO")
	return buf.Bytes(), nil
}