package main

import (
	"mail_service/api/middlewares"
	"mail_service/api/routes"
	grpcclient "mail_service/infrastructure/grpc_client"
	mailsending "mail_service/infrastructure/mail_sending"
	"mail_service/internal/handler"
	"mail_service/internal/repository"
	"mail_service/internal/service"
	"net/http"
	"os"
	"path/filepath"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize logger for mail_service
	currentPath, _ := os.Getwd()
	mailServiceLogPath := filepath.Join(currentPath, "logs", "mail_service.log")
	logging.InitLogger("mail_service", mailServiceLogPath, 10, 5, 30)

	// Load running environment variable
	environment := env.GetEnv("RUNNING_ENVIRONMENT", "local")
	logging.LogMessage("mail_service", "Running in "+environment+" environment", "INFO")

	// Load environment variables from the .env file
	environmentFilePath := filepath.Join(currentPath, "configs", environment+".env")
	if err := env.LoadEnv(environmentFilePath); err != nil {
		logging.LogMessage("mail_service", "Failed to load environment variables from "+environmentFilePath+": "+err.Error(), "FATAL")
		logging.LogMessage("mail_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("mail_service", "Environment variables loaded successfully from "+environmentFilePath, "INFO")
	}

	mailGRPCClient, err := grpcclient.StartGRPCClient()
	if err != nil {
		logging.LogMessage("mail_service", "Failed to connect to Server Administration's GRPC server, err: " + err.Error(), "ERROR")
		logging.LogMessage("mail_service", "Exiting ...", "FATAL")
		os.Exit(1)
	}

	mailSending := mailsending.NewMailSending(env.GetEnv("SENDER_EMAIL", ""), env.GetEnv("SENDER_PASSWORD", ""))
	mailGRPCClientRepository := repository.NewMailGRPCClientRepository(mailGRPCClient)
	mailService := service.NewMailService(mailSending, mailGRPCClientRepository)
	mailHandler := handler.NewMailHandler(mailService)

	mailServerHost := env.GetEnv("MAIL_SERVICE_HOST", "localhost")
	mailServerPort := env.GetEnv("MAIL_SERVICE_PORT", "10003")

	r := mux.NewRouter()
	routes.RegisterRoutes(r, mailHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, change this for security
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	logging.LogMessage("mail_service", "Starting HTTP server on port " + mailServerPort, "INFO")
	if err := http.ListenAndServe(mailServerHost + ":" + mailServerPort, middlewares.CorsMiddleware(corsHandler)); err != nil {
		logging.LogMessage("mail_service", "Failed to start HTTP server: "+err.Error(), "FATAL")
		logging.LogMessage("mail_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}
	logging.LogMessage("mail_service", "HTTP server stopped", "INFO")
	logging.LogMessage("mail_service", "Exiting the program...", "INFO")
	os.Exit(0)
	
	// startTime := "2025-06-24T00:00:00Z"
	// endTime := "2025-06-24T23:59:59Z"
}