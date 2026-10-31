// Package cmdssh — ssh_deploy_router.go routes ssh deploy subcommands (keys, node-config, ide, bin).
package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunSSHDeployRouterCLI routes `gitmap ssh deploy` to keys, node-config, IDE, or binary deployment.
func RunSSHDeployRouterCLI(args []string) error {
	if len(args) == 0 {
		printDeployHelp()

		return nil
	}

	sub := strings.ToLower(args[0])

	return dispatchDeploySubcommand(sub, args)
}

func dispatchDeploySubcommand(sub string, args []string) error {
	if isDeployIdeCommand(sub) {
		return routeDeployIdeOrAgy(args[1:])
	}

	if isDeployConfigCommand(sub) {
		return dispatchDeployConfig(sub, args)
	}

	if isDeployKeysCommand(sub) {
		return dispatchDeployKeys(sub, args)
	}

	return dispatchDeployOperations(sub, args)
}

func isDeployIdeCommand(sub string) bool {
	isIde := sub == "ide" || sub == "agy" || sub == "antigravity"

	return isIde
}

func isDeployConfigCommand(sub string) bool {
	isConfig := sub == "config" || sub == "config-ssh" || sub == "ssh-config" || sub == "ssh" || sub == "node-config" || sub == "nodeconfig" || sub == "nc" || sub == "nodes"

	return isConfig
}

func dispatchDeployConfig(sub string, args []string) error {
	if sub == "ssh" {
		return routeDeploySSHSub(args)
	}

	if sub == "node-config" || sub == "nodeconfig" || sub == "nc" || sub == "nodes" {
		return RunSSHDeployConfigSSHCLI(args[1:])
	}

	return routeDeployConfigSSH(args)
}

func isDeployKeysCommand(sub string) bool {
	isKeys := sub == "all-keys" || sub == "all-key" || sub == "allkeys" || sub == "all" ||
		sub == "keys" || sub == "key" || sub == "k" || sub == "keys-all" || sub == "deploy-keys" ||
		sub == "deploy-keys-all" || sub == "keys-hyphen-all"

	return isKeys
}

func dispatchDeployKeys(sub string, args []string) error {
	if sub == "all-keys" || sub == "all-key" || sub == "allkeys" {
		return RunSSHDeployKeysCLI(args[1:])
	}

	if sub == "all" {
		return routeDeployAll(args)
	}

	return routeDeployKeysSub(args)
}

func dispatchDeployOperations(sub string, args []string) error {
	if isDeployBinFallback(sub) {
		return RunSSHDeployBinCLI(args[1:])
	}

	if sub == "export" || sub == "export-json" || sub == "exportjson" {
		return RunSSHNodesExportJSON(args[1:])
	}

	if sub == "import" || sub == "import-json" || sub == "importjson" {
		return RunSSHNodesImportJSON(args[1:])
	}

	return dispatchDeployTail(sub, args)
}

func dispatchDeployTail(sub string, args []string) error {
	if sub == "macro" || sub == "macros" {
		return cmdmacro.ExecuteMacroDeploySSH(args[1:])
	}

	if sub == "help" || sub == "--help" || sub == "-h" {
		printDeployHelp()

		return nil
	}

	return routeFallbackDeploy(args)
}

func routeDeployConfigSSH(args []string) error {
	if len(args) > 1 && strings.EqualFold(args[1], "ssh") {
		return RunSSHDeployConfigSSHCLI(args[2:])
	}

	return RunSSHDeployConfigSSHCLI(args[1:])
}

func routeDeploySSHSub(args []string) error {
	if len(args) > 1 && strings.EqualFold(args[1], "config") {
		return RunSSHDeployConfigSSHCLI(args[2:])
	}

	return RunSSHDeployConfigSSHCLI(args[1:])
}

func routeDeployAll(args []string) error {
	if len(args) > 1 && isDeployKeysSubToken(strings.ToLower(args[1])) {
		return RunSSHDeployKeysCLI(args[2:])
	}

	return RunSSHDeployConfigSSHCLI(args)
}

func routeDeployKeysSub(args []string) error {
	if len(args) > 1 && strings.EqualFold(args[1], "all") {
		return RunSSHDeployKeysCLI(args[2:])
	}

	return RunSSHDeployKeysCLI(args[1:])
}

