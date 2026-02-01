package health

import (
	"fmt"
	"os"
)

// CheckReadiness checks if the service is ready by reading the status file
func CheckReadiness(statusFile string) error {
	if _, err := os.Stat(statusFile); os.IsNotExist(err) {
		return fmt.Errorf("service not ready")
	}
	return nil
}
