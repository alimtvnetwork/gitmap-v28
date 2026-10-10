package cmdupdate

import (
	"context"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
	"strings"
	"sync"
	"time"
)

func executeParallelFleetUpdate(targets []FleetTarget, opts FleetUpdateOptions) []FleetUpdateNodeResult {
	if opts.IsZip {
		clearZipCache()
	}
	results := make([]FleetUpdateNodeResult, len(targets))
	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Printf("\n%s[FLEET UPDATE]%s Updating '%s' across %d cluster node(s) in parallel...\n\n",
		constants.ColorCyan, constants.ColorReset, opts.Pkg, len(targets))

	for idx, t := range targets {
		wg.Add(1)
		go func(i int, target FleetTarget) {
			defer wg.Done()
			res := runSingleFleetUpdate(target, opts)
			mu.Lock()
			results[i] = res
			mu.Unlock()
		}(idx, t)
	}
	wg.Wait()
	return results
}

func runSingleFleetUpdate(target FleetTarget, opts FleetUpdateOptions) FleetUpdateNodeResult {
	if opts.IsDryRun {
		return FleetUpdateNodeResult{
			Alias:      target.Alias,
			IP:         target.IP,
			IsSuccess:  true,
			Status:     "SUCCESS",
			DurationMs: 0,
			Details:    fmt.Sprintf("[DRY-RUN] Would update %s", opts.Pkg),
		}
	}

	isOnline, reason := CheckConnLivenessFn(context.Background(), target.IP, target.Port, 1000*time.Millisecond)
	if !isOnline {
		res := FleetUpdateNodeResult{
			Alias:      target.Alias,
			IP:         target.IP,
			IsSuccess:  false,
			IsOffline:  true,
			Status:     "OFFLINE",
			DurationMs: 15,
			Details:    fmt.Sprintf("machine is off or unreachable (%s)", reason),
		}
		printSingleFleetUpdateProgress(res)
		return res
	}

	start := time.Now()
	rawOutput, err := ExecuteRemoteUpdateFn(target, opts)
	dur := time.Since(start).Milliseconds()

	telemetry := ParseFleetUpdateTelemetry(rawOutput, target, err)
	isSuccess := err == nil && telemetry.Success
	status := "FAILED"
	if isSuccess {
		status = "SUCCESS"
	}

	details := telemetry.Details
	if details == "" {
		details = formatTelemetryDetails(telemetry)
	}

	res := FleetUpdateNodeResult{
		Alias:      target.Alias,
		IP:         target.IP,
		IsSuccess:  isSuccess,
		Status:     status,
		Telemetry:  telemetry,
		DurationMs: dur,
		Details:    details,
		Error:      err,
	}
	printSingleFleetUpdateProgress(res)
	return res
}

func formatTelemetryDetails(t FleetUpdateTelemetry) string {
	if len(t.Updated) > 0 {
		return fmt.Sprintf("Updated: %s", strings.Join(t.Updated, ", "))
	}
	if t.CurrentVersion != "" {
		return fmt.Sprintf("Version: %s", t.CurrentVersion)
	}
	return "OK"
}

func printSingleFleetUpdateProgress(res FleetUpdateNodeResult) {
	if res.IsSuccess {
		fmt.Printf("  %s✓%s [%s|%s] %s (%dms)\n",
			constants.ColorGreen, constants.ColorReset, res.Alias, res.IP, res.Details, res.DurationMs)
		return
	}
	if res.IsOffline {
		fmt.Printf("  %s●%s [%s|%s] OFFLINE: %s (%dms)\n",
			constants.ColorYellow, constants.ColorReset, res.Alias, res.IP, res.Details, res.DurationMs)
		return
	}
	fmt.Printf("  %s✖%s [%s|%s] FAILED: %s (%dms)\n",
		constants.ColorRed, constants.ColorReset, res.Alias, res.IP, res.Details, res.DurationMs)
}

func printFleetNoTargetsBanner(pkg string, excludedCount int) {
	fmt.Printf("\n%s[FLEET UPDATE]%s No matching active targets found for package '%s' (excluded: %d).\n\n",
		constants.ColorYellow, constants.ColorReset, pkg, excludedCount)
}

func renderFleetUpdateSummary(results []FleetUpdateNodeResult, pkg string, excludedCount int) {
	if len(results) == 0 {
		return
	}
	successCount, failCount, offlineCount := calculateFleetMetrics(results)
	fmt.Printf("\n%s================================================================================%s\n",
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" %sSSH Fleet Update Summary [%s]:%s Total: %d | Succeeded: %d | Failed: %d | Offline: %d | Excluded: %d\n",
		constants.ColorBold, pkg, constants.ColorReset, len(results), successCount, failCount, offlineCount, excludedCount)
	fmt.Printf("%s--------------------------------------------------------------------------------%s\n",
		constants.ColorDim, constants.ColorReset)

	cfg := termout.TableConfig{
		Columns: []termout.Column{
			{Title: "ALIAS", Align: termout.AlignLeft, MinWidth: 15},
			{Title: "IP", Align: termout.AlignLeft, MinWidth: 16},
			{Title: "STATUS", Align: termout.AlignLeft, MinWidth: 10},
			{Title: "DURATION", Align: termout.AlignRight, MinWidth: 10},
			{Title: "DETAILS", Align: termout.AlignLeft, MinWidth: 25},
		},
		Rows: buildFleetUpdateRows(results),
	}
	termout.PrintTable(cfg)
	fmt.Printf("%s================================================================================%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func calculateFleetMetrics(results []FleetUpdateNodeResult) (int, int, int) {
	succeeded := 0
	failed := 0
	offline := 0
	for _, r := range results {
		if r.IsSuccess {
			succeeded++
			continue
		}
		if r.IsOffline {
			offline++
			continue
		}
		failed++
	}
	return succeeded, failed, offline
}

func buildFleetUpdateRows(results []FleetUpdateNodeResult) []termout.Row {
	rows := make([]termout.Row, 0, len(results))
	for _, r := range results {
		statusStr := constants.ColorRed + "FAILED" + constants.ColorReset
		if r.IsSuccess {
			statusStr = constants.ColorGreen + "SUCCESS" + constants.ColorReset
		} else if r.IsOffline {
			statusStr = constants.ColorYellow + "OFFLINE" + constants.ColorReset
		}
		detail := sanitizeTableRowDetail(r.Details)
		rows = append(rows, termout.Row{
			Cells: []string{
				r.Alias,
				r.IP,
				statusStr,
				fmt.Sprintf("%dms", r.DurationMs),
				detail,
			},
		})
	}
	return rows
}

func sanitizeTableRowDetail(raw string) string {
	clean := strings.ReplaceAll(raw, "\r", "")
	clean = strings.ReplaceAll(clean, "\n", " ")
	clean = strings.TrimSpace(clean)
	if len(clean) > 50 {
		return clean[:47] + "..."
	}
	if clean == "" {
		return "OK"
	}
	return clean
}
