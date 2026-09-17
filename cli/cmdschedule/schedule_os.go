package cmdschedule

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// OSActionType identifies the power operation to schedule.
type OSActionType string

const (
	OSActionShutdown OSActionType = "shutdown"
	OSActionRestart  OSActionType = "restart"
)

// OSActionParams specifies the execution context for an OS power action.
type OSActionParams struct {
	Action     OSActionType
	Delay      time.Duration
	Seconds    int64
	Executable string
	Args       []string
	OnTick     func(remaining time.Duration)
}

// OSActionExecutor executes the configured OS power action.
type OSActionExecutor func(params OSActionParams) error

// DefaultOSActionExecutor is the injectable power command executor.
var DefaultOSActionExecutor OSActionExecutor = executeNativeOSAction

// BuildOSActionCommand constructs the OS command line for action and countdown seconds.
func BuildOSActionCommand(action OSActionType, seconds int64) (string, []string) {
	switch runtime.GOOS {
	case "windows":
		return buildWindowsActionCommand(action, seconds)
	case "darwin":
		return buildDarwinActionCommand(action)
	default:
		return buildLinuxActionCommand(action)
	}
}

func buildWindowsActionCommand(action OSActionType, seconds int64) (string, []string) {
	flag := "/s"
	if action == OSActionRestart {
		flag = "/r"
	}
	comment := "GitMap scheduled " + string(action)
	return "shutdown", []string{flag, "/t", fmt.Sprintf("%d", seconds), "/c", comment}
}

func buildDarwinActionCommand(action OSActionType) (string, []string) {
	flag := "-h"
	if action == OSActionRestart {
		flag = "-r"
	}
	return "sudo", []string{"shutdown", flag, "now"}
}

func buildLinuxActionCommand(action OSActionType) (string, []string) {
	sub := "poweroff"
	if action == OSActionRestart {
		sub = "reboot"
	}
	return "sudo", []string{sub}
}

func executeNativeOSAction(params OSActionParams) error {
	performDelayCountdown(params)
	if params.Executable == "" {
		return nil
	}
	cmd := exec.Command(params.Executable, params.Args...)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "execute "+string(params.Action))
	}
	return nil
}

func performDelayCountdown(params OSActionParams) {
	if params.Delay <= 0 || runtime.GOOS == "windows" {
		return
	}
	remaining := params.Delay
	for remaining > 0 {
		if params.OnTick != nil {
			params.OnTick(remaining)
		}
		step := resolveCountdownStep(remaining)
		time.Sleep(step)
		remaining -= step
	}
}

func resolveCountdownStep(remaining time.Duration) time.Duration {
	if remaining < time.Second {
		return remaining
	}
	return time.Second
}

// RunSchedulePowerCLI handles the schedule power subcommands from the CLI.
func RunSchedulePowerCLI(action OSActionType, args []string) error {
	if isHelpArg(args) {
		return printSchedulePowerHelp(action)
	}
	dur, err := resolveSchedulePowerDuration(args)
	if err != nil {
		return err
	}
	params := buildPowerActionParams(action, dur)
	if execErr := DefaultOSActionExecutor(params); execErr != nil {
		return apperror.WrapSimple(execErr, "schedule "+string(action))
	}
	printPowerScheduleConfirmation(action, dur)
	return nil
}

func isHelpArg(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}
	return false
}

func printSchedulePowerHelp(action OSActionType) error {
	fmt.Printf("Usage: gitmap schedule %s [duration]\n\n", action)
	fmt.Printf("Schedules a system %s with flexible duration formats.\n\n", action)
	fmt.Println("Supported duration formats:")
	fmt.Println("  1:45hr, 1:45h, 01:30:00   Colon notation (hours:minutes[:seconds])")
	fmt.Println("  1day, 2d                  Day notation")
	fmt.Println("  2h, 120m, 1s              Hour, minute, second notation")
	fmt.Println("  now, 0                    Execute immediately")
	return nil
}

func resolveSchedulePowerDuration(args []string) (time.Duration, error) {
	durStr := "0"
	if len(args) > 0 {
		durStr = args[0]
	}
	res := ParseScheduleDuration(durStr)
	if res.IsFailure() {
		return 0, res.AppError()
	}
	return res.Value, nil
}

func buildPowerActionParams(action OSActionType, dur time.Duration) OSActionParams {
	secs := int64(dur / time.Second)
	exe, cmdArgs := BuildOSActionCommand(action, secs)
	return OSActionParams{
		Action:     action,
		Delay:      dur,
		Seconds:    secs,
		Executable: exe,
		Args:       cmdArgs,
	}
}

func printPowerScheduleConfirmation(action OSActionType, dur time.Duration) {
	secs := int64(dur / time.Second)
	fmt.Printf("✓ Scheduled system %s in %v (%ds)\n", action, dur, secs)
}

// runScheduleRestart runs native OS restart without arguments.
func runScheduleRestart() error {
	return RunSchedulePowerCLI(OSActionRestart, nil)
}

// runScheduleShutdown runs native OS shutdown without arguments.
func runScheduleShutdown() error {
	return RunSchedulePowerCLI(OSActionShutdown, nil)
}
