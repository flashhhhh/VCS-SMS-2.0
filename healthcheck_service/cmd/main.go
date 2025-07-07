package main

import (
	"context"
	grpcclient "healthcheck_service/infrastructure/grpc_client"
	"healthcheck_service/infrastructure/healthcheck"
	"healthcheck_service/proto"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
)

func main() {
	// Initialize logger for healthcheck_service
	currentPath, _ := os.Getwd()
	serverServiceLogPath := filepath.Join(currentPath, "logs", "healthcheck_service.log")
	logging.InitLogger("healthcheck_service", serverServiceLogPath, 10, 5, 30)

	// Load running environment variable
	environment := env.GetEnv("RUNNING_ENVIRONMENT", "local")
	logging.LogMessage("healthcheck_service", "Running in "+environment+" environment", "INFO")

	// Load environment variables from the .env file
	environmentFilePath := filepath.Join(currentPath, "configs", environment+".env")
	if err := env.LoadEnv(environmentFilePath); err != nil {
		logging.LogMessage("healthcheck_service", "Failed to load environment variables from "+environmentFilePath+": "+err.Error(), "FATAL")
		logging.LogMessage("healthcheck_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("healthcheck_service", "Environment variables loaded successfully from "+environmentFilePath, "INFO")
	}

	serverAdministrationGRPCClient, err := grpcclient.StartGRPCClient()
	if err != nil {
		logging.LogMessage("healthcheck_service", "Failed to connect to Server Administration's GRPC server, err: " + err.Error(), "ERROR")
		logging.LogMessage("healthcheck_service", "Exiting ...", "FATAL")
		os.Exit(1)
	}

	healthcheckPeriodStr := env.GetEnv("HEALTHCHECK_PERIOD", "0")
	healthcheckPeriod, _ := strconv.Atoi(healthcheckPeriodStr)

	for {
		logging.LogMessage("healthcheck_service", "Get all addresses of all servers", "INFO")

		serverAddressesList, err := serverAdministrationGRPCClient.GetAddressAndStatus(context.Background(), &proto.EmptyRequest{})
		if err != nil {
			logging.LogMessage("healthcheck_service", "Failed to receive addresses and status of all servers, err: " + err.Error(), "ERROR")
		} else {
			serverStatusList, err := healthcheck.CheckAllServers(serverAddressesList)
			if err != nil {
				logging.LogMessage("healthcheck_service", "Failed to healthcheck all servers, err: " + err.Error(), "ERROR")
			}

			serverAdministrationGRPCClient.UpdateStatus(context.Background(), &serverStatusList)
		}

		if healthcheckPeriod == 0 {
			return
		}

		time.Sleep(time.Duration(healthcheckPeriod))
	}
}