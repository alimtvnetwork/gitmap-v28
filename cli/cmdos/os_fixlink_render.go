package cmdos

import (
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// FixLinkJSONPayload wraps fix-link results for --json formatting.
type FixLinkJSONPayload struct {
	Status        string       `json:"status"`
	IsDryRun      bool         `json:"is_dry_run"`
	TotalChecked  int          `json:"total_checked"`
	TotalHealthy  int          `json:"total_healthy"`
	TotalRepaired int          `json:"total_repaired"`
	TotalBroken   int          `json:"total_broken"`
	Links         []LinkResult `json:"links"`
}

func renderLinkResults(results []LinkResult, opts FixLinkOptions) error {
	if opts.IsJSON {
		return renderLinkResultsJSON(results, opts.IsDryRun)
	}

	renderLinkResultsTerminal(results, opts.IsDryRun)

	return nil
}

func renderLinkResultsTerminal(results []LinkResult, isDryRun bool) {
	printHeaderTerminal(isDryRun)

	for _, r := range results {
		renderLinkItemTerminal(r)
	}

	printSummaryTerminal(results)
}

func printHeaderTerminal(isDryRun bool) {
	if isDryRun {
		fmt.Printf("▶ %sInspecting symlinks (dry run)...%s\n", constants.ColorYellow, constants.ColorReset)

		return
	}

	fmt.Printf("▶ %sInspecting and repairing symlinks...%s\n", constants.ColorCyan, constants.ColorReset)
}

func renderLinkItemTerminal(r LinkResult) {
	if r.IsRepaired {
		fmt.Printf("  %s✓ Repaired:%s %s -> %s (%s)\n", constants.ColorGreen, constants.ColorReset, r.Path, r.Target, r.Message)

		return
	}

	if r.IsHealthy {
		fmt.Printf("  %s✓ Healthy:%s  %s -> %s\n", constants.ColorGreen, constants.ColorReset, r.Path, r.Target)

		return
	}

	fmt.Printf("  %s✖ Broken:%s   %s -> %s (%s)\n", constants.ColorRed, constants.ColorReset, r.Path, r.Target, r.Message)
}

func printSummaryTerminal(results []LinkResult) {
	healthy, repaired, broken := countResultStats(results)
	total := len(results)

	fmt.Println()
	fmt.Printf("  • Summary: %d checked, %d healthy, %d repaired, %d broken\n",
		total, healthy, repaired, broken)
}

func countResultStats(results []LinkResult) (int, int, int) {
	var healthy, repaired, broken int
	for _, r := range results {
		switch {
		case r.IsRepaired:
			repaired++
		case r.IsHealthy:
			healthy++
		default:
			broken++
		}
	}

	return healthy, repaired, broken
}

func renderLinkResultsJSON(results []LinkResult, isDryRun bool) error {
	h, r, b := countResultStats(results)
	status := "success"
	if b > 0 {
		status = "degraded"
	}

	payload := FixLinkJSONPayload{
		Status:        status,
		IsDryRun:      isDryRun,
		TotalChecked:  len(results),
		TotalHealthy:  h,
		TotalRepaired: r,
		TotalBroken:   b,
		Links:         results,
	}

	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
