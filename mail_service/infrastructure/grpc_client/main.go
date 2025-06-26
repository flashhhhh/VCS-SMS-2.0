package grpcclient

import (
	"context"
	"mail_service/proto"

	"github.com/flashhhhh/pkg/env"
	"google.golang.org/grpc"
)

type MailGRPCClient interface {
	GetServersInformation(ctx context.Context, req *proto.TimeRequest) (*proto.ServersInformationResponse, error)
}

type mailGRPCClientWrapper struct {
	client proto.ServerAdministrationServiceClient
}

func (w *mailGRPCClientWrapper) GetServersInformation(ctx context.Context, req *proto.TimeRequest) (*proto.ServersInformationResponse, error) {
	return w.client.GetServersInformation(ctx, req)
}

func StartGRPCClient() (MailGRPCClient, error) {
	// Create a connection to the server.
	conn, err := grpc.Dial(env.GetEnv("GRPC_SERVER_ADMINISTRATION_SERVER", "localhost") + ":" + env.GetEnv("GRPC_SERVER_ADMINISTRATION_PORT", "50052"), grpc.WithInsecure())
	println(env.GetEnv("GRPC_SERVER_ADMINISTRATION_SERVER", "localhost") + ":" + env.GetEnv("GRPC_SERVER_ADMINISTRATION_PORT", "50052"))
	if err != nil {
		return nil, err
	}

	// Create a new client
	client := proto.NewServerAdministrationServiceClient(conn)

	return &mailGRPCClientWrapper{client: client}, nil
}