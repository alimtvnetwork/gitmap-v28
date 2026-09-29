// Package cmdssh — ssh_deploy_router.go routes ssh deploy subcommands (keys, node-config).
package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
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
	case "config", "config-ssh", "ssh-config":
		if len(args) > 1 && strings.EqualFold(args[1], "ssh") {
			return RunSSHDeployConfigSSHCLI(args[2:])
		}
		return RunSSHDeployConfigSSHCLI(args[1:])
	case "ssh":
		if len(args) > 1 && strings.EqualFold(args[1], "config") {
			return RunSSHDeployConfigSSHCLI(args[2:])
		}
		return RunSSHDeployConfigSSHCLI(args[1:])
	case "keys", "key", "k", "keys-all", "deploy-keys", "deploy-keys-all", "keys-hyphen-all":
		return RunSSHDeployKeysCLI(args[1:])
	case "export", "export-json", "exportjson":
		return RunSSHNodesExportJSON(args[1:])
	case "bin", "binary", "exe", "gitmap":
		return RunSSHDeployBinCLI(args[1:])
	case "node-config", "nodeconfig", "nc", "nodes":
		return RunSSHDeployConfigSSHCLI(args[1:])
	case "macro", "macros":
		return cmdmacro.ExecuteMacroDeploySSH(args[1:])
	case "import", "import-json", "importjson":
		return RunSSHNodesImportJSON(args[1:])
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
	fmt.Println("    gitmap deploy keys [all] [--except <id,ip,alias>] (alias: deploy-keys-all)")
	fmt.Println("        Gather local and remote public keys, deduplicate unique keys,")
	fmt.Println("        and deploy into authorized_keys across all nodes for passwordless SSH.")
	fmt.Println()
	fmt.Println("    gitmap deploy config ssh [all|<target>] [--file <path>] [--except <id,ip,alias>]")
	fmt.Println("        Deploy node topology, IP addresses, aliases, and credentials across fleet nodes.")
	fmt.Println()
	fmt.Println("    gitmap ssh export-json [file.json] (alias: nodes export-json)")
	fmt.Println("        Export all registered SSH fleet nodes, IPs, and credentials to a portable JSON file.")
	fmt.Println()
	fmt.Println("    gitmap ssh import-json [file.json] (alias: deploy import)")
	fmt.Println("        Import fleet nodes and credentials from a JSON file into the local registry.")
	fmt.Println()
	fmt.Println("    gitmap ssh deploy bin [all|<target>] [--file <path>]")
	fmt.Println("        Deploy local GitMap binary across remote fleet nodes and verify version.")
	fmt.Println()
	fmt.Println("  💡 Suggestions & Cross-Section Guidance:")
	fmt.Println("    • Export nodes on this machine:  gitmap ssh export-json")
	fmt.Println("    • Import nodes on other machine: gitmap ssh import-json gitmap-ssh-nodes.json")
	fmt.Println("    • Deploy keys to all nodes:      gitmap deploy keys all  (or: gitmap deploy-keys-all)")
	fmt.Println("    • Verify passwordless status:    gitmap ssh check")
	fmt.Println()
}
