// Package cmd: emit + report helpers for `gitmap audit-legacy`.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// emitAuditLegacy prints results in JSON or human format.
func emitAuditLegacy(opts auditLegacyOpts, hits []auditLegacyHit, fileCount int) {
	if opts.AsJSON {
		emitAuditLegacyJSON(opts, hits, fileCount)

		return
	}
	emitAuditLegacyText(opts, hits)
}

// emitAuditLegacyJSON prints a JSON report to stdout.
func emitAuditLegacyJSON(opts auditLegacyOpts, hits []auditLegacyHit, fileCount int) {
	report := map[string]any{
		"root":           opts.Root,
		"patterns":       opts.Raw,
		"filesScanned":   fileCount,
		"matchCount":     len(hits),
		"matches":        hits,
		"filesWithMatch": uniqueAuditFiles(hits),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "audit-legacy: json encode failed: %v\n", err)
	}
}

// emitAuditLegacyText prints a human-readable report.
func emitAuditLegacyText(opts auditLegacyOpts, hits []auditLegacyHit) {
	if len(hits) == 0 {
		fmt.Fprintf(os.Stdout, constants.MsgAuditLegacyClean, opts.Root)

		return
	}
	files := uniqueAuditFiles(hits)
	fmt.Fprintf(os.Stdout, constants.MsgAuditLegacyHeader, len(hits), len(files), opts.Raw)
	for _, h := range hits {
		fmt.Fprintf(os.Stdout, constants.MsgAuditLegacyHit, h.File, h.Line, h.Text)
	}
}

// uniqueAuditFiles returns the deduped file list for the report.
func uniqueAuditFiles(hits []auditLegacyHit) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0)
	for _, h := range hits {
		if _, ok := seen[h.File]; ok {
			continue
		}
		seen[h.File] = struct{}{}
		out = append(out, h.File)
	}

	return out
}