func routeDeployIdeOrAgy(args []string) error {
	opts := parseAgyDeployArgs(args)
	if opts.TargetNode == "" || opts.TargetNode == "help" || opts.TargetNode == "-h" || opts.TargetNode == "--help" {
		printDeployHelp()

		return nil
	}

	return executeAgyDeployWorkflow(opts)
}

func executeAgyDeployWorkflow(opts cmdagy.AgyDeployOptions) error {
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

func parseAgyDeployArgs(args []string) cmdagy.AgyDeployOptions {
	opts := cmdagy.AgyDeployOptions{IsAll: true}
	for i := 0; i < len(args); i++ {
		i = parseSingleAgyArg(args, i, &opts)
	}

	return opts
}

func parseSingleAgyArg(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	arg := args[i]
	if isAgyModifierFlag(arg) {
		applyAgyModifier(arg, opts)

		return i
	}

	return parseAgyValueArg(args, i, opts)
}

func isAgyModifierFlag(arg string) bool {
	isMod := argMatches(arg, "--all") ||
		argMatches(arg, "--plugins") ||
		argMatches(arg, "--skills") ||
		argMatches(arg, "--binaries") ||
		argMatches(arg, "--projects") ||
		argMatches(arg, "--dry-run") ||
		argMatches(arg, "--json") ||
		argMatches(arg, "--force") ||
		argMatches(arg, "--restart")

	return isMod
}

func applyAgyModifier(arg string, opts *cmdagy.AgyDeployOptions) {
	applyAgyResourceFlags(arg, opts)
	applyAgyControlFlags(arg, opts)
}

func applyAgyResourceFlags(arg string, opts *cmdagy.AgyDeployOptions) {
	if argMatches(arg, "--all") {
		opts.IsAll = true
	}

	if argMatches(arg, "--plugins") {
		opts.HasPlugins = true
	}

	applyAgyResourceFlagsPart2(arg, opts)
}

func applyAgyResourceFlagsPart2(arg string, opts *cmdagy.AgyDeployOptions) {
	if argMatches(arg, "--skills") {
		opts.HasSkills = true
	}

	if argMatches(arg, "--binaries") {
		opts.HasBinaries = true
	}

	if argMatches(arg, "--projects") {
		opts.HasProjects = true
	}
}

func applyAgyControlFlags(arg string, opts *cmdagy.AgyDeployOptions) {
	if argMatches(arg, "--dry-run") {
		opts.IsDryRun = true
	}

	if argMatches(arg, "--json") {
		opts.IsJSON = true
	}

	applyAgyControlFlagsPart2(arg, opts)
}

func applyAgyControlFlagsPart2(arg string, opts *cmdagy.AgyDeployOptions) {
	if argMatches(arg, "--force") {
		opts.IsForce = true
	}

	if argMatches(arg, "--restart") {
		opts.IsRestart = true
	}
}

func parseAgyValueArg(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	arg := args[i]
	if argMatches(arg, "--preset") && i+1 < len(args) {
		opts.Preset = args[i+1]

		return i + 1
	}

	return parseAgyPresetDirectOrNext(args, i, opts)
}

func parseAgyPresetDirectOrNext(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--preset=") {
		opts.Preset = strings.TrimPrefix(arg, "--preset=")

		return i
	}

	return parseAgyPresetEagerOrNext(args, i, opts)
}

func parseAgyPresetEagerOrNext(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	if argMatches(args[i], "--preset") {
		opts.Preset = "eager"

		return i
	}

	return parseAgyThemeOrTarget(args, i, opts)
}

func parseAgyThemeOrTarget(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	arg := args[i]
	if argMatches(arg, "--theme") && i+1 < len(args) {
		opts.Theme = args[i+1]

		return i + 1
	}

	return parseAgyThemeDirectOrNext(args, i, opts)
}

func parseAgyThemeDirectOrNext(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--theme=") {
		opts.Theme = strings.TrimPrefix(arg, "--theme=")

		return i
	}

	return parseAgyThemeDarkOrNext(args, i, opts)
}

func parseAgyThemeDarkOrNext(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	if argMatches(args[i], "--theme") {
		opts.Theme = "dark"

		return i
	}

	return parseAgyTargetArg(args, i, opts)
}

func parseAgyTargetArg(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	arg := args[i]
	if argMatches(arg, "--target") && i+1 < len(args) {
		opts.TargetNode = args[i+1]

		return i + 1
	}

	return parseAgyTargetDirectOrPositional(args, i, opts)
}

