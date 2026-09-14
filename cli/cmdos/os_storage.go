package cmdos

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// StorageRunnerFn is registered by package cmd to invoke the full storage suite.
var StorageRunnerFn func(args []string) error

func runOSStorage(args []string) error {
	if StorageRunnerFn != nil {
		return StorageRunnerFn(args)
	}

	return runDefaultOSStorage()
}

func runDefaultOSStorage() error {
	fmt.Println(constants.ColorCyan + "▶ OS Storage & Drive Status:" + constants.ColorReset)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Printf("  • Current Workspace: %s\n", wd)
	fmt.Printf("  • System Temp:       %s\n", os.TempDir())

	return nil
}
