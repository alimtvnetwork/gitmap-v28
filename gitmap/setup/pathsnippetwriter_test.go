package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWritePathSnippet_AppendsToFreshProfile(t *testing.T) {
	dir := t.TempDir()
	profile := filepath.Join(dir, ".bashrc")
	res, err := WritePathSnippet("bash", "/opt/gitmap", "test", profile)
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if res.Action != "appended" {
		t.Errorf("first write should be appended, got %q", res.Action)
	}
	got, _ := os.ReadFile(profile)
	if !strings.Contains(string(got), "managed by test") {
		t.Errorf("profile missing marker; got:\n%s", got)
	}
}

func TestWritePathSnippet_RewritesExistingBlock(t *testing.T) {
	dir := t.TempDir()
	profile := filepath.Join(dir, ".bashrc")
	original := "alias ll='ls -la'\n"
	os.WriteFile(profile, []byte(original), 0o644)

	// First write: append.
	_, err := WritePathSnippet("bash", "/opt/old", "test", profile)
	if err != nil {
		t.Fatalf("first write: %v", err)
	}
	// Second write with different dir: rewrite.
	res, err := WritePathSnippet("bash", "/opt/new", "test", profile)
	if err != nil {
		t.Fatalf("second write: %v", err)
	}
	if res.Action != "rewritten" {
		t.Errorf("second write should be rewritten, got %q", res.Action)
	}

	got, _ := os.ReadFile(profile)
	if strings.Contains(string(got), "/opt/old") {
		t.Errorf("old dir should be gone; got:\n%s", got)
	}
	if !strings.Contains(string(got), "/opt/new") {
		t.Errorf("new dir should be present; got:\n%s", got)
	}
	if !strings.HasPrefix(string(got), "alias ll='ls -la'") {
		t.Errorf("user's alias was clobbered; got:\n%s", got)
	}
}

func TestWritePathSnippet_NoopOnIdenticalSecondWrite(t *testing.T) {
	dir := t.TempDir()
	profile := filepath.Join(dir, ".bashrc")
	_, err := WritePathSnippet("bash", "/opt/gitmap", "test", profile)
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	res, err := WritePathSnippet("bash", "/opt/gitmap", "test", profile)
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	if res.Action != "noop" {
		t.Errorf("identical second write should be noop, got %q", res.Action)
	}
}

func TestWritePathSnippet_DifferentManagersCoexist(t *testing.T) {
	dir := t.TempDir()
	profile := filepath.Join(dir, ".bashrc")
	if _, err := WritePathSnippet("bash", "/opt/a", "run.sh", profile); err != nil {
		t.Fatalf("run.sh write: %v", err)
	}
	if _, err := WritePathSnippet("bash", "/opt/b", "installer", profile); err != nil {
		t.Fatalf("installer write: %v", err)
	}
	got, _ := os.ReadFile(profile)
	if !strings.Contains(string(got), "managed by run.sh") {
		t.Errorf("run.sh block missing")
	}
	if !strings.Contains(string(got), "managed by installer") {
		t.Errorf("installer block missing")
	}
}
