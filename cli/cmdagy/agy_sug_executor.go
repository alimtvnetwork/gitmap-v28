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
func ExecuteSUGWatch(projects []string, interval time.Duration) error {
	printSUGStartBanner(projects, interval)
	for {
		pendingCount := evaluateAllProjectsStatus(projects)
		if pendingCount == 0 {
			return triggerSystemShutdown(len(projects))
		}
		printSUGPendingWait(pendingCount, interval)
		time.Sleep(interval)
	}
}

func printSUGStartBanner(projects []string, interval time.Duration) {
	fmt.Printf("\n  %s[SUG]%s Monitoring %d project(s) until green for automated shutdown (interval: %v)...\n\n",
		constants.ColorCyan, constants.ColorReset, len(projects), interval)
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

func triggerSystemShutdown(total int) error {
	fmt.Println()
	fmt.Printf("  %s🚀 All %d monitored projects are green. Initiating system shutdown...%s\n\n",
		constants.ColorGreen+"\033[1m", total, constants.ColorReset)
	return OSShutdownExecutorFn(runtime.GOOS)
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
