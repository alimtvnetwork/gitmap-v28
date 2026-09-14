package cmdvmware

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// Run dispatches gitmap vmware CLI commands.
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
		return auditVmwareCommand("shared", func() error { return runVmwareShared(rest) })
	case "shared-enable", "mount":
		return auditVmwareCommand("shared-enable", func() error { return runVmwareSharedEnable(rest) })
	case "shared-status":
		return auditVmwareCommand("shared-status", func() error { return runVmwareStatus(rest) })
	case constants.CmdInstall, constants.CmdInstallAlias:
		return auditVmwareCommand("install", func() error { return runVmwareInstall(rest) })
	case constants.SubCmdSharedStatus:
		return auditVmwareCommand("status", func() error { return runVmwareStatus(rest) })
	default:
		return unknownVmwareSubcommandError(subCmd)
	}
}

func auditVmwareCommand(action string, runFn func() error) error {
	start := time.Now()
	runErr := runFn()
	durationMs := time.Since(start).Milliseconds()
	isSuccess := runErr == nil
	recordVmwareAudit(action, durationMs, isSuccess, runErr)

	return runErr
}

func recordVmwareAudit(action string, durationMs int64, isSuccess bool, runErr error) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return
	}

	defer splitDB.Close()

	exitCode, errMsg := resolveAuditErrorDetails(runErr)
	_ = splitDB.RecordExecution(
		"vmware", action, "", "system", durationMs,
		isSuccess, exitCode, "", errMsg, "gitmap vmware "+action, "", "",
	)
}

func resolveAuditErrorDetails(err error) (int, string) {
	if err == nil {
		return 0, ""
	}

	return 1, err.Error()
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
	fmt.Println("  install (in)           Install or verify VMware tools (apt on Linux, VMTools on Windows)")
	fmt.Println("  shared enable (mount)  Mount /mnt/hgfs (Linux) or create Desktop shortcut to UNC share (Windows)")
	fmt.Println("  shared status          Check /mnt/hgfs mount (Linux) or UNC share access (Windows)")
	fmt.Println("  status                 Display VMware hypervisor and guest detection status")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -y, --yes              Auto-confirm package installation without prompting")
	fmt.Println("  -n, --dry-run          Show actions without mounting or modifying system")
	fmt.Println("  -h, --help             Show this help message")
}
