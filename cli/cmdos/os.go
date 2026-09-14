package cmdos

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvmware"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	case constants.SubCmdOSIP:
		return runOSIP(subArgs)
	case constants.SubCmdOSZsh:
		return runOSZsh(subArgs)
	case "fix", "fixes":
		return runOSFix(subArgs)
	case "clean", "clear":
		return runOSClean(subArgs)
	case constants.SubCmdOSUser:
		return runOSUser(subArgs)
	case "vmware":
		return cmdvmware.Run(subArgs)
	case "group", "groups", "user-group", "usergroup":
		return runOSGroup(subArgs)
	case "cron", "crontab":
		return runOSCron(subArgs)
	case "storage", "disk", "space":
		return runOSStorage(subArgs)
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

	desktopDir := cmdvmware.ResolveUserDesktopDir()
	fmt.Printf("  • Desktop Dir:      %s\n", desktopDir)

	return inspectStandardLinksStatus(desktopDir)
}

const osUsageText = `Usage: gitmap os [subcommand] [flags]

Commands:
  ip                  Inspect, set, change, switch, or revert network IP configuration
  fix                 Register, edit, run, export, and import system repair scripts
  clean (clear)       Clean temporary and ephemeral system cache directories
  zsh                 Install, theme, switch, profile, and clean ZSH & Oh-My-Zsh
  user                Add, edit, export, import, or remove operating system users
  group (user-group)  List, create, edit, export, import, and remove user groups
  vmware              Discover and mount VMware shared folders (/mnt/hgfs)
  cron                Inspect, append, and remove crontab scheduled jobs
  storage (disk)      Inspect disk drive capacities, partitions, and storage metrics
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
