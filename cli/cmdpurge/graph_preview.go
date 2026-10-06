package cmdpurge

import (
	"fmt"
	"io"
	"strings"
)

const defaultBoxWidth = 76

// RenderPreflightCard returns the visual ASCII/Unicode boxed before-and-after graph diff.
func RenderPreflightCard(report *PreflightReport) string {
	var sb strings.Builder
	renderCardTop(&sb, defaultBoxWidth)
	renderHeaderSection(&sb, report, defaultBoxWidth)
	renderDivider(&sb, defaultBoxWidth)
	renderMetricsSection(&sb, report, defaultBoxWidth)
	renderDivider(&sb, defaultBoxWidth)
	renderTopologyDiff(&sb, report, defaultBoxWidth)
	renderDivider(&sb, defaultBoxWidth)
	renderAdvisoryFooter(&sb, defaultBoxWidth)
	renderCardBottom(&sb, defaultBoxWidth)

	return sb.String()
}

// RenderPreflightBox writes a structured impact summary box to the output writer.
func RenderPreflightBox(report *PreflightReport, out io.Writer) error {
	card := RenderPreflightCard(report)
	_, err := fmt.Fprint(out, card)

	return err
}

// RenderGraphComparison renders before and after commit graph topologies.
func RenderGraphComparison(report *PreflightReport, out io.Writer) error {
	var sb strings.Builder
	sb.WriteString("Commit Topology Diff:\n")
	renderBeforeGraph(&sb, report)
	sb.WriteString("\n")
	renderAfterGraph(&sb, report)
	_, err := fmt.Fprint(out, sb.String())

	return err
}

func renderCardTop(sb *strings.Builder, width int) {
	title := " GITMAP HISTORY PURGE: PRE-FLIGHT IMPACT REPORT "
	dashes := width - 2 - len(title)
	if dashes < 2 {
		dashes = 2
	}
	sb.WriteString("┌─" + title + strings.Repeat("─", dashes-2) + "┐\n")
}

func renderDivider(sb *strings.Builder, width int) {
	sb.WriteString("├" + strings.Repeat("─", width-2) + "┤\n")
}

func renderCardBottom(sb *strings.Builder, width int) {
	sb.WriteString("└" + strings.Repeat("─", width-2) + "┘\n")
}

func padBoxLine(content string, width int) string {
	innerCap := width - 4
	if innerCap < 10 {
		innerCap = 10
	}
	if len(content) > innerCap {
		content = content[:innerCap-3] + "..."
	}
	padding := innerCap - len(content)
	if padding < 0 {
		padding = 0
	}

	return "│ " + content + strings.Repeat(" ", padding) + " │\n"
}

func renderHeaderSection(sb *strings.Builder, r *PreflightReport, width int) {
	sb.WriteString(padBoxLine("Target Repo:        "+r.RepoSlug, width))
	sb.WriteString(padBoxLine("Target Mode:        "+formatModeLabel(r.TargetType), width))
	sb.WriteString(padBoxLine("Target Path:        "+r.TargetPath, width))
	sb.WriteString(padBoxLine("Working Tree:       "+formatDirtyLabel(r.IsWorkingTreeDirty), width))
}

func renderMetricsSection(sb *strings.Builder, r *PreflightReport, width int) {
	affCount := len(r.AffectedCommits)
	sb.WriteString(padBoxLine(fmt.Sprintf("Commits Scanned:    %d commits", r.TotalCommitsScanned), width))
	sb.WriteString(padBoxLine(fmt.Sprintf("Affected Commits:   %d commits contain target data", affCount), width))
	sb.WriteString(padBoxLine(fmt.Sprintf("Distinct Files:     %d unique file paths", r.DistinctFilesCount), width))
	sb.WriteString(padBoxLine(fmt.Sprintf("Estimated Blobs:    %s to stage in temp backup", formatBytes(r.TotalEstimatedBytes)), width))
	sb.WriteString(padBoxLine("Impacted Branches:  "+joinOrNone(r.AffectedBranches), width))
	sb.WriteString(padBoxLine("Impacted Tags:      "+joinOrNone(r.AffectedTags), width))
	sb.WriteString(padBoxLine("Remote Status:      "+formatRemoteLabel(r.HasPushedCommits), width))
}

