package cmdpending

import (
	"strings"
	"testing"
)

func TestLooksLikeHashOrOpaqueID(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"aaa07946b60797706c131ca50e50ca526a44b073", true},
		{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", true},
		{"c9a646d3-9c61-4cb7-8014-442e88f87245", true},
		{"scripts-fixer", false},
		{"sdists-v9", false},
		{"main", false},
		{"", false},
		{"12345", false},
	}

	for _, tc := range tests {
		got := looksLikeHashOrOpaqueID(tc.input)
		if got != tc.want {
			t.Errorf("looksLikeHashOrOpaqueID(%q) = %v; want %v", tc.input, got, tc.want)
		}
	}
}

func TestExtractRepoNameFromRemoteURL(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://github.com/alimtvnetwork/seo-packages-v2.git", "seo-packages-v2"},
		{"https://github.com/alimtvnetwork/seo-packages-v2", "seo-packages-v2"},
		{"git@github.com:alimtvnetwork/image-generate-v2.git", "image-generate-v2"},
		{"https://token@github.com/owner/repo.git/", "repo"},
		{"ssh://git@github.com/owner/repo", "repo"},
	}

	for _, tc := range tests {
		got := extractRepoNameFromRemoteURL(tc.url)
		if got != tc.want {
			t.Errorf("extractRepoNameFromRemoteURL(%q) = %q; want %q", tc.url, got, tc.want)
		}
	}
}

func TestIsGitErrorOutput(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"fatal: not a git repository", true},
		{"fatal: ambiguous argument 'HEAD'", true},
		{"error: pathspec 'foo' did not match any file(s) known to git", true},
		{"main", false},
		{"feature/user-auth", false},
		{"v1.0.0", false},
	}

	for _, tc := range tests {
		got := isGitErrorOutput(tc.input)
		if got != tc.want {
			t.Errorf("isGitErrorOutput(%q) = %v; want %v", tc.input, got, tc.want)
		}
	}
}

func TestFormatTreeLine_ExactWidth(t *testing.T) {
	shortLine := "   ├── Option 1: gitmap sends cp repo \"wip: save changes\""
	formattedShort := formatTreeLine(shortLine)
	if len([]rune(formattedShort)) != 78 {
		t.Errorf("expected rune length 78 for short line, got %d", len([]rune(formattedShort)))
	}

	longLine := "   ├── Option 1: gitmap sends cp extremely-long-repository-name-that-definitely-exceeds-the-maximum-table-width-of-seventy-eight-characters \"wip: save changes\""
	formattedLong := formatTreeLine(longLine)
	if len([]rune(formattedLong)) != 78 {
		t.Errorf("expected rune length 78 for long line, got %d", len([]rune(formattedLong)))
	}
	if !strings.HasSuffix(formattedLong, "…") {
		t.Errorf("expected long line to end with ellipsis, got %q", formattedLong)
	}
}

func TestSafeQueryRepoBranch_NoFatalBranch(t *testing.T) {
	oldExec := currentPendingCommitsGitExecutor
	defer func() { currentPendingCommitsGitExecutor = oldExec }()

	currentPendingCommitsGitExecutor = func(dir string, args ...string) (string, error) {
		return "fatal: ambiguous argument 'HEAD': unknown revision", nil
	}

	branch := safeQueryRepoBranch(".", "main")
	if branch != "main" {
		t.Errorf("safeQueryRepoBranch returned %q; want 'main'", branch)
	}
	if isGitErrorOutput(branch) {
		t.Errorf("safeQueryRepoBranch returned git error string: %q", branch)
	}
}
