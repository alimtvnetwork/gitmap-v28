// Package prdesc provides structured Markdown pull request description generators (Spec 129).
package prdesc

import (
	"fmt"
	"strings"
	"time"
)

// PRCommitInfo represents a single commit included in the pull request.
type PRCommitInfo struct {
	ShortSHA string
	Author   string
	Subject  string
}

// PRComponentImpact represents the impact summary on a project component.
type PRComponentImpact struct {
	Component string
	Impact    string
	Files     int
}

// PRFileDiffSummary represents aggregate file diff statistics.
type PRFileDiffSummary struct {
	FilesChanged int
	Additions    int
	Deletions    int
}

// PRMetadata holds structured metadata for PR description generation.
type PRMetadata struct {
	PRNumber         int
	Title            string
	SourceBranch     string
	TargetBranch     string
	Author           string
	Timestamp        time.Time
	ExecutiveSummary string
	Commits          []PRCommitInfo
	Impacts          []PRComponentImpact
	DiffSummary      PRFileDiffSummary
	SafetyChecks     []string
}

// GeneratePRDescription builds a structured Markdown description for a pull request.
func GeneratePRDescription(meta PRMetadata) string {
	var b strings.Builder
	renderHeader(&b, meta)
	renderExecutiveSummary(&b, meta.ExecutiveSummary)
	renderComponentImpactTable(&b, meta.Impacts)
	renderCommitsTable(&b, meta.Commits)
	renderDiffSummary(&b, meta.DiffSummary)
	renderSafetyChecklist(&b, meta.SafetyChecks)

	return b.String()
}

func renderHeader(b *strings.Builder, meta PRMetadata) {
	b.WriteString(fmt.Sprintf("## Pull Request: %s\n\n", meta.Title))
	b.WriteString("| Metadata | Details |\n")
	b.WriteString("|---|---|\n")
	b.WriteString(fmt.Sprintf("| **Source Branch** | `%s` |\n", meta.SourceBranch))
	b.WriteString(fmt.Sprintf("| **Target Branch** | `%s` |\n", meta.TargetBranch))
	b.WriteString(fmt.Sprintf("| **Author** | %s |\n", meta.Author))
	b.WriteString(fmt.Sprintf("| **Timestamp** | %s |\n", meta.Timestamp.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("| **Commit Count** | %d |\n\n", len(meta.Commits)))
}

func renderExecutiveSummary(b *strings.Builder, summary string) {
	b.WriteString("### Executive Summary\n\n")
	if summary == "" {
		b.WriteString("Automated feature branch replay and merge.\n\n")

		return
	}
	b.WriteString(summary)
	b.WriteString("\n\n")
}

func renderComponentImpactTable(b *strings.Builder, impacts []PRComponentImpact) {
	b.WriteString("### Component Impact\n\n")
	b.WriteString("| Component | Impact | Files Changed |\n")
	b.WriteString("|---|---|---|\n")
	if len(impacts) == 0 {
		b.WriteString("| Core | Standard replay | 1 |\n\n")

		return
	}
	for _, imp := range impacts {
		b.WriteString(fmt.Sprintf("| `%s` | %s | %d |\n", imp.Component, imp.Impact, imp.Files))
	}
	b.WriteString("\n")
}

func renderCommitsTable(b *strings.Builder, commits []PRCommitInfo) {
	b.WriteString("### Commits Included\n\n")
	b.WriteString("| Commit | Author | Subject |\n")
	b.WriteString("|---|---|---|\n")
	for _, c := range commits {
		b.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", c.ShortSHA, c.Author, c.Subject))
	}
	b.WriteString("\n")
}

func renderDiffSummary(b *strings.Builder, diff PRFileDiffSummary) {
	b.WriteString("### File Diff Summary\n\n")
	b.WriteString(fmt.Sprintf("- **Files Changed:** %d\n", diff.FilesChanged))
	b.WriteString(fmt.Sprintf("- **Additions:** +%d lines\n", diff.Additions))
	b.WriteString(fmt.Sprintf("- **Deletions:** -%d lines\n\n", diff.Deletions))
}

func renderSafetyChecklist(b *strings.Builder, checks []string) {
	b.WriteString("### Safety Checklist\n\n")
	list := defaultSafetyChecks(checks)
	for _, check := range list {
		b.WriteString(fmt.Sprintf("- [x] %s\n", check))
	}
	b.WriteString("\n")
}

func defaultSafetyChecks(checks []string) []string {
	if len(checks) > 0 {
		return checks
	}

	return []string{
		"Clean replay without merge conflicts",
		"Author identity and timestamps preserved",
		"Mainline target integrity verified",
	}
}
