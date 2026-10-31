// Package cmdcontextmenu - contextmenu_help.go: status and help printers for context-menu.
// Split from contextmenu.go on 31-Oct-26 so contextmenu.go stays under the 300-line file-size gate.
package cmdcontextmenu

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

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
