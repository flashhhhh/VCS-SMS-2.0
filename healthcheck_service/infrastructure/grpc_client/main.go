package grpcclient

import (
	"context"
	"healthcheck_service/proto"

	"github.com/flashhhhh/pkg/env"
	"google.golang.org/grpc"
)

type ServerAdministrationGRPCClient interface {
	GetAddressAndStatus(ctx context.Context, req *proto.EmptyRequest) (*proto.IDAddressAndStatusList, error)
	UpdateStatus(ctx context.Context, req *proto.ServerStatusList) (*proto.EmptyResponse, error)
}

type serverAdministrationGRPCClientWrapper struct {
	client proto.ServerAdministrationServiceClient
}

func (w *serverAdministrationGRPCClientWrapper) GetAddressAndStatus(ctx context.Context, req *proto.EmptyRequest) (*proto.IDAddressAndStatusList, error) {
	return w.client.GetAddressAndStatus(ctx, req)
}

func (w *serverAdministrationGRPCClientWrapper) UpdateStatus(ctx context.Context, req *proto.ServerStatusList) (*proto.EmptyResponse, error) {
	return w.client.UpdateStatus(ctx, req)
}

func StartGRPCClient() (ServerAdministrationGRPCClient, error) {
	// Create a connection to the server.
	conn, err := grpc.Dial(env.GetEnv("GRPC_SERVER_ADMINISTRATION_SERVER", "localhost") + ":" + env.GetEnv("GRPC_SERVER_ADMINISTRATION_PORT", "50052"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// Create a new client
	client := proto.NewServerAdministrationServiceClient(conn)

	return &serverAdministrationGRPCClientWrapper{client: client}, nil
}