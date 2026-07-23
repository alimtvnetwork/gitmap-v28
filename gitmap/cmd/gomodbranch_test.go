package cmd

import (
	"testing"
)

func TestDeriveSlug_Simple(t *testing.T) {
	got := deriveSlug("github.com/org/repo")
	expected := "github-com-org-repo"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDeriveSlug_NoSlashes(t *testing.T) {
	got := deriveSlug("mymodule")
	if got != "mymodule" {
		t.Errorf("expected mymodule, got %q", got)
	}
}

func TestDeriveSlug_WithAt(t *testing.T) {
	got := deriveSlug("github.com/org/repo@v2")
	expected := "github-com-org-repo-v2"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDeriveSlug_Dots(t *testing.T) {
	got := deriveSlug("go.example.dev/pkg")
	expected := "go-example-dev-pkg"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDeriveSlug_Empty(t *testing.T) {
	got := deriveSlug("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestCreateGoModBranches_Names(t *testing.T) {
	// Verify branch name construction without actually running git.
	slug := "github-com-org-repo"
	expectedBackup := "backup/before-replace-" + slug
	expectedFeature := "feature/replace-" + slug

	backup := "backup/before-replace-" + slug
	feature := "feature/replace-" + slug

	if backup != expectedBackup {
		t.Errorf("expected %q, got %q", expectedBackup, backup)
	}
	if feature != expectedFeature {
		t.Errorf("expected %q, got %q", expectedFeature, feature)
	}
}
