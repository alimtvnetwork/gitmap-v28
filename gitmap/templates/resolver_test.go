package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// withTempHome points os.UserHomeDir at a temp dir for the duration of a test.
func withTempHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp) // Windows

	return tmp
}

func TestUserDirUnderHome(t *testing.T) {
	home := withTempHome(t)
	got, err := UserDir()
	if err != nil {
		t.Fatalf("UserDir error: %v", err)
	}
	want := filepath.Join(home, userTemplatesDirName, userTemplatesSubdir)
	if got != want {
		t.Fatalf("UserDir = %q, want %q", got, want)
	}
}

func TestEnsureUserDirCreates(t *testing.T) {
	withTempHome(t)
	dir, err := EnsureUserDir()
	if err != nil {
		t.Fatalf("EnsureUserDir error: %v", err)
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Fatalf("EnsureUserDir did not create %s: %v", dir, err)
	}
}

func TestResolveUserOverlayWinsOverEmbed(t *testing.T) {
	withTempHome(t)
	dir, err := EnsureUserDir()
	if err != nil {
		t.Fatalf("EnsureUserDir: %v", err)
	}

	// Place a fake user-overlay template.
	overlay := filepath.Join(dir, kindIgnore, "overlay-only"+templateExtIgnore)
	if mkErr := os.MkdirAll(filepath.Dir(overlay), 0o755); mkErr != nil {
		t.Fatalf("mkdir: %v", mkErr)
	}
	body := []byte("# user override\nfoo\n")
	if wErr := os.WriteFile(overlay, body, 0o644); wErr != nil {
		t.Fatalf("write: %v", wErr)
	}

	r, err := Resolve(kindIgnore, "overlay-only")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if r.Source != SourceUser {
		t.Fatalf("Source = %v, want SourceUser", r.Source)
	}
	if !bytes.Equal(r.Content, body) {
		t.Fatalf("Content mismatch: got %q", r.Content)
	}
}

func TestResolveMissingReturnsNotFound(t *testing.T) {
	withTempHome(t)
	_, err := Resolve(kindIgnore, "definitely-not-a-language")
	if err == nil {
		t.Fatal("expected error for missing template, got nil")
	}
}
