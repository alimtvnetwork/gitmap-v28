package clonefrom

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestRenderSummary_TalliesAllStatuses asserts the header counts
// every status correctly. Pure assembly test — uses synthesized
// Result values, no git.
func TestRenderSummary_TalliesAllStatuses(t *testing.T) {
	results := []Result{
		{Status: constants.CloneFromStatusOK, Row: Row{URL: "a"}},
		{Status: constants.CloneFromStatusOK, Row: Row{URL: "b"}},
		{Status: constants.CloneFromStatusSkipped, Row: Row{URL: "c"}, Detail: "dest exists"},
		{Status: constants.CloneFromStatusFailed, Row: Row{URL: "d"}, Detail: "boom"},
	}

	var buf bytes.Buffer
	if err := RenderSummary(&buf, results, "/tmp/r.csv"); err != nil {
		t.Fatalf("RenderSummary: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "2 ok, 1 skipped, 1 failed (4 total)") {
		t.Errorf("tally line missing or wrong:\n%s", out)
	}

	if !strings.Contains(out, "report: /tmp/r.csv") {
		t.Errorf("report path line missing")
	}
}
