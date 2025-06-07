// Package system contains functions that works with OS.
package system

import (
	"fmt"
	"os/exec"
)

// RestartService restarts linux service.
func RestartService(serviceName string) error {
	cmd := exec.Command("systemctl", "restart", serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restart %s service: %w", serviceName, err)
	}
	return nil
}
