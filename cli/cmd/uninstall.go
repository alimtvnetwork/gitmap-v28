package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwinutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type uninstallFlags struct {
	isDryRun        bool
	isForce         bool
	isPurge         bool
	isPurgeWebView2 bool
	backupPath      string
}

func parseUninstallFlags(args []string) (*flag.FlagSet, *uninstallFlags) {
	fs := flag.NewFlagSet("uninstall", flag.ExitOnError)
	f := &uninstallFlags{}
	bindCoreUninstallFlags(fs, f)
	bindPurgeAndBackupFlags(fs, f)
	fs.Parse(reorderFlagsBeforeArgs(args))

	return fs, f
}

func bindCoreUninstallFlags(fs *flag.FlagSet, f *uninstallFlags) {
	fs.BoolVar(&f.isDryRun, constants.FlagUninstallDryRun, false, constants.FlagDescUninstallDryRun)
	fs.BoolVar(&f.isDryRun, "n", false, "Dry run alias")
	fs.BoolVar(&f.isForce, constants.FlagUninstallForce, false, constants.FlagDescUninstallForce)
	fs.BoolVar(&f.isForce, "yes", false, "Force alias")
	fs.BoolVar(&f.isForce, "y", false, "Force alias")
}

func bindPurgeAndBackupFlags(fs *flag.FlagSet, f *uninstallFlags) {
	fs.BoolVar(&f.isPurge, constants.FlagUninstallPurge, false, constants.FlagDescUninstallPurge)
	fs.BoolVar(&f.isPurge, "all", false, "Purge all files and configurations")
	fs.BoolVar(&f.isPurge, "a", false, "Purge alias")
	fs.BoolVar(&f.isPurgeWebView2, "purge-webview2", false, "Purge WebView2 runtime when removing Edge")
	fs.StringVar(&f.backupPath, "backup", "", "Backup file path for state snapshot")
	fs.StringVar(&f.backupPath, "b", "", "Backup file path alias")
}

func isSelfUninstallTool(tool string) bool {
	return tool == "" || tool == "gitmap" || tool == "gitmap-cli"
}

func isCtxTool(tool, canonical string) bool {
	return tool == constants.ToolCtx || canonical == constants.ToolCtx
}

// runUninstall handles the "uninstall" command.
func runUninstall(args []string) error {
	checkHelp("uninstall", args)
	if !hasPositionalToolArg(args) {
		runSelfUninstall(args)

		return nil
	}

	return executeUninstallFlow(args)
}

func executeUninstallFlow(args []string) error {
	fs, flags := parseUninstallFlags(args)
	tool := fs.Arg(0)
	if strings.EqualFold(tool, "app") || strings.EqualFold(tool, "apps") {
		return runAppsUninstall(fs.Args()[1:])
	}
	if isSelfUninstallTool(tool) {
		runSelfUninstall(args)

		return nil
	}

	return processToolUninstall(tool, flags)
}

func processToolUninstall(tool string, flags *uninstallFlags) error {
	validateToolName(tool)
	canonical := resolveToolAlias(tool)
	if isCtxTool(tool, canonical) {
		runUninstallCtx()
		return nil
	}
	if handled, err := dispatchAgyAgmUninstall(tool, canonical, flags); handled {
		return err
	}
	if handled, err := dispatchWinUtilToolUninstall(tool, canonical, flags); handled {
		return err
	}
	return executeValidatedUninstall(tool, canonical, flags)
}

func dispatchAgyAgmUninstall(tool, canonical string, flags *uninstallFlags) (bool, error) {
	if isAgyUninstallTarget(tool, canonical) {
		return true, runAgyUninstallFlow(tool, canonical, flags)
	}
	if isAgmUninstallTarget(tool, canonical) {
		return true, runAgmUninstallFlow(tool, canonical, flags)
	}
	return false, nil
}

func dispatchWinUtilToolUninstall(tool, canonical string, flags *uninstallFlags) (bool, error) {
	if isCopilotUninstallTarget(tool, canonical) {
		return true, runCopilotUninstallFlow(flags)
	}
	if isEdgeUninstallTarget(tool, canonical) {
		return true, runEdgeUninstallFlow(tool, flags)
	}
	return false, nil
}

