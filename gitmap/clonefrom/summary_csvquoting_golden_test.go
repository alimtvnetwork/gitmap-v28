package clonefrom

// summary_csvquoting_golden_test.go — pins encoding/csv's quoting
// behavior for adversarial Detail/URL strings (commas, double
// quotes, embedded LF, embedded CRLF, leading whitespace). The
// underlying csv.Writer rules are stable across Go versions, but
// this test guards against:
//
//   - accidental column reorder (would shift quoted cells)
//   - drop/add of UseCRLF (would change line terminators)
//   - replacing csv.Writer with a hand-rolled emitter that mishandles
//     quoting
//   - schema additions that re-flow per-row formatting
//
// Encoding/csv quoting contract recap (RFC 4180 + Go specifics):
//
//   - field containing `,`  → wrapped in `"…"`
//   - field containing `"`  → wrapped + every `"` doubled to `""`
//   - field containing LF or CR → wrapped in `"…"`, embedded byte
//     preserved verbatim (CR/LF are NOT escaped, only quoted)
//   - leading/trailing space is NOT quoted by csv.Writer (spec
//     allows but doesn't require it) — pinning catches if Go ever
//     tightens that
//   - record terminator with UseCRLF=true is `\r\n`, even on Unix
//
// Regenerate after deliberate schema changes:
//
//	GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 \
//	  go test ./gitmap/clonefrom/ -run TestCloneFromReport_Golden_Quoting

import (
	"bytes"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// quotingEdgeCaseResults builds a 4-row fixture where every row
// stresses a different csv.Writer quoting branch. Hand-constructed
// (vs. derived from a real Execute run) so the test stays
// hermetic — no git, no filesystem, no network.
//
// Row breakdown:
//
//  1. URL contains a comma — `,` in any field forces quoting of
//     that field only. Detail also contains a comma to verify
//     multiple quoted cells in one record.
//  2. Detail contains an embedded double-quote — must be doubled
//     (`"foo"` → `"""foo"""`). Tests the `"` → `""` escape.
//  3. Detail contains an embedded LF — must be quoted, LF
//     preserved verbatim inside the quotes. This produces a
//     multi-line CSV record that downstream parsers MUST handle.
//  4. Detail contains an embedded CRLF + leading whitespace —
//     guards against csv.Writer ever changing how it handles CR
//     and confirms leading-space is NOT auto-quoted.
func quotingEdgeCaseResults() []Result {
	return []Result{
		{
			Row: Row{URL: "https://example.com/repo,with,commas.git",
				Branch: "main", Depth: 0},
			Dest:     "repo-commas",
			Status:   constants.CloneFromStatusFailed,
			Detail:   "fatal: bad url, see man git-clone",
			Duration: 100 * time.Millisecond,
		},
		{
			Row:      Row{URL: "https://example.com/quoted.git"},
			Dest:     "quoted",
			Status:   constants.CloneFromStatusFailed,
			Detail:   `error: ref "refs/heads/x" not found`,
			Duration: 50 * time.Millisecond,
		},
		{
			Row:      Row{URL: "https://example.com/multiline.git"},
			Dest:     "multiline",
			Status:   constants.CloneFromStatusFailed,
			Detail:   "line one\nline two",
			Duration: 25 * time.Millisecond,
		},
		{
			Row:      Row{URL: "https://example.com/crlf.git"},
			Dest:     "crlf",
			Status:   constants.CloneFromStatusFailed,
			Detail:   "  hint: leading space\r\nplus crlf inside",
			Duration: 10 * time.Millisecond,
		},
	}
}

// TestCloneFromReport_Golden_Quoting pins the bytes for the
// quoting edge-case fixture. Catches drift in csv.Writer escaping
// behavior, record terminator (CRLF), and column ordering when
// quoted cells are present. Reuses assertReportGolden from
// summary_golden_test.go so the regenerate flow is identical.
func TestCloneFromReport_Golden_Quoting(t *testing.T) {
	var buf bytes.Buffer
	if err := writeReportRows(&buf, quotingEdgeCaseResults()); err != nil {
		t.Fatalf("writeReportRows: %v", err)
	}
	assertReportGolden(t, "clonefrom_report_quoting.csv", buf.Bytes())
}
