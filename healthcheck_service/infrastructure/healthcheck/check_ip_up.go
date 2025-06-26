package healthcheck

import (
	"errors"
	"time"

	"github.com/prometheus-community/pro-bing"
)

func IsHostUp(ipv4 string) (bool, error) {
	pinger, err := probing.NewPinger(ipv4)
	if err != nil {
		return false, err
	}

	pinger.Count = 3
	pinger.Timeout = 5 * time.Second // set timeout to 5 seconds

	err = pinger.Run()
	if err != nil {
		return false, err
	}

	stats := pinger.Statistics()
	println("Packet recv: ", stats.PacketsRecv)
	if stats.PacketsRecv > 0 {
		return true, nil
	} else {
		return false, errors.New("Don't receive any packet from address: " + ipv4)
	}
}