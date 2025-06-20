package main

import (
	"net/http"
	"os"
	"path/filepath"
	"server_administration_service/api/routes"
	"server_administration_service/infrastructure/postgres"
	"server_administration_service/infrastructure/redis"
	"server_administration_service/internal/handler"
	"server_administration_service/internal/repository"
	"server_administration_service/internal/service"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	// Initialize Redis client
	redisAddress := env.GetEnv("REDIS_HOST", "localhost") + 
				":" + env.GetEnv("REDIS_PORT", "6379")
	redis := redis.NewRedisClient(redisAddress)

	// Initialize the server
	serverRepository := repository.NewServerCRUDRepository(db, redis)
	serverService := service.NewServerCRUDService(serverRepository)
	serverHandler := handler.NewServerRestHandler(serverService)

	// Initialize the HTTP server
	serverPort := env.GetEnv("SERVER_ADMINISTRATION_PORT", "10002")
	
	r := mux.NewRouter()
	routes.RegisterRoutes(r, serverHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, change this for security
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	logging.LogMessage("server_administration_service", "Starting server on port "+serverPort, "INFO")
	if err := http.ListenAndServe(":"+serverPort, corsHandler); err != nil {
		logging.LogMessage("server_administration_service", "Failed to start server: "+err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}
	logging.LogMessage("user_service", "HTTP server stopped", "INFO")
	logging.LogMessage("user_service", "Exiting the program...", "INFO")
	os.Exit(0)
}