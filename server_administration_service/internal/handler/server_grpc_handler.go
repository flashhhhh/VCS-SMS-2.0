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
	proto.UnimplementedServerAdministrationServiceServer
}

func NewServerGRPCHandler(serverGRPCService service.ServerGRPCService) *ServerGRPCHandler {
	return &ServerGRPCHandler{
		serverGRPCService: serverGRPCService,
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