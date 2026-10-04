package cmdcursor

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func isCursorHelpRequested(args []string) bool {
	if len(args) == 0 {
		return true
	}
	arg := strings.ToLower(args[0])
	return arg == "help" || arg == "--help" || arg == "-h"
}

// RenderCursorHelp prints the Cursor CLI reference.
func RenderCursorHelp() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║            gitmap cursor / cur - Cursor IDE Integration          ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println()
	fmt.Println("  Usage:")
	fmt.Println("    gitmap cursor <subcommand> [arguments]")
	fmt.Println("    gitmap cur <subcommand> [arguments]")
	fmt.Println()
	fmt.Println("  Subcommands:")
	fmt.Println("    open [target]                 Open target repository or path in Cursor")
	fmt.Println("    sync                          Sync GitMap repos with Cursor Project Manager")
	fmt.Println("    list-projects, lp, ls         List projects registered in Cursor Project Manager")
	fmt.Println("    install                       Install Cursor or run Ubuntu provisioning script")
	fmt.Println("    help                          Show this help menu")
	fmt.Println()
}

func isTargetDirectory(arg string) bool {
	info, err := os.Stat(arg)
	return err == nil && info.IsDir()
}

func routeCursorSubcommand(sub string, args []string) error {
	switch strings.ToLower(sub) {
	case "open", "o":
		return RunCursorOpen(args[1:])
	case "sync", "s":
		return RunCursorSync(args[1:])
	case "list-projects", "lp", "ls", "list":
		return RunCursorListProjects(args[1:])
	case "install", "in", "i":
		return RunCursorInstall(args[1:])
	default:
		if isTargetDirectory(sub) {
			return RunCursorOpen(args)
		}
		RenderCursorHelp()
		return apperror.NewWithDetails("cmd.cursor.dispatch", "E1030", fmt.Sprintf("unknown cursor subcommand '%s'", sub), "cmdcursor", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
	}
}

// RunCursor handles `gitmap cursor` / `gitmap cur` execution.
func RunCursor(args []string) error {
	if isCursorHelpRequested(args) {
		RenderCursorHelp()
		return nil
	}
	return routeCursorSubcommand(args[0], args)
}
