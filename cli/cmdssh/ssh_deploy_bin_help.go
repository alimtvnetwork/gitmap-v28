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
	fmt.Printf("  %sTargets (<ip, alias, seq, id>):%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    all                         Deploy to all online registered fleet nodes (default)")
	fmt.Println("    <alias>                     Deploy by node alias (e.g. w1, w2, w3, alpha-win)")
	fmt.Println("    <ip>                        Deploy by node IP address (e.g. 192.168.1.3)")
	fmt.Println("    <seq>                       Deploy by sequence number from 'ssh ls' (e.g. 1, 2, 4)")
	fmt.Println("    <id>                        Deploy by database host ID (e.g. host-1, h-node-2)")
	fmt.Println("    <t1,t2,...>                 Deploy to multiple targets (e.g. w1,w2 or 1,2)")
	fmt.Println()
	fmt.Printf("  %sFlags:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    --file, -f <path>           Local binary path to deploy (default: active binary)")
	fmt.Println("    --dry-run                   Preview deployment steps without modifying files")
	fmt.Println("    --timeout, -t <duration>    Per-node transfer timeout (default: 60s)")
	fmt.Println("    --help, -h                  Show this help menu")
	fmt.Println()
	fmt.Printf("  %sExamples:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    gitmap deploy-bin w1                        # By alias: deploy to node w1")
	fmt.Println("    gitmap deploy-bin 192.168.1.3               # By IP: deploy to 192.168.1.3")
	fmt.Println("    gitmap deploy-bin 1                         # By sequence: deploy to 1st node in 'ssh ls'")
	fmt.Println("    gitmap deploy-bin 4                         # By sequence: deploy to 4th node in 'ssh ls'")
	fmt.Println("    gitmap deploy-bin host-1                    # By host ID: deploy to database host ID")
	fmt.Println("    gitmap deploy-bin w1,w2                     # Multiple targets: deploy to w1 and w2")
	fmt.Println("    gitmap deploy-bin all                       # Broadcast to all fleet machines")
	fmt.Println("    gitmap deploy-bin                           # Default: deploy current binary to all nodes")
	fmt.Println("    gitmap deploy-bin 1 --dry-run               # Preview deployment to sequence #1")
	fmt.Println("    gitmap deploy-bin w1 --file ./gitmap.exe    # Deploy specific local build artifact to w1")
	fmt.Println()
}
