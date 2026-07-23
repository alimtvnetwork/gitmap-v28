package formatter

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// TestValidateRecords_TableDriven exercises every issue path the writer
// will surface, plus the happy path.
func TestValidateRecords_TableDriven(t *testing.T) {
	cases := []struct {
		name         string
		records      []model.ScanRecord
		wantIssueCnt int
		wantField    string // first issue's Field, "" when wantIssueCnt == 0
	}{
		{
			name: "happy_path_no_issues",
			records: []model.ScanRecord{{
				Slug:         "my-repo",
				RepoName:     "my-repo",
				HTTPSUrl:     "https://github.com/u/my-repo.git",
				RelativePath: "my-repo",
			}},
			wantIssueCnt: 0,
		},
		{
			name: "missing_repo_name",
			records: []model.ScanRecord{{
				HTTPSUrl:     "https://github.com/u/r.git",
				RelativePath: "r",
			}},
			wantIssueCnt: 1,
			wantField:    "RepoName",
		},
		{
			name: "missing_relative_path",
			records: []model.ScanRecord{{
				RepoName: "r",
				HTTPSUrl: "https://github.com/u/r.git",
			}},
			wantIssueCnt: 1,
			wantField:    "RelativePath",
		},
		{
			name: "missing_both_urls",
			records: []model.ScanRecord{{
				RepoName:     "r",
				RelativePath: "r",
			}},
			wantIssueCnt: 1,
			wantField:    "HTTPSUrl|SSHUrl",
		},
		{
			name: "ssh_url_alone_is_fine",
			records: []model.ScanRecord{{
				RepoName:     "r",
				RelativePath: "r",
				SSHUrl:       "git@github.com:u/r.git",
			}},
			wantIssueCnt: 0,
		},
		{
			name: "slug_mismatch",
			records: []model.ScanRecord{{
				Slug:         "wrong-slug",
				RepoName:     "MyRepo",
				HTTPSUrl:     "https://github.com/u/MyRepo.git",
				RelativePath: "MyRepo",
			}},
			wantIssueCnt: 1,
			wantField:    "Slug",
		},
		{
			name:    "multiple_issues_one_record",
			records: []model.ScanRecord{{
				// Missing RepoName, RelativePath, AND any URL → 3 issues.
			}},
			wantIssueCnt: 3,
			wantField:    "RepoName",
		},
		{
			name:         "empty_input_no_issues",
			records:      nil,
			wantIssueCnt: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateRecords(tc.records)
			if len(got) != tc.wantIssueCnt {
				t.Fatalf("got %d issues, want %d (issues=%v)", len(got), tc.wantIssueCnt, got)
			}
			if tc.wantIssueCnt > 0 && got[0].Field != tc.wantField {
				t.Errorf("first issue Field = %q, want %q", got[0].Field, tc.wantField)
			}
		})
	}
}

// TestValidationIssue_StringHasContext makes sure operators see row index,
// repo identifier, and the offending field at a glance.
func TestValidationIssue_StringHasContext(t *testing.T) {
	issue := ValidationIssue{
		RowIndex: 3,
		RepoName: "alpha",
		Field:    "RelativePath",
		Reason:   "required field is empty",
	}
	out := issue.String()
	for _, want := range []string{"row 3", "alpha", "RelativePath", "empty"} {
		if !strings.Contains(out, want) {
			t.Errorf("issue.String() = %q, missing %q", out, want)
		}
	}
}

// TestValidationIssue_StringFallsBackToUnnamed verifies the renderer
// stays useful when the offending record has no RepoName.
func TestValidationIssue_StringFallsBackToUnnamed(t *testing.T) {
	issue := ValidationIssue{RowIndex: 0, Field: "RepoName", Reason: "required field is empty"}
	if !strings.Contains(issue.String(), "<unnamed>") {
		t.Errorf("expected <unnamed> placeholder in %q", issue.String())
	}
}
