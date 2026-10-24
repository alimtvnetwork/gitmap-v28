package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvmware"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runOS dispatches gitmap os subcommands with strict UX ordering.
func runOS(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) {
		renderOSHelpWithIdentityFooter()

		return nil
	}

	if len(args) == 1 && hasOSJSONArg(args) {
		return renderSingleMachineIdentity(resolveLocalMachineIdentity(), true)
	}

	return dispatchOSSubcommand(strings.ToLower(args[0]), args[1:])
}

// renderOSHelpWithIdentityFooter enforces the canonical 3-part layout:
// 1. OS Commands Menu
// 2. Active Host Machine Identity
// 3. Real-World Sample Commands at the bottom
func renderOSHelpWithIdentityFooter() {
	RenderModernOSHelpMenuOnly()
	renderOSMachineIdentityFooter()
	PrintCrossPlatformOSExamples()
}

func hasOSJSONArg(args []string) bool {
	for _, arg := range args {
		normalized := strings.ToLower(strings.TrimSpace(arg))
		if normalized == "--json" || normalized == "-j" || normalized == "json" {
			return true
		}
	}
	return false
}

func renderOSMachineIdentityFooter() {
	id := resolveLocalMachineIdentity()
	fmt.Printf("\n%s  ● Active Host Identity:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    ├─ %-16s%s\n", "Machine Name:", id.MachineName)
	fmt.Printf("    ├─ %-16s%s\n", "Fleet Alias:", id.Alias)
	fmt.Printf("    ├─ %-16s%s\n", "Local IPv4:", id.IPAddress)
	fmt.Printf("    ├─ %-16s%s\n", "OS Platform:", id.OSPlatform)
	archDesc := formatArchWithCores(id.Architecture, id.CPUCores)
	fmt.Printf("    ├─ %-16s%s\n", "Architecture:", archDesc)
	fmt.Printf("    └─ %-16s%s\n\n", "Active User:", id.CurrentUser)
}

func isOSHelpArg(arg string) bool {
	return argMatches(arg, "-h", "--help", "help")
}

func dispatchOSSubcommand(subCmd string, subArgs []string) error {
	switch subCmd {
	case "--json", "-j", "json":
		return renderSingleMachineIdentity(resolveLocalMachineIdentity(), true)
	case "network", "net", "adapter", "adapters", "interface", "interfaces":
		return RunOSNetworkCLI(subArgs)
	case "machine", "machines", "machine-name", "hostname":
		return RunMachineCLI(subArgs)
	case "alias", "aliases":
		return RunAliasCLI(subArgs)
	case constants.SubCmdOSDisplay, constants.SubCmdOSDisplayAlias, constants.SubCmdOSDisplayAlias2:
		return runOSDisplay(subArgs)
	case constants.SubCmdFixLink, constants.SubCmdFixLinkAlias, constants.SubCmdFixLinkAlias2:
		return runOSFixLink(subArgs)
	case "info", "os-info", "osinfo", "sysinfo", "system-info":
		return RunOSInfoCLI(subArgs)
	case constants.SubCmdOSStatus, "st":
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
	case "devtool", "devtools", "dev-tool", "dev-tools", "dt":
		return RunDevToolCLI(subArgs)
	case "dev-clean", "dev-cleanup", "cleandev", "devcleanup", "clean-dev", "devtools-cache", "dev-tools-cache":
		return RunOSDevClean(subArgs)
	case "terminal", "terminal-clear", "clear-terminal", "clean-terminal":
		return RunTerminalCleanCLI(subArgs)
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
	case "cron", "crontab", "schedule", "schedules", "sched":
		return runOSCron(subArgs)
	case "storage", "disk", "space":
		return runOSStorage(subArgs)
	case "autologin", "auto-login", "al":
		return runOSAutoLogin(subArgs)
	case "change-password", "passwd", "chpasswd", "password":
		return RunChangePasswordCLI(subArgs)
	case "tweak", "tweaks", "twk":
		return runOSTweak(subArgs)
	case "dm", "display-manager":
		return runOSDMCommand(subArgs)
	case "dns", "nameserver":
		return runOSDNSCommand(subArgs)
	case "theme", "colorscheme":
		return runOSThemeCommand(subArgs)
	case "dock", "panel", "start-menu", "dock-position", "taskbar":
		return RunOSDockCommand(subArgs)
	case "update", "upgrade":
		return runOSUpdateCommand(subCmd == "upgrade", subArgs)
	case "tui", "menu", "gui", "dashboard":
		return RunOSTUICommand(subArgs)
	case constants.SubCmdOSHelp:
		printOSUsage()
		return nil
	default:
		return unknownOSSubcommandError(subCmd)
	}
}

func unknownOSSubcommandError(subCmd string) error {
	return apperror.NewSimple("unknown os subcommand '"+subCmd+"' (see 'gitmap os --help')", "E_INVALID_OS_SUBCMD")
}
