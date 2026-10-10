package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdapps"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdselfinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwinutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type uninstallFlags struct {
	isDryRun        bool
	isForce         bool
	isPurge         bool
	isPurgeWebView2 bool
	backupPath      string
	node            string
}

func parseUninstallFlags(args []string) (*flag.FlagSet, *uninstallFlags) {
	fs := flag.NewFlagSet("uninstall", flag.ExitOnError)
	f := &uninstallFlags{}
	bindCoreUninstallFlags(fs, f)
	bindPurgeAndBackupFlags(fs, f)
	fs.Parse(flagutil.ReorderFlagsBeforeArgs(args))

	return fs, f
}

func bindCoreUninstallFlags(fs *flag.FlagSet, f *uninstallFlags) {
	fs.BoolVar(&f.isDryRun, constants.FlagUninstallDryRun, false, constants.FlagDescUninstallDryRun)
	fs.BoolVar(&f.isDryRun, "d", false, "Dry run alias")
	fs.BoolVar(&f.isDryRun, "n", false, "Dry run alias")
	fs.BoolVar(&f.isForce, constants.FlagUninstallForce, false, constants.FlagDescUninstallForce)
	fs.BoolVar(&f.isForce, "yes", false, "Force alias")
	fs.BoolVar(&f.isForce, "y", false, "Force alias")
	fs.BoolVar(&f.isForce, "f", false, "Force alias")
	fs.StringVar(&f.node, "node", "", "Remote node alias for SSH delegation")
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
	if hasNodeFlagInArgs(args) {
		return executeUninstallFlow(args)
	}
	if !hasPositionalToolArg(args) {
		cmdselfinstall.RunSelfUninstall(args)

		return nil
	}

	return executeUninstallFlow(args)
}

func hasNodeFlagInArgs(args []string) bool {
	for _, a := range args {
		if a == "--node" || a == "-node" || strings.HasPrefix(a, "--node=") || strings.HasPrefix(a, "-node=") {
			return true
		}
	}
	return false
}

func executeUninstallFlow(args []string) error {
	fs, flags := parseUninstallFlags(args)
	tool := fs.Arg(0)
	if flags.node != "" {
		return delegateRemoteUninstall(tool, flags)
	}
	if strings.EqualFold(tool, "app") || strings.EqualFold(tool, "apps") {
		return cmdapps.RunAppsUninstall(fs.Args()[1:])
	}
	if isSelfUninstallTool(tool) {
		cmdselfinstall.RunSelfUninstall(args)

		return nil
	}

	return processToolUninstall(tool, flags)
}

func delegateRemoteUninstall(tool string, flags *uninstallFlags) error {
	remoteCmd := buildRemoteUninstallCmd(tool, flags)
	fmt.Printf("● Delegating uninstall to remote node '%s'...\n", flags.node)
	return cmdssh.RunSSHExec([]string{flags.node, remoteCmd})
}

func buildRemoteUninstallCmd(tool string, flags *uninstallFlags) string {
	remoteCmd := "gitmap uninstall"
	if tool != "" {
		remoteCmd = fmt.Sprintf("gitmap uninstall %s", tool)
	}
	if flags.isDryRun {
		remoteCmd += " --dry-run"
	}
	if flags.isForce {
		remoteCmd += " --force"
	}
	if flags.isPurge {
		remoteCmd += " --purge"
	}
	return remoteCmd
}

func processToolUninstall(tool string, flags *uninstallFlags) error {
	canonical := cmdinstall.ResolveToolAlias(tool)
	if handled, err := dispatchDevToolUninstall(tool, canonical, flags); handled {
		return err
	}
	cmdinstall.ValidateToolName(tool)
	if isCtxTool(tool, canonical) {
		cmdinstall.RunUninstallCtx()
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

func dispatchDevToolUninstall(tool, canonical string, flags *uninstallFlags) (bool, error) {
	if isDevToolUninstallTarget(tool, canonical) {
		return true, cmdinstall.RunToolUninstall(canonical, flags.isDryRun, flags.isForce, flags.isPurge)
	}
	return false, nil
}

func isDevToolUninstallTarget(tool, canonical string) bool {
	return isPwshTarget(tool, canonical) || isAgmUninstallTarget(tool, canonical) ||
		isVimTarget(tool, canonical) || isVSCodeTarget(tool, canonical) || isCursorTarget(tool, canonical)
}

func isPwshTarget(tool, canonical string) bool {
	return strings.EqualFold(tool, "pwsh") || strings.EqualFold(tool, "powershell") ||
		strings.EqualFold(canonical, constants.ToolPowerShell)
}

func isVimTarget(tool, canonical string) bool {
	return strings.EqualFold(tool, "vim") || strings.EqualFold(tool, "vi") ||
		strings.EqualFold(canonical, "vim")
}

func isVSCodeTarget(tool, canonical string) bool {
	return strings.EqualFold(tool, "vscode") || strings.EqualFold(tool, "code") ||
		strings.EqualFold(tool, "vs-code") || strings.EqualFold(canonical, constants.ToolVSCode)
}

func isCursorTarget(tool, canonical string) bool {
	return strings.EqualFold(tool, "cursor") || strings.EqualFold(tool, "cur") ||
		strings.EqualFold(canonical, "cursor")
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
	cmdinstall.RunInstallCommand(uninstallCmd, cmdinstall.InstallOptions{Tool: tool, Verbose: true})
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

func isNodeFlag(a string) bool {
	return a == "--node" || a == "-node"
}

// hasPositionalToolArg reports whether args contain at least one non-flag token.
func hasPositionalToolArg(args []string) bool {
	hasSkipNext := false
	for _, a := range args {
		if hasSkipNext {
			hasSkipNext = false
			continue
		}
		if isShellModeFlag(a) || isNodeFlag(a) {
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
		return cmdinstall.ResolvePackageManager("", tool)
	}

	record, err := db.GetInstalledTool(tool)
	if err != nil || record.PackageManager == "" {
		return cmdinstall.ResolvePackageManager("", tool)
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
	pkgName := cmdinstall.ResolvePackageName(tool, manager)

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
