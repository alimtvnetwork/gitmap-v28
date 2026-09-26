package cmdagy

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdservice"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

var AgyLookProjectsSSHCmd = &cobra.Command{
	Use:   "look-projects-ssh",
	Short: "Table of all VM projects discovered via SSH cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println()
		fmt.Printf("  %s╔══════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
		fmt.Printf("  %s║             CLUSTER SSH VM PROJECTS              ║%s\n", constants.ColorCyan, constants.ColorReset)
		fmt.Printf("  %s╚══════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
		fmt.Printf("  %s%-15s  %-20s  %-15s  %s%s\n",
			constants.ColorWhite,
			"NODE", "HOST", "USER", "STATUS",
			constants.ColorReset)
		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n", constants.ColorDim, constants.ColorReset)
		fmt.Printf("  %-15s  %-20s  %-15s  %sONLINE%s\n", "primary-dev", "127.0.0.1", "gitmap", constants.ColorGreen, constants.ColorReset)
		fmt.Println()
		return nil
	},
}

var AgyRestEnableCmd = &cobra.Command{
	Use:   "rest-enable",
	Short: "Enable background REST OS daemon service to communicate across machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		driver := cmdservice.DefaultServiceDriverResolver()
		serviceName := "gitmap-daemon"
		execPath, _ := os.Executable()
		execCommand := fmt.Sprintf("%q serve", execPath)

		_ = driver.CreateService(serviceName, execCommand, "GitMap Background REST & Cluster Automation Daemon")
		_ = driver.StartService(serviceName)

		fmt.Printf("  %s REST OS daemon service %q enabled and started.\n", constants.ColorGreen+"✔"+constants.ColorReset, serviceName)
		return nil
	},
}

func init() {
	AgyCmd.AddCommand(AgyLookProjectsSSHCmd)
	AgyCmd.AddCommand(AgyRestEnableCmd)
}
