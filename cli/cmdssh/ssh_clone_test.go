package cmdssh

import (
	"errors"
	"strings"
	"testing"
)

func TestParseSSHCloneOptions(t *testing.T) {
	args := []string{"my-repo", "git", "--target", "worker-1", "--ssh"}
	opts := parseSSHCloneOptions(args)

	if opts.Repo != "my-repo" {
		t.Fatalf("expected repo 'my-repo', got %q", opts.Repo)
	}
	if opts.DestPath != "git" {
		t.Fatalf("expected destPath 'git', got %q", opts.DestPath)
	}
	if opts.Target != "worker-1" {
		t.Fatalf("expected target 'worker-1', got %q", opts.Target)
	}
	if !opts.IsSSH {
		t.Fatalf("expected isSSH true")
	}
}

func TestResolveRepoName(t *testing.T) {
	cases := []struct {
		url      string
		expected string
	}{
		{"https://github.com/alimtvnetwork/gitmap-v28.git", "gitmap-v28"},
		{"git@github.com:user/service.git", "service"},
		{"https://gitlab.com/group/sub/my-project", "my-project"},
		{"simple-name", "simple-name"},
	}

	for _, c := range cases {
		got := resolveRepoName(c.url)
		if got != c.expected {
			t.Errorf("resolveRepoName(%q) = %q; want %q", c.url, got, c.expected)
		}
	}
}

func TestExpandRepoTokenToURL(t *testing.T) {
	url1 := expandRepoTokenToURL("https://github.com/foo/bar.git")
	if url1 != "https://github.com/foo/bar.git" {
		t.Fatalf("unexpected url: %s", url1)
	}

	url2 := expandRepoTokenToURL("org/repo")
	if url2 != "https://github.com/org/repo.git" {
		t.Fatalf("unexpected url: %s", url2)
	}
}

func TestResolveRemoteDestPathForNode(t *testing.T) {
	p1 := resolveRemoteDestPathForNode("my-repo", "", "linux")
	if p1 != "~/git/my-repo" {
		t.Fatalf("expected ~/git/my-repo, got %q", p1)
	}

	p2 := resolveRemoteDestPathForNode("my-repo", "git", "windows")
	if p2 != "~\\git\\my-repo" {
		t.Fatalf("expected ~\\git\\my-repo, got %q", p2)
	}

	p3 := resolveRemoteDestPathForNode("my-repo", "/opt/apps/", "linux")
	if !strings.HasSuffix(p3, "my-repo") {
		t.Fatalf("expected path ending in my-repo, got %q", p3)
	}
}

func TestIsAuthFailure(t *testing.T) {
	if !isAuthFailure("Permission denied (publickey).", nil) {
		t.Errorf("expected true for permission denied")
	}
	if !isAuthFailure("fatal: repository 'https://...' not found", nil) {
		t.Errorf("expected true for repository not found")
	}
	if !isAuthFailure("", errors.New("HTTP 401 Unauthorized")) {
		t.Errorf("expected true for HTTP 401")
	}
	if isAuthFailure("Everything up-to-date", nil) {
		t.Errorf("expected false for success output")
	}
}

func TestExtractGitHost(t *testing.T) {
	h1 := extractGitHost("https://gitlab.com/org/repo.git")
	if h1 != "gitlab.com" {
		t.Fatalf("expected gitlab.com, got %s", h1)
	}

	h2 := extractGitHost("git@github.com:org/repo.git")
	if h2 != "github.com" {
		t.Fatalf("expected github.com, got %s", h2)
	}
}
