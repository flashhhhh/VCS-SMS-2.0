package redis

import (
	"context"
	"os"

	"github.com/flashhhhh/pkg/logging"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) *redis.Client {
	logging.LogMessage("server_administration_service", "Connecting to Redis at "+addr, "INFO")
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		logging.LogMessage("server_administration_service", "Failed to connect to Redis: "+err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}

	logging.LogMessage("server_administration_service", "Connected to Redis successfully", "INFO")
	return client
}