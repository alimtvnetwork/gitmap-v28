package committransfer

import (
	"bytes"
	"strings"
	"testing"
)

// TestPrintPlanNoticeV6 confirms the v6.0.0 stderr notice inversion:
// when IncludeMerges is false and MergeExcluded > 0, the message
// confirms the opt-out rather than advising an opt-in.
func TestPrintPlanNoticeV6(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	// Case 1: IncludeMerges=true, MergeExcluded=0 → no notice
	planTrue := ReplayPlan{
		Commits:       []SourceCommit{{ShortSHA: "abc1234", Cleaned: "feat: x"}},
		MergeExcluded: 0,
		IncludeMerges: true,
	}
	PrintPlan(&buf, planTrue, "[test]")
	out := buf.String()
	if strings.Contains(out, "merge commits") {
		t.Errorf("IncludeMerges=true with no excluded merges should not emit merge notice; got:\n%s", out)
	}

	// Case 2: IncludeMerges=false, MergeExcluded=2 → confirmation notice
	buf.Reset()
	planFalse := ReplayPlan{
		Commits:       []SourceCommit{{ShortSHA: "def5678", Cleaned: "feat: y"}},
		MergeExcluded: 2,
		IncludeMerges: false,
	}
	PrintPlan(&buf, planFalse, "[test]")
	out = buf.String()
	want := "note: 2 merge commits excluded by --no-include-merges"
	if !strings.Contains(out, want) {
		t.Errorf("IncludeMerges=false with excluded merges should emit %q; got:\n%s", want, out)
	}
	// Ensure the OLD advisory message is gone
	if strings.Contains(out, "pass --include-merges") {
		t.Errorf("old advisory message should not appear in v6; got:\n%s", out)
	}
}
