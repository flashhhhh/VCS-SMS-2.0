package healthcheck

import (
	"os/exec"
)

func IsHostUp(ipv4 string) (bool, error) {
	cmd := exec.Command("ping", "-c", "3", "-w", "5", ipv4)
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}