package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunSSHMacroCLI executes the ssh macro command dispatcher.
func RunSSHMacroCLI(args []string) error {
	return runSSHMacroCLI(args)
}

func runSSHMacroCLI(args []string) error {
	if isSSHMacroHelpRequested(args) {
		printSSHMacroHelp()
		return nil
	}
	return dispatchSSHMacroSubcommand(args[0], args[1:])
}

func isSSHMacroHelpRequested(args []string) bool {
	if len(args) == 0 {
		return true
	}
	sub := args[0]
	return sub == "-h" || sub == "--help" || sub == "help"
}

func dispatchSSHMacroSubcommand(sub string, args []string) error {
	switch sub {
	case "sync":
		return runSSHMacroSyncCLI(args)
	case "export", "exp", "dump", "pull":
		return runSSHMacroExportCLI(args)
	case "import", "imp", "load", "push":
		return runSSHMacroImportCLI(args)
	default:
		printSSHMacroHelp()
		return nil
	}
}

func printSSHMacroHelp() {
	printSSHMacroHelpHeader()
	printSSHMacroHelpUsage()
}

func printSSHMacroHelpHeader() {
	fmt.Printf("\n  %sGitMap SSH Macro Engine%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Synchronize, export, and import macro workflows across SSH nodes.")
	fmt.Println()
}

func printSSHMacroHelpUsage() {
	fmt.Println("Usage:")
	fmt.Println("  gitmap ssh macro sync [name|all] [flags]   Push macros to remote SSH nodes")
	fmt.Println("  gitmap ssh macro export <name> [flags]    Pull macro from remote node to local")
	fmt.Println("  gitmap ssh macro import <file> [flags]    Push macro file to remote nodes")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: all online)")
	fmt.Println("      --except string     Exclude machines by alias, IP, or ID")
	fmt.Println("  -f, --force             Overwrite existing macros")
	fmt.Println("      --dry-run           Simulate action without writing files")
	fmt.Println()
}
