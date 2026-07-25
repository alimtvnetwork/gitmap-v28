package cmd

import "testing"

// TestResolveCloneFolderPreservesVersionSuffix locks the v6.83.0
// rule: base `clone` never rewrites the destination folder name.
// A trailing `-vN` MUST survive; only an explicit folder argument
// may override the derived name. Version flattening belongs to
// `clone-next`, not `clone`.
func TestResolveCloneFolderPreservesVersionSuffix(t *testing.T) {
	cases := []struct {
		name     string
		repoName string
		folder   string
		want     string
	}{
		{"plain repo", "scripts-fixer", "", "scripts-fixer"},
		{"v1 preserved", "codex-june-6-v1", "", "codex-june-6-v1"},
		{"v13 preserved", "wp-onboarding-v13", "", "wp-onboarding-v13"},
		{"explicit folder wins", "wp-onboarding-v13", "custom", "custom"},
		{"explicit folder over plain", "scripts-fixer", "other", "other"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := resolveCloneFolder(tc.repoName, tc.folder); got != tc.want {
				t.Fatalf("resolveCloneFolder(%q,%q): got %q want %q",
					tc.repoName, tc.folder, got, tc.want)
			}
		})
	}
}

// TestRepoNameFromURLKeepsVersionSuffix guards the upstream half of
// the same promise: URL parsing must not strip `-vN` either, and
// `.git` / trailing slashes still normalise away.
func TestRepoNameFromURLKeepsVersionSuffix(t *testing.T) {
	cases := []struct{ url, want string }{
		{"https://github.com/owner/codex-june-6-v2.git", "codex-june-6-v2"},
		{"https://github.com/owner/codex-june-6-v2/", "codex-june-6-v2"},
		{"git@github.com:owner/wp-onboarding-v13.git", "wp-onboarding-v13"},
		{"https://github.com/owner/scripts-fixer", "scripts-fixer"},
	}
	for _, tc := range cases {
		got := resolveCloneFolder(repoNameFromURL(tc.url), "")
		if got != tc.want {
			t.Fatalf("folder for %q: got %q want %q", tc.url, got, tc.want)
		}
	}
}
