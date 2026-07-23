package cloner

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

func TestPickCloneStrategy(t *testing.T) {
	cases := []struct {
		name      string
		rec       model.ScanRecord
		useBranch bool
		branch    string
	}{
		{
			name:      "HEAD with valid branch checks out branch",
			rec:       model.ScanRecord{Branch: "main", BranchSource: gitutil.BranchSourceHEAD},
			useBranch: true,
			branch:    "main",
		},
		{
			name:      "HEAD with literal HEAD branch falls back",
			rec:       model.ScanRecord{Branch: "HEAD", BranchSource: gitutil.BranchSourceHEAD},
			useBranch: false,
		},
		{
			name:      "remote-tracking with branch checks out branch",
			rec:       model.ScanRecord{Branch: "develop", BranchSource: gitutil.BranchSourceRemoteTracking},
			useBranch: true,
			branch:    "develop",
		},
		{
			name:      "default with branch checks out branch",
			rec:       model.ScanRecord{Branch: "main", BranchSource: gitutil.BranchSourceDefault},
			useBranch: true,
			branch:    "main",
		},
		{
			name:      "detached falls back to remote default",
			rec:       model.ScanRecord{Branch: "abc1234", BranchSource: gitutil.BranchSourceDetached},
			useBranch: false,
		},
		{
			name:      "unknown falls back to remote default",
			rec:       model.ScanRecord{Branch: "main", BranchSource: gitutil.BranchSourceUnknown},
			useBranch: false,
		},
		{
			name:      "empty branchSource falls back to remote default",
			rec:       model.ScanRecord{Branch: "main", BranchSource: ""},
			useBranch: false,
		},
		{
			name:      "unrecognized branchSource falls back",
			rec:       model.ScanRecord{Branch: "main", BranchSource: "weird"},
			useBranch: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pickCloneStrategy(tc.rec)
			if got.useBranch != tc.useBranch {
				t.Fatalf("useBranch: want %v, got %v (reason=%q)", tc.useBranch, got.useBranch, got.reason)
			}
			if tc.useBranch && got.branch != tc.branch {
				t.Fatalf("branch: want %q, got %q", tc.branch, got.branch)
			}
			if got.reason == "" {
				t.Fatalf("expected non-empty reason")
			}
		})
	}
}
