package handler

import (
	"encoding/json"
	"server_administration_service/internal/service"

	"github.com/IBM/sarama"
	"github.com/flashhhhh/pkg/logging"
)

type ServerConsumerHandler struct {
	serverKafkaService service.ServerKafkaService
}

func NewServerConsumerHandler(serverKafkaService service.ServerKafkaService) *ServerConsumerHandler {
	return &ServerConsumerHandler{
		serverKafkaService: serverKafkaService,
	}
}

func (h ServerConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h ServerConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h ServerConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Create a semaphore channel to limit concurrent goroutines
	const maxWorkers = 10
	semaphore := make(chan struct{}, maxWorkers)
	
	for message := range claim.Messages() {
		// Acquire semaphore
		semaphore <- struct{}{}
		
		go func(message *sarama.ConsumerMessage) {
			// Release semaphore when done
			defer func() { <-semaphore }()
			
			logging.LogMessage("server_administration_service", "Received message: " + string(message.Value), "INFO")
			session.MarkMessage(message, "")

			// Parse the message
			var serverMessage struct {
				ServerID 	string	`json:"server_id"`
				Status   	string	`json:"status"`
			}
			
			if err := json.Unmarshal(message.Value, &serverMessage); err != nil {
				logging.LogMessage("server_administration_service", "Error parsing message: "+err.Error(), "ERROR")
				return
			}
			
			// Now you can use the parsed message
			logging.LogMessage("server_administration_service", "Updating server status for server_id: " + serverMessage.ServerID, "INFO")

			err := h.serverKafkaService.UpdateStatus(serverMessage.ServerID, serverMessage.Status)
			if err != nil {
				logging.LogMessage("server_administration_service", "Failed to update status: " + serverMessage.Status + 
																	" for server id: " + serverMessage.ServerID + 
																	" , err: " + err.Error(), "ERROR")
			}

			logging.LogMessage("server_administration_service", "Message processed: " + string(message.Value), "INFO")
		}(message)
	}

	return nil
}