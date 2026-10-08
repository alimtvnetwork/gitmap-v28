package cmdupdate

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
