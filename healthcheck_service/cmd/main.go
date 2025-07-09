package main

import (
	"context"
	"encoding/json"
	grpcclient "healthcheck_service/infrastructure/grpc_client"
	"healthcheck_service/infrastructure/healthcheck"
	"healthcheck_service/proto"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/kafka"
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
	} else {
		logging.LogMessage("healthcheck_service", "Successfully connect to GRPC server", "INFO")
	}

	// Initialize Kafka producer
	kafka_address := env.GetEnv("KAFKA_HOST", "kafka") + ":" + env.GetEnv("KAFKA_PORT", "9092")
	kafkaProducer, err := kafka.NewKafkaProducer([]string{kafka_address})
	if err != nil {
		logging.LogMessage("healthcheck_service", "Failed to connect to Kafka Server, err: " + err.Error(), "ERROR")
		logging.LogMessage("healthcheck_service", "Exiting ...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("healthcheck_service", "Successfully connect to Kafka address: " + kafka_address, "INFO")
	}

	kafka_topic := env.GetEnv("KAFKA_TOPIC", "healthcheck_topic")

	healthcheckPeriodStr := env.GetEnv("HEALTHCHECK_PERIOD", "0")
	healthcheckPeriod, _ := strconv.Atoi(healthcheckPeriodStr)

	for {
		logging.LogMessage("healthcheck_service", "Get all addresses of all servers", "INFO")

		serverAddressesList, err := serverAdministrationGRPCClient.GetAddressAndStatus(context.Background(), &proto.EmptyRequest{})
		if err != nil {
			logging.LogMessage("healthcheck_service", "Failed to receive addresses and status of all servers, err: " + err.Error(), "ERROR")
		} else {
			for _, serverAddress := range serverAddressesList.ServerList {
				server_id := serverAddress.ServerId
				address := serverAddress.Address
				status := serverAddress.Status

				logging.LogMessage("healthcheck_service", "Pinging server " + server_id + " at address " + address, "INFO")
				newStatusBool, err := healthcheck.IsHostUp(address)
				if err != nil {
					logging.LogMessage("healthcheck_service", "Pinging server " + server_id + " at address " + address + " has error: " + err.Error(), "ERROR")
				}

				newStatus := "Off"
				if newStatusBool {
					newStatus = "On"
				}

				logging.LogMessage("healthcheck_service", "Pinging server " + server_id + " at address " + address + " has status: " + newStatus, "INFO")

				// Send message to Kafka if newStatus != status
				if status != newStatus {
					healthcheckResult := map[string]interface{}{
						"server_id":	server_id,
						"status":		newStatus,
					}

					logging.LogMessage("healthcheck_service", "Sending server " + server_id + " at address " + address +
															" with status: " + newStatus + " to Kafka server", "INFO")

					healthcheckMessage, _ := json.Marshal(healthcheckResult)
					err := kafkaProducer.SendMessage(kafka_topic, healthcheckMessage)

					if err != nil {
						logging.LogMessage("healthcheck_service", "Failed to send server " + server_id + " at address " + address +
															" with status " + newStatus + " to Kafka server, err: " + err.Error(), "ERROR")
					} else {
						logging.LogMessage("healthcheck_service", "Sending server " + server_id + " at address " + address +
															" with status: " + newStatus + " to Kafka server successfully!", "INFO")
					}
				}
			}
		}

		if healthcheckPeriod == 0 {
			return
		}

		time.Sleep(time.Duration(healthcheckPeriod))
	}
}