package cmdos

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func runOSDMCommand(args []string) error {
	if len(args) == 0 {
		return runOSDMStatus()
	}

	subCmd := strings.ToLower(args[0])
	if isOSHelpArg(subCmd) {
		printOSDMHelp()

		return nil
	}

	return dispatchOSDMSubcommand(subCmd, args[1:])
}

func dispatchOSDMSubcommand(subCmd string, subArgs []string) error {
	switch subCmd {
	case "status", "st", "info":
		return runOSDMStatus()
	case "wayland":
		return handleOSDMWayland(subArgs)
	case "restart":
		return handleOSDMRestart()
	default:
		msg := fmt.Sprintf("unknown os dm subcommand %q (see 'gitmap os dm --help')", subCmd)

		return apperror.NewSimple(msg, "E_INVALID_OS_DM_SUBCMD")
	}
}

func runOSDMStatus() error {
	mgr := newLinuxDMManager()
	status, err := mgr.GetStatus()
	if err != nil {
		return err
	}

	renderOSDMStatus(status)

	return nil
}

func handleOSDMWayland(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("usage: gitmap os dm wayland <enable|disable>", "E_MISSING_ARG")
	}

	action := strings.ToLower(args[0])
	enableWayland := action == "enable" || action == "on" || action == "true"
	mgr := newLinuxDMManager()
	if err := mgr.SetWayland(enableWayland); err != nil {
		return err
	}

	state := "disabled (forcing X11)"
	if enableWayland {
		state = "enabled"
	}

	fmt.Printf("✔ Display Manager: Wayland %s successfully.\n", state)

	return nil
}

func handleOSDMRestart() error {
	mgr := newLinuxDMManager()
	if err := mgr.RestartService(); err != nil {
		return err
	}

	fmt.Println("✔ Display Manager service restarted successfully.")

	return nil
}
