package cmdlogviewer

import (
	"reflect"
	"testing"
)

func TestParseCommitLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected CommitItem
		isOk     bool
	}{
		{
			name:     "valid commit line",
			input:    "90f75b081057019d3d246e33f6db41ed02b09df5|90f75b081|Alim Karim|2026-10-25|fix: test commit subject",
			expected: CommitItem{
				Hash:      "90f75b081057019d3d246e33f6db41ed02b09df5",
				ShortHash: "90f75b081",
				Author:    "Alim Karim",
				Date:      "2026-10-25",
				Subject:   "fix: test commit subject",
			},
			isOk: true,
		},
		{
			name:     "commit subject with pipe",
			input:    "1111111111111111111111111111111111111111|1111111|Author|2026-01-01|feat: foo | bar | baz",
			expected: CommitItem{
				Hash:      "1111111111111111111111111111111111111111",
				ShortHash: "1111111",
				Author:    "Author",
				Date:      "2026-01-01",
				Subject:   "feat: foo | bar | baz",
			},
			isOk: true,
		},
		{
			name:     "empty line",
			input:    "",
			expected: CommitItem{},
			isOk:     false,
		},
		{
			name:     "too few fields",
			input:    "hash|short",
			expected: CommitItem{},
			isOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, isOk := parseCommitLine(tt.input)
			if isOk != tt.isOk {
				t.Fatalf("parseCommitLine() isOk = %v, want %v", isOk, tt.isOk)
			}
			if isOk && !reflect.DeepEqual(item, tt.expected) {
				t.Errorf("parseCommitLine() = %+v, want %+v", item, tt.expected)
			}
		})
	}
}

func TestFormatAction(t *testing.T) {
	cases := []struct {
		code     string
		expected string
	}{
		{"A", "Added"},
		{"M", "Modified"},
		{"D", "Deleted"},
		{"R100", "Renamed"},
		{"C75", "Copied"},
		{"unknown", "Modified"},
	}

	for _, c := range cases {
		got := formatAction(c.code)
		if got != c.expected {
			t.Errorf("formatAction(%q) = %q, want %q", c.code, got, c.expected)
		}
	}
}

func TestParseNameStatusLines(t *testing.T) {
	raw := "A\tnew_file.go\nM\texisting.go\nR100\told_name.go\tnew_name.go\nD\tdeleted.go"
	numMap := map[string]string{
		"new_file.go": "+10, -0",
		"existing.go": "+5, -2",
		"new_name.go": "+1, -1",
		"deleted.go":  "+0, -20",
	}

	files := parseNameStatusLines(raw, numMap)
	if len(files) != 4 {
		t.Fatalf("expected 4 files, got %d", len(files))
	}

	if files[0].Action != "Added" || files[0].Path != "new_file.go" || files[0].Stats != "+10, -0" {
		t.Errorf("file[0] mismatch: %+v", files[0])
	}
	if files[2].Action != "Renamed" || files[2].Path != "new_name.go" || files[2].OldPath != "old_name.go" {
		t.Errorf("file[2] rename mismatch: %+v", files[2])
	}
}

func TestBuildLogArgs(t *testing.T) {
	argsAll := buildLogArgs(50, "")
	if len(argsAll) < 4 || argsAll[1] != "-n50" {
		t.Errorf("unexpected buildLogArgs for all files: %v", argsAll)
	}

	argsPath := buildLogArgs(25, "some/file.go")
	if len(argsPath) < 6 || argsPath[4] != "--" || argsPath[5] != "some/file.go" {
		t.Errorf("unexpected buildLogArgs with path: %v", argsPath)
	}
}

func TestLiveGitLogQueries(t *testing.T) {
	repoRoot, err := findRepoRoot(".")
	if err != nil {
		t.Skip("skipping live test: not in a git repo")
	}

	commits, err := QueryCommits(repoRoot, ".", 5)
	if err != nil {
		t.Fatalf("QueryCommits failed: %v", err)
	}

	if len(commits) == 0 {
		t.Fatal("expected at least 1 commit in current repo")
	}

	first := commits[0]
	if first.Hash == "" || first.ShortHash == "" || first.Author == "" {
		t.Errorf("commit missing expected fields: %+v", first)
	}

	detail, err := QueryCommitDetail(repoRoot, first.Hash)
	if err != nil {
		t.Fatalf("QueryCommitDetail failed: %v", err)
	}

	if detail.Hash != first.Hash {
		t.Errorf("expected detail hash %q, got %q", first.Hash, detail.Hash)
	}
}
