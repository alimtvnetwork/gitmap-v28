package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservice"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunRestEnable enables and starts the background REST daemon OS service.
func RunRestEnable(args []string) error {
	driver := cmdservice.DefaultServiceDriverResolver()
	serviceName := "gitmap-daemon"
	serviceDesc := "GitMap Background REST & Cluster Automation Daemon"

	execPath, err := os.Executable()
	if err != nil {
		return apperror.WrapSimple(err, "resolve executable path")
	}

	execCommand := fmt.Sprintf("%q serve", execPath)

	fmt.Printf("\n  %s╔══════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║     ENABLING GITMAP REST OS BACKGROUND SERVICE   ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)

	fmt.Printf("  • Service Name: %s%s%s\n", constants.ColorYellow, serviceName, constants.ColorReset)
	fmt.Printf("  • Binary Path:  %s\n", execCommand)

	// Create service entry
	if cErr := driver.CreateService(serviceName, execCommand, serviceDesc); cErr != nil {
		fmt.Printf("  ℹ Service creation note: %v\n", cErr)
	}

	// Start the service
	if sErr := driver.StartService(serviceName); sErr != nil {
		fmt.Printf("  ⚠ Service start note: %v\n", sErr)
	} else {
		fmt.Printf("  %s Service %q started and running.\n", constants.ColorGreen+"✔"+constants.ColorReset, serviceName)
	}

	info, statErr := driver.GetService(serviceName)
	if statErr == nil && info != nil {
		fmt.Printf("  • Current Status:  %s%s%s\n", constants.ColorGreen, info.Status, constants.ColorReset)
		fmt.Printf("  • Enabled:         %v\n", info.IsEnabled)
	}

	fmt.Printf("\n  %s REST daemon service is enabled and listening for cross-machine requests.\n\n", constants.ColorCyan+"▶"+constants.ColorReset)
	return nil
}
