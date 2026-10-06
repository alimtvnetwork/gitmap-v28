//go:build !windows

package cmdpurge

import (
	"os"
)

// SendToRecycleBin moves a file or directory to the platform trash or removes it.
func SendToRecycleBin(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(path)
}
