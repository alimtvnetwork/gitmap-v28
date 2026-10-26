package cmdclone

import (
	"strings"
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

func TestBuildSlugCandidates(t *testing.T) {
	candidatesVersioned := buildSlugCandidates("macro-ahk-v1")
	if len(candidatesVersioned) != 2 || candidatesVersioned[1] != "macro-ahk" {
		t.Errorf("unexpected candidates for versioned slug: %v", candidatesVersioned)
	}

	candidatesUnversioned := buildSlugCandidates("macro-ahk")
	if len(candidatesUnversioned) != 3 || candidatesUnversioned[1] != "macro-ahk-v1" {
		t.Errorf("unexpected candidates for unversioned slug: %v", candidatesUnversioned)
	}

	candidatesScoped := buildSlugCandidates("alimtvnetwork/macro-ahk-v2")
	if len(candidatesScoped) != 2 || candidatesScoped[1] != "alimtvnetwork/macro-ahk" {
		t.Errorf("unexpected candidates for scoped slug: %v", candidatesScoped)
	}
}

func TestIsRecognizedCloneShorthand(t *testing.T) {
	if !isRecognizedCloneShorthand("json") {
		t.Errorf("expected 'json' to be recognized shorthand")
	}

	if !isRecognizedCloneShorthand("csv") {
		t.Errorf("expected 'csv' to be recognized shorthand")
	}

	if !isRecognizedCloneShorthand("text") {
		t.Errorf("expected 'text' to be recognized shorthand")
	}

	if isRecognizedCloneShorthand("my-repo") {
		t.Errorf("did not expect 'my-repo' to be recognized shorthand")
	}
}

func TestDispatchCloneExecution_BareSlugNotFound(t *testing.T) {
	cf := CloneFlags{
		Source: "non-existent-test-repo-slug-xyz987",
	}
	err := dispatchCloneExecution(cf)
	if err == nil {
		t.Fatalf("expected error for non-existent bare slug, got nil")
	}

	if !strings.Contains(err.Error(), "not found or resolved") {
		t.Errorf("expected error message to mention not found or resolved, got %v", err)
	}
}

func TestValidateCFRTarget_Empty(t *testing.T) {
	_, err := validateCFRTarget("")
	if err == nil {
		t.Fatalf("expected error for empty CFR target, got nil")
	}
}

func TestValidateCFRTarget_NotFound(t *testing.T) {
	_, err := validateCFRTarget("non-existent-cfr-repo-slug-xyz987")
	if err == nil {
		t.Fatalf("expected error for non-existent CFR target, got nil")
	}
}
