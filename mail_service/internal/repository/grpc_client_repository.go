package repository

import (
	"context"
	grpcclient "mail_service/infrastructure/grpc_client"
	"mail_service/proto"

	"github.com/flashhhhh/pkg/logging"
)

type MailGRPCClientRepository interface {
	GetServersInformation(startTime, endTime string) (int, int, int, float64, error)
}

type mailGRPCClientRepository struct {
	mailGRPCClient grpcclient.MailGRPCClient
}

func NewMailGRPCClientRepository(mailGRPCClient grpcclient.MailGRPCClient) MailGRPCClientRepository {
	return &mailGRPCClientRepository {
		mailGRPCClient: mailGRPCClient,
	}
}

func (r *mailGRPCClientRepository) GetServersInformation(startTime, endTime string) (int, int, int, float64, error) {
	resp, err := r.mailGRPCClient.GetServersInformation(context.Background(), &proto.TimeRequest{
		StartTime: startTime,
		EndTime: endTime,
	})

	if err != nil {
		logging.LogMessage("mail_service", "Cannot get server information from Server Administration's GRPC server. Err: " + err.Error(), "ERROR")
		return 0, 0, 0, 0, err
	}

	logging.LogMessage("mail_service", "Get server information from Server Administration's GRPC server successfully!", "INFO")
	return int(resp.NumServers), int(resp.NumOnServers), int(resp.NumOffServers), resp.MeanUpTimeRatio, nil
}