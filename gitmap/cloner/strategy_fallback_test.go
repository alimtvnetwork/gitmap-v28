package cloner

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// TestApplyDefaultBranchFallback_EmptyFallbackIsNoop pins the
// backwards-compat guarantee: no --default-branch flag means the
// returned slice is the SAME slice (not a copy), every record is
// identical, and BranchSource=unknown still routes through the
// "remote default HEAD" path. Existing manifests must be unaffected.
func TestApplyDefaultBranchFallback_EmptyFallbackIsNoop(t *testing.T) {
	in := []model.ScanRecord{
		{Branch: "", BranchSource: gitutil.BranchSourceUnknown},
		{Branch: "feat-x", BranchSource: gitutil.BranchSourceHEAD},
	}
	out := applyDefaultBranchFallback(in, "")
	if &out[0] != &in[0] {
		t.Fatalf("empty fallback must return the same slice header (no copy)")
	}
}

// TestApplyDefaultBranchFallback_RewritesUntrustedRows verifies the
// core behavior: every record that pickCloneStrategy would skip -b for
// (empty Branch, detached, unknown, empty source, unrecognized source,
// HEAD-with-literal-HEAD) is rewritten to the trusted default path.
func TestApplyDefaultBranchFallback_RewritesUntrustedRows(t *testing.T) {
	cases := []struct {
		name string
		rec  model.ScanRecord
	}{
		{"empty branch + unknown source", model.ScanRecord{Branch: "", BranchSource: gitutil.BranchSourceUnknown}},
		{"detached", model.ScanRecord{Branch: "abc1234", BranchSource: gitutil.BranchSourceDetached}},
		{"empty source", model.ScanRecord{Branch: "main", BranchSource: ""}},
		{"unrecognized source", model.ScanRecord{Branch: "main", BranchSource: "weird"}},
		{"HEAD with literal HEAD branch", model.ScanRecord{Branch: "HEAD", BranchSource: gitutil.BranchSourceHEAD}},
		{"HEAD with empty branch", model.ScanRecord{Branch: "", BranchSource: gitutil.BranchSourceHEAD}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := applyDefaultBranchFallback([]model.ScanRecord{tc.rec}, "trunk")
			if len(out) != 1 {
				t.Fatalf("expected 1 record, got %d", len(out))
			}
			if out[0].Branch != "trunk" {
				t.Errorf("Branch: want %q, got %q", "trunk", out[0].Branch)
			}
			if out[0].BranchSource != gitutil.BranchSourceDefault {
				t.Errorf("BranchSource: want %q, got %q",
					gitutil.BranchSourceDefault, out[0].BranchSource)
			}
			// And the strategy now picks -b trunk, closing the loop.
			strategy := pickCloneStrategy(out[0])
			if !strategy.useBranch || strategy.branch != "trunk" {
				t.Errorf("strategy after fallback: want useBranch=true branch=trunk, got %+v", strategy)
			}
			// Breadcrumb must be present so audits can tell what happened.
			if !strings.Contains(out[0].Notes, "default-branch fallback applied: trunk") {
				t.Errorf("Notes missing breadcrumb: %q", out[0].Notes)
			}
		})
	}
}

// TestApplyDefaultBranchFallback_LeavesTrustedRowsUntouched confirms
// records that already pick -b are returned bit-identical: no Branch
// rewrite, no Notes mutation. The fallback must be additive only.
func TestApplyDefaultBranchFallback_LeavesTrustedRowsUntouched(t *testing.T) {
	in := []model.ScanRecord{
		{Branch: "main", BranchSource: gitutil.BranchSourceHEAD, Notes: "scanned"},
		{Branch: "develop", BranchSource: gitutil.BranchSourceRemoteTracking},
		{Branch: "release", BranchSource: gitutil.BranchSourceDefault},
	}
	out := applyDefaultBranchFallback(in, "trunk")
	for i, rec := range out {
		if rec.Branch != in[i].Branch {
			t.Errorf("[%d] trusted Branch mutated: %q → %q", i, in[i].Branch, rec.Branch)
		}
		if rec.BranchSource != in[i].BranchSource {
			t.Errorf("[%d] trusted BranchSource mutated: %q → %q",
				i, in[i].BranchSource, rec.BranchSource)
		}
		if rec.Notes != in[i].Notes {
			t.Errorf("[%d] trusted Notes mutated: %q → %q", i, in[i].Notes, rec.Notes)
		}
	}
}

// TestApplyDefaultBranchFallback_DoesNotMutateInputSlice guards the
// "no aliasing the caller's data" promise. The input must be safe to
// reuse after the call (the cloner historically passes the slice into
// progress / cache / summary in parallel).
func TestApplyDefaultBranchFallback_DoesNotMutateInputSlice(t *testing.T) {
	in := []model.ScanRecord{
		{Branch: "", BranchSource: gitutil.BranchSourceUnknown, Notes: "orig"},
	}
	_ = applyDefaultBranchFallback(in, "trunk")
	if in[0].Branch != "" || in[0].BranchSource != gitutil.BranchSourceUnknown || in[0].Notes != "orig" {
		t.Fatalf("input slice mutated: %+v", in[0])
	}
}
