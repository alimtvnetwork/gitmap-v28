package committransfer

// Compile-time + plumbing assertions for the v5.83.0 --max-history-scan
// escape hatch (spec 114 Gap A follow-up). The runtime behaviour
// (bounded vs unbounded `git log`) is covered indirectly by the
// existing TestPlanIdempotenceBeyond200Commits regression guard, which
// continues to pass because the default (0) preserves the unbounded
// scan.

import "testing"

// TestOptions_MaxHistoryScanZeroIsUnbounded pins the contract that
// the zero value of Options.MaxHistoryScan means "unbounded" so
// existing callers (and Go zero-value semantics) keep the v5.78.0+
// behaviour automatically.
func TestOptions_MaxHistoryScanZeroIsUnbounded(t *testing.T) {
	var opts Options
	if opts.MaxHistoryScan != 0 {
		t.Fatalf("Options.MaxHistoryScan zero value = %d, want 0 (unbounded)", opts.MaxHistoryScan)
	}
}

// TestOptions_MaxHistoryScanRoundTrip exercises the struct field so
// any rename or removal triggers a compile failure here BEFORE it
// reaches CLI wiring or plan.go.
func TestOptions_MaxHistoryScanRoundTrip(t *testing.T) {
	opts := Options{MaxHistoryScan: 500}
	if opts.MaxHistoryScan != 500 {
		t.Fatalf("Options.MaxHistoryScan = %d, want 500", opts.MaxHistoryScan)
	}
}
