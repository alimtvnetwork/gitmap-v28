package cmd

import (
	"path/filepath"
	"testing"
)

func TestParseCreateParams_SingleNameWithSpaces(t *testing.T) {
	args := []string{"My Cool Project"}
	p, err := parseCreateParams(args, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "My Cool Project" {
		t.Errorf("got Name %q, want 'My Cool Project'", p.Name)
	}
	if p.Slug != "my-cool-project" {
		t.Errorf("got Slug %q, want 'my-cool-project'", p.Slug)
	}
	wantDir := filepath.Join(".", "my-cool-project")
	if p.LocalDir != wantDir {
		t.Errorf("got LocalDir %q, want %q", p.LocalDir, wantDir)
	}
	if p.IsSkipRemote {
		t.Errorf("expected IsSkipRemote=false")
	}
}

func TestParseCreateParams_NameAndSlug(t *testing.T) {
	args := []string{"My Project", "custom-slug"}
	p, err := parseCreateParams(args, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "My Project" {
		t.Errorf("got Name %q, want 'My Project'", p.Name)
	}
	if p.Slug != "custom-slug" {
		t.Errorf("got Slug %q, want 'custom-slug'", p.Slug)
	}
	wantDir := filepath.Join(".", "custom-slug")
	if p.LocalDir != wantDir {
		t.Errorf("got LocalDir %q, want %q", p.LocalDir, wantDir)
	}
}

func TestParseCreateParams_NameAndFolder(t *testing.T) {
	args := []string{"My Project", "./target/dir"}
	p, err := parseCreateParams(args, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "My Project" {
		t.Errorf("got Name %q, want 'My Project'", p.Name)
	}
	if p.LocalDir != "./target/dir" {
		t.Errorf("got LocalDir %q, want './target/dir'", p.LocalDir)
	}
	if p.Slug != "my-project" {
		t.Errorf("got Slug %q, want 'my-project'", p.Slug)
	}
}

func TestParseCreateParams_ThreeArgs(t *testing.T) {
	args := []string{"My Project", "./target/dir", "explicit-slug"}
	p, err := parseCreateParams(args, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "My Project" {
		t.Errorf("got Name %q, want 'My Project'", p.Name)
	}
	if p.LocalDir != "./target/dir" {
		t.Errorf("got LocalDir %q, want './target/dir'", p.LocalDir)
	}
	if p.Slug != "explicit-slug" {
		t.Errorf("got Slug %q, want 'explicit-slug'", p.Slug)
	}
}

func TestParseCreateParams_LocalMode(t *testing.T) {
	p1, _ := parseCreateParams([]string{"Local Repo", "--local"}, false)
	if !p1.IsSkipRemote {
		t.Errorf("expected IsSkipRemote=true with --local flag")
	}

	p2, _ := parseCreateParams([]string{"Local Repo"}, true)
	if !p2.IsSkipRemote {
		t.Errorf("expected IsSkipRemote=true with defaultLocal=true")
	}
}
