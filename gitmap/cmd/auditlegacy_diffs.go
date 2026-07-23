// Package cmd: per-file unified-diff writer for `gitmap audit-legacy`.
//
// When --diffs is set alongside --report, every file with at least one
// match gets its own `<reportDir>/diffs/<sanitized-path>.diff` artifact
// containing a minimal unified diff that previews the legacy → v8
// substitution. The Markdown report links each diff so a reviewer can
// click straight from the file-counts table to the proposed change.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// auditDiffPlan is everything the report renderer needs to link a diff.
type auditDiffPlan struct {
	SourceFile  string // original repo-relative file path
	DiffPath    string // absolute on-disk path of the .diff artifact
	DiffRelLink string // path relative to the report file (for Markdown links)
}

// writeAuditLegacyDiffs writes one unified-diff artifact per file with hits.
// Returns plans the report renderer uses to insert hyperlinks. No-op when
// --diffs wasn't set or the report path is empty.
func writeAuditLegacyDiffs(opts auditLegacyOpts, hits []auditLegacyHit) []auditDiffPlan {
	if !opts.WriteDiffs || opts.ReportPath == "" || len(hits) == 0 {
		return nil
	}
	diffsDir := auditDiffsDir(opts.ReportPath)
	if err := os.MkdirAll(diffsDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAuditLegacyDiffWrite, diffsDir, err)

		return nil
	}
	plans := buildAuditDiffPlans(opts, hits, diffsDir)
	if len(plans) > 0 {
		fmt.Fprintf(os.Stdout, constants.MsgAuditLegacyDiffsWrote, len(plans), diffsDir)
	}

	return plans
}

// auditDiffsDir resolves <reportDir>/diffs/.
func auditDiffsDir(reportPath string) string {
	return filepath.Join(filepath.Dir(reportPath), constants.DefaultAuditLegacyDiffsDir)
}

// buildAuditDiffPlans iterates unique files and writes one diff each.
func buildAuditDiffPlans(opts auditLegacyOpts, hits []auditLegacyHit, diffsDir string) []auditDiffPlan {
	files := uniqueAuditFiles(hits)
	plans := make([]auditDiffPlan, 0, len(files))
	for _, file := range files {
		plan, ok := writeOneAuditDiff(opts, file, diffsDir)
		if ok {
			plans = append(plans, plan)
		}
	}

	return plans
}

// writeOneAuditDiff renders + persists the diff for a single source file.
func writeOneAuditDiff(opts auditLegacyOpts, file, diffsDir string) (auditDiffPlan, bool) {
	body, err := renderAuditDiff(file, opts.Patterns)
	if err != nil || body == "" {
		return auditDiffPlan{}, false
	}
	diffPath := filepath.Join(diffsDir, sanitizeAuditDiffName(file)+".diff")
	if err := os.WriteFile(diffPath, []byte(body), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAuditLegacyDiffWrite, diffPath, err)

		return auditDiffPlan{}, false
	}
	rel := relAuditDiffLink(opts.ReportPath, diffPath)

	return auditDiffPlan{SourceFile: file, DiffPath: diffPath, DiffRelLink: rel}, true
}

// sanitizeAuditDiffName flattens a path into a single safe filename so
// diffs/<name>.diff is collision-resistant across nested directories.
func sanitizeAuditDiffName(file string) string {
	clean := filepath.ToSlash(filepath.Clean(file))
	clean = strings.TrimPrefix(clean, "./")
	clean = strings.ReplaceAll(clean, "/", "__")
	clean = strings.ReplaceAll(clean, ":", "_")

	return clean
}

// relAuditDiffLink computes a report-relative path for Markdown links.
// Falls back to the absolute diff path when filepath.Rel fails.
func relAuditDiffLink(reportPath, diffPath string) string {
	rel, err := filepath.Rel(filepath.Dir(reportPath), diffPath)
	if err != nil {
		return filepath.ToSlash(diffPath)
	}

	return filepath.ToSlash(rel)
}
