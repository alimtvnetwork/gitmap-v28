package cmdos

import (
	"fmt"
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
	case "dev":
		return runOSDevSubcommand(subArgs)
	case "dev-clean", "dev-cleanup", "cleandev", "devcleanup", "clean-dev":
		return RunOSDevClean(subArgs)
	case "cleanup":
		return runOSCleanup(subArgs)
	case "ai-clean", "aiclean", "clean-ai":
		return RunOSAICleanCLI(subArgs)
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
	case "autologin", "auto-login", "al":
		return runOSAutoLogin(subArgs)
	case "tweak", "tweaks", "twk":
		return runOSTweak(subArgs)
	case "dm", "display-manager":
		return runOSDMCommand(subArgs)
	case "dns", "nameserver":
		return runOSDNSCommand(subArgs)
	case "theme", "colorscheme":
		return runOSThemeCommand(subArgs)
	case "update":
		return runOSUpdateCommand(false, subArgs)
	case "upgrade":
		return runOSUpdateCommand(true, subArgs)
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
