package cmdagy

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// OSShutdownExecutorFn is an injectable cross-platform shutdown runner to prevent accidental OS shutdown during testing.
var OSShutdownExecutorFn = defaultCrossPlatformShutdown

// ExecuteSUGWatch monitors registered projects and executes OS shutdown when all turn green.
func ExecuteSUGWatch(projects []string, interval time.Duration, isDryRun bool, isOnce bool) error {
	RecordSUGWatchStarted(projects, interval)
	defer RecordSUGWatchStopped()

	printSUGStartBanner(projects, interval, isDryRun)
	for {

		pendingCount := evaluateAllProjectsStatus(projects)
		if pendingCount == 0 {
			return triggerSystemShutdown(len(projects), isDryRun)
		}
		if isOnce {
			printSUGOnceFinished(pendingCount)
			return nil
		}
		printSUGPendingWait(pendingCount, interval)
		time.Sleep(interval)
	}
}

func printSUGStartBanner(projects []string, interval time.Duration, isDryRun bool) {
	mode := ""
	if isDryRun {
		mode = " " + constants.ColorYellow + "[DRY-RUN]" + constants.ColorCyan
	}
	fmt.Printf("\n  %s[SUG]%s%s Monitoring %d project(s) until green for automated shutdown (interval: %v)...\n\n",
		constants.ColorCyan, mode, constants.ColorReset, len(projects), interval)
}

func printSUGOnceFinished(pending int) {
	fmt.Printf("\n  %sℹ [SUG] Single evaluation pass completed (%d pending).%s\n\n",
		constants.ColorCyan, pending, constants.ColorReset)
}

func printSUGPendingWait(pending int, interval time.Duration) {
	fmt.Printf("\n  %s⏳ [SUG] %d project(s) not yet green. Next evaluation in %v...%s\n\n",
		constants.ColorYellow, pending, interval, constants.ColorReset)
}

func evaluateAllProjectsStatus(projects []string) int {
	pending := 0
	for _, p := range projects {
		isGreen := checkSingleProjectPipelineGreen(p)
		printProjectEvaluationRow(p, isGreen)
		if !isGreen {
			pending++
		}
	}
	return pending
}

func checkSingleProjectPipelineGreen(project string) bool {
	payload, _, hasFailures := cmdpipeline.FetchPipelineErrorReportWithMeta(project, false)
	return !hasFailures && !payload.IsRunning
}

func printProjectEvaluationRow(project string, isGreen bool) {
	status := constants.ColorGreen + "✔ GREEN" + constants.ColorReset
	if !isGreen {
		status = constants.ColorYellow + "⏳ RUNNING/FAILING" + constants.ColorReset
	}
	fmt.Printf("    • %-30s : %s\n", project, status)
}

func triggerSystemShutdown(total int, isDryRun bool) error {
	fmt.Println()
	if isDryRun {
		cmdStr := getShutdownCommandStr(runtime.GOOS)
		fmt.Printf("  %s✔ [DRY-RUN] All %d monitored projects are green. OS shutdown command would be executed (%s).%s\n\n",
			constants.ColorGreen+"\033[1m", total, cmdStr, constants.ColorReset)
		return nil
	}
	cfg := loadSUGConfig()
	if cfg.IsClearOnShutdown() {
		cfg.ProjectTargets = []string{}
		_ = saveSUGConfig(cfg)
	}
	fmt.Printf("  %s🚀 All %d monitored projects are green. Initiating system shutdown...%s\n\n",
		constants.ColorGreen+"\033[1m", total, constants.ColorReset)
	return OSShutdownExecutorFn(runtime.GOOS)
}

func getShutdownCommandStr(goos string) string {
	switch strings.ToLower(goos) {
	case "windows":
		return "shutdown /s /t 60"
	case "darwin":
		return `osascript -e 'tell app "System Events" to shut down'`
	default:
		return "shutdown -h +1"
	}
}

func defaultCrossPlatformShutdown(goos string) error {
	switch strings.ToLower(goos) {
	case "windows":
		return exec.Command("shutdown", "/s", "/t", "60").Run()
	case "darwin":
		return execMacOSShutdown()
	default:
		return exec.Command("shutdown", "-h", "+1").Run()
	}
}

func execMacOSShutdown() error {
	err := exec.Command("osascript", "-e", "tell app \"System Events\" to shut down").Run()
	if err != nil {
		return exec.Command("shutdown", "-h", "+1").Run()
	}
	return nil
}
