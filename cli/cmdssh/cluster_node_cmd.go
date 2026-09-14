package cmdssh

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var executeNodeCmdFn = ExecuteNodeCommand

// ClusterNodeCmd represents the gitmap cluster node command.
var ClusterNodeCmd = &cobra.Command{
	Use:   "node <subcommand> [args...]",
	Short: "Provision and configure Ubuntu cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClusterNodeCLI(args)
	},
}

func showClusterNodeHelp() error {
	fmt.Println("Usage: gitmap cluster node <subcommand> [args...]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  set-ip <target> <new-ip> [--route-ip <gw>]           Configure static Netplan IP and gateway")
	fmt.Println("  install-base <target>                                Install base packages (curl, git, zsh, etc.)")
	fmt.Println("  create-user <target> <username> [pass] [--theme <t>] Create sudo user with zsh and oh-my-zsh")
	fmt.Println("  set-theme <target> <theme>                           Update ZSH_THEME in ~/.zshrc")
	fmt.Println("  purge <target>                                       Purge unneeded packages and clean apt cache")
	fmt.Println("\nFlags:")
	fmt.Println("  -h, --help    Show this help message")

	return nil
}

func isClusterNodeHelp(args []string) bool {
	isEmpty := len(args) == 0
	if isEmpty {
		return true
	}
	for _, arg := range args {
		if isClusterHelpFlag(arg) {
			return true
		}
	}

	return false
}

func getNextArg(args []string, idx int) string {
	hasMore := idx+1 < len(args)
	if hasMore {
		return args[idx+1]
	}

	return ""
}

func parseRouteIPFlag(arg, nextArg string) (string, int) {
	if strings.HasPrefix(arg, "--route-ip=") {
		return strings.TrimPrefix(arg, "--route-ip="), 1
	}
	hasFlag := arg == "--route-ip"
	hasVal := nextArg != ""
	if hasFlag && hasVal {
		return nextArg, 2
	}

	return "", 0
}

func buildSetIPParams(pos []string, routeIP string) (string, string, string, error) {
	isShort := len(pos) < 2
	if isShort {
		return "", "", "", apperror.NewValidationError("usage: gitmap cluster node set-ip <target> <new-ip> [--route-ip <gw>]")
	}

	return pos[0], pos[1], routeIP, nil
}

func parseSetIPArgs(args []string) (string, string, string, error) {
	var positional []string
	var routeIP string
	idx := 0
	for idx < len(args) {
		val, consumed := parseRouteIPFlag(args[idx], getNextArg(args, idx))
		if consumed > 0 {
			routeIP = val
			idx += consumed
			continue
		}
		positional = append(positional, args[idx])
		idx++
	}

	return buildSetIPParams(positional, routeIP)
}

func parseThemeFlag(arg, nextArg string) (string, int) {
	if strings.HasPrefix(arg, "--theme=") {
		return strings.TrimPrefix(arg, "--theme="), 1
	}
	hasFlag := arg == "--theme" || arg == "-t"
	hasVal := nextArg != ""
	if hasFlag && hasVal {
		return nextArg, 2
	}

	return "", 0
}

func extractUserPassword(pos []string) string {
	hasPass := len(pos) >= 3
	if hasPass {
		return pos[2]
	}

	return ""
}

func buildCreateUserParams(pos []string, theme string) (string, string, string, string, error) {
	isShort := len(pos) < 2
	if isShort {
		return "", "", "", "", apperror.NewValidationError("usage: gitmap cluster node create-user <target> <username> [password] [--theme <theme>]")
	}
	password := extractUserPassword(pos)

	return pos[0], pos[1], password, theme, nil
}

func parseCreateUserArgs(args []string) (string, string, string, string, error) {
	var positional []string
	var theme string
	idx := 0
	for idx < len(args) {
		val, consumed := parseThemeFlag(args[idx], getNextArg(args, idx))
		if consumed > 0 {
			theme = val
			idx += consumed
			continue
		}
		positional = append(positional, args[idx])
		idx++
	}

	return buildCreateUserParams(positional, theme)
}

