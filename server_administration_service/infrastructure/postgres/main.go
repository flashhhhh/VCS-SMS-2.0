package postgres

import (
	"os"
	"server_administration_service/internal/domain"

	"github.com/flashhhhh/pkg/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(dsn string) *gorm.DB {
	logging.LogMessage("server_administration_service", "Connecting to the database...", "INFO")
	logging.LogMessage("server_administration_service", "Database connection string: "+dsn, "DEBUG")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to connect to the database: "+err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}

	logging.LogMessage("server_administration_service", "Connected to the database successfully", "INFO")
	return db
}

func Migrate(db *gorm.DB) {
	logging.LogMessage("server_administration_service", "Migrating the database...", "INFO")

	// Check if the table exists
	tableExists := db.Migrator().HasTable(&domain.Server{})
	if !tableExists {
		logging.LogMessage("server_administration_service", "Tables don't exist, migrating...", "INFO")
		
		err := db.AutoMigrate(&domain.Server{})
		if err != nil {
			logging.LogMessage("server_administration_service", "Failed to migrate the database: "+err.Error(), "FATAL")
			logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
			os.Exit(1)
		}
		
		logging.LogMessage("server_administration_service", "Database migrated successfully", "INFO")
	} else {
		logging.LogMessage("server_administration_service", "Tables already exist, skipping migration", "INFO")
	}
}