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
	case "add", "create", "new":
		return runSSHMacroAddCLI(args)
	case "rm", "remove", "delete":
		return runSSHMacroRmCLI(args)
	case "ls", "list":
		return runSSHMacroLsCLI(args)
	case "edit", "modify":
		return runSSHMacroEditCLI(args)
	case "run", "exec":
		return runSSHMacroRunCLI(args)
	case "deploy":
		return runSSHMacroDeployCLI(args)
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
	fmt.Println("  Create locally, verify, and deploy macros across the SSH fleet.")
	fmt.Println()
}

func printSSHMacroHelpUsage() {
	fmt.Println("Usage:")
	fmt.Println("  gitmap ssh macro add <name> <steps...>     Create a macro locally for SSH fleet")
	fmt.Println("  gitmap ssh macro rm <name>                 Delete a macro locally")
	fmt.Println("  gitmap ssh macro ls                        List all saved macros")
	fmt.Println("  gitmap ssh macro edit <name>               Edit macro steps locally")
	fmt.Println("  gitmap ssh macro run <name> [flags]        Run macro locally or across remote nodes")
	fmt.Println("  gitmap ssh macro deploy <name> [flags]     Deploy macro JSON to remote fleet nodes")
	fmt.Println("  gitmap ssh macro sync [name|all] [flags]   Push macros to remote SSH nodes")
	fmt.Println("  gitmap ssh macro export <name> [flags]     Pull macro from remote node to local")
	fmt.Println("  gitmap ssh macro import <file> [flags]     Push macro file to remote nodes")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: all online)")
	fmt.Println("      --except string     Exclude machines by alias, IP, or ID")
	fmt.Println("      --deploy            Deploy before remote running")
	fmt.Println("  -f, --force             Overwrite existing macros")
	fmt.Println("      --dry-run           Simulate action without writing files")
	fmt.Println()
}
