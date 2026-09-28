package cmdwinutil

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RunWinUtil dispatches winutil commands.
func RunWinUtil(args []string) error {
	if len(args) == 0 || isWinUtilHelp(args[0]) {
		printWinUtilHelp()
		return nil
	}

	subCmd := strings.ToLower(args[0])
	return dispatchWinUtilSubcommand(subCmd, args[1:])
}

func isWinUtilHelp(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func dispatchWinUtilSubcommand(subCmd string, subArgs []string) error {
	switch subCmd {
	case "copilot":
		return RunWinUtilCopilotCLI(subArgs)
	case "edge":
		return RunWinUtilEdgeCLI(subArgs)
	default:
		printWinUtilHelp()
		return nil
	}
}

// RunWinUtilCopilotCLI executes Copilot removal CLI.
func RunWinUtilCopilotCLI(args []string) error {
	opts := parseWinUtilOptions(args)
	res, err := RunCopilotUninstall(opts.IsDryRun)
	if err != nil {
		return err
	}

	renderRemovalResult(res, opts)
	return nil
}

// RunWinUtilEdgeCLI executes Edge removal CLI.
func RunWinUtilEdgeCLI(args []string) error {
	opts := parseWinUtilOptions(args)
	res, err := RunEdgeUninstall(opts.HasKeepWebView2, opts.IsDryRun)
	if err != nil {
		return err
	}

	renderRemovalResult(res, opts)
	return nil
}

func parseWinUtilOptions(args []string) WinUtilOptions {
	opts := WinUtilOptions{HasKeepWebView2: true}
	for i := 0; i < len(args); i++ {
		parseSingleWinUtilFlag(args[i], &opts)
	}
	return opts
}

func parseSingleWinUtilFlag(arg string, opts *WinUtilOptions) {
	low := strings.ToLower(arg)
	switch {
	case low == "--dry-run" || low == "-n" || low == "dry-run":
		opts.IsDryRun = true
	case low == "--yes" || low == "-y" || low == "yes":
		opts.HasAutoYes = true
	case low == "--purge-webview2":
		opts.HasKeepWebView2 = false
	case low == "--json":
		opts.IsJSON = true
	}
}

func renderRemovalResult(res WinRemovalResult, opts WinUtilOptions) {
	if opts.IsJSON {
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
		return
	}

	renderResultBanner(res)
	renderResultActions(res.Actions)
	renderResultWarnings(res.Warnings)
}

func renderResultBanner(res WinRemovalResult) {
	fmt.Println()
	mode := "Completed"
	if res.IsDryRun {
		mode = "Dry-Run Preview"
	}
	fmt.Printf("  WinUtil: %s uninstallation (%s)\n", res.Target, mode)
	fmt.Println("  " + strings.Repeat("═", 50))
}

func renderResultActions(actions []string) {
	for _, act := range actions {
		fmt.Printf("  ✔ %s\n", act)
	}
}

func renderResultWarnings(warnings []string) {
	for _, warn := range warnings {
		fmt.Printf("  ⚠ %s\n", warn)
	}
	fmt.Println()
}

func printWinUtilHelp() {
	fmt.Println("Usage: gitmap winutil <target> uninstall [flags]")
	fmt.Println("       gitmap uninstall <target> [flags]")
	fmt.Println()
	fmt.Println("Targets:")
	fmt.Println("  copilot          Remove Windows Copilot Appx and set disable policies")
	fmt.Println("  edge             Remove Microsoft Edge (Chris Titus WinUtil parity)")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -n, --dry-run          Preview actions without modifying system")
	fmt.Println("      --purge-webview2   Also remove WebView2 runtime (Edge only)")
	fmt.Println("  -y, --yes              Bypass confirmation")
	fmt.Println("      --json             Output structured JSON")
}
