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
	if len(args) == 0 || isOSHelpArg(args[0]) {
		checkHelp("os", args)
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
	case constants.SubCmdOSDisplay, constants.SubCmdOSDisplayAlias, constants.SubCmdOSDisplayAlias2:
		return runOSDisplay(subArgs)
	case constants.SubCmdFixLink, constants.SubCmdFixLinkAlias, constants.SubCmdFixLinkAlias2:
		return runOSFixLink(subArgs)
	case constants.SubCmdOSStatus, "st", "info":
		return runOSStatus(subArgs)
	case constants.SubCmdOSHelp:
		return handleOSHelp()
	default:
		return unknownOSSubcommandError(subCmd)
	}
}

func handleOSHelp() error {
	printOSUsage()

	return nil
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

const osUsageText = `Usage: gitmap os [subcommand] [flags]

Commands:
  display (disp)      Inspect and configure OS display settings, resolution & timeouts
  fix-link (fixlink)  Inspect and repair broken symlinks and shared directories
  status (st)         Display operating system environment and link diagnostics
  help                Show this help message

Flags:
  --target <path>     Explicit target for symlink repair
  --force (-f)        Recreate symlinks even if target missing
  --recursive (-r)    Recursively inspect directories
  --dry-run (-n)      Inspect without modifying disk
  --json              Output as JSON`

func printOSUsage() {
	fmt.Println(osUsageText)
}

