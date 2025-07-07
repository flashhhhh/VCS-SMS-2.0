package elasticsearch

import (
	"os"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/flashhhhh/pkg/logging"
)

func ConnectES(dsn string) (*elasticsearch.Client) {
	client, connectES_err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{dsn},
	})

	if connectES_err != nil {
		logging.LogMessage("server_administration_service", "Error creating Elasticsearch client: " + connectES_err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}

	logging.LogMessage("server_administration_service", "Connected to Elasticsearch at "+dsn, "INFO")
	return client
}