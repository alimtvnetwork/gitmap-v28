package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runOS dispatches gitmap os subcommands.
func runOS(args []string) error {
	checkHelp("os", args)

	if len(args) == 0 || isOSHelpArg(args[0]) {
		printOSUsage()

		return nil
	}

	subCmd := strings.ToLower(args[0])

	return dispatchOSSubcommand(subCmd, args[1:])
}

func isOSHelpArg(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func dispatchOSSubcommand(subCmd string, subArgs []string) error {
	switch subCmd {
	case constants.SubCmdFixLink, constants.SubCmdFixLinkAlias, constants.SubCmdFixLinkAlias2:
		return runOSFixLink(subArgs)
	case constants.SubCmdOSStatus, "st", "info":
		return runOSStatus(subArgs)
	case constants.SubCmdOSHelp:
		printOSUsage()

		return nil
	default:
		return unknownOSSubcommandError(subCmd)
	}
}

func unknownOSSubcommandError(subCmd string) error {
	msg := fmt.Sprintf("unknown os subcommand %q (see 'gitmap os --help')", subCmd)

	return apperror.NewSimple(msg, "E_INVALID_OS_SUBCMD")
}

func runOSStatus(args []string) error {
	checkHelp("os", args)
	fmt.Println("▶ gitmap os status")
	fmt.Printf("  • Operating System: %s (%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  • User Home:        %s\n", expandHome("~"))

	desktopDir := resolveUserDesktopDir()
	fmt.Printf("  • Desktop Dir:      %s\n", desktopDir)

	return inspectStandardLinksStatus(desktopDir)
}

func printOSUsage() {
	fmt.Println("Usage: gitmap os [subcommand] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  fix-link (fixlink)  Inspect and repair broken symlinks and shared directories")
	fmt.Println("  status (st)         Display operating system environment and link diagnostics")
	fmt.Println("  help                Show this help message")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --target <path>     Explicit target for symlink repair")
	fmt.Println("  --force (-f)        Recreate symlinks even if target missing")
	fmt.Println("  --recursive (-r)    Recursively inspect directories")
	fmt.Println("  --dry-run (-n)      Inspect without modifying disk")
	fmt.Println("  --json              Output as JSON")
}
