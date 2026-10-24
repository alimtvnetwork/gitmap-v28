package cmdschedule

import (
	"fmt"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule/enums"
)

// OSActionType is the canonical power operation type, defined in enums.
type OSActionType = enums.OSActionType

const (
	OSActionShutdown = enums.OSActionShutdown
	OSActionRestart  = enums.OSActionRestart
	OSActionCancel   = enums.OSActionCancel
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

var (
	DefaultOSActionExecutor OSActionExecutor = executeNativeOSAction
	countdownOS                              = runtime.GOOS
	countdownSleepFn                         = time.Sleep
)

func resolveTargetOS(targetOS ...string) string {
	if len(targetOS) > 0 && targetOS[0] != "" {
		return targetOS[0]
	}

	return runtime.GOOS
}

// BuildOSActionCommand constructs the OS command line for action and countdown seconds.
func BuildOSActionCommand(action OSActionType, seconds int64, targetOS ...string) (string, []string) {
	switch resolveTargetOS(targetOS...) {
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

	return "shutdown", []string{flag, "/t", fmt.Sprintf("%d", seconds), "/c", "GitMap scheduled " + action.String()}
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
		return apperror.WrapSimple(err, "execute "+params.Action.String())
	}

	return nil
}

func hasActivePowerSchedule() bool {
	active, err := GetActivePowerSchedule()
	if err != nil || active == nil {
		return false
	}

	return active.Status == "ARMED"
}

func performDelayCountdown(params OSActionParams) {
	if params.Action == OSActionCancel || params.Delay <= 0 || countdownOS == "windows" {
		return
	}

	runCountdownLoop(params)
}

func runCountdownLoop(params OSActionParams) {
	rem := params.Delay
	for rem > 0 && hasActivePowerSchedule() {
		if params.OnTick != nil {
			params.OnTick(rem)
		}
		step := resolveCountdownStep(rem)
		countdownSleepFn(step)
		rem -= step
	}
}

func resolveCountdownStep(remaining time.Duration) time.Duration {
	if remaining < time.Second {
		return remaining
	}

	return time.Second
}

func isPowerStatusSubcmd(arg string) bool {
	return slices.Contains([]string{"status", "ls", "check"}, arg)
}

func isPowerCancelSubcmd(arg string) bool {
	return slices.Contains([]string{"cancel", "abort", "stop", "rm", "remove"}, arg)
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

	return armAndDispatchPowerSchedule(action, dur)
}

func armAndDispatchPowerSchedule(action OSActionType, dur time.Duration) error {
	params := buildPowerActionParams(action, dur)
	recordPowerSchedule(action, dur, params.Args)
	if execErr := DefaultOSActionExecutor(params); execErr != nil {
		_ = CancelActivePowerSchedule()

		return apperror.WrapSimple(execErr, "schedule "+action.String())
	}

	printPowerScheduleConfirmation(action, dur)

	return nil
}

func recordPowerSchedule(action OSActionType, dur time.Duration, cmdArgs []string) {
	now := time.Now()
	_ = SavePowerScheduleState(PowerScheduleState{
		Action: action, Status: "ARMED", ScheduledAt: now, TriggerAt: now.Add(dur),
		DurationSeconds: int64(dur / time.Second), RawDuration: dur.String(),
		Command: strings.Join(cmdArgs, " "),
	})
}

// ShowSchedulePowerStatus prints the current status of a scheduled power action.
func ShowSchedulePowerStatus(action OSActionType) error {
	active, err := GetActivePowerSchedule()
	if err != nil {
		return err
	}

	return renderActivePowerOrNone(active, action)
}

func renderActivePowerOrNone(active *PowerScheduleState, action OSActionType) error {
	if active == nil {
		fmt.Printf("No scheduled %s is currently active.\n", action)

		return nil
	}

	renderActivePowerStatus(active)

	return nil
}

func renderActivePowerStatus(active *PowerScheduleState) {
	fmt.Printf("\n  \033[1;96m⚡ Active Power Schedule:\033[0m\n"+
		"    • Action:     \033[1m%s\033[0m\n"+
		"    • Status:     %s\n"+
		"    • Scheduled:  %s\n"+
		"    • Target:     %s\n"+
		"    • Remaining:  \033[1;92m%s\033[0m\n"+
		"    • To cancel:  gitmap schedule %s cancel\n\n",
		active.Action, active.Status, active.ScheduledAt.Format("2006-01-02 15:04:05"),
		active.TriggerAt.Format("2006-01-02 15:04:05"), FormatRemainingDuration(active.TriggerAt), active.Action)
}

// BuildCancelOSActionCommand returns the native command to abort an OS shutdown.
func BuildCancelOSActionCommand(targetOS ...string) (string, []string) {
	if resolveTargetOS(targetOS...) == "windows" {
		return "shutdown", []string{"/a"}
	}

	return "shutdown", []string{"-c"}
}

// CancelSchedulePowerCLI cancels any armed power schedule and invokes native abort.
func CancelSchedulePowerCLI(action OSActionType) error {
	_ = CancelActivePowerSchedule()
	exe, cmdArgs := BuildCancelOSActionCommand()
	params := OSActionParams{Action: OSActionCancel, Executable: exe, Args: cmdArgs}
	if err := DefaultOSActionExecutor(params); err != nil {
		return err
	}

	fmt.Printf("✓ Canceled active scheduled %s.\n", action)

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
	fmt.Printf("Usage: gitmap schedule %s [duration|status|cancel]\n\n"+
		"Examples: gitmap schedule %s 1:45hr | status | cancel\n", action, action)

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

	return OSActionParams{Action: action, Delay: dur, Seconds: secs, Executable: exe, Args: cmdArgs}
}

func printPowerScheduleConfirmation(action OSActionType, dur time.Duration) {
	target := time.Now().Add(dur).Format("15:04:05")
	fmt.Printf("✓ Scheduled system %s in %v (%ds). Target time: %s\n"+
		"  • Check status: gitmap schedule %s status\n"+
		"  • Cancel:       gitmap schedule %s cancel\n",
		action, dur, int64(dur/time.Second), target, action, action)
}
