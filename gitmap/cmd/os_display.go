package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/power"
)

func runOSDisplay(args []string) error {
	if len(args) == 0 {
		return runOSDisplayStatus()
	}

	if isHelpArg(args[0]) {
		printOSDisplayUsage()

		return nil
	}

	subCmd := strings.ToLower(args[0])

	return dispatchOSDisplaySubcommand(subCmd, args[1:])
}

func dispatchOSDisplaySubcommand(subCmd string, subArgs []string) error {
	switch subCmd {
	case "status", "st", "info":
		return runOSDisplayStatus()
	case "never-sleep", "never", "ns", "off":
		return runPowerNeverSleep()
	case "set":
		return runPowerSet(subArgs)
	case "reset", "restore":
		return runPowerReset()
	default:
		return handleOSDisplayFallback(subCmd, subArgs)
	}
}

func handleOSDisplayFallback(subCmd string, subArgs []string) error {
	if _, err := strconv.Atoi(subCmd); err == nil {
		setArgs := append([]string{subCmd}, subArgs...)

		return runPowerSet(setArgs)
	}

	msg := fmt.Sprintf("unknown os display subcommand %q (see 'gitmap os display --help')", subCmd)

	return apperror.NewSimple(msg, "E_INVALID_OS_DISPLAY_SUBCMD")
}

func runOSDisplayStatus() error {
	mgr, err := power.NewManager()
	if err != nil {
		return err
	}

	current, err := mgr.GetStatus()
	if err != nil {
		return err
	}

	dbSetting := queryActivePowerSetting()
	printOSDisplayDetails(current, dbSetting)

	return nil
}

func queryActivePowerSetting() power.Settings {
	db, err := openDB()
	if err != nil || db == nil {
		return power.Settings{}
	}

	defer db.Close()
	setting, _ := db.GetActivePowerSetting()

	return setting
}

func detectDisplayServer() string {
	if runtime.GOOS == "windows" {
		return "Windows Desktop Window Manager (DWM)"
	}

	if runtime.GOOS == "darwin" {
		return "macOS Quartz / WindowServer"
	}

	return detectLinuxDisplayServer()
}

func detectLinuxDisplayServer() string {
	sessType := strings.ToLower(os.Getenv("XDG_SESSION_TYPE"))
	if sessType != "" {
		return sessType
	}

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "wayland"
	}

	if os.Getenv("DISPLAY") != "" {
		return "x11"
	}

	return "headless / tty"
}

func detectDesktopSession() string {
	if runtime.GOOS == "windows" {
		return "Windows Shell"
	}

	if runtime.GOOS == "darwin" {
		return "Aqua"
	}

	return detectLinuxDesktopSession()
}

func detectLinuxDesktopSession() string {
	desktop := os.Getenv("XDG_CURRENT_DESKTOP")
	if desktop != "" {
		return desktop
	}

	sess := os.Getenv("DESKTOP_SESSION")
	if sess != "" {
		return sess
	}

	return "none (terminal/headless)"
}

func printOSDisplayDetails(current, dbSetting power.Settings) {
	fmt.Printf("▶ OS Display & Screen Status (%s)\n", current.Platform)
	fmt.Printf("  • Operating System: %s (%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  • Display Server:   %s\n", detectDisplayServer())
	fmt.Printf("  • Desktop Session:  %s\n", detectDesktopSession())
	printOSDisplayPowerMetrics(current, dbSetting)
}

func printOSDisplayPowerMetrics(current, dbSetting power.Settings) {
	disp := formatDisplayTimeout(current.DisplayTimeoutMinutes, current.IsNeverSleep)
	sleep := formatDisplayTimeout(current.SleepTimeoutMinutes, current.IsNeverSleep)
	fmt.Printf("  • Display Timeout:  %s\n", disp)
	fmt.Printf("  • Sleep Timeout:    %s\n", sleep)
	printOSDisplayModeAndProfile(current.IsNeverSleep, dbSetting.Source)
}

func printOSDisplayModeAndProfile(isNeverSleep bool, profileSource string) {
	if isNeverSleep {
		fmt.Println("  • Mode:             Never-Sleep (inhibited)")
	} else {
		fmt.Println("  • Mode:             Standard timeouts")
	}

	if profileSource != "" {
		fmt.Printf("  • SQLite Profile:   %s\n", profileSource)
	}
}

func formatDisplayTimeout(minutes int, isNeverSleep bool) string {
	if isNeverSleep || minutes == 0 {
		return "Never (inhibited)"
	}

	return fmt.Sprintf("%d minutes", minutes)
}

const osDisplayUsageText = `Usage: gitmap os display [subcommand] [flags]

Commands:
  status (st)           Display desktop session, display server & screen blanking status
  never-sleep (never)   Disable screen blanking and sleep timeouts
  set <minutes>         Configure screen and display idle timeouts
  reset (restore)       Restore previous display timeout settings from snapshot
  help                  Show this help message

Examples:
  gitmap os display
  gitmap os display never-sleep
  gitmap os display set 15
  gitmap os display reset`

func printOSDisplayUsage() {
	fmt.Println(osDisplayUsageText)
}
