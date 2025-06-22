package healthcheck

import (
	"errors"

	"github.com/prometheus-community/pro-bing"
)

func IsHostUp(ipv4 string) (bool, error) {
	pinger, err := probing.NewPinger(ipv4)
	if err != nil {
		return false, err
	}

	pinger.Count = 3
	err = pinger.Run()

	if err != nil {
		return false, err
	}

	stats := pinger.Statistics()
	if (stats.PacketsRecv > 0) {
		return true, nil
	} else {
		return false, errors.New("Don't receive any packet from address: " + ipv4)
	}
}