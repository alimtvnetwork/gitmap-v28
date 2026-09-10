//go:build !windows

package cmd

import (
	"os"
)

func sendToRecycleBin(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(path)
}