func renderTopologyDiff(sb *strings.Builder, r *PreflightReport, width int) {
	sb.WriteString(padBoxLine("Commit Topology Diff Preview:", width))
	sb.WriteString(padBoxLine("  BEFORE (Current History):", width))
	renderTopologyBeforeNodes(sb, r.AffectedCommits, width)
	sb.WriteString(padBoxLine("  AFTER (Proposed Rewrite):", width))
	renderTopologyAfterNodes(sb, r.AffectedCommits, r.TargetType, width)
}

func renderTopologyBeforeNodes(sb *strings.Builder, nodes []AffectedCommitNode, width int) {
	displayNodes := limitDisplayNodes(nodes, 4)
	for _, n := range displayNodes {
		badge := "          "
		if n.IsAffected {
			badge = "[AFFECTED]"
		}
		sha := shortShaOrFallback(n)
		msg := truncateString(n.CommitMessage, 28)
		sb.WriteString(padBoxLine(fmt.Sprintf("    * %s %s %s", sha, badge, msg), width))
	}
}

func renderTopologyAfterNodes(sb *strings.Builder, nodes []AffectedCommitNode, t TargetType, width int) {
	displayNodes := limitDisplayNodes(nodes, 4)
	for _, n := range displayNodes {
		sha := shortShaOrFallback(n)
		if n.IsAffected && t == TargetTypeCommit {
			sb.WriteString(padBoxLine("    * [EXCISED] commit removed from tree", width))
			continue
		}

		if n.IsAffected {
			sb.WriteString(padBoxLine(fmt.Sprintf("    * %s (cln)   %s", sha, truncateString(n.CommitMessage, 26)), width))
			continue
		}

		sb.WriteString(padBoxLine(fmt.Sprintf("    * %s          %s", sha, truncateString(n.CommitMessage, 28)), width))
	}
}

func renderBeforeGraph(sb *strings.Builder, r *PreflightReport) {
	sb.WriteString("  BEFORE (Current History):\n")
	for _, n := range limitDisplayNodes(r.AffectedCommits, 6) {
		status := ""
		if n.IsAffected {
			status = " [AFFECTED] <-- (target data removed)"
		}
		sb.WriteString(fmt.Sprintf("    * %s %s%s\n", shortShaOrFallback(n), n.CommitMessage, status))
	}
}

func renderAfterGraph(sb *strings.Builder, r *PreflightReport) {
	sb.WriteString("  AFTER (Proposed Rewrite):\n")
	for _, n := range limitDisplayNodes(r.AffectedCommits, 6) {
		if n.IsAffected && r.TargetType == TargetTypeCommit {
			sb.WriteString("    * [EXCISED] commit dropped\n")
		} else if n.IsAffected {
			sb.WriteString(fmt.Sprintf("    * %s (cln) %s <-- [re-parented]\n", shortShaOrFallback(n), n.CommitMessage))
		} else {
			sb.WriteString(fmt.Sprintf("    * %s %s\n", shortShaOrFallback(n), n.CommitMessage))
		}
	}
}

func renderAdvisoryFooter(sb *strings.Builder, width int) {
	sb.WriteString(padBoxLine("You can undo this if you wanted to. We do not confirm this,", width))
	sb.WriteString(padBoxLine("but you can try: gitmap history undo <operation-id>", width))
}

func limitDisplayNodes(nodes []AffectedCommitNode, limit int) []AffectedCommitNode {
	if len(nodes) <= limit {
		return nodes
	}

	return nodes[:limit]
}

func shortShaOrFallback(n AffectedCommitNode) string {
	if n.ShortSha != "" {
		return n.ShortSha
	}
	if len(n.CommitSha) >= 7 {
		return n.CommitSha[:7]
	}

	return n.CommitSha
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max-3] + "..."
}

func formatModeLabel(t TargetType) string {
	switch t {
	case TargetTypeFolder:
		return "Folder Purge"
	case TargetTypeFile:
		return "File Purge"
	case TargetTypeCommit:
		return "Commit Purge"
	default:
		return "History Purge"
	}
}

func formatDirtyLabel(isDirty bool) string {
	if isDirty {
		return "DIRTY (uncommitted changes detected)"
	}

	return "CLEAN"
}

func formatRemoteLabel(hasPushed bool) string {
	if hasPushed {
		return "COMMITS ALREADY PUSHED (force push required)"
	}

	return "LOCAL ONLY (safe to purge)"
}

func joinOrNone(items []string) string {
	if len(items) == 0 {
		return "none"
	}

	return strings.Join(items, ", ")
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024.0)
	}

	return fmt.Sprintf("%.2f MB", float64(bytes)/(1024.0*1024.0))
}