func parseAgyTargetDirectOrPositional(args []string, i int, opts *cmdagy.AgyDeployOptions) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--target=") {
		opts.TargetNode = strings.TrimPrefix(arg, "--target=")

		return i
	}

	if !strings.HasPrefix(arg, "-") && opts.TargetNode == "" {
		opts.TargetNode = arg
	}

	return i
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
	isToken := s == "keys" || s == "key" || s == "k" || s == "all-keys" || s == "keys-all"

	return isToken
}

func isDeployAllKeysArg(args []string) bool {
	hasKeys := len(args) > 1 && isDeployKeysSubToken(strings.ToLower(args[1]))

	return hasKeys
}

func routeFallbackDeploy(args []string) error {
	first := strings.ToLower(args[0])
	if isDeployBinFallback(first) {
		return RunSSHDeployBinCLI(args[1:])
	}

	if isDeployKeysFallback(args, first) {
		return routeDeployKeysFallback(args, first)
	}

	return routeFallbackDeployRest(args, first)
}

func isDeployBinFallback(first string) bool {
	isBin := first == "bin" || first == "binary" || first == "exe" || first == "deploy-bin" || first == "push-bin" || first == "sync-bin"

	return isBin
}

func isDeployKeysFallback(args []string, first string) bool {
	isKeys := first == "all-keys" || first == "allkeys" || (first == "all" && isDeployAllKeysArg(args))

	return isKeys
}

func routeDeployKeysFallback(args []string, first string) error {
	if first == "all" {
		return RunSSHDeployKeysCLI(args[2:])
	}

	return RunSSHDeployKeysCLI(args[1:])
}

func routeFallbackDeployRest(args []string, first string) error {
	if first == "all" || strings.HasPrefix(first, "-") {
		return RunSSHDeployConfigSSHCLI(args)
	}

	fmt.Printf("\n  %sUnknown deploy target: %q%s\n", constants.ColorRed, args[0], constants.ColorReset)
	printDeployHelp()

	return nil
}

func printDeployHelp() {
	printDeployHelpBanner()
	printDeployHelpSubcommands()
	printDeployHelpSuggestions()
}

func printDeployHelpBanner() {
	fmt.Printf("\n  %s🚀 GitMap SSH Fleet Deploy Commands%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    gitmap deploy ide <target> [flags] (alias: gitmap deploy agy)")
	fmt.Println("        Deploy Antigravity IDE themes, presets, 4 official plugins, and 43 skills across fleet.")
	fmt.Println()
	fmt.Println("    gitmap deploy keys [all] [--except <id,ip,alias>] (alias: deploy-keys-all)")
	fmt.Println("        Gather local and remote public keys, deduplicate unique keys,")
	fmt.Println("        and deploy into authorized_keys across all nodes for passwordless SSH.")
	fmt.Println()
}

func printDeployHelpSubcommands() {
	fmt.Println("    gitmap deploy config ssh [all|<target>] [--file <path>] [--except <id,ip,alias>]")
	fmt.Println("        Deploy node topology, IP addresses, aliases, and credentials across fleet nodes.")
	fmt.Println()
	fmt.Println("    gitmap ssh export-json [file.json] (alias: nodes export-json)")
	fmt.Println("        Export all registered SSH fleet nodes, IPs, and credentials to a portable JSON file.")
	fmt.Println()
	fmt.Println("    gitmap ssh import-json [file.json] (alias: deploy import)")
	fmt.Println("        Import fleet nodes and credentials from a JSON file into the local registry.")
	fmt.Println()
	fmt.Println("    gitmap ssh deploy bin [all|<target>] [flags] (alias: push-bin, sync-bin)")
	fmt.Println("        Deploy local GitMap binary across remote fleet nodes and verify version.")
	fmt.Println("        Flags: --file <path>, --dest <path>, --force, --install, --dry-run, --timeout <dur>")
	fmt.Println()
}

func printDeployHelpSuggestions() {
	fmt.Println("  💡 Suggestions & Cross-Section Guidance:")
	fmt.Println("    • Deploy Antigravity to fleet node: gitmap deploy ide u1 --all")
	fmt.Println("    • Deploy binary across all nodes:   gitmap ssh deploy bin all")
	fmt.Println("    • Deploy binary with force:         gitmap ssh deploy bin u1 --file ./gitmap --force")
	fmt.Println("    • Deploy keys to all nodes:         gitmap deploy keys all  (or: gitmap deploy-keys-all)")
	fmt.Println("    • Verify passwordless status:       gitmap ssh check")
	fmt.Println()
}
