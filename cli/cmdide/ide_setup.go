package cmdide

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/setup"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runIDESetup(args []string) error {
	opts, err := parseSetupOptions(args)
	if err != nil {
		return err
	}

	res, runErr := setup.RunSetup(opts)
	if runErr != nil {
		return runErr
	}

	return outputSetupResult(res, opts)
}

func outputSetupResult(res *setup.IDESetupResult, opts setup.IDESetupOptions) error {
	if opts.IsJSON {
		return printJSONResult(res)
	}

	printTerminalResult(res, opts)

	return nil
}

func defaultSetupOptions() setup.IDESetupOptions {
	return setup.IDESetupOptions{
		Target:          setup.TargetAll,
		Provider:        setup.ProviderKimi,
		IsBackupEnabled: true,
	}
}

func parseSetupOptions(args []string) (setup.IDESetupOptions, error) {
	opts := defaultSetupOptions()
	if err := populateSetupOptions(&opts, args); err != nil {
		return opts, err
	}

	setup.ApplyPresetDefaults(&opts)

	return opts, nil
}

func populateSetupOptions(opts *setup.IDESetupOptions, args []string) error {
	for i := 0; i < len(args); i++ {
		step, isFlag, err := parseSetupArg(opts, args, i)
		if err != nil {
			return err
		}

		if isFlag {
			i += step
			continue
		}

		opts.Target = resolveTarget(args[i])
	}

	return nil
}

func parseSetupArg(opts *setup.IDESetupOptions, args []string, idx int) (int, bool, error) {
	if step, hasVal, err := parseValueFlags(opts, args, idx); hasVal {
		return step, true, err
	}

	return parseBoolFlags(opts, args[idx])
}

func parseValueFlags(opts *setup.IDESetupOptions, args []string, idx int) (int, bool, error) {
	if idx+1 >= len(args) {
		return 0, false, nil
	}

	arg, val := args[idx], args[idx+1]
	if isProviderFlag(arg) {
		p, err := setup.ParseIDEProviderType(val)
		if err != nil {
			return 0, true, err
		}
		opts.Provider = p

		return 1, true, nil
	}

	return parseOverrideValueFlags(opts, arg, val)
}

func isProviderFlag(arg string) bool {
	return arg == "-p" || arg == "--provider"
}

func parseOverrideValueFlags(opts *setup.IDESetupOptions, arg, val string) (int, bool, error) {
	switch arg {
	case "-m", "--model":
		opts.Model = val
		return 1, true, nil
	case "-b", "--base-url", "-u":
		opts.BaseURL = val
		return 1, true, nil
	case "-k", "--api-key":
		opts.APIKey = val
		return 1, true, nil
	}

	return 0, false, nil
}

func parseBoolFlags(opts *setup.IDESetupOptions, arg string) (int, bool, error) {
	if isDryRunOrJSONFlag(opts, arg) || isTestOrQuietFlag(opts, arg) {
		return 0, true, nil
	}

	if arg == "-f" || arg == "--force" {
		opts.IsForceCreate = true

		return 0, true, nil
	}

	return 0, false, nil
}

func isDryRunOrJSONFlag(opts *setup.IDESetupOptions, arg string) bool {
	if arg == "-n" || arg == "--dry-run" {
		opts.IsDryRun = true
		return true
	}
	if arg == "-j" || arg == "--json" {
		opts.IsJSON = true
		return true
	}

	return false
}

func isTestOrQuietFlag(opts *setup.IDESetupOptions, arg string) bool {
	if arg == "--test" {
		opts.IsTestEnabled = true
		return true
	}
	if arg == "-q" || arg == "--quiet" {
		opts.IsQuiet = true
		return true
	}

	return false
}

func resolveTarget(raw string) setup.IDETargetType {
	clean := strings.ToLower(strings.TrimSpace(raw))
	switch clean {
	case "cline", "composer":
		return setup.TargetCline
	case "continue", "tab":
		return setup.TargetContinue
	case "zed":
		return setup.TargetZed
	default:
		return setup.TargetAll
	}
}

func printJSONResult(res *setup.IDESetupResult) error {
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func printTerminalResult(res *setup.IDESetupResult, opts setup.IDESetupOptions) {
	fmt.Printf("\n%s● GitMap AI IDE Setup%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  • Provider: %s (%s)\n", res.Provider, opts.Model)
	if opts.BaseURL != "" {
		fmt.Printf("  • Base URL: %s\n", opts.BaseURL)
	}
	printConnectionStatus(res.ConnectionTest)
	printTargetResults(res.TargetResults)
	printSummaryFooter(res, opts.IsDryRun)
}

func printConnectionStatus(test setup.ConnectionTestResult) {
	if test.LatencyMs <= 0 && !test.IsSuccess {
		return
	}
	if test.IsSuccess {
		fmt.Printf("  • Probe:    %sConnected%s (%dms)\n", constants.ColorGreen, constants.ColorReset, test.LatencyMs)
		return
	}

	fmt.Printf("  • Probe:    %sFailed%s: %s\n", constants.ColorRed, constants.ColorReset, test.ErrorMessage)
}

func printTargetResults(targets map[string]setup.IDETargetResult) {
	fmt.Printf("\n  %sConfigured Targets:%s\n", constants.ColorBold, constants.ColorReset)
	for name, tr := range targets {
		status := constants.ColorGreen + "Configured" + constants.ColorReset
		if !tr.IsConfigured {
			status = constants.ColorYellow + "Skipped / Error" + constants.ColorReset
		}
		fmt.Printf("    %-10s %s (%s)\n", name, status, tr.ConfigPath)
	}
}

func printSummaryFooter(res *setup.IDESetupResult, isDryRun bool) {
	if isDryRun {
		fmt.Printf("\n%s[dry-run] No configuration files were modified.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	if res.IsSuccess {
		fmt.Printf("\n%s✓ Setup completed successfully.%s\n\n", constants.ColorGreen, constants.ColorReset)
		return
	}

	fmt.Printf("\n%s⚠ Setup completed with warnings.%s\n\n", constants.ColorYellow, constants.ColorReset)
}
