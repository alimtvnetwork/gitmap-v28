package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func TestIsFixLsRequest_Variations(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"bare_ls", []string{"ls"}, true},
		{"bare_list", []string{"list"}, true},
		{"short_flag_l", []string{"-l"}, true},
		{"long_flag_list", []string{"--list"}, true},
		{"uppercase_ls", []string{"LS"}, true},
		{"ls_with_extra", []string{"ls", "more"}, true},
		{"normal_repo", []string{"my-repo"}, false},
		{"all_cmd", []string{"all"}, false},
		{"empty_args", []string{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isFixLsRequest(tc.args)
			if actual != tc.expected {
				t.Fatalf("expected %v for %v, got %v", tc.expected, tc.args, actual)
			}
		})
	}
}

func TestIsRepoIssue_Conditions(t *testing.T) {
	testCases := []struct {
		name       string
		diag       gitutil.DirtyDiagnosis
		rs         gitutil.RepoStatus
		isConflict bool
		isLocked   bool
		expected   bool
	}{
		{
			name:     "clean_repo",
			diag:     gitutil.DirtyDiagnosis{IsDirty: false},
			rs:       gitutil.RepoStatus{},
			expected: false,
		},
		{
			name:     "dirty_repo",
			diag:     gitutil.DirtyDiagnosis{IsDirty: true, ModifiedCount: 2},
			rs:       gitutil.RepoStatus{},
			expected: true,
		},
		{
			name:       "conflict_repo",
			diag:       gitutil.DirtyDiagnosis{IsDirty: false},
			rs:         gitutil.RepoStatus{},
			isConflict: true,
			expected:   true,
		},
		{
			name:     "locked_repo",
			diag:     gitutil.DirtyDiagnosis{IsDirty: false},
			rs:       gitutil.RepoStatus{},
			isLocked: true,
			expected: true,
		},
		{
			name:     "behind_repo",
			diag:     gitutil.DirtyDiagnosis{IsDirty: false},
			rs:       gitutil.RepoStatus{Behind: 3},
			expected: true,
		},
		{
			name:     "stash_repo",
			diag:     gitutil.DirtyDiagnosis{IsDirty: false},
			rs:       gitutil.RepoStatus{StashCount: 1},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isRepoIssue(tc.diag, tc.rs, tc.isConflict, tc.isLocked)
			if actual != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestBuildIssueSummary_Formatting(t *testing.T) {
	diag := gitutil.DirtyDiagnosis{IsDirty: true, SummaryReason: "+1 modified, +2 untracked"}
	rs := gitutil.RepoStatus{Behind: 2, StashCount: 1}
	summary := buildIssueSummary(diag, rs, true, false)

	if !strings.Contains(summary, "merge conflict") {
		t.Fatalf("expected merge conflict in %s", summary)
	}
	if !strings.Contains(summary, "+1 modified") {
		t.Fatalf("expected +1 modified in %s", summary)
	}
	if !strings.Contains(summary, "behind (2)") {
		t.Fatalf("expected behind in %s", summary)
	}
}