func isAgyUninstallTarget(tool, canonical string) bool {
	return strings.EqualFold(tool, "agy") || strings.EqualFold(tool, "agy-all") ||
		strings.EqualFold(tool, "antigravity") || strings.EqualFold(tool, "antigravity-all") ||
		strings.EqualFold(canonical, constants.ToolAgy) || strings.EqualFold(canonical, constants.ToolAntigravity)
}

func isAgmUninstallTarget(tool, canonical string) bool {
	return strings.EqualFold(tool, "agm") || strings.EqualFold(tool, "agm-all") ||
		strings.EqualFold(tool, "ag-manager") || strings.EqualFold(tool, "ag-manager-all") ||
		strings.EqualFold(tool, "antigravity-manager") || strings.EqualFold(tool, "antigravity-manager-all") ||
		strings.EqualFold(canonical, constants.ToolAgManager)
}

func runAgyUninstallFlow(tool, canonical string, flags *uninstallFlags) error {
	isFullPurge := flags.isPurge || strings.EqualFold(tool, "agy-all") || strings.EqualFold(tool, "antigravity-all")
	return cmdagy.RunAGYUninstall(isFullPurge, flags.isForce, flags.isDryRun, flags.backupPath)
}

func runAgmUninstallFlow(tool, canonical string, flags *uninstallFlags) error {
	isFullPurge := flags.isPurge || strings.EqualFold(tool, "agm-all") || strings.EqualFold(tool, "ag-manager-all") || strings.EqualFold(tool, "antigravity-manager-all")
	return cmdinstall.RunAGMUninstall(isFullPurge, flags.isForce, flags.isDryRun)
}

func isCopilotUninstallTarget(tool, canonical string) bool {
	return strings.EqualFold(tool, "copilot") || strings.EqualFold(tool, "copilot-all") ||
		strings.EqualFold(canonical, "copilot") || strings.EqualFold(tool, "windows-copilot")
}

func isEdgeUninstallTarget(tool, canonical string) bool {
	return strings.EqualFold(tool, "edge") || strings.EqualFold(tool, "edge-all") ||
		strings.EqualFold(canonical, "edge") || strings.EqualFold(tool, "msedge")
}

func runCopilotUninstallFlow(flags *uninstallFlags) error {
	return cmdwinutil.RunWinUtilCopilotCLI(buildWinUtilFlags(flags))
}

func runEdgeUninstallFlow(tool string, flags *uninstallFlags) error {
	return cmdwinutil.RunWinUtilEdgeCLI(buildWinUtilEdgeFlags(tool, flags))
}

func buildWinUtilFlags(flags *uninstallFlags) []string {
	var args []string
	if flags.isDryRun {
		args = append(args, "--dry-run")
	}
	if flags.isForce {
		args = append(args, "--yes")
	}
	return args
}

func buildWinUtilEdgeFlags(tool string, flags *uninstallFlags) []string {
	args := buildWinUtilFlags(flags)
	if flags.isPurge || flags.isPurgeWebView2 || strings.EqualFold(tool, "edge-all") {
		args = append(args, "--purge-webview2")
	}
	return args
}

func executeValidatedUninstall(tool, canonical string, flags *uninstallFlags) error {
	db, _ := openDB()
	if db != nil {
		defer db.Close()
	}

	if err := checkUninstallEligibility(db, tool, canonical, flags.isForce); err != nil {
		return err
	}

	if !flags.isForce && !confirmUninstall(tool) {
		return nil
	}

	return dispatchUninstallExecution(db, tool, canonical, flags.isDryRun, flags.isPurge)
}

