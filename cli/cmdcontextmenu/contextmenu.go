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
	"install":          runContextMenuInstall,
	"in":               runContextMenuInstall,
	"add":              runContextMenuInstall,
	"uninstall":        runUninstallAdapter,
	"remove":           runUninstallAdapter,
	"rm":               runUninstallAdapter,
	"un":               runUninstallAdapter,
	"status":           runStatusAdapter,
	"check":            runStatusAdapter,
	"ls":               runStatusAdapter,
	"list":             runStatusAdapter,
	"clean-powershell": runContextMenuCleanPowerShell,
	"clean-pwsh":       runContextMenuCleanPowerShell,
	"pwsh-clean":       runContextMenuCleanPowerShell,
	"customize":        runContextMenuCustomize,
	"custom":           runContextMenuCustomize,
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

func printCustomizeStatus() error {
	fmt.Println("GitMap Context Menu Customizer")
	fmt.Println("==============================")
	printIntegrationStatus()
	printPowerShellIntegrationStatus()
	printCustomizeActions()

	return nil
}

func printIntegrationStatus() {
	isActive := cmdinstall.IsContextMenuActive()
	if isActive {
		fmt.Printf("• GitMap Context Menu  : %sActive%s\n", constants.ColorGreen, constants.ColorReset)

		return
	}
	fmt.Printf("• GitMap Context Menu  : %sInactive%s\n", constants.ColorYellow, constants.ColorReset)
}

func printPowerShellIntegrationStatus() {
	st := DetectPowerShellStatus()
	if st.IsCleaned {
		fmt.Printf("• PowerShell Duplicates: %sCleaned (Legacy hidden)%s\n", constants.ColorGreen, constants.ColorReset)

		return
	}
	if st.IsDuplicate {
		fmt.Printf("• PowerShell Duplicates: %sDetected (Legacy + PS7 both visible)%s\n", constants.ColorYellow, constants.ColorReset)

		return
	}
	fmt.Printf("• PowerShell Duplicates: %sNone%s\n", constants.ColorCyan, constants.ColorReset)
}

func printCustomizeActions() {
	fmt.Println("\nAvailable Customizations:")
	fmt.Println("  gitmap context-menu clean-powershell          Hide legacy PowerShell duplicate")
	fmt.Println("  gitmap context-menu clean-powershell --restore Restore legacy PowerShell menu")
	fmt.Println("  gitmap context-menu hide-git-actions          Hide Git direct actions from context menu")
	fmt.Println("  gitmap context-menu show-git-actions          Show Git direct actions in context menu")
	fmt.Println("  gitmap context-menu install                   Register GitMap context menus")
	fmt.Println("  gitmap context-menu uninstall                 Remove GitMap context menus")
	fmt.Println("  gitmap context-menu wizard                    Open install wizard UI")
}

func printCustomizeHelp() {
	fmt.Println("Usage: gitmap context-menu customize [flags]")
	fmt.Println("\nView and customize OS context menu entries and third-party integrations.")
	fmt.Println("\nFlags:")
	fmt.Println("  --help, -h    Show this help message")
}

func runContextMenuIDESetup(args []string) error {
	if IDESetupRunner != nil {
		return IDESetupRunner(args)
	}

	return fmt.Errorf("IDE setup runner is not initialized")
}

func printCleanPowerShellHelp() {
	fmt.Println("Usage: gitmap context-menu clean-powershell [flags]")
	fmt.Println("\nSuppress redundant legacy PowerShell context menu when PowerShell 7 is installed.")
	fmt.Println("\nFlags:")
	fmt.Println("  --restore     Remove HKCU overrides and restore default menus")
	fmt.Println("  --dry-run     Preview registry changes without writing")
	fmt.Println("  --help, -h    Show this help message")
}

func printContextMenuHelp() {
	fmt.Println("Usage: gitmap context-menu <command> [flags]")
	fmt.Println("\nManage OS right-click context menu integrations.")
	fmt.Println("\nCommands:")
	fmt.Println("  install          Register GitMap in Windows / Linux / macOS context menus")
	fmt.Println("  uninstall        Remove GitMap from OS context menus")
	fmt.Println("  status           Check whether context menu integration is active")
	fmt.Println("  clean-powershell Clean duplicate legacy PowerShell context menu entries")
	fmt.Println("  hide-git-actions Hide Git direct actions from folder context menu")
	fmt.Println("  show-git-actions Restore Git direct actions in folder context menu")
	fmt.Println("  customize        Overview and configure context menu customizations")
	fmt.Println("  ide-setup        Configure Cursor-like IDEs (Cline, Continue, Zed)")
	fmt.Println("  wizard           Open the interactive Install Wizard overview")
	fmt.Println("\nFlags:")
	fmt.Println("  --hide-git       Omit Git direct actions during install")
	fmt.Println("  --explain        Print command templates before execution")
	fmt.Println("  --help, -h       Show this help message")
}