func parseSetThemeArgs(args []string) (string, string, error) {
	isShort := len(args) < 2
	if isShort {
		return "", "", apperror.NewValidationError("usage: gitmap cluster node set-theme <target> <theme>")
	}

	return args[0], args[1], nil
}

func parseSingleTargetArg(subcmd string, args []string) (string, error) {
	isShort := len(args) < 1
	if isShort {
		msg := fmt.Sprintf("usage: gitmap cluster node %s <target>", subcmd)
		return "", apperror.NewValidationError(msg)
	}

	return args[0], nil
}

func runNodeSetIP(ctx context.Context, args []string) error {
	target, ip, routeIP, err := parseSetIPArgs(args)
	if err != nil {
		return err
	}
	script := GenerateNetplanScript(ip, routeIP)

	return executeClusterNodeScript(ctx, target, script)
}

func runNodeInstallBase(ctx context.Context, args []string) error {
	target, err := parseSingleTargetArg("install-base", args)
	if err != nil {
		return err
	}
	script := GenerateBasePackagesScript()

	return executeClusterNodeScript(ctx, target, script)
}

func runNodeCreateUser(ctx context.Context, args []string) error {
	target, user, pass, theme, err := parseCreateUserArgs(args)
	if err != nil {
		return err
	}
	script := GenerateCreateUserScript(user, pass, theme)

	return executeClusterNodeScript(ctx, target, script)
}

func runNodeSetTheme(ctx context.Context, args []string) error {
	target, theme, err := parseSetThemeArgs(args)
	if err != nil {
		return err
	}
	script := GenerateThemeScript(theme)

	return executeClusterNodeScript(ctx, target, script)
}

func runNodePurge(ctx context.Context, args []string) error {
	target, err := parseSingleTargetArg("purge", args)
	if err != nil {
		return err
	}
	script := GeneratePurgeScript()

	return executeClusterNodeScript(ctx, target, script)
}

func runScriptOnHosts(ctx context.Context, hosts []store.SSHHost, script string) []ClusterRunResult {
	results := make([]ClusterRunResult, 0, len(hosts))
	for _, host := range hosts {
		res := executeNodeCmdFn(ctx, host, script, true)
		results = append(results, res)
	}

	return results
}

func executeClusterNodeScript(ctx context.Context, target string, script string) error {
	hosts, err := resolveClusterHosts(ctx, target)
	if err != nil {
		return err
	}
	results := runScriptOnHosts(ctx, hosts, script)
	PrintClusterSummaryTable(results)
	if hasExecutionFailures(results) {
		return apperror.NewExecutionError("one or more cluster nodes failed recipe execution")
	}

	return nil
}

func routeClusterNodeCommand(ctx context.Context, sub string, rest []string) (error, bool) {
	switch sub {
	case "set-ip":
		return runNodeSetIP(ctx, rest), true
	case "install-base":
		return runNodeInstallBase(ctx, rest), true
	case "create-user":
		return runNodeCreateUser(ctx, rest), true
	case "set-theme":
		return runNodeSetTheme(ctx, rest), true
	case "purge":
		return runNodePurge(ctx, rest), true
	default:
		return nil, false
	}
}

// RunClusterNodeCLI dispatches cluster node provisioning recipes across target hosts.
func RunClusterNodeCLI(args []string) error {
	isHelp := isClusterNodeHelp(args)
	if isHelp {
		return showClusterNodeHelp()
	}
	ctx := context.Background()
	err, isMatched := routeClusterNodeCommand(ctx, args[0], args[1:])
	if isMatched {
		return err
	}

	return apperror.NewValidationError(fmt.Sprintf("unknown cluster node subcommand: %s", args[0]))
}