func executeStandardUninstall(db *store.DB, tool string, isDryRun, isPurge bool) error {
	manager := resolveUninstallManager(db, tool)
	uninstallCmd := buildUninstallCommand(manager, tool, isPurge)
	if isDryRun {
		fmt.Printf(constants.MsgUninstallDryCmd, strings.Join(uninstallCmd, " "))

		return nil
	}

	fmt.Printf(constants.MsgUninstallRemoving, tool)
	runInstallCommand(uninstallCmd, installOptions{Tool: tool, Verbose: true})
	removeToolFromDatabase(db, tool)

	fmt.Printf(constants.MsgUninstallSuccess, tool)

	return nil
}

func isAntigravityTool(tool string) bool {
	return tool == constants.ToolAntigravity || tool == constants.ToolAgy
}

func purgeDualAntigravityDatabases(db *store.DB) {
	_ = cmdinstall.PurgeAntigravityDualDatabase()
	if db == nil || db.Conn() == nil {
		return
	}

	_, _ = db.Conn().Exec(constants.SQLDeleteInstalledTool, constants.ToolAntigravity)
	_, _ = db.Conn().Exec(constants.SQLDeleteInstalledTool, constants.ToolAgy)
	_ = db.SyncKnownSplitDatabases()
}

func removeToolFromDatabase(db *store.DB, tool string) {
	if isAntigravityTool(tool) {
		purgeDualAntigravityDatabases(db)

		return
	}

	removeGenericToolFromDB(db, tool)
}

func removeGenericToolFromDB(db *store.DB, tool string) {
	if db == nil {
		return
	}

	if errRemove := db.RemoveInstalledTool(tool); errRemove != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUninstallDBRemove, tool, errRemove)
	}
}

// confirmUninstall prompts the user for confirmation.
func confirmUninstall(tool string) bool {
	fmt.Printf(constants.MsgUninstallConfirm, tool)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	trimmed := strings.TrimSpace(input)

	return strings.EqualFold(trimmed, "y") || strings.EqualFold(trimmed, "yes")
}

func isShellModeFlag(a string) bool {
	return isFlagToken(a) && (a == "--shell-mode" || a == "-shell-mode")
}

// hasPositionalToolArg reports whether args contain at least one non-flag token.
func hasPositionalToolArg(args []string) bool {
	hasSkipNext := false
	for _, a := range args {
		if hasSkipNext {
			hasSkipNext = false

			continue
		}

		if isShellModeFlag(a) {
			hasSkipNext = true

			continue
		}

		if !isFlagToken(a) {
			return true
		}
	}

	return false
}

// resolveUninstallManager determines which manager was used to install.
func resolveUninstallManager(db *store.DB, tool string) string {
	if db == nil {
		return resolvePackageManager("", tool)
	}

	record, err := db.GetInstalledTool(tool)
	if err != nil || record.PackageManager == "" {
		return resolvePackageManager("", tool)
	}

	return record.PackageManager
}

func buildPkgMgrCommand(manager, pkgName string, isPurge bool) []string {
	switch manager {
	case constants.PkgMgrChocolatey:
		return buildChocoUninstall(pkgName, isPurge)
	case constants.PkgMgrWinget:
		return []string{"winget", "uninstall", pkgName}
	case constants.PkgMgrApt:
		return buildAptUninstall(pkgName, isPurge)
	case constants.PkgMgrBrew:
		return []string{"brew", "uninstall", pkgName}
	case constants.PkgMgrSnap:
		return []string{"sudo", "snap", "remove", pkgName}
	}

	return buildChocoUninstall(pkgName, isPurge)
}

// buildUninstallCommand builds the uninstall command for a manager.
func buildUninstallCommand(manager, tool string, isPurge bool) []string {
	pkgName := resolvePackageName(tool, manager)

	return buildPkgMgrCommand(manager, pkgName, isPurge)
}

// buildChocoUninstall builds a Chocolatey uninstall command.
func buildChocoUninstall(pkg string, isPurge bool) []string {
	args := []string{"choco", "uninstall", pkg, "-y"}
	if isPurge {
		args = append(args, "-x")
	}

	return args
}

// buildAptUninstall builds an apt uninstall command.
func buildAptUninstall(pkg string, isPurge bool) []string {
	action := "remove"
	if isPurge {
		action = "purge"
	}

	return []string{"sudo", "apt", action, "-y", pkg}
}
