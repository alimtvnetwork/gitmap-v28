package cmdcontextmenu

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

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
	switch strings.ToLower(subCmd) {
	case "install", "in", "add":
		return runContextMenuInstall(rest)
	case "uninstall", "remove", "rm", "un":
		return runContextMenuUninstall()
	case "status", "check", "ls", "list":
		return runContextMenuStatus()
	case "wizard", "ui":
		return cmdinstall.RunInstallWizard()
	default:
		printContextMenuHelp()

		return nil
	}
}

func hasExplainFlag(args []string) bool {
	for _, a := range args {
		if a == "--explain" || a == "-explain" {
			return true
		}
	}

	return false
}

func runContextMenuInstall(args []string) error {
	hasExplain := hasExplainFlag(args)

	return cmdinstall.RunInstallCtx(hasExplain)
}

func runContextMenuUninstall() error {
	cmdinstall.RunUninstallCtx()

	return nil
}

func runContextMenuStatus() error {
	isActive := cmdinstall.IsContextMenuActive()
	if isActive {
		fmt.Printf("✓ %sGitMap context menu is active in the OS.%s\n", constants.ColorGreen, constants.ColorReset)
		fmt.Println("  Use 'gitmap context-menu uninstall' to remove it.")

		return nil
	}

	fmt.Printf("○ %sGitMap context menu is inactive.%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("  Use 'gitmap context-menu install' to register it.")

	return nil
}

func printContextMenuHelp() {
	fmt.Println("Usage: gitmap context-menu <command> [flags]")
	fmt.Println("\nManage OS right-click context menu integrations.")
	fmt.Println("\nCommands:")
	fmt.Println("  install       Register GitMap in Windows / Linux / macOS context menus")
	fmt.Println("  uninstall     Remove GitMap from OS context menus")
	fmt.Println("  status        Check whether context menu integration is active")
	fmt.Println("  wizard        Open the interactive Install Wizard overview")
	fmt.Println("\nFlags:")
	fmt.Println("  --explain     Print command templates before execution")
	fmt.Println("  --help, -h    Show this help message")
}
