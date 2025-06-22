package grpc

import (
	"net"
	"os"
	"server_administration_service/internal/handler"
	"server_administration_service/proto"

	"github.com/flashhhhh/pkg/logging"
	"google.golang.org/grpc"
)

func StartGRPCServer(serverGRPCHandler *handler.ServerGRPCHandler, port string) {
	lis, err := net.Listen("tcp", ":" + port)
	if (err != nil) {
		logging.LogMessage("server_administration_service", "Failed to start GRPC server", "ERROR")
		logging.LogMessage("server_administration_service", "Shutting down ...", "FATAL")
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterServerAdministrationServiceServer(grpcServer, serverGRPCHandler)

	logging.LogMessage("server_administration_service", "gRPC server is running on port: "+port, "INFO")
	if err := grpcServer.Serve(lis); err != nil {
		logging.LogMessage("server_administration_service", "Failed to serve: "+err.Error(), "ERROR")
	}
}