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
	fmt.Fprintf(out, "  ║ GITMAP FLEET NODES %-86s║\n", strings.ToUpper(string(opts.Kind))+" DISPATCH")
	fmt.Fprintln(out, "  ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")
	fileMsg := ""
	if opts.HasFile {
		fileMsg = fmt.Sprintf(" (staged '%s' to remote work directories)", opts.DetectedFile)
	}
	fmt.Fprintf(out, "  ▸ Dispatching '%s' across local host and %d remote fleet node(s)%s...\n\n", opts.Kind, nodeCount, fileMsg)
}

func renderFleetResultsTable(out io.Writer, results []RemoteCloneNodeResult, isLocalSuccess bool) {
	fmt.Fprintln(out, "  NODE (ALIAS)     HOST                   ROLE       STATUS        DURATION   DETAILS")
	fmt.Fprintln(out, "  --------------------------------------------------------------------------------------------------------------")
	renderLocalRow(out, isLocalSuccess)
	for _, r := range results {
		renderRemoteRow(out, r)
	}
	fmt.Fprintln(out, "  --------------------------------------------------------------------------------------------------------------")
	renderFleetSummaryFooter(out, results, isLocalSuccess)
}

func renderLocalRow(out io.Writer, isLocalSuccess bool) {
	statusTag := constants.ColorGreen + "● success" + constants.ColorReset
	if !isLocalSuccess {
		statusTag = constants.ColorRed + "✗ failed" + constants.ColorReset
	}
	fmt.Fprintf(out, "  %-16s %-22s %-10s %-20s %-10s %s\n",
		"local (current)", "127.0.0.1", "master", statusTag, "in-process", "executed directly on host machine")
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

func truncateFirstLine(s string) string {
	firstLine := strings.Split(s, "\n")[0]
	if len(firstLine) > 45 {
		return firstLine[:42] + "..."
	}
	return firstLine
}

func formatDetails(r RemoteCloneNodeResult) string {
	if r.Error != "" {
		return r.Error
	}
	if r.Stdout != "" {
		return truncateFirstLine(r.Stdout)
	}
	return "done"
}

func renderFleetSummaryFooter(out io.Writer, results []RemoteCloneNodeResult, isLocalSuccess bool) {
	var succCount int
	var failCount int
	for _, r := range results {
		if r.Status == "success" {
			succCount++
		} else {
			failCount++
		}
	}
	if isLocalSuccess {
		succCount++
	} else {
		failCount++
	}
	total := len(results) + 1
	fmt.Fprintf(out, "\n  ✔ Fleet Clone Summary: %d/%d node(s) completed successfully (%d failed)\n\n",
		succCount, total, failCount)
}
