// Package cmdagy — agy_wpr_cmd.go defines Cobra commands, flag binding, and CLI dispatch for WPR.
package cmdagy

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	wprIntervalStr string
	wprPrefixTpl   string
	wprSuffixTpl   string
	wprFileTarget  string
	isWPRSSH       bool
	isWPRJSON      bool
	isWPRDryRun    bool
	isWPRRestart   bool
	isWPROnce      bool
)

// AgyWPRCmd is the root Cobra command for watch-prompts-running.
var AgyWPRCmd = &cobra.Command{
	Use:     "watch-prompts-running [command]",
	Aliases: []string{"wpr", "watch-running-prompts"},
	Short:   "Watch, backup, auto-recover, and deploy running prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunWPRCLI(args)
	},
}

func init() {
	setupWPRFlags(AgyWPRCmd)
	AgyWPRCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		RenderWPRHelp()
	})
	AgyWPRCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return wprCompletionEntries(), cobra.ShellCompDirectiveNoFileComp
	}
	AgyCmd.AddCommand(AgyWPRCmd)
}

func setupWPRFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&wprIntervalStr, "time", "t", "2m", "Polling interval (minimum 10s, default 2m)")
	cmd.Flags().StringVarP(&wprPrefixTpl, "prefix", "p", "default", "Prefix template for re-injected prompts")
	cmd.Flags().StringVarP(&wprSuffixTpl, "suffix", "s", "", "Suffix template for re-injected prompts")
	cmd.Flags().StringVarP(&wprFileTarget, "file", "f", "", "Custom target Split-DB file path")
	cmd.Flags().BoolVar(&isWPRSSH, "ssh", false, "Execute across cluster SSH fleet")
	cmd.Flags().BoolVarP(&isWPRJSON, "json", "j", false, "Output machine-readable JSON metrics")
	cmd.Flags().BoolVarP(&isWPRDryRun, "dry-run", "n", false, "Simulate watchdog execution")
	cmd.Flags().BoolVar(&isWPRRestart, "restart", false, "Restart IDE process immediately")
	cmd.Flags().BoolVar(&isWPROnce, "once", false, "Run single verification pass and exit")
}

func wprCompletionEntries() []string {
	return []string{
		"all\tSnapshot all running/queued prompts to Split-DB",
		"ls\tInspect watched projects and schedule status",
		"start\tStart watch loop and auto-recovery",
		"disable\tPause watch loop without removing projects",
		"shutdown\tStop watch loop and terminate IDE process",
		"restart\tRestart IDE process and re-inject prompts",
		"switch-account\tSwitch account with prompt preservation",
		"fast-forward\tFast-forward account switch and IDE restart",
		"logs\tInspect watcher event history and recovery logs",
		"status\tInspect watchdog runtime state and cached fleet",
		"remove\tRemove project from watch list",
		"deploy\tDeploy Split-DB and enqueue WPR task via SSH",
		"help\tShow command usage guide",
	}
}

// RunWPRCLI routes watch-prompts-running CLI arguments to specific subcommands.
func RunWPRCLI(args []string) error {
	opts := parseWPROptions(args)
	if opts.Subcommand == "" || opts.Subcommand == "help" {
		RenderWPRHelp()
		return nil
	}

	return routeWPRSubcommand(opts, args)
}

func parseWPROptions(args []string) WPROptions {
	opts := WPROptions{
		IntervalStr:    "2m",
		Interval:       2 * time.Minute,
		PrefixTemplate: "default",
	}

	var positional []string
	for i := 0; i < len(args); i++ {
		i = parseSingleWPRArg(args, i, &opts, &positional)
	}

	extractWPRPositional(positional, &opts)
	return opts
}

func parseSingleWPRArg(args []string, i int, opts *WPROptions, positional *[]string) int {
	arg := args[i]
	switch {
	case (arg == "-t" || arg == "--time") && i+1 < len(args):
		opts.IntervalStr = args[i+1]
		opts.Interval, _ = time.ParseDuration(args[i+1])
		return i + 1
	case (arg == "-p" || arg == "--prefix") && i+1 < len(args):
		opts.PrefixTemplate = args[i+1]
		return i + 1
	case (arg == "-s" || arg == "--suffix") && i+1 < len(args):
		opts.SuffixTemplate = args[i+1]
		return i + 1
	case arg == "--ssh":
		opts.IsSSH = true
	case arg == "-j" || arg == "--json":
		opts.IsJSON = true
	case arg == "-n" || arg == "--dry-run":
		opts.IsDryRun = true
	case arg == "--restart":
		opts.IsRestart = true
	case arg == "--once" || arg == "-1":
		opts.IsOnce = true
	case !strings.HasPrefix(arg, "-"):
		*positional = append(*positional, arg)
	}

	return i
}

func extractWPRPositional(pos []string, opts *WPROptions) {
	if len(pos) > 0 {
		opts.Subcommand = strings.ToLower(pos[0])
	}
	if len(pos) > 1 {
		opts.Target = pos[1]
		opts.Email = pos[1]
	}
	if len(pos) > 2 {
		opts.TargetProject = pos[2]
	}
}

func routeWPRSubcommand(opts WPROptions, rawArgs []string) error {
	switch opts.Subcommand {
	case "all":
		return RunWPRAll(opts)
	case "ls", "list":
		return RunWPRLs(opts)
	case "start":
		return RunWPRStart(opts)
	case "disable", "pause":
		return RunWPRDisable(opts)
	case "shutdown", "sd", "stop":
		return RunWPRShutdown(opts)
	case "restart":
		return RunWPRRestart(opts)
	case "switch-account", "sa":
		return RunWPRSwitchAccount(opts)
	case "fast-forward", "ff":
		return RunWPRFastForward(opts)
	case "logs", "log":
		return RunWPRLogs(opts)
	case "status", "st":
		return RunWPRStatus(opts)
	case "remove", "rm", "del", "delete":
		return RunWPRRemove(opts)
	case "deploy":
		return RunWPRDeploy(filterDeployArgs(rawArgs), opts)
	default:
		RenderWPRHelp()
		return nil
	}
}

func filterDeployArgs(args []string) []string {
	var clean []string
	foundDeploy := false
	for _, a := range args {
		if !foundDeploy && strings.EqualFold(a, "deploy") {
			foundDeploy = true
			continue
		}
		if foundDeploy && !strings.HasPrefix(a, "-") {
			clean = append(clean, a)
		}
	}

	return clean
}
