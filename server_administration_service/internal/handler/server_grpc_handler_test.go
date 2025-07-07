package handler_test

import (
	"context"
	"errors"
	"testing"

	"server_administration_service/internal/dto"
	"server_administration_service/internal/handler"
	"server_administration_service/proto"

	"github.com/stretchr/testify/mock"
)

// Mock implementations for service interfaces

type mockServerGRPCService struct {
	mock.Mock
}

func (m *mockServerGRPCService) GetServerAddresses() ([]dto.ServerAddress, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ServerAddress), args.Error(1)
}
func (m *mockServerGRPCService) UpdateStatus(serverID, status string) error {
	args := m.Called(serverID, status)
	return args.Error(0)
}

type mockServerInfoService struct {
	mock.Mock
}

func (m *mockServerInfoService) GetNumServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}
func (m *mockServerInfoService) GetNumOnServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}
func (m *mockServerInfoService) GetNumOffServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}
func (m *mockServerInfoService) GetServerMeanUpTimeRatio(startTime, endTime string) (float64, error) {
	args := m.Called(startTime, endTime)
	return args.Get(0).(float64), args.Error(1)
}

func TestGetAddressAndStatus_Success(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	addresses := []dto.ServerAddress{
		{ServerID: "1", IPv4: "10.0.0.1", Status: "On"},
		{ServerID: "2", IPv4: "10.0.0.2", Status: "Off"},
	}
	mockGRPC.On("GetServerAddresses").Return(addresses, nil)

	resp, err := handler.GetAddressAndStatus(context.Background(), &proto.EmptyRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.ServerList) != 2 {
		t.Errorf("expected 2 servers, got %d", len(resp.ServerList))
	}
	if resp.ServerList[0].ServerId != "1" || resp.ServerList[1].ServerId != "2" {
		t.Errorf("unexpected server IDs: %+v", resp.ServerList)
	}
}

func TestGetAddressAndStatus_Error(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	mockGRPC.On("GetServerAddresses").Return(nil, errors.New("db error"))

	resp, err := handler.GetAddressAndStatus(context.Background(), &proto.EmptyRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if resp != nil {
		t.Errorf("expected nil response, got %+v", resp)
	}
}

func TestUpdateStatus_AllSuccess(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	statusList := []*proto.ServerStatus{
		{ServerId: "1", Status: "On"},
		{ServerId: "2", Status: "Off"},
	}
	mockGRPC.On("UpdateStatus", mock.Anything, mock.Anything).Return(nil)

	resp, err := handler.UpdateStatus(context.Background(), &proto.ServerStatusList{StatusList: statusList})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

func TestUpdateStatus_WithError(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	statusList := []*proto.ServerStatus{
		{ServerId: "1", Status: "On"},
		{ServerId: "2", Status: "Off"},
	}
	// First call returns nil, second returns error
	mockGRPC.On("UpdateStatus", "1", "On").Return(nil)
	mockGRPC.On("UpdateStatus", "2", "Off").Return(errors.New("update error"))

	resp, err := handler.UpdateStatus(context.Background(), &proto.ServerStatusList{StatusList: statusList})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

func TestGetServersInformation_Success(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	mockInfo.On("GetNumServers").Return(5, nil)
	mockInfo.On("GetNumOnServers").Return(3, nil)
	mockInfo.On("GetNumOffServers").Return(2, nil)
	mockInfo.On("GetServerMeanUpTimeRatio", "2025-06-24T00:00:00Z", "2025-06-24T23:59:59Z").Return(0.75, nil)

	req := &proto.TimeRequest{StartTime: "2025-06-24T00:00:00Z", EndTime: "2025-06-24T23:59:59Z"}
	resp, err := handler.GetServersInformation(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.NumServers != 5 || resp.NumOnServers != 3 || resp.NumOffServers != 2 || resp.MeanUpTimeRatio != 0.75 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetServersInformation_NumServersError(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	mockInfo.On("GetNumServers").Return(0, errors.New("fail"))
	req := &proto.TimeRequest{StartTime: "2025-06-24T00:00:00Z", EndTime: "2025-06-24T23:59:59Z"}
	resp, err := handler.GetServersInformation(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

func TestGetServersInformation_NumOnServersError(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	mockInfo.On("GetNumServers").Return(5, nil)
	mockInfo.On("GetNumOnServers").Return(0, errors.New("fail"))
	req := &proto.TimeRequest{StartTime: "2025-06-24T00:00:00Z", EndTime: "2025-06-24T23:59:59Z"}
	resp, err := handler.GetServersInformation(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

func TestGetServersInformation_NumOffServersError(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	mockInfo.On("GetNumServers").Return(5, nil)
	mockInfo.On("GetNumOnServers").Return(3, nil)
	mockInfo.On("GetNumOffServers").Return(0, errors.New("fail"))
	req := &proto.TimeRequest{StartTime: "2025-06-24T00:00:00Z", EndTime: "2025-06-24T23:59:59Z"}
	resp, err := handler.GetServersInformation(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

func TestGetServersInformation_MeanUpTimeRatioError(t *testing.T) {
	mockGRPC := new(mockServerGRPCService)
	mockInfo := new(mockServerInfoService)
	handler := handler.NewServerGRPCHandler(mockGRPC, mockInfo)

	mockInfo.On("GetNumServers").Return(5, nil)
	mockInfo.On("GetNumOnServers").Return(3, nil)
	mockInfo.On("GetNumOffServers").Return(2, nil)
	mockInfo.On("GetServerMeanUpTimeRatio", "2025-06-24T00:00:00Z", "2025-06-24T23:59:59Z").Return(0.0, errors.New("fail"))
	req := &proto.TimeRequest{StartTime: "2025-06-24T00:00:00Z", EndTime: "2025-06-24T23:59:59Z"}
	resp, err := handler.GetServersInformation(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}