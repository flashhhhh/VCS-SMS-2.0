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
	statusChanged := make(chan *proto.ServerStatus, len(serverAddresses.ServerList))

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

			newStatusBool, err := IsHostUp(address)
			if err != nil {
				logging.LogMessage("healthcheck_service", "Pinging to address " + address + " makes error: " + err.Error(), "DEBUG")
			}

			newStatus := "Off"
			if newStatusBool {
				newStatus = "On"
			}

			if status != newStatus {
				statusChanged <- &proto.ServerStatus{
					Id:     id,
					Status: newStatus,
				}
			}
		}(serverAddress)
	}

	go func() {
		wg.Wait()
		close(statusChanged)
	}()

	for s := range statusChanged {
		mu.Lock()
		serverStatusList = append(serverStatusList, s)
		mu.Unlock()
	}

	wg.Wait()

	// If you need to return []proto.ServerStatus, convert pointers to values here (without copying the mutex)
	statusList := make([]proto.ServerStatus, len(serverStatusList))
	for i, s := range serverStatusList {
		statusList[i] = proto.ServerStatus{
			Id:     s.Id,
			Status: s.Status,
		}
	}
	return proto.ServerStatusList{StatusList: serverStatusList}, nil
}