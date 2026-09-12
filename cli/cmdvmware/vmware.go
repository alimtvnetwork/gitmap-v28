package cmdvmware

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runVmware dispatches gitmap vmware CLI commands.
func Run(args []string) error {
	checkHelp(constants.CmdVmware, args)
	if len(args) == 0 {
		printVmwareUsage()

		return nil
	}

	subCmd := strings.ToLower(args[0])

	return dispatchVmwareSubcommand(subCmd, args[1:])
}

func dispatchVmwareSubcommand(subCmd string, rest []string) error {
	switch subCmd {
	case constants.SubCmdVmwareShared:
		return runVmwareShared(rest)
	case "shared-enable", "mount":
		return runVmwareSharedEnable(rest)
	case constants.SubCmdSharedStatus:
		return runVmwareStatus(rest)
	case constants.CmdInstall, constants.CmdInstallAlias:
		return runVmwareInstall(rest)
	default:
		return unknownVmwareSubcommandError(subCmd)
	}
}

func unknownVmwareSubcommandError(subCmd string) error {
	msg := fmt.Sprintf("unknown vmware subcommand %q; see 'gitmap vmware --help'", subCmd)

	return apperror.NewWithDetails(
		"cmd.runVmware",
		"E4001",
		msg,
		"cmd.vmware",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		nil,
	)
}

func printVmwareUsage() {
	fmt.Println("Usage: gitmap vmware [subcommand] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  install (in)           Install open-vm-tools & desktop packages via apt")
	fmt.Println("  shared enable (mount)  Install open-vm-tools, mount /mnt/hgfs & persist in crontab")
	fmt.Println("  shared status          Check /mnt/hgfs mount, desktop symlink & crontab persistence")
	fmt.Println("  status                 Display VMware guest environment detection status")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -y, --yes              Auto-confirm package installation without prompting")
	fmt.Println("  -n, --dry-run          Show actions without mounting or writing crontab")
	fmt.Println("  -h, --help             Show this help message")
}
