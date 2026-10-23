package cmdnodes

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderResultsTableHeader(out io.Writer) {
	fmt.Fprintf(out, "    %-16s %-18s %-10s %-14s %-12s %-12s %s\n",
		"NODE (ALIAS)", "HOST", "ROLE", "STATUS", "VERSION", "DURATION", "DETAILS")
	fmt.Fprintln(out, "    --------------------------------------------------------------------------------------------------------")
}

func renderSkippedLocalRow(out io.Writer) {
	tag := constants.ColorCyan + "○ skipped" + constants.ColorReset
	fmt.Fprintf(out, "    %-16s %-18s %-10s %s %-12s %-12s %s\n",
		"local (current)", "127.0.0.1", "master", padVisual(tag, 14), constants.Version, "-", "skipped local execution (except-self)")
}

func resolveLocalDetails(localDetails string) string {
	if localDetails != "" {
		return localDetails
	}

	return "executed directly on host machine"
}

func renderExecutedLocalRow(out io.Writer, isLocalSuccess bool, localDetails string, localDur time.Duration) {
	tag := constants.ColorGreen + "● success" + constants.ColorReset
	if !isLocalSuccess {
		tag = constants.ColorRed + "✗ failed" + constants.ColorReset
	}

	durStr := fmt.Sprintf("%dms", localDur.Milliseconds())
	fmt.Fprintf(out, "    %-16s %-18s %-10s %s %-12s %-12s %s\n",
		"local (current)", "127.0.0.1", "master", padVisual(tag, 14), constants.Version, durStr, resolveLocalDetails(localDetails))
}

func renderLocalResultRow(out io.Writer, isLocalSuccess, isSkipLocal bool, localDetails string, localDur time.Duration) {
	if isSkipLocal {
		renderSkippedLocalRow(out)
		return
	}

	renderExecutedLocalRow(out, isLocalSuccess, localDetails, localDur)
}

func formatDurationMs(ms int64) string {
	if ms == 0 {
		return "-"
	}

	return fmt.Sprintf("%dms", ms)
}

func renderRemoteResultRow(out io.Writer, r RemoteCloneNodeResult) {
	statusTag := resolveStatusTag(r.Status)
	durStr := formatDurationMs(r.DurationMs)
	ver := resolveDisplayValue(getNodeVersion(r.Alias))
	fmt.Fprintf(out, "    %-16s %-18s %-10s %s %-12s %-12s %s\n",
		r.Alias, r.Host, r.Role, padVisual(statusTag, 14), ver, durStr, formatDetails(r))
}

func resolveStatusTag(status string) string {
	switch status {
	case "success", "cloned":
		return constants.ColorGreen + "● success" + constants.ColorReset
	case "auth_failed":
		return constants.ColorMagenta + "▲ auth_failed" + constants.ColorReset
	case "offline":
		return constants.ColorYellow + "○ offline" + constants.ColorReset
	default:
		return constants.ColorRed + "✗ failed" + constants.ColorReset
	}
}

func cleanNewlines(s string) string {
	clean := strings.ReplaceAll(s, "\r\n", " ")
	clean = strings.ReplaceAll(clean, "\n", " ")
	return strings.ReplaceAll(clean, "\t", " ")
}

func truncateDetail(s string) string {
	if len(s) > 55 {
		return s[:52] + "..."
	}

	return s
}

func sanitizeError(errStr string) string {
	clean := cleanNewlines(errStr)
	if idx := strings.Index(clean, "output: "); idx != -1 {
		clean = strings.TrimSpace(clean[idx+8:])
	}

	clean = strings.TrimSpace(strings.TrimSuffix(clean, ")"))
	return truncateDetail(clean)
}

func isIgnoredDetailLine(line string) bool {
	low := strings.ToLower(line)
	if strings.HasPrefix(line, "=") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "┌") || strings.HasPrefix(line, "│") || strings.HasPrefix(line, "└") {
		return true
	}

	if strings.HasPrefix(low, "at ") || strings.HasPrefix(low, "origin:") || strings.HasPrefix(low, "stack trace:") {
		return true
	}

	return strings.HasPrefix(low, "pending task already exists")
}

func sanitizeStdout(stdout string) string {
	lines := strings.Split(stdout, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" || isIgnoredDetailLine(line) {
			continue
		}

		return truncateDetail(line)
	}

	return "done"
}

func formatDetails(r RemoteCloneNodeResult) string {
	res := "done"
	switch {
	case r.Details != "":
		res = r.Details
	case r.Error != "":
		res = sanitizeError(r.Error)
	case r.Stdout != "":
		res = sanitizeStdout(r.Stdout)
	}

	return strings.ReplaceAll(res, "(machine is off)", "(unreachable or port 22 closed)")
}

func renderFleetSummaryFooter(out io.Writer, results []RemoteCloneNodeResult, isLocalSuccess bool, isSkipLocal bool) {
	succCount, failCount := countFleetResults(results)
	if isSkipLocal {
		renderSkippedLocalSummary(out, len(results), succCount, failCount)
		return
	}

	succCount, failCount = adjustForLocalResult(succCount, failCount, isLocalSuccess)
	total := len(results) + 1
	fmt.Fprintf(out, "\n  ✔ Fleet CFR Summary: %d/%d node(s) completed successfully (%d failed)\n\n",
		succCount, total, failCount)
}

func countFleetResults(results []RemoteCloneNodeResult) (int, int) {
	var succCount int
	var failCount int
	for _, r := range results {
		if r.Status == "success" || r.Status == "cloned" {
			succCount++
			continue
		}

		failCount++
	}

	return succCount, failCount
}

func adjustForLocalResult(succCount, failCount int, isLocalSuccess bool) (int, int) {
	if isLocalSuccess {
		return succCount + 1, failCount
	}

	return succCount, failCount + 1
}

func renderSkippedLocalSummary(out io.Writer, total, succCount, failCount int) {
	fmt.Fprintf(out, "\n  ✔ Fleet CFR Summary: %d/%d remote node(s) completed successfully (%d failed, local skipped)\n\n",
		succCount, total, failCount)
}

func renderFleetFooterSuggestions(out io.Writer, opts NodesCloneOptions) {
	target := "repo"
	if len(opts.PassArgs) > 0 {
		target = opts.PassArgs[0]
	}

	fmt.Fprintln(out, "  [tip] Fleet Operations & Suggested Commands:")
	fmt.Fprintln(out, "    • Ping Fleet Nodes:          gitmap nodes ping")
	fmt.Fprintln(out, "    • Inspect Node Connection:   gitmap ssh test <alias>")
	fmt.Fprintln(out, "    • Query Machine Telemetry:   gitmap machine --ssh")
	fmt.Fprintf(out, "    • Rerun Except Local Host:   gitmap nodes %s %s --except-self\n\n", opts.Kind, target)
}

func renderFleetResultsTable(out io.Writer, results []RemoteCloneNodeResult, isLocalSuccess bool, localDetails string, localDuration time.Duration, opts NodesCloneOptions) {
	renderResultsTableHeader(out)
	renderLocalResultRow(out, isLocalSuccess, opts.IsSkipLocal, localDetails, localDuration)
	for _, r := range results {
		renderRemoteResultRow(out, r)
	}

	fmt.Fprintln(out, "    --------------------------------------------------------------------------------------------------------")
	renderFleetSummaryFooter(out, results, isLocalSuccess, opts.IsSkipLocal)
	renderFleetFooterSuggestions(out, opts)
}
