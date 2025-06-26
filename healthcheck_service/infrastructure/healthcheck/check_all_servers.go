package healthcheck

import (
	"healthcheck/proto"
	"strconv"
	"sync"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
)

func CheckAllServers(serverAddresses *proto.IDAddressAndStatusList) (proto.ServerStatusList, error) {
	maxGoroutinesStr := env.GetEnv("MAX_GOROUTINES", "10")
	maxGoroutines, _ := strconv.Atoi(maxGoroutinesStr)

	semaphore := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	var serverStatusList []*proto.ServerStatus

	var mu sync.Mutex

	for _, serverAddress := range serverAddresses.ServerList {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(serverAddress *proto.IDAddressAndStatus) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()

			id := serverAddress.Id
			address := serverAddress.Address
			status := serverAddress.Status

			logging.LogMessage("healthcheck_service", "Pinging to address " + address, "INFO")

			newStatusBool, err := IsHostUp(address)
			if err != nil {
				logging.LogMessage("healthcheck_service", "Pinging to address "+address+" makes error: "+err.Error(), "DEBUG")
			}

			newStatus := "Off"
			if newStatusBool {
				newStatus = "On"
			}

			if status != newStatus {
				mu.Lock()
				serverStatusList = append(serverStatusList, &proto.ServerStatus{
					Id:     id,
					Status: newStatus,
				})
				mu.Unlock()
			}

			logging.LogMessage("healthcheck_service", "Pinging to address " + address + " with status " + newStatus, "INFO")
		}(serverAddress)
	}

	wg.Wait()
	return proto.ServerStatusList{StatusList: serverStatusList}, nil
}