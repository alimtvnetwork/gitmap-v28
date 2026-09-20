package prclean

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/prdb"
)

func TestFormatMergedTime(t *testing.T) {
	if got := formatMergedTime(0); got != "N/A" {
		t.Errorf("expected 'N/A' for zero timestamp, got %q", got)
	}

	ts := int64(1710000000)
	if got := formatMergedTime(ts); got != "1710000000 (epoch)" {
		t.Errorf("expected '1710000000 (epoch)', got %q", got)
	}
}

func TestProcessCleanBranches_Empty(t *testing.T) {
	err := processCleanBranches(nil, ".", []prdb.PrBranchRecord{}, true)
	if err != nil {
		t.Errorf("expected nil error for empty branch list, got %v", err)
	}
}
