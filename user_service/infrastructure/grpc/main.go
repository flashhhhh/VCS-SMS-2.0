package grpc

import (
	"log"
	"net"
	"user_service/internal/handler"
	"user_service/pb"

	"google.golang.org/grpc"
)

func StartGRPCServer(userHandler *handler.GrpcHandler, port string) {
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	log.Println("gRPC server is running on port: ", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}