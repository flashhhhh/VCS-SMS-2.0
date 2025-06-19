package postgres

import (
	"os"
	"time"
	"user_service/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/flashhhhh/pkg/logging"
)

func ConnectDB(dsn string) (*gorm.DB) {
	logging.LogMessage("user_service", "Connecting to Postgres Database...", "INFO")
	logging.LogMessage("user_service", "DSN = " + dsn, "DEBUG")
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.LogMessage("user_service", "Failed to connect to Postgres Database: "+err.Error(), "FATAL")
		logging.LogMessage("user_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}

	logging.LogMessage("user_service", "Connected to Postgres Database", "INFO")
	return db
}

func Migrate(db *gorm.DB) {
	logging.LogMessage("user_service", "Running migrations...", "INFO")

	for {
		if db == nil {
			logging.LogMessage("user_service", "Database connection is nil, retrying migration in 10 seconds...", "ERROR")
			time.Sleep(10 * time.Second)
			continue
		}

		// Migrate the schema
		if err := db.AutoMigrate(&domain.User{}); err != nil {
			// Fatal error, exit the program
			logging.LogMessage("user_service", "Failed to run migrations: "+err.Error(), "ERROR")
			logging.LogMessage("user_service", "Exiting the program...", "FATAL")
			os.Exit(1)
		}

		logging.LogMessage("user_service", "Migrations completed successfully", "INFO")
		return
	}
}