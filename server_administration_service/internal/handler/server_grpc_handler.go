package handler

import (
	"context"
	"server_administration_service/internal/service"
	"server_administration_service/proto"
	"strconv"

	"github.com/flashhhhh/pkg/logging"
)

type ServerGRPCHandler struct {
	serverGRPCService service.ServerGRPCService
	serverInfoService service.ServerInfoService
	proto.UnimplementedServerAdministrationServiceServer
}

func NewServerGRPCHandler(serverGRPCService service.ServerGRPCService, serverInfoService service.ServerInfoService) *ServerGRPCHandler {
	return &ServerGRPCHandler{
		serverGRPCService: serverGRPCService,
		serverInfoService: serverInfoService,
	}
}

func (h *ServerGRPCHandler) GetAddressAndStatus(ctx context.Context, req *proto.EmptyRequest) (*proto.IDAddressAndStatusList, error) {
	logging.LogMessage("server_administration_service", "Get Address and current status list of all servers", "INFO")

	serverAddresses, err := h.serverGRPCService.GetServerAddresses()
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to get address and current status list, err: " + err.Error(), "INFO")
		return nil, err
	}

	logging.LogMessage("server_administration_service", "Addresses were retrieved successfully", "INFO")

	idAddressAndStatusList := &proto.IDAddressAndStatusList{}
	for _, serverAddress := range serverAddresses {
		idAddressAndStatusList.ServerList = append(idAddressAndStatusList.ServerList, &proto.IDAddressAndStatus{
			Id: int64(serverAddress.ID),
			Address: serverAddress.IPv4,
			Status: serverAddress.Status,
		})
	}

	return idAddressAndStatusList, nil
}

func (h *ServerGRPCHandler) UpdateStatus(ctx context.Context, req *proto.ServerStatusList) (*proto.EmptyResponse, error) {
	for _, serverStatus := range req.StatusList {
		id := int(serverStatus.Id)
		status := serverStatus.Status

		logging.LogMessage("server_administration_service", "Update status " + status + " for id " + strconv.Itoa(id), "INFO")

		err := h.serverGRPCService.UpdateStatus(id, status)
		if err != nil {
			logging.LogMessage("server_administration_service", "Failed to update status " + status + " for id " + strconv.Itoa(id) + ", err: " + err.Error(), "INFO")
			// return &proto.EmptyResponse{}, err
		}

		logging.LogMessage("server_administration_service", "Update status " + status + " for id " + strconv.Itoa(id) + " successfully!", "INFO")
	}

	return &proto.EmptyResponse{}, nil
}

func (h *ServerGRPCHandler) GetServersInformation(ctx context.Context, req *proto.TimeRequest) (*proto.ServersInformationResponse, error) {
	numServers, err := h.serverInfoService.GetNumServers()
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to get the number of servers, err: " + err.Error(), "ERROR")
		return &proto.ServersInformationResponse{}, err
	}

	numOnServers, err := h.serverInfoService.GetNumOnServers()
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to get the number of ON servers, err: " + err.Error(), "ERROR")
		return &proto.ServersInformationResponse{}, err
	}

	numOffServers, err := h.serverInfoService.GetNumOffServers()
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to get the number of OFF servers, err: " + err.Error(), "ERROR")
		return &proto.ServersInformationResponse{}, err
	}

	meanUpTimeRatio, err := h.serverInfoService.GetServerMeanUpTimeRatio(req.StartTime, req.EndTime)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to get the mean of uptime ratio, err: " + err.Error(), "ERROR")
		return &proto.ServersInformationResponse{}, err
	}

	return &proto.ServersInformationResponse{
		NumServers: int64(numServers),
		NumOnServers: int64(numOnServers),
		NumOffServers: int64(numOffServers),
		MeanUpTimeRatio: meanUpTimeRatio,
	}, nil
}