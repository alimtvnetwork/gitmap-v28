package cmdclone

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// TestIsDirectURL_HTTPS verifies HTTPS URLs are detected.
func TestIsDirectURL_HTTPS(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"https://github.com/user/repo.git", true},
		{"https://github.com/user/repo", true},
		{"https://gitlab.com/org/project.git", true},
		{"HTTPS://GITHUB.COM/USER/REPO.git", true},
	}

	for _, tc := range cases {
		if got := isDirectURL(tc.input); got != tc.want {
			t.Errorf("isDirectURL(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

// TestIsDirectURL_HTTP verifies plain HTTP URLs are detected.
func TestIsDirectURL_HTTP(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"http://github.com/user/repo.git", true},
		{"HTTP://EXAMPLE.COM/repo.git", true},
	}

	for _, tc := range cases {
		if got := isDirectURL(tc.input); got != tc.want {
			t.Errorf("isDirectURL(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

// TestIsDirectURL_SSH verifies SSH URLs are detected.
func TestIsDirectURL_SSH(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"git@github.com:user/repo.git", true},
		{"git@gitlab.com:org/project.git", true},
	}

	for _, tc := range cases {
		if got := isDirectURL(tc.input); got != tc.want {
			t.Errorf("isDirectURL(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

// TestIsDirectURL_NonURL verifies file paths and shorthands are rejected.
func TestIsDirectURL_NonURL(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"json", false},
		{"csv", false},
		{"text", false},
		{"gitmap.json", false},
		{"./output/gitmap.csv", false},
		{".gitmap/output/gitmap.json", false},
		{"C:\\repos\\output.json", false},
		{"/home/user/repos.txt", false},
		{"", false},
	}

	for _, tc := range cases {
		if got := isDirectURL(tc.input); got != tc.want {
			t.Errorf("isDirectURL(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

// TestRepoNameFromURL_HTTPS verifies name extraction from HTTPS URLs.
func TestRepoNameFromURL_HTTPS(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://github.com/alimtvnetwork/wp-alim.git", "wp-alim"},
		{"https://github.com/user/my-repo.git", "my-repo"},
		{"https://github.com/user/repo", "repo"},
		{"https://gitlab.com/org/sub/project.git", "project"},
	}

	for _, tc := range cases {
		if got := repoNameFromURL(tc.input); got != tc.want {
			t.Errorf("repoNameFromURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// TestRepoNameFromURL_SSH verifies name extraction from SSH URLs.
func TestRepoNameFromURL_SSH(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"git@github.com:user/my-repo.git", "my-repo"},
		{"git@github.com:org/project.git", "project"},
		{"git@gitlab.com:group/sub/repo.git", "repo"},
	}

	for _, tc := range cases {
		if got := repoNameFromURL(tc.input); got != tc.want {
			t.Errorf("repoNameFromURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// TestRepoNameFromURL_EdgeCases verifies edge case handling.
func TestRepoNameFromURL_EdgeCases(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://github.com/user/repo.git.git", "repo.git"},
		{"my-repo.git", "my-repo"},
		{"my-repo", "my-repo"},
		{"", ""},
	}

	for _, tc := range cases {
		if got := repoNameFromURL(tc.input); got != tc.want {
			t.Errorf("repoNameFromURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// TestRepoNameFromURL_TrailingSlash is a regression guard for the
// "downloader gets confused with trailing slash" bug: a URL ending
// in `/` (or `\`, or `.git/`) used to collapse the basename to ""
// which made the clone target equal to CWD and triggered the
// destructive "target exists" replace flow.
func TestRepoNameFromURL_TrailingSlash(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://github.com/alimtvnetwork/gitmap-v28/", "gitmap-v28"},
		{"https://github.com/alimtvnetwork/gitmap-v28.git/", "gitmap-v28"},
		{"https://github.com/alimtvnetwork/gitmap-v28///", "gitmap-v28"},
		{"https://github.com/alimtvnetwork/gitmap-v28\\", "gitmap-v28"},
		{"git@github.com:owner/repo.git/", "repo"},
		{"git@github.com:owner/repo/", "repo"},
		{"ssh://git@github.com/owner/repo/", "repo"},
	}

	for _, tc := range cases {
		got := repoNameFromURL(tc.input)
		if got != tc.want {
			t.Errorf("repoNameFromURL(%q) = %q, want %q (regression: trailing slash collapses basename)", tc.input, got, tc.want)
		}

		if got == "" {
			t.Errorf("repoNameFromURL(%q) returned empty — would target CWD and trigger replace flow", tc.input)
		}
	}
}

func TestParseCloneFlags_ExcludeAndList(t *testing.T) {
	cf := ParseCloneFlags([]string{"a.json", "ls", "--exclude", "repoStarts,repostarts2nd"})
	if !cf.IsListOnly {
		t.Errorf("expected IsListOnly=true, got false")
	}
	if cf.Source != "a.json" {
		t.Errorf("expected Source='a.json', got %q", cf.Source)
	}
	if cf.ExcludeFilter != "repoStarts,repostarts2nd" {
		t.Errorf("expected ExcludeFilter='repoStarts,repostarts2nd', got %q", cf.ExcludeFilter)
	}
}

func TestParseCloneFlags_SSHAndPositional(t *testing.T) {
	cf := ParseCloneFlags([]string{"--ssh", "my-manifest.json"})
	if !cf.UseSSH {
		t.Errorf("expected UseSSH=true, got false")
	}
	if cf.Source != "my-manifest.json" {
		t.Errorf("expected Source='my-manifest.json', got %q", cf.Source)
	}
}

func TestParseCloneFlags_BareListOnly(t *testing.T) {
	cf := ParseCloneFlags([]string{"ls"})
	if !cf.IsListOnly {
		t.Errorf("expected IsListOnly=true, got false")
	}
	if cf.Source != "" {
		t.Errorf("expected Source='', got %q", cf.Source)
	}
}

func TestFilterRecordsByExclude(t *testing.T) {
	records := []model.ScanRecord{
		{RepoName: "repoStartsOne", RelativePath: "repoStartsOne"},
		{RepoName: "repostarts2nd", RelativePath: "repostarts2nd"},
		{RepoName: "keepThisRepo", RelativePath: "keepThisRepo"},
	}
	filtered := filterRecordsByExclude(records, "repoStarts,repostarts2nd")
	if len(filtered) != 1 {
		t.Fatalf("expected 1 record after exclude, got %d", len(filtered))
	}
	if filtered[0].RepoName != "keepThisRepo" {
		t.Errorf("expected keepThisRepo, got %s", filtered[0].RepoName)
	}
}

func TestConvertRecordsToSSH(t *testing.T) {
	records := []model.ScanRecord{
		{HTTPSUrl: "https://github.com/user/myrepo.git", RelativePath: "myrepo"},
	}
	converted := convertRecordsToSSH(records)
	if converted[0].Transport != "ssh" {
		t.Errorf("expected Transport='ssh', got %s", converted[0].Transport)
	}
	if converted[0].SSHUrl != "git@github.com:user/myrepo.git" {
		t.Errorf("expected git@github.com:user/myrepo.git, got %s", converted[0].SSHUrl)
	}
}
