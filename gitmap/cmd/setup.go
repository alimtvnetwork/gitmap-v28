package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/completion"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/setup"
)

// runSetup handles the "setup" subcommand.
func runSetup(args []string) {
	// Subcommand: `gitmap setup print-path-snippet ...`
	// Used by run.sh + install.sh to fetch the canonical marker-block
	// snippet so all three drivers emit byte-identical output.
	if len(args) > 0 && args[0] == "print-path-snippet" {
		runPrintPathSnippet(args[1:])

		return
	}

	checkHelp("setup", args)
	configPath, dryRun, hasConfig := parseSetupFlags(args)
	configPath = resolveSetupConfigPath(configPath, hasConfig)
	cfg := mustLoadSetupConfig(configPath)
	printSetupBanner(dryRun)
	result := setup.Apply(cfg, dryRun)
	installShellCompletion(dryRun)
	installCDFunction(dryRun)
	installPathSnippet(dryRun)
	ensureGitignoreStep(dryRun)
	verifyShellWrapper(dryRun)
	printSetupSummary(result)
}

// installPathSnippet writes the canonical marker-block PATH snippet to
// the user's profile so future shells pick up the gitmap install dir.
// Idempotent: rewrites the existing block if present, otherwise appends.
func installPathSnippet(dryRun bool) {
	shell := completion.DetectShell()
	fmt.Printf("\n  %s%s%s\n", constants.ColorYellow, "PATH snippet:", constants.ColorReset)

	dir := resolveActiveBinaryDir()
	if len(dir) == 0 {
		fmt.Fprintf(os.Stderr, "  %sskipped: could not resolve active gitmap directory%s\n",
			constants.ColorYellow, constants.ColorReset)

		return
	}

	if dryRun {
		fmt.Printf("  %s[dry-run]%s would write PATH snippet for %s -> %s\n",
			constants.ColorDim, constants.ColorReset, shell, dir)

		return
	}

	res, err := setup.WritePathSnippet(shell, dir, "gitmap setup", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s%v%s\n", constants.ColorYellow, err, constants.ColorReset)

		return
	}
	fmt.Printf("  %s%s%s -> %s\n", constants.ColorGreen, res.Action, constants.ColorReset, res.Profile)
}

// resolveActiveBinaryDir returns the directory containing the running
// gitmap binary, used as the PATH entry to inject.
func resolveActiveBinaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		resolved = exe
	}

	return filepath.Dir(resolved)
}

// installShellCompletion detects the shell and installs completions.
func installShellCompletion(dryRun bool) {
	shell := completion.DetectShell()
	fmt.Printf("\n  %s%s %s%s\n", constants.ColorYellow, constants.SetupSectionComp, shell, constants.ColorReset)

	if dryRun {
		fmt.Printf("  %s[dry-run]%s would install %s completion\n",
			constants.ColorDim, constants.ColorReset, shell)

		return
	}

	err := completion.Install(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s%s\n", constants.ColorYellow, err, constants.ColorReset)

		return
	}

	fmt.Fprintf(os.Stderr, constants.MsgCompInstalled, shell)
}

// installCDFunction detects the shell and installs the gcd wrapper.
func installCDFunction(dryRun bool) {
	shell := completion.DetectShell()
	fmt.Printf("\n  %s%s %s%s\n", constants.ColorYellow, "cd function:", shell, constants.ColorReset)

	if dryRun {
		fmt.Printf("  %s[dry-run]%s would install gcd function for %s\n",
			constants.ColorDim, constants.ColorReset, shell)

		return
	}

	err := completion.InstallCDFunction(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s%s\n", constants.ColorYellow, err, constants.ColorReset)
	}
}

// ensureGitignoreStep adds release-related paths to .gitignore during setup.
func ensureGitignoreStep(dryRun bool) {
	fmt.Printf("\n  %s■ Gitignore —%s\n", constants.ColorYellow, constants.ColorReset)

	if dryRun {
		fmt.Printf("  %s[dry-run]%s would ensure release paths are in .gitignore\n",
			constants.ColorDim, constants.ColorReset)

		return
	}

	release.EnsureGitignore()
	fmt.Printf("  %s✓%s Release paths verified in .gitignore\n", constants.ColorGreen, constants.ColorReset)
}

// parseSetupFlags parses flags for the setup command.
func parseSetupFlags(args []string) (configPath string, dryRun, hasConfig bool) {
	fs := flag.NewFlagSet(constants.CmdSetup, flag.ExitOnError)
	cfgFlag := fs.String("config", constants.DefaultSetupConfigPath, constants.FlagDescSetupConfig)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	fs.Parse(args)
	fs.Visit(func(f *flag.Flag) {
		hasConfig = hasConfig || f.Name == "config"
	})

	return *cfgFlag, *dryRunFlag, hasConfig
}

// printSetupBanner shows the setup header.
func printSetupBanner(dryRun bool) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.SetupBannerTop, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.SetupBannerTitle, constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorCyan, constants.SetupBannerBottom, constants.ColorReset)
	if dryRun {
		fmt.Printf("\n  %s%s%s\n", constants.ColorYellow, constants.SetupDryRunFmt, constants.ColorReset)
	}
}

// printSetupSummary shows the final results.
func printSetupSummary(r setup.SetupResult) {
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermSeparator, constants.ColorReset)
	_, _ = filepath.Abs(".")
	printSetupCounts(r)
	printSetupErrors(r)
	fmt.Println()
}

// printSetupCounts prints applied/skipped/failed counts.
func printSetupCounts(r setup.SetupResult) {
	if r.Applied > 0 {
		fmt.Printf("  %s"+constants.SetupAppliedFmt+"%s\n", constants.ColorGreen, r.Applied, constants.ColorReset)
	}
	if r.Skipped > 0 {
		fmt.Printf("  %s"+constants.SetupSkippedFmt+"%s\n", constants.ColorDim, r.Skipped, constants.ColorReset)
	}
	if r.Failed > 0 {
		fmt.Printf("  %s"+constants.SetupFailedFmt+"%s\n", constants.ColorYellow, r.Failed, constants.ColorReset)
	}
}

// printSetupErrors prints each failed setting detail.
func printSetupErrors(r setup.SetupResult) {
	if r.Failed == 0 {
		return
	}
	for _, e := range r.Errors {
		fmt.Printf("    %s"+constants.SetupErrorEntryFmt+"%s\n", constants.ColorYellow, e, constants.ColorReset)
	}
}
