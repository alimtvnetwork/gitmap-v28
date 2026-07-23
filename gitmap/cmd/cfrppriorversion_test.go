package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestCFRPPriorMaxLookbackSane guards spec/01-app/113 §2.3 — the
// lookback must stay bounded so a fresh `cfrp` on vN never fans out
// dozens of API calls. Adjust this test deliberately if the spec
// changes.
func TestCFRPPriorMaxLookbackSane(t *testing.T) {
	if constants.CFRPPriorMaxLookback < 1 || constants.CFRPPriorMaxLookback > 50 {
		t.Fatalf("CFRPPriorMaxLookback=%d out of sane [1,50] range",
			constants.CFRPPriorMaxLookback)
	}
}

// TestResolvePriorScanIdentity_NotAGitRepo confirms a non-git path
// yields the "nothing to scan" sentinel ("", 0) instead of panicking.
func TestResolvePriorScanIdentity_NotAGitRepo(t *testing.T) {
	dir := t.TempDir()

	base, current := resolvePriorScanIdentity(dir)
	if base != "" || current != 0 {
		t.Fatalf("non-git dir should yield (\"\",0); got (%q,%d)", base, current)
	}
}
