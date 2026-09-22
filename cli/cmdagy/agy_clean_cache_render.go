// Package cmdagy — agy_clean_cache_render.go renders console outputs and reports for cache cleaning.
package cmdagy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// FormatBytes converts a byte count to a human-readable string (KB, MB, GB).
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// RenderCleanPreview prints target directories and detected processes to standard output.
func RenderCleanPreview(targets []AgyCacheTarget, procs []AgyProcessInfo) {
	fmt.Printf("\n%s⚡ Antigravity Cache & Process Cleanup%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%sTarget cleanup locations:%s\n", constants.ColorWhite, constants.ColorReset)

	for _, t := range targets {
		if t.Exists {
			fmt.Printf("  %s✓%s %-48s %s%s, %d files%s\n",
				constants.ColorGreen, constants.ColorReset,
				t.Path, constants.ColorDim, FormatBytes(t.SizeBytes), t.FileCount, constants.ColorReset)
		} else {
			fmt.Printf("  %s•%s %-48s %s(not found)%s\n",
				constants.ColorDim, constants.ColorReset, t.Path, constants.ColorDim, constants.ColorReset)
		}
	}

	fmt.Printf("\n%sProcesses to terminate:%s\n", constants.ColorWhite, constants.ColorReset)
	if len(procs) == 0 {
		fmt.Printf("  %s• None running%s\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, p := range procs {
		fmt.Printf("  %s• PID %-6d %s%s\n", constants.ColorYellow, p.PID, p.Name, constants.ColorReset)
	}
}

// AskProceedConfirmation prompts the user interactively to confirm cleanup.
func AskProceedConfirmation() bool {
	fmt.Printf("\n%sProceed with cleanup? [y/N]: %s", constants.ColorYellow, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	ans = strings.TrimSpace(strings.ToLower(ans))
	return ans == "y" || ans == "yes"
}

// RenderCleanSuccess prints execution metrics and safety confirmation.
func RenderCleanSuccess(report CleanCacheReport) {
	fmt.Printf("\n%s✓ Cleanup completed successfully!%s\n", constants.ColorGreen, constants.ColorReset)
	if report.ProcessesTerminated > 0 {
		fmt.Printf("  • Terminated %d process(es)\n", report.ProcessesTerminated)
	}
	fmt.Printf("  • Removed %d file(s) (%s freed)\n", report.FilesDeleted, report.HumanFreed)
	if report.ConvsDeleted > 0 {
		fmt.Printf("  • Pruned %d conversation(s) exceeding retention limit\n", report.ConvsDeleted)
	}
	fmt.Printf("  • Completed in %dms\n", report.DurationMs)
}

// RenderCleanJSON marshals and prints the CleanCacheReport in indented JSON format.
func RenderCleanJSON(report CleanCacheReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
