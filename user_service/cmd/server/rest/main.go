package main

import (
	"net/http"
	"os"
	"path/filepath"
	"user_service/api/middlewares"
	api "user_service/api/routes"
	"user_service/infrastructure/postgres"
	"user_service/internal/handler"
	"user_service/internal/repository"
	"user_service/internal/service"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize the logger for the user_service
	currentPath, _ := os.Getwd()
	userServiceLogPath := filepath.Join(currentPath, "logs", "user_service.log")
	logging.InitLogger("user_service", userServiceLogPath, 10, 5, 30)

	// Load running environment variable
	environment := env.GetEnv("RUNNING_ENVIRONMENT", "local")
	logging.LogMessage("user_service", "Running in "+environment+" environment", "INFO")

	// Load environment variables from the .env file
	environmentFilePath := filepath.Join(currentPath, "configs", environment+".env")
	if err := env.LoadEnv(environmentFilePath); err != nil {
		logging.LogMessage("user_service", "Failed to load environment variables from "+environmentFilePath+": "+err.Error(), "FATAL")
		logging.LogMessage("user_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("user_service", "Environment variables loaded successfully from "+environmentFilePath, "INFO")
	}

	// Connect to the database
	dsn := "host=" + env.GetEnv("PG_HOST", "localhost") +
		" user=" + env.GetEnv("PG_USER", "postgres") +
		" password=" + env.GetEnv("PG_PASSWORD", "password") +
		" dbname=" + env.GetEnv("PG_NAME", "user_service") +
		" port=" + env.GetEnv("PG_PORT", "5432") +
		" sslmode=disable"
	db := postgres.ConnectDB(dsn)

	// Migrate the database
	postgres.Migrate(db)

	// Initialize internal services
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewRestHandler(userService)

	// Start the HTTP server
	user_service_port := env.GetEnv("USER_SERVICE_PORT", "10001")

	r := mux.NewRouter()
	api.RegisterRoutes(r, userHandler)
	
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, change this for security
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	logging.LogMessage("user_service", "Starting HTTP server on port "+user_service_port, "INFO")
	if err := http.ListenAndServe(":"+user_service_port, middlewares.CorsMiddleware(corsHandler)); err != nil {
		logging.LogMessage("user_service", "Failed to start HTTP server: "+err.Error(), "FATAL")
		logging.LogMessage("user_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}
	logging.LogMessage("user_service", "HTTP server stopped", "INFO")
	logging.LogMessage("user_service", "Exiting the program...", "INFO")
	os.Exit(0)
}