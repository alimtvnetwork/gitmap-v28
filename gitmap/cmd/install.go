package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runInstall handles the "install" command.
func runInstall(args []string) {
	checkHelp("install", args)

	fs := flag.NewFlagSet("install", flag.ExitOnError)

	var manager, version string
	var verbose, dryRun, check, list, yes, explain bool

	fs.StringVar(&manager, constants.FlagInstallManager, "", constants.FlagDescInstallManager)
	fs.StringVar(&version, constants.FlagInstallVersion, "", constants.FlagDescInstallVersion)
	fs.BoolVar(&verbose, constants.FlagInstallVerbose, false, constants.FlagDescInstallVerbose)
	fs.BoolVar(&dryRun, constants.FlagInstallDryRun, false, constants.FlagDescInstallDryRun)
	fs.BoolVar(&check, constants.FlagInstallCheck, false, constants.FlagDescInstallCheck)
	fs.BoolVar(&list, constants.FlagInstallList, false, constants.FlagDescInstallList)
	fs.BoolVar(&yes, constants.FlagInstallYes, false, constants.FlagDescInstallYes)
	fs.BoolVar(&yes, "y", false, constants.FlagDescInstallYes)
	fs.BoolVar(&explain, constants.FlagInstallExplain, false, constants.FlagDescInstallExplain)

	reordered := reorderFlagsBeforeArgs(args)
	fs.Parse(reordered)

	if list {
		printInstallListGrouped()

		return
	}

	tool := fs.Arg(0)
	if tool == "" {
		fmt.Fprint(os.Stderr, constants.ErrInstallToolRequired)
		os.Exit(1)
	}

	validateToolName(tool)

	opts := installOptions{
		Tool:    tool,
		Manager: manager,
		Version: version,
		Verbose: verbose,
		DryRun:  dryRun,
		Check:   check,
		Yes:     yes,
		Explain: explain,
	}

	executeInstall(opts)
}

// installOptions holds parsed install flags.
type installOptions struct {
	Tool    string
	Manager string
	Version string
	Verbose bool
	DryRun  bool
	Check   bool
	Yes     bool
	Explain bool
}

// validateToolName checks if the tool is supported.
func validateToolName(tool string) {
	// Clean-code installer aliases (clean-code/code-guide/cg/cc) are not
	// in the standard InstallToolDescriptions map; they dispatch to a
	// dedicated PowerShell IRM | IEX flow. Allow them through.
	if isCleanCodeAlias(tool) {
		return
	}

	_, exists := constants.InstallToolDescriptions[tool]
	if exists {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrInstallUnknownTool, tool)
	os.Exit(1)
}

// executeInstall runs the install flow for a tool.
//
// Refactored to keep gocyclo below the 15 threshold:
//
//   - The 9-way "special tool" branch (clean-code, scripts, npp-settings,
//     vscode-sync, obs-sync, wt-sync, vscode-ctx, pwsh-ctx, all-dev-tools)
//     is now a single table lookup via specialInstallHandler().
//   - The "generic install" path (detect → check-only → confirm → install
//     → post-install npp settings) is broken into focused subroutines so
//     each function reads top-to-bottom with one responsibility.
func executeInstall(opts installOptions) {
	if handler := specialInstallHandler(opts.Tool); handler != nil {
		handler(opts)

		return
	}

	executeGenericInstall(opts)
}

