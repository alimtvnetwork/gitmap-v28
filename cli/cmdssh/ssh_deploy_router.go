// Package cmdssh — ssh_deploy_router.go routes ssh deploy subcommands (keys, node-config, ide).
package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunSSHDeployRouterCLI routes `gitmap ssh deploy` to keys, node-config, or IDE configuration.
func RunSSHDeployRouterCLI(args []string) error {
	if len(args) == 0 {
		printDeployHelp()
		return nil
	}
	sub := strings.ToLower(args[0])
	switch sub {
	case "ide", "agy", "antigravity":
		return routeDeployIdeOrAgy(args[1:])
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
	case "all-keys", "all-key", "allkeys":
		return RunSSHDeployKeysCLI(args[1:])
	case "all":
		if len(args) > 1 && isDeployKeysSubToken(strings.ToLower(args[1])) {
			return RunSSHDeployKeysCLI(args[2:])
		}
		return RunSSHDeployConfigSSHCLI(args)
	case "keys", "key", "k", "keys-all", "deploy-keys", "deploy-keys-all", "keys-hyphen-all":
		if len(args) > 1 && strings.EqualFold(args[1], "all") {
			return RunSSHDeployKeysCLI(args[2:])
		}
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

func routeDeployIdeOrAgy(args []string) error {
	opts := cmdagy.AgyDeployOptions{
		IsAll: true,
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--all":
			opts.IsAll = true
		case arg == "--plugins":
			opts.HasPlugins = true
		case arg == "--skills":
			opts.HasSkills = true
		case arg == "--binaries":
			opts.HasBinaries = true
		case arg == "--projects":
			opts.HasProjects = true
		case arg == "--dry-run":
			opts.IsDryRun = true
		case arg == "--json":
			opts.IsJSON = true
		case arg == "--force":
			opts.IsForce = true
		case arg == "--restart":
			opts.IsRestart = true
		case arg == "--preset" && i+1 < len(args):
			i++
			opts.Preset = args[i]
		case strings.HasPrefix(arg, "--preset="):
			opts.Preset = strings.TrimPrefix(arg, "--preset=")
		case arg == "--preset":
			opts.Preset = "eager"
		case arg == "--theme" && i+1 < len(args):
			i++
			opts.Theme = args[i]
		case strings.HasPrefix(arg, "--theme="):
			opts.Theme = strings.TrimPrefix(arg, "--theme=")
		case arg == "--theme":
			opts.Theme = "dark"
		case arg == "--target" && i+1 < len(args):
			i++
			opts.TargetNode = args[i]
		case strings.HasPrefix(arg, "--target="):
			opts.TargetNode = strings.TrimPrefix(arg, "--target=")
		case !strings.HasPrefix(arg, "-") && opts.TargetNode == "":
			opts.TargetNode = arg
		}
	}

	if opts.TargetNode == "" || opts.TargetNode == "help" || opts.TargetNode == "-h" || opts.TargetNode == "--help" {
		printDeployHelp()
		return nil
	}

	res, err := cmdagy.ExecuteAgyDeploy(opts)
	if err != nil {
		return err
	}

	if opts.IsJSON {
		return printDeployResultJSON(res)
	}

	cmdagy.RenderAgyDeploySummary(res)
	return nil
}

func printDeployResultJSON(res *cmdagy.AgyDeployResultJSON) error {
	payload, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(payload))

	return nil
}

func isDeployKeysSubToken(s string) bool {
	return s == "keys" || s == "key" || s == "k" || s == "all-keys" || s == "keys-all"
}

func isDeployAllKeysArg(args []string) bool {
	return len(args) > 1 && isDeployKeysSubToken(strings.ToLower(args[1]))
}

func routeFallbackDeploy(args []string) error {
	first := strings.ToLower(args[0])
	if first == "all-keys" || first == "allkeys" {
		return RunSSHDeployKeysCLI(args[1:])
	}
	if first == "all" && isDeployAllKeysArg(args) {
		return RunSSHDeployKeysCLI(args[2:])
	}
	if first == "all" {
		return RunSSHDeployConfigSSHCLI(args)
	}
	if strings.HasPrefix(first, "-") {
		return RunSSHDeployConfigSSHCLI(args)
	}
	fmt.Printf("\n  %sUnknown deploy target: %q%s\n", constants.ColorRed, args[0], constants.ColorReset)
	printDeployHelp()
	return nil
}

func printDeployHelp() {
	fmt.Printf("\n  %s🚀 GitMap SSH Fleet Deploy Commands%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    gitmap deploy ide <target> [flags] (alias: gitmap deploy agy)")
	fmt.Println("        Deploy Antigravity IDE themes, presets, 4 official plugins, and 43 skills across fleet.")
	fmt.Println()
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
	fmt.Println("    • Deploy Antigravity to fleet node: gitmap deploy ide u1 --all")
	fmt.Println("    • Export nodes on this machine:     gitmap ssh export-json")
	fmt.Println("    • Import nodes on other machine:    gitmap ssh import-json gitmap-ssh-nodes.json")
	fmt.Println("    • Deploy keys to all nodes:         gitmap deploy keys all  (or: gitmap deploy-keys-all)")
	fmt.Println("    • Verify passwordless status:       gitmap ssh check")
	fmt.Println()
}
