// Package cmdssh — ssh_deploy_router.go routes ssh deploy subcommands (keys, node-config).
package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunSSHDeployRouterCLI routes `gitmap ssh deploy` to keys or node-config.
func RunSSHDeployRouterCLI(args []string) error {
	if len(args) == 0 {
		printDeployHelp()
		return nil
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "config", "config-ssh":
		if len(args) > 1 && strings.ToLower(args[1]) == "ssh" {
			return RunSSHDeployConfigSSHCLI(args[2:])
		}
		return RunSSHDeployConfigSSHCLI(args[1:])
	case "ssh":
		return RunSSHDeployConfigSSHCLI(args[1:])
	case "keys", "key", "k":
		return RunSSHDeployKeysCLI(args[1:])
	case "bin", "binary", "exe", "gitmap":
		return RunSSHDeployBinCLI(args[1:])
	case "node-config", "nodeconfig", "nc", "nodes":
		return RunSSHDeployConfigSSHCLI(args[1:])
	case "help", "--help", "-h":
		printDeployHelp()
		return nil
	default:
		return routeFallbackDeploy(args)
	}
}

func routeFallbackDeploy(args []string) error {
	first := strings.ToLower(args[0])
	if first == "all" || strings.HasPrefix(first, "-") {
		return RunSSHDeployConfigSSHCLI(args)
	}
	fmt.Printf("\n  %sUnknown deploy target: %q%s\n", constants.ColorRed, args[0], constants.ColorReset)
	printDeployHelp()
	return nil
}

func printDeployHelp() {
	fmt.Printf("\n  %s🚀 GitMap SSH Fleet Deploy Commands%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    gitmap deploy config ssh [all|<target>] [--file <path>] [--except <id,ip,alias>]")
	fmt.Println("        Deploy node topology, IP addresses, aliases, and credentials across fleet nodes.")
	fmt.Println("        (Deploys from CLI config or JSON file; updates existing records without duplicates)")
	fmt.Println()
	fmt.Println("    gitmap ssh deploy node-config [all] [--except <id,ip,alias>] (alias: nc)")
	fmt.Println("        Deploy node topology, IP addresses, and aliases across all target fleet nodes.")
	fmt.Println("        (Recommended to run FIRST so all nodes know cluster topology)")
	fmt.Println()
	fmt.Println("    gitmap ssh deploy bin [all|<target>] [--file <path>]")
	fmt.Println("        Deploy local GitMap binary across remote fleet nodes and verify version.")
	fmt.Println("        (Autonomous one-liner replacement for manual scp commands)")
	fmt.Println()
	fmt.Println("    gitmap ssh deploy keys [all] [--except <id,ip,alias>]")
	fmt.Println("        Gather local and remote public keys, deduplicate unique keys,")
	fmt.Println("        and deploy into ~/.ssh/authorized_keys across all nodes for passwordless SSH.")
	fmt.Println()
	fmt.Println("    Flags:")
	fmt.Println("        -f, --file <path>      Deploy from specified JSON file instead of local CLI config")
	fmt.Println("        --except, -e <tokens>  Exclude nodes by numeric ID, worker ID, IP, or alias")
	fmt.Println("        --dry-run, -n          Preview deployment actions without modifying files")
	fmt.Println("        --json                 Output machine-readable JSON metrics")
	fmt.Println()
}
