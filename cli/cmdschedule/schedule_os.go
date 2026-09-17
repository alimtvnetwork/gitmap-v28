package cmdschedule

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
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

func isPowerStatusSubcmd(arg string) bool {
	return arg == "status" || arg == "ls" || arg == "check"
}

func isPowerCancelSubcmd(arg string) bool {
	return arg == "cancel" || arg == "abort" || arg == "stop" || arg == "rm"
}

// RunSchedulePowerCLI handles the schedule power subcommands from the CLI.
func RunSchedulePowerCLI(action OSActionType, args []string) error {
	if isHelpArg(args) {
		return printSchedulePowerHelp(action)
	}
	if len(args) > 0 && isPowerStatusSubcmd(args[0]) {
		return ShowSchedulePowerStatus(action)
	}
	if len(args) > 0 && isPowerCancelSubcmd(args[0]) {
		return CancelSchedulePowerCLI(action)
	}
	return executeSchedulePowerCLI(action, args)
}

func executeSchedulePowerCLI(action OSActionType, args []string) error {
	dur, err := resolveSchedulePowerDuration(args)
	if err != nil {
		return err
	}
	params := buildPowerActionParams(action, dur)
	if execErr := DefaultOSActionExecutor(params); execErr != nil {
		return apperror.WrapSimple(execErr, "schedule "+string(action))
	}
	recordPowerSchedule(action, dur, params.Args)
	printPowerScheduleConfirmation(action, dur)
	return nil
}

func recordPowerSchedule(action OSActionType, dur time.Duration, cmdArgs []string) {
	now := time.Now()
	state := PowerScheduleState{
		Action:          action,
		Status:          "ARMED",
		ScheduledAt:     now,
		TriggerAt:       now.Add(dur),
		DurationSeconds: int64(dur / time.Second),
		Command:         strings.Join(cmdArgs, " "),
	}
	_ = SavePowerScheduleState(state)
}

// ShowSchedulePowerStatus prints the current status of a scheduled power action.
func ShowSchedulePowerStatus(action OSActionType) error {
	active, err := GetActivePowerSchedule()
	if err != nil {
		return err
	}
	if active == nil {
		fmt.Printf("No scheduled %s is currently active.\n", action)
		return nil
	}
	renderActivePowerStatus(active)
	return nil
}

func renderActivePowerStatus(active *PowerScheduleState) {
	rem := FormatRemainingDuration(active.TriggerAt)
	fmt.Printf("\n  \033[1;96m⚡ Active Power Schedule:\033[0m\n")
	fmt.Printf("    • Action:     \033[1m%s\033[0m\n", active.Action)
	fmt.Printf("    • Status:     %s\n", active.Status)
	fmt.Printf("    • Scheduled:  %s\n", active.ScheduledAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("    • Target:     %s\n", active.TriggerAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("    • Remaining:  \033[1;92m%s\033[0m\n", rem)
	fmt.Printf("    • To cancel:  gitmap schedule %s cancel\n\n", active.Action)
}

// BuildCancelOSActionCommand returns the native command to abort an OS shutdown.
func BuildCancelOSActionCommand() (string, []string) {
	if runtime.GOOS == "windows" {
		return "shutdown", []string{"/a"}
	}
	return "shutdown", []string{"-c"}
}

// CancelSchedulePowerCLI cancels any armed power schedule and invokes native abort.
func CancelSchedulePowerCLI(action OSActionType) error {
	_ = CancelActivePowerSchedule()
	exe, cmdArgs := BuildCancelOSActionCommand()
	cmd := exec.Command(exe, cmdArgs...)
	_ = cmd.Run()
	fmt.Printf("✓ Cancelled active scheduled %s.\n", action)
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
	fmt.Printf("Usage: gitmap schedule %s [duration|status|cancel]\n\n", action)
	fmt.Printf("Schedules a system %s with flexible duration formats.\n\n", action)
	fmt.Println("Supported duration formats:")
	fmt.Println("  1:45hr, 1:45h, 01:30:00   Colon notation (hours:minutes[:seconds])")
	fmt.Println("  1day, 1d, 2days, 2d       Day notation")
	fmt.Println("  2h, 120m, 1s              Hour, minute, second notation")
	fmt.Println("  now, 0                    Execute immediately")
	printSchedulePowerExamples(action)
	return nil
}

func printSchedulePowerExamples(action OSActionType) {
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  gitmap schedule %s 1:45hr\n", action)
	fmt.Printf("  gitmap schedule %s 2h\n", action)
	fmt.Printf("  gitmap schedule %s 120m\n", action)
	fmt.Printf("  gitmap schedule %s 1day\n", action)
	fmt.Printf("  gitmap schedule %s 1d\n", action)
	fmt.Printf("  gitmap schedule %s status    # View active timer and remaining time\n", action)
	fmt.Printf("  gitmap schedule %s cancel    # Abort active scheduled %s\n", action, action)
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
	target := time.Now().Add(dur).Format("15:04:05")
	fmt.Printf("✓ Scheduled system %s in %v (%ds). Target time: %s\n", action, dur, secs, target)
	fmt.Printf("  • Check status: gitmap schedule %s status\n", action)
	fmt.Printf("  • Cancel:       gitmap schedule %s cancel\n", action)
}

// runScheduleRestart runs native OS restart without arguments.
func runScheduleRestart() error {
	return RunSchedulePowerCLI(OSActionRestart, nil)
}

// runScheduleShutdown runs native OS shutdown without arguments.
func runScheduleShutdown() error {
	return RunSchedulePowerCLI(OSActionShutdown, nil)
}
