package main

import (
	"os"
	"path/filepath"
	"server_administration_service/infrastructure/elasticsearch"
	"server_administration_service/infrastructure/grpc"
	"server_administration_service/infrastructure/postgres"
	"server_administration_service/internal/handler"
	"server_administration_service/internal/repository"
	"server_administration_service/internal/service"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
)

func main() {
	// Initialize logger for server_administration_service
	currentPath, _ := os.Getwd()
	serverServiceLogPath := filepath.Join(currentPath, "logs", "server_administration_service.log")
	logging.InitLogger("server_administration_service", serverServiceLogPath, 10, 5, 30)

	// Load running environment variable
	environment := env.GetEnv("RUNNING_ENVIRONMENT", "local")
	logging.LogMessage("server_administration_service", "Running in "+environment+" environment", "INFO")

	// Load environment variables from the .env file
	environmentFilePath := filepath.Join(currentPath, "configs", environment+".env")
	if err := env.LoadEnv(environmentFilePath); err != nil {
		logging.LogMessage("server_administration_service", "Failed to load environment variables from "+environmentFilePath+": "+err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("server_administration_service", "Environment variables loaded successfully from "+environmentFilePath, "INFO")
	}

	// Connect to the database
	dsn := "host=" + env.GetEnv("PG_HOST", "localhost") +
		" user=" + env.GetEnv("PG_USER", "postgres") +
		" password=" + env.GetEnv("PG_PASSWORD", "password") +
		" dbname=" + env.GetEnv("PG_NAME", "server_administration_service") +
		" port=" + env.GetEnv("PG_PORT", "5432") +
		" sslmode=disable"
	db := postgres.ConnectDB(dsn)

	// Migrate the database
	if environment == "local" {
		logging.LogMessage("server_administration_service", "Running database migrations in local environment", "INFO")
		postgres.Migrate(db)
	} else {
		logging.LogMessage("server_administration_service", "Skipping database migrations in non-local environment", "INFO")
	}

	// Initialize ES client
	esAddress := env.GetEnv("ES_HOST", "http://localhost") +
				":" + env.GetEnv("ES_PORT", "9200")
	es := elasticsearch.ConnectES(esAddress)
	esc := elasticsearch.NewElasticsearchClient(es)

	// Initialize the server
	serverGRPCRepository := repository.NewServerGRPCRepository(db)
	serverGRPCService := service.NewServerGRPCService(serverGRPCRepository)

	serverInfoRepository := repository.NewServerInfoRepository(db, esc)
	serverInfoService := service.NewServerInfoService(serverInfoRepository)
	serverGRPCHandler := handler.NewServerGRPCHandler(serverGRPCService, serverInfoService)

	serverGRPCPort := env.GetEnv("SERVER_ADMINISTRATION_GPRC_PORT", "50051")
	logging.LogMessage("server_administration_service", "Starting gRPC server on port " + serverGRPCPort, "INFO")
	grpc.StartGRPCServer(serverGRPCHandler, serverGRPCPort)
}