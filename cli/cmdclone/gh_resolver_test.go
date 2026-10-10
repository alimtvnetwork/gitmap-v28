package cmdclone

import (
	"testing"
)

func TestIsGitHubOwnerRepo(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"alimtvnetwork/seo-packages-v2", true},
		{"owner/repo", true},
		{"owner/repo.git", true},
		{"owner-name/repo_name.v2", true},
		{"seo-packages-v2", false},
		{"https://github.com/owner/repo", false},
		{"./local/path", false},
		{"../local/path", false},
		{"/abs/path", false},
		{"C:/work/repo", false},
		{"", false},
		{"a/b/c", false},
	}

	for _, tc := range tests {
		got := isGitHubOwnerRepo(tc.input)
		if got != tc.want {
			t.Errorf("isGitHubOwnerRepo(%q) = %v; want %v", tc.input, got, tc.want)
		}
	}
}

func TestIsFullURLOrPath(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"https://github.com/owner/repo", true},
		{"http://github.com/owner/repo", true},
		{"git@github.com:owner/repo.git", true},
		{"ssh://git@github.com/owner/repo", true},
		{"file:///path/to/repo", true},
		{"./local/repo", true},
		{"../local/repo", true},
		{"/var/git/repo", true},
		{"\\server\\share\\repo", true},
		{"alimtvnetwork/seo-packages-v2", false},
		{"seo-packages-v2", false},
	}

	for _, tc := range tests {
		got := isFullURLOrPath(tc.input)
		if got != tc.want {
			t.Errorf("isFullURLOrPath(%q) = %v; want %v", tc.input, got, tc.want)
		}
	}
}

func TestResolveRepoSlug_OwnerRepoFallback(t *testing.T) {
	input := "custom-org/some-unknown-repo-xyz-123"
	got := ResolveRepoSlug(input)
	want := "https://github.com/custom-org/some-unknown-repo-xyz-123"
	if got != want {
		t.Errorf("ResolveRepoSlug(%q) = %q; want %q", input, got, want)
	}
}

func TestFilterRedundantSuggestions(t *testing.T) {
	suggs := []string{"seo-packages-v2", "seo-packages-v1", "seo-other"}
	filtered := filterRedundantSuggestions(suggs, "seo-packages-v2", "alimtvnetwork/seo-packages-v2")
	for _, s := range filtered {
		if s == "seo-packages-v2" {
			t.Errorf("expected seo-packages-v2 to be filtered out, got %v", filtered)
		}
	}
	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered suggestions, got %d (%v)", len(filtered), filtered)
	}
}
