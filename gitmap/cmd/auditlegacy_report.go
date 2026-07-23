// Package cmd: Markdown report writer for `gitmap audit-legacy`.
//
// Splits per-pattern + per-file counts and the full hit list into a
// human-readable Markdown file. Used so CI / contributors can attach
// a single artifact to a PR or share an audit summary without piping
// JSON through a formatter.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// writeAuditLegacyReport renders a Markdown audit report to disk.
// No-op when ReportPath is empty (the user didn't pass --report).
func writeAuditLegacyReport(opts auditLegacyOpts, hits []auditLegacyHit, fileCount int, plans []auditDiffPlan) {
	if opts.ReportPath == "" {
		return
	}
	body := renderAuditMarkdown(opts, hits, fileCount, plans)
	if err := writeAuditReportFile(opts.ReportPath, body); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAuditLegacyReportWrite, opts.ReportPath, err)

		return
	}
	fmt.Fprintf(os.Stdout, constants.MsgAuditLegacyReportWrote, opts.ReportPath)
}

// writeAuditReportFile creates parent dirs then writes body to path.
func writeAuditReportFile(path, body string) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(path, []byte(body), 0o644)
}

// renderAuditMarkdown builds the full report body.
func renderAuditMarkdown(opts auditLegacyOpts, hits []auditLegacyHit, fileCount int, plans []auditDiffPlan) string {
	var b strings.Builder
	writeAuditMDHeader(&b, opts, hits, fileCount)
	writeAuditMDPatternCounts(&b, opts, hits)
	writeAuditMDFileCounts(&b, hits, plans)
	writeAuditMDDiffArtifacts(&b, plans)
	writeAuditMDHitList(&b, hits)

	return b.String()
}

// writeAuditMDHeader writes the title + summary block.
func writeAuditMDHeader(b *strings.Builder, opts auditLegacyOpts, hits []auditLegacyHit, fileCount int) {
	files := uniqueAuditFiles(hits)
	fmt.Fprintf(b, "# Legacy Reference Audit\n\n")
	fmt.Fprintf(b, "- Root scanned: `%s`\n", opts.Root)
	fmt.Fprintf(b, "- Patterns: `%s`\n", strings.Join(opts.Raw, "`, `"))
	fmt.Fprintf(b, "- Files scanned: **%d**\n", fileCount)
	fmt.Fprintf(b, "- Total matches: **%d**\n", len(hits))
	fmt.Fprintf(b, "- Files with matches: **%d**\n\n", len(files))
}

// writeAuditMDPatternCounts writes a per-pattern hit-count table.
func writeAuditMDPatternCounts(b *strings.Builder, opts auditLegacyOpts, hits []auditLegacyHit) {
	fmt.Fprintf(b, "## Counts by pattern\n\n")
	fmt.Fprintf(b, "| Pattern | Matches |\n|---|---:|\n")
	counts := countAuditByPattern(opts.Raw, hits)
	for _, p := range opts.Raw {
		fmt.Fprintf(b, "| `%s` | %d |\n", p, counts[p])
	}
	fmt.Fprintln(b)
}

// writeAuditMDFileCounts writes a per-file hit-count table. When diff
// plans are present, an extra "Diff" column links to each artifact.
func writeAuditMDFileCounts(b *strings.Builder, hits []auditLegacyHit, plans []auditDiffPlan) {
	fmt.Fprintf(b, "## Counts by file\n\n")
	if len(hits) == 0 {
		fmt.Fprintln(b, "_None — repo is clean._")
		fmt.Fprintln(b)

		return
	}
	links := indexAuditDiffPlans(plans)
	writeAuditMDFileCountsTable(b, hits, links)
}

// writeAuditMDFileCountsTable renders the actual table rows.
func writeAuditMDFileCountsTable(b *strings.Builder, hits []auditLegacyHit, links map[string]string) {
	hasDiffs := len(links) > 0
	if hasDiffs {
		fmt.Fprintf(b, "| File | Matches | Diff |\n|---|---:|---|\n")
	} else {
		fmt.Fprintf(b, "| File | Matches |\n|---|---:|\n")
	}
	for _, row := range sortedFileCounts(hits) {
		writeAuditMDFileCountRow(b, row, links, hasDiffs)
	}
	fmt.Fprintln(b)
}

// writeAuditMDFileCountRow emits one row of the file-counts table.
func writeAuditMDFileCountRow(b *strings.Builder, row auditFileCount, links map[string]string, hasDiffs bool) {
	if !hasDiffs {
		fmt.Fprintf(b, "| `%s` | %d |\n", row.file, row.count)

		return
	}
	link := links[row.file]
	if link == "" {
		fmt.Fprintf(b, "| `%s` | %d | — |\n", row.file, row.count)

		return
	}
	fmt.Fprintf(b, "| `%s` | %d | [view](%s) |\n", row.file, row.count, link)
}

// writeAuditMDDiffArtifacts lists every diff artifact in a dedicated
// section so reviewers can grab them all from one place.
func writeAuditMDDiffArtifacts(b *strings.Builder, plans []auditDiffPlan) {
	if len(plans) == 0 {
		return
	}
	fmt.Fprintf(b, "## Per-file diffs\n\n")
	fmt.Fprintf(b, "Each diff previews the legacy → `gitmap-v28` substitution for one file. ")
	fmt.Fprintf(b, "Apply with `patch -p0 < <file>.diff` from the repo root.\n\n")
	for _, p := range plans {
		fmt.Fprintf(b, "- [`%s`](%s)\n", p.SourceFile, p.DiffRelLink)
	}
	fmt.Fprintln(b)
}

// indexAuditDiffPlans builds a source-file → relative-link map.
func indexAuditDiffPlans(plans []auditDiffPlan) map[string]string {
	out := make(map[string]string, len(plans))
	for _, p := range plans {
		out[p.SourceFile] = p.DiffRelLink
	}

	return out
}

// writeAuditMDHitList writes every match as `file:line: text`.
func writeAuditMDHitList(b *strings.Builder, hits []auditLegacyHit) {
	if len(hits) == 0 {
		return
	}
	fmt.Fprintf(b, "## All matches\n\n```\n")
	for _, h := range hits {
		fmt.Fprintf(b, "%s:%d: %s\n", h.File, h.Line, h.Text)
	}
	fmt.Fprintf(b, "```\n")
}

// countAuditByPattern returns per-raw-pattern hit counts.
func countAuditByPattern(raws []string, hits []auditLegacyHit) map[string]int {
	out := map[string]int{}
	for _, p := range raws {
		out[p] = 0
	}
	for _, h := range hits {
		out[h.Pattern]++
	}

	return out
}

// auditFileCount is one row of the per-file count table.
type auditFileCount struct {
	file  string
	count int
}

// sortedFileCounts returns per-file counts sorted by count desc, file asc.
func sortedFileCounts(hits []auditLegacyHit) []auditFileCount {
	counts := map[string]int{}
	for _, h := range hits {
		counts[h.File]++
	}
	out := make([]auditFileCount, 0, len(counts))
	for f, c := range counts {
		out = append(out, auditFileCount{file: f, count: c})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].count != out[j].count {
			return out[i].count > out[j].count
		}

		return out[i].file < out[j].file
	})

	return out
}
