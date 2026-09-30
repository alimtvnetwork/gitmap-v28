// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderFleetStartBanner(out io.Writer, opts NodesCloneOptions, nodeCount int) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  ╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintf(out, "  ║ GITMAP FLEET NODES %-94s║\n", strings.ToUpper(string(opts.Kind))+" DISPATCH")
	fmt.Fprintln(out, "  ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")
	fileMsg := ""
	if opts.HasFile {
		fileMsg = fmt.Sprintf(" (staged '%s' to remote work directories)", opts.DetectedFile)
	}
	destMsg := ""
	if opts.TargetDir != "" {
		destMsg = fmt.Sprintf(" [dest: %s]", opts.TargetDir)
	}
	execScope := "local host and"
	if opts.IsSkipLocal {
		execScope = "remote-only (except-self) across"
	}
	fmt.Fprintf(out, "  ▸ Dispatching '%s'%s %s %d remote fleet node(s)%s...\n\n",
		opts.Kind, destMsg, execScope, nodeCount, fileMsg)
}

func renderFleetResultsTable(out io.Writer, results []RemoteCloneNodeResult, isLocalSuccess bool, localDetails string, opts NodesCloneOptions) {
	fmt.Fprintln(out, "  NODE (ALIAS)     HOST                   ROLE       STATUS        DURATION   DETAILS")
	fmt.Fprintln(out, "  --------------------------------------------------------------------------------------------------------------")
	renderLocalRow(out, isLocalSuccess, opts.IsSkipLocal, localDetails)
	for _, r := range results {
		renderRemoteRow(out, r)
	}
	fmt.Fprintln(out, "  --------------------------------------------------------------------------------------------------------------")
	renderFleetSummaryFooter(out, results, isLocalSuccess, opts.IsSkipLocal)
}

func renderLocalRow(out io.Writer, isLocalSuccess bool, isSkipLocal bool, localDetails string) {
	if isSkipLocal {
		statusTag := constants.ColorCyan + "○ skipped" + constants.ColorReset
		fmt.Fprintf(out, "  %-16s %-22s %-10s %-20s %-10s %s\n",
			"local (current)", "127.0.0.1", "master", statusTag, "-", "skipped local execution (except-self)")
		return
	}
	statusTag := constants.ColorGreen + "● success" + constants.ColorReset
	if !isLocalSuccess {
		statusTag = constants.ColorRed + "✗ failed" + constants.ColorReset
	}
	details := "executed directly on host machine"
	if localDetails != "" {
		details = localDetails
	}
	fmt.Fprintf(out, "  %-16s %-22s %-10s %-20s %-10s %s\n",
		"local (current)", "127.0.0.1", "master", statusTag, "in-process", details)
}

func renderRemoteRow(out io.Writer, r RemoteCloneNodeResult) {
	statusTag := resolveStatusTag(r.Status)
	details := formatDetails(r)
	durStr := fmt.Sprintf("%dms", r.DurationMs)
	if r.DurationMs == 0 {
		durStr = "-"
	}
	fmt.Fprintf(out, "  %-16s %-22s %-10s %-20s %-10s %s\n",
		r.Alias, r.Host, r.Role, statusTag, durStr, details)
}

func resolveStatusTag(status string) string {
	switch status {
	case "success":
		return constants.ColorGreen + "● success" + constants.ColorReset
	case "auth_failed":
		return constants.ColorYellow + "○ auth_failed" + constants.ColorReset
	case "offline":
		return constants.ColorYellow + "○ offline" + constants.ColorReset
	default:
		return constants.ColorRed + "✗ failed" + constants.ColorReset
	}
}

func sanitizeError(errStr string) string {
	clean := strings.ReplaceAll(errStr, "\r\n", " ")
	clean = strings.ReplaceAll(clean, "\n", " ")
	clean = strings.ReplaceAll(clean, "\t", " ")
	if idx := strings.Index(clean, "output: "); idx != -1 {
		clean = strings.TrimSpace(clean[idx+8:])
	}
	clean = strings.TrimSuffix(clean, ")")
	clean = strings.TrimSpace(clean)
	if len(clean) > 55 {
		return clean[:52] + "..."
	}
	return clean
}

func isIgnoredDetailLine(line string) bool {
	low := strings.ToLower(line)
	if strings.HasPrefix(line, "=") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "╔") || strings.HasPrefix(line, "║") || strings.HasPrefix(line, "╚") {
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
		if len(line) > 55 {
			return line[:52] + "..."
		}
		return line
	}
	return "done"
}

func formatDetails(r RemoteCloneNodeResult) string {
	if r.Details != "" {
		return r.Details
	}
	if r.Error != "" {
		return sanitizeError(r.Error)
	}
	if r.Stdout != "" {
		return sanitizeStdout(r.Stdout)
	}
	return "done"
}

func renderFleetSummaryFooter(out io.Writer, results []RemoteCloneNodeResult, isLocalSuccess bool, isSkipLocal bool) {
	succCount, failCount := countFleetResults(results)
	if isSkipLocal {
		renderSkippedLocalSummary(out, len(results), succCount, failCount)
		return
	}

	succCount, failCount = adjustForLocalResult(succCount, failCount, isLocalSuccess)
	total := len(results) + 1
	fmt.Fprintf(out, "\n  ✔ Fleet Clone Summary: %d/%d node(s) completed successfully (%d failed)\n\n",
		succCount, total, failCount)
}

func countFleetResults(results []RemoteCloneNodeResult) (int, int) {
	var succCount int
	var failCount int
	for _, r := range results {
		if r.Status == "success" {
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
	fmt.Fprintf(out, "\n  ✔ Fleet Clone Summary: %d/%d remote node(s) completed successfully (%d failed, local skipped)\n\n",
		succCount, total, failCount)
}
