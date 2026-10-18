package cmdlemon

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func runLemonStatusCmd(args []string) error {
	isJSON := hasFlag(args, "--json")
	var state DaemonState
	data, err := os.ReadFile(StateFilePath)
	if err == nil {
		_ = json.Unmarshal(data, &state)
	}
	if isJSON {
		out, _ := json.MarshalIndent(state, "", "  ")
		fmt.Println(string(out))
		return nil
	}
	fmt.Println("🍋 GitMap Lemon Daemon Status")
	fmt.Printf("  Running:      %v\n", state.IsRunning)
	fmt.Printf("  PID:          %d\n", state.Pid)
	fmt.Printf("  Poll Count:   %d\n", state.PollCount)
	fmt.Printf("  Last Polled:  %s\n", state.LastPollTime.Format(time.RFC3339))
	return nil
}

func runLemonServiceCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing service action: install, uninstall, start, stop, restart, status")
	}
	mgr := NewServiceManager()
	details, err := mgr.Execute(ServiceActionType(args[0]))
	if err != nil {
		return err
	}
	if details != nil {
		fmt.Printf("✔ Service action %q applied (name=%s, status=%s, platform=%s)\n",
			args[0], details.Name, details.StatusType, details.Platform)
	}
	return nil
}

func parseIntervalFlag(args []string) time.Duration {
	for i, arg := range args {
		if dur, ok := tryParseIntervalArg(args, i, arg); ok {
			return dur
		}
	}
	return DefaultPollingInterval
}

func tryParseIntervalArg(args []string, i int, arg string) (time.Duration, bool) {
	if (arg == "--interval" || arg == "-i") && i+1 < len(args) {
		return parseDurationOrFallback(args[i+1])
	}
	if strings.HasPrefix(arg, "--interval=") {
		return parseDurationOrFallback(strings.TrimPrefix(arg, "--interval="))
	}
	return 0, false
}

func parseDurationOrFallback(val string) (time.Duration, bool) {
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0, false
	}
	return d, true
}

func hasFlag(args []string, flags ...string) bool {
	for _, arg := range args {
		if containsFlag(flags, arg) {
			return true
		}
	}
	return false
}

func containsFlag(flags []string, target string) bool {
	for _, f := range flags {
		if f == target {
			return true
		}
	}
	return false
}