// specialInstallHandler returns the dedicated handler for tools that
// bypass the generic detect → confirm → install pipeline (settings-only
// syncs, context-menu installers, and the meta "all-dev-tools" target).
// Returns nil for any tool that should follow the generic flow.
//
// Each entry takes installOptions even when most handlers ignore them, so
// the table is uniform and future handlers can opt-in to flag access
// (e.g. --dry-run, --yes) without re-shaping the dispatcher.
func specialInstallHandler(tool string) func(installOptions) {
	if isCleanCodeAlias(tool) {
		return func(installOptions) { runInstallCleanCode() }
	}

	switch tool {
	case constants.ToolScripts:
		return func(installOptions) { runInstallScripts() }
	case constants.ToolNppSettings:
		return func(installOptions) { runNppSettingsOnly() }
	case constants.ToolVSCodeSync:
		return func(installOptions) { runVSCodeSettingsOnly() }
	case constants.ToolOBSSync:
		return func(installOptions) { runOBSSettingsOnly() }
	case constants.ToolWTSync:
		return func(installOptions) { runWTSettingsOnly() }
	case constants.ToolVSCodeCtx:
		return func(installOptions) { runVSCodeContextMenu() }
	case constants.ToolPwshCtx:
		return func(installOptions) { runPwshContextMenu() }
	case constants.ToolCtx:
		return func(opts installOptions) { runInstallCtx(opts.Explain) }
	case constants.ToolAllDevTools:
		return func(opts installOptions) { runAllDevTools(opts) }
	case constants.ToolGitmapOneliner:
		return func(installOptions) { runInstallGitmapOneliner() }
	}

	return nil
}

// executeGenericInstall runs the standard package-manager-backed install
// pipeline: detect existing version → honor --check → resolve manager →
// announce → confirm → install → post-install settings sync.
func executeGenericInstall(opts installOptions) {
	originalTool := opts.Tool
	installName := resolveNppInstallName(opts.Tool)

	if alreadyInstalled(installName) {
		return
	}

	if opts.Check {
		fmt.Printf(constants.MsgInstallNotFound, installName)

		return
	}

	opts.Tool = installName
	manager := resolvePackageManager(opts.Manager)
	announceInstallPlan(opts.Version, manager)

	if !confirmInstallIfNeeded(opts, installName, manager) {
		return
	}

	installTool(opts)
	postInstallSettingsSync(originalTool)
}

// alreadyInstalled prints the "checking" line, probes the installed
// version, and returns true (with a "found" message) when the tool is
// already present. Returns false when the tool is missing and the
// caller should proceed with the install pipeline.
func alreadyInstalled(installName string) bool {
	fmt.Printf(constants.MsgInstallChecking, installName)

	existingVersion := detectInstalledVersion(installName)
	if existingVersion == "" {
		return false
	}

	fmt.Printf(constants.MsgInstallFound, installName, existingVersion)

	return true
}

// announceInstallPlan prints the version and manager that will be used,
// matching the previous two-line pre-confirmation banner.
func announceInstallPlan(version, manager string) {
	if version != "" {
		fmt.Printf(constants.MsgInstallVersion, version)
	} else {
		fmt.Print(constants.MsgInstallVersionLabel)
	}

	fmt.Printf(constants.MsgInstallManager, manager)
}

// confirmInstallIfNeeded honors the --yes and --dry-run flags. Returns
// true when the install should proceed, false when the user declined and
// the caller has already printed the abort message.
func confirmInstallIfNeeded(opts installOptions, installName, manager string) bool {
	if opts.Yes || opts.DryRun {
		return true
	}

	if confirmInstall(installName, opts.Version, manager) {
		return true
	}

	fmt.Print(constants.MsgInstallAborted)

	return false
}

// postInstallSettingsSync handles the npp-specific post-install behavior:
// `install npp` syncs Notepad++ settings after the binary is installed,
// while `install install-npp` explicitly skips the settings step and
// prints a hint to the user.
func postInstallSettingsSync(originalTool string) {
	switch originalTool {
	case constants.ToolNpp:
		runNppSettings()
	case constants.ToolNppInstall:
		fmt.Print(constants.MsgInstallNppSkipSet)
	}
}

// confirmInstall prompts the user for install confirmation.
func confirmInstall(tool, version, manager string) bool {
	if version != "" {
		fmt.Printf(constants.MsgInstallPrompt, tool, version, manager)
	} else {
		fmt.Printf(constants.MsgInstallPromptNoVer, tool, manager)
	}

	var answer string

	fmt.Scanln(&answer)

	return answer == "y" || answer == "Y"
}
