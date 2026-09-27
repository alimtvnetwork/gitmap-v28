// Package cmdssh — ssh_deploy_bin_help.go renders UI help menus for binary deployment.
package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RenderDeployBinHelp prints the formatted terminal help guide for deploy-bin.
func RenderDeployBinHelp() {
	fmt.Printf("\n%s%s🚀 GitMap Fleet Binary Deployment (gitmap ssh deploy-bin)%s\n\n", constants.ColorCyan, constants.ColorBold, constants.ColorReset)
	fmt.Println("  Autonomous one-liner deployment of GitMap binary across remote SSH fleet nodes.")
	fmt.Println("  Eliminates manual scp commands, handles OS-specific install paths, and verifies version.")
	fmt.Println()
	fmt.Printf("  %sUsage:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    gitmap ssh deploy-bin [target] [flags]")
	fmt.Println("    gitmap ssh push-bin [target] [flags]        (alias)")
	fmt.Println("    gitmap ssh sync-bin [target] [flags]        (alias)")
	fmt.Println("    gitmap deploy-bin [target] [flags]          (top-level alias)")
	fmt.Println()
	fmt.Printf("  %sTargets:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    all                         Deploy to all online registered fleet nodes (default)")
	fmt.Println("    <alias>                     Deploy to specific node alias (e.g. w1, w2, w3)")
	fmt.Println("    <ip>                        Deploy to specific node IP address (e.g. 192.168.1.3)")
	fmt.Println()
	fmt.Printf("  %sFlags:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    --file, -f <path>           Local binary path to deploy (default: active binary)")
	fmt.Println("    --dry-run                   Preview deployment steps without modifying files")
	fmt.Println("    --timeout, -t <duration>    Per-node transfer timeout (default: 60s)")
	fmt.Println("    --help, -h                  Show this help menu")
	fmt.Println()
	fmt.Printf("  %sExamples:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    gitmap ssh deploy-bin w1                    # Deploy to node w1 (replaces raw scp)")
	fmt.Println("    gitmap ssh deploy-bin all                   # Broadcast to all fleet machines")
	fmt.Println("    gitmap ssh deploy-bin w2 --file ./gitmap.exe# Deploy custom local build to w2")
	fmt.Println("    gitmap deploy-bin                           # Deploy current binary to all nodes")
	fmt.Println()
}
