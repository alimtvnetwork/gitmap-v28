package clonenext

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadLocalRepoState_PlainRepo(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(filepath.Join(gitDir, "refs", "heads"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeFile(t, filepath.Join(gitDir, "HEAD"), "ref: refs/heads/main\n")
	writeFile(t, filepath.Join(gitDir, "refs", "heads", "main"), "abc123def456\n")
	writeFile(t, filepath.Join(gitDir, "config"), `[core]
	repositoryformatversion = 0
[remote "origin"]
	url = https://github.com/acme/alpha-v3.git
	fetch = +refs/heads/*:refs/remotes/origin/*
`)

	state, err := ReadLocalRepoState(root)
	if err != nil {
		t.Fatalf("ReadLocalRepoState: %v", err)
	}
	if state.OriginURL != "https://github.com/acme/alpha-v3.git" {
		t.Errorf("OriginURL = %q", state.OriginURL)
	}
	if state.HeadSHA != "abc123def456" {
		t.Errorf("HeadSHA = %q, want abc123def456", state.HeadSHA)
	}
}

func TestReadLocalRepoState_NotARepo(t *testing.T) {
	root := t.TempDir()
	_, err := ReadLocalRepoState(root)
	if !errors.Is(err, ErrNotAGitRepo) {
		t.Errorf("err = %v, want ErrNotAGitRepo", err)
	}
}

func TestReadLocalRepoState_DetachedHead(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeFile(t, filepath.Join(gitDir, "HEAD"), "deadbeefcafe1234\n")
	writeFile(t, filepath.Join(gitDir, "config"), `[remote "origin"]
	url = git@github.com:acme/alpha-v3.git
`)

	state, err := ReadLocalRepoState(root)
	if err != nil {
		t.Fatalf("ReadLocalRepoState: %v", err)
	}
	if state.HeadSHA != "deadbeefcafe1234" {
		t.Errorf("HeadSHA = %q, want detached SHA", state.HeadSHA)
	}
	if state.OriginURL != "git@github.com:acme/alpha-v3.git" {
		t.Errorf("OriginURL = %q", state.OriginURL)
	}
}

func TestReadLocalRepoState_NoOriginRemote(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeFile(t, filepath.Join(gitDir, "HEAD"), "ref: refs/heads/main\n")
	writeFile(t, filepath.Join(gitDir, "config"), `[core]
	repositoryformatversion = 0
[remote "upstream"]
	url = https://github.com/acme/other.git
`)

	state, err := ReadLocalRepoState(root)
	if err != nil {
		t.Fatalf("ReadLocalRepoState: %v", err)
	}
	if state.OriginURL != "" {
		t.Errorf("OriginURL = %q, want empty (no origin remote)", state.OriginURL)
	}
}

func TestExtractOriginURL_IgnoresOtherRemotes(t *testing.T) {
	cfg := `[remote "fork"]
	url = https://github.com/me/x.git
[remote "origin"]
	url = https://github.com/them/x.git
[remote "upstream"]
	url = https://github.com/upstream/x.git
`
	got := extractOriginURL(cfg)
	want := "https://github.com/them/x.git"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
