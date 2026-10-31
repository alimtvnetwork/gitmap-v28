package cmdcontextmenu

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

type ctxSubcmdHandler func(args []string) error

// IDESetupRunner allows hooking IDE setup from root command dispatcher without heavy dependency imports.
var IDESetupRunner func(args []string) error

var ctxCmdDispatch = map[string]ctxSubcmdHandler{
	"install":            runContextMenuInstall,
	"in":                 runContextMenuInstall,
	"add":                runContextMenuInstall,
	"uninstall":          runUninstallAdapter,
	"remove":             runUninstallAdapter,
	"rm":                 runUninstallAdapter,
	"un":                 runUninstallAdapter,
	"status":             runStatusAdapter,
	"check":              runStatusAdapter,
	"ls":                 runStatusAdapter,
	"list":               runStatusAdapter,
	"clean-powershell":   runContextMenuCleanPowerShell,
	"clean-pwsh":         runContextMenuCleanPowerShell,
	"pwsh-clean":         runContextMenuCleanPowerShell,
	"customize":          runContextMenuCustomize,
	"custom":             runContextMenuCustomize,
	"wizard":             runWizardAdapter,
	"ui":                 runWizardAdapter,
	"ide-setup":          runContextMenuIDESetup,
	"setup-ide":          runContextMenuIDESetup,
	"cursor-setup":       runContextMenuIDESetup,
	"hide-git-actions":   runContextMenuHideGitActions,
	"hide-git":           runContextMenuHideGitActions,
	"show-git-actions":   runContextMenuShowGitActions,
	"show-git":           runContextMenuShowGitActions,
	"toggle-git-actions": runContextMenuToggleGitActions,
	"toggle-git":         runContextMenuToggleGitActions,
}

// RunContextMenu handles the 'context-menu' command and its subcommands.
func RunContextMenu(args []string) error {
	if isContextMenuHelp(args) {
		printContextMenuHelp()

		return nil
	}
	if len(args) == 0 {
		return runContextMenuStatus()
	}

	return dispatchContextMenuSubcommand(args[0], args[1:])
}

func isContextMenuHelp(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])

	return strutil.EqualFoldAny(first, "--help", "-h", "help")
}

func dispatchContextMenuSubcommand(subCmd string, rest []string) error {
	handler, hasHandler := ctxCmdDispatch[strings.ToLower(subCmd)]
	if !hasHandler {
		printContextMenuHelp()

		return nil
	}

	return handler(rest)
}

func runUninstallAdapter(_ []string) error {
	return runContextMenuUninstall()
}

func runStatusAdapter(_ []string) error {
	return runContextMenuStatus()
}

func runWizardAdapter(_ []string) error {
	return cmdinstall.RunInstallWizard()
}

func hasExplainFlag(args []string) bool {
	for _, a := range args {
		if a == "--explain" || a == "-explain" {
			return true
		}
	}

	return false
}

func hasRestoreFlag(args []string) bool {
	for _, a := range args {
		if strutil.EqualFoldAny(a, "--restore", "-restore") {
			return true
		}
	}

	return false
}

func hasDryRunFlag(args []string) bool {
	for _, a := range args {
		if strutil.EqualFoldAny(a, "--dry-run", "-dry-run", "--dryrun") {
			return true
		}
	}

	return false
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if strutil.EqualFoldAny(a, "--help", "-h", "help") {
			return true
		}
	}

	return false
}

func runContextMenuInstall(args []string) error {
	hasExplain := hasExplainFlag(args)
	if hasHideGitFlag(args) {
		_ = cmdinstall.SetGitActionsHidden(true)
	}

	return cmdinstall.RunInstallCtx(hasExplain)
}

func hasHideGitFlag(args []string) bool {
	for _, a := range args {
		if strutil.EqualFoldAny(a, "--hide-git", "--hide-git-actions", "-hide-git") {
			return true
		}
	}

	return false
}

func runContextMenuHideGitActions(_ []string) error {
	if err := cmdinstall.SetGitActionsHidden(true); err != nil {
		return err
	}
	fmt.Printf("✓ %sGit actions hidden from context menu.%s\n", constants.ColorGreen, constants.ColorReset)

	return cmdinstall.RunInstallCtx(false)
}

func runContextMenuShowGitActions(_ []string) error {
	if err := cmdinstall.SetGitActionsHidden(false); err != nil {
		return err
	}
	fmt.Printf("✓ %sGit actions enabled in context menu.%s\n", constants.ColorGreen, constants.ColorReset)

	return cmdinstall.RunInstallCtx(false)
}

func runContextMenuToggleGitActions(_ []string) error {
	isHidden := cmdinstall.IsGitActionsHidden()

	return dispatchToggleGitActions(!isHidden)
}

func dispatchToggleGitActions(hide bool) error {
	if hide {
		return runContextMenuHideGitActions(nil)
	}

	return runContextMenuShowGitActions(nil)
}

func runContextMenuUninstall() error {
	cmdinstall.RunUninstallCtx()

	return nil
}

func runContextMenuStatus() error {
	isActive := cmdinstall.IsContextMenuActive()
	if isActive {
		fmt.Printf("✓ %sGitMap context menu is active in the OS.%s\n", constants.ColorGreen, constants.ColorReset)
		printGitActionsVisibility()
		fmt.Println("  Use 'gitmap context-menu uninstall' to remove it.")

		return nil
	}

	fmt.Printf("○ %sGitMap context menu is inactive.%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("  Use 'gitmap context-menu install' to register it.")

	return nil
}

func printGitActionsVisibility() {
	if cmdinstall.IsGitActionsHidden() {
		fmt.Printf("  • Git Actions: %sHidden%s\n", constants.ColorYellow, constants.ColorReset)

		return
	}
	fmt.Printf("  • Git Actions: %sVisible%s\n", constants.ColorGreen, constants.ColorReset)
}

func runContextMenuCleanPowerShell(args []string) error {
	if hasHelpFlag(args) {
		printCleanPowerShellHelp()

		return nil
	}
	isRestore := hasRestoreFlag(args)
	isDryRun := hasDryRunFlag(args)

	return CleanDuplicatePowerShell(isRestore, isDryRun)
}

func runContextMenuCustomize(args []string) error {
	if hasHelpFlag(args) {
		printCustomizeHelp()

		return nil
	}

	return printCustomizeStatus()
}
