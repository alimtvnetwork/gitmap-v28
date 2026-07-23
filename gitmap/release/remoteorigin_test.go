package release_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// TestParseHTTPSURL verifies HTTPS remote URLs are parsed correctly.
func TestParseHTTPSURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		owner string
		repo  string
	}{
		{"standard", "https://github.com/octocat/hello-world.git", "octocat", "hello-world"},
		{"no .git suffix", "https://github.com/octocat/hello-world", "octocat", "hello-world"},
		{"deep host", "https://gitlab.example.com/org/project.git", "org", "project"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, err := release.ParseGitURLExported(tc.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if owner != tc.owner {
				t.Errorf("owner: expected %q, got %q", tc.owner, owner)
			}
			if repo != tc.repo {
				t.Errorf("repo: expected %q, got %q", tc.repo, repo)
			}
		})
	}
}

// TestParseSSHURL verifies SSH remote URLs are parsed correctly.
func TestParseSSHURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		owner string
		repo  string
	}{
		{"standard", "git@github.com:octocat/hello-world.git", "octocat", "hello-world"},
		{"no .git suffix", "git@github.com:octocat/hello-world", "octocat", "hello-world"},
		{"custom host", "git@gitlab.example.com:org/project.git", "org", "project"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, err := release.ParseGitURLExported(tc.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if owner != tc.owner {
				t.Errorf("owner: expected %q, got %q", tc.owner, owner)
			}
			if repo != tc.repo {
				t.Errorf("repo: expected %q, got %q", tc.repo, repo)
			}
		})
	}
}

// TestParseInvalidURL verifies invalid URLs return errors.
func TestParseInvalidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty string", ""},
		{"plain text", "not-a-url"},
		{"ftp protocol", "ftp://github.com/owner/repo"},
		{"https too few parts", "https://github.com"},
		{"ssh no colon", "git@github.com"},
		{"ssh no slash", "git@github.com:noslash"},
		{"whitespace only", "   "},
		{"https bare scheme", "https://"},
		{"https scheme + slash", "https:/"},
		{"http (not https)", "http://github.com/owner/repo"},
		{"file scheme", "file:///srv/git/repo.git"},
		{"https just host slash", "https://github.com/"},
		{"ssh empty path", "git@github.com:"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := release.ParseGitURLExported(tc.url)
			if err == nil {
				t.Error("expected error for invalid URL, got nil")
			}
		})
	}
}

// TestParseHTTPSURLEdgeCases locks valid HTTPS variations the
// parser must continue to accept (port suffix on host, multi-level
// path / subgroups, hyphenated and dot-bearing repo names). Each
// case asserts the exact owner/repo extracted so a future refactor
// of parseHTTPSURL cannot quietly shift the trailing-segment rule.
func TestParseHTTPSURLEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		owner, repo string
	}{
		{"explicit port", "https://github.com:443/octocat/hello-world.git", "octocat", "hello-world"},
		{"gitlab subgroup", "https://gitlab.com/group/sub/repo.git", "sub", "repo"},
		{"deep subgroup", "https://gitlab.com/a/b/c/d/repo.git", "d", "repo"},
		{"hyphen + underscore in repo", "https://github.com/my-org/my_repo-name.git", "my-org", "my_repo-name"},
		{"dot in repo (non .git)", "https://github.com/octocat/site.io.git", "octocat", "site.io"},
		{"uppercase segments", "https://GitHub.com/Octo/Hello.git", "Octo", "Hello"},
		{"no .git, deep host", "https://git.self-hosted.example/team/proj", "team", "proj"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, err := release.ParseGitURLExported(tc.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if owner != tc.owner || repo != tc.repo {
				t.Errorf("got (owner=%q, repo=%q), want (%q, %q)",
					owner, repo, tc.owner, tc.repo)
			}
		})
	}
}

// TestParseSSHURLEdgeCases locks valid SSH variations: non-default
// users (CI bots, mirroring accounts), custom hosts, and repo names
// with hyphens/underscores/dots. The parser splits on the LAST
// colon, so an unusual user prefix (e.g. `bot-account@`) must still
// surface owner/repo from the path that follows the colon.
func TestParseSSHURLEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		owner, repo string
	}{
		{"non-git user", "deploy@github.com:octocat/hello.git", "octocat", "hello"},
		{"hyphenated user", "ci-bot@gitlab.example.com:org/proj.git", "org", "proj"},
		{"dotted host", "git@source.code.example.io:team/tool.git", "team", "tool"},
		{"underscored repo", "git@github.com:org/my_tool.git", "org", "my_tool"},
		{"dot in repo (non .git)", "git@github.com:octocat/site.io.git", "octocat", "site.io"},
		{"uppercase repo", "git@github.com:Org/RepoName.git", "Org", "RepoName"},
		{"no .git suffix + host alias", "git@gh:org/repo", "org", "repo"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, err := release.ParseGitURLExported(tc.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if owner != tc.owner || repo != tc.repo {
				t.Errorf("got (owner=%q, repo=%q), want (%q, %q)",
					owner, repo, tc.owner, tc.repo)
			}
		})
	}
}
