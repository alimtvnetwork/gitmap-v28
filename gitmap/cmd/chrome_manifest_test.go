// Package cmd — tests for chrome backup manifest + Local State parser.
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- manifest -----------------------------------------------------------

func writeTestTarball(t *testing.T, dir string, files map[string]string) string {
	t.Helper()
	path := filepath.Join(dir, "snap.tar.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	for name, body := range files {
		_ = tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(body))})
		_, _ = tw.Write([]byte(body))
	}
	_ = tw.Close()
	_ = gz.Close()
	_ = f.Close()
	return path
}

func TestBuildAndVerifyChromeManifestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	tar1 := writeTestTarball(t, dir, map[string]string{
		"Default/Bookmarks":   `{"v":1}`,
		"Default/Preferences": `{"profile":{}}`,
	})
	if _, err := writeChromeManifest(tar1); err != nil {
		t.Fatalf("writeChromeManifest: %v", err)
	}
	ok, miss, err := verifyChromeManifest(tar1)
	if err != nil {
		t.Fatalf("verify err: %v", err)
	}
	if !ok || len(miss) != 0 {
		t.Fatalf("expected clean verify; got ok=%v miss=%v", ok, miss)
	}
}

func TestVerifyChromeManifestDetectsTampering(t *testing.T) {
	dir := t.TempDir()
	tarPath := writeTestTarball(t, dir, map[string]string{"Default/Bookmarks": "v1"})
	if _, err := writeChromeManifest(tarPath); err != nil {
		t.Fatal(err)
	}
	// Rewrite tarball with different content but same member name.
	_ = os.Remove(tarPath)
	writeTestTarball(t, dir, map[string]string{"Default/Bookmarks": "TAMPERED"})
	ok, miss, err := verifyChromeManifest(tarPath)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if ok || len(miss) == 0 {
		t.Fatalf("expected mismatch detection; got ok=%v miss=%v", ok, miss)
	}
}

func TestVerifyChromeManifestMissingSidecar(t *testing.T) {
	dir := t.TempDir()
	tarPath := writeTestTarball(t, dir, map[string]string{"a": "b"})
	_, _, err := verifyChromeManifest(tarPath)
	if err == nil || !strings.Contains(err.Error(), "manifest missing") {
		t.Fatalf("expected manifest-missing error, got %v", err)
	}
}

// --- Local State parser -------------------------------------------------

func TestParseChromeLocalStateFullDoc(t *testing.T) {
	raw := []byte(`{
		"profile": {
			"last_used": "Profile 2",
			"last_active_profiles": ["Profile 2", "Default"],
			"info_cache": {
				"Default":   {"name": "Personal", "gaia_name": "Alice"},
				"Profile 2": {"name": "Work",     "user_name": "alice@work.com"}
			}
		}
	}`)
	s, err := ParseChromeLocalState(raw)
	if err != nil {
		t.Fatal(err)
	}
	if s.LastUsed != "Profile 2" {
		t.Errorf("last_used: %q", s.LastUsed)
	}
	if got := s.DisplayNameFor("Profile 2"); got != "Work" {
		t.Errorf("display: %q", got)
	}
	if !s.Profiles["Default"].IsActive || !s.Profiles["Profile 2"].IsActive {
		t.Errorf("active flags wrong: %+v", s.Profiles)
	}
	if s.Profiles["Default"].GAIAName != "Alice" {
		t.Errorf("gaia_name lost: %+v", s.Profiles["Default"])
	}
}

func TestParseChromeLocalStateMissingInfoCache(t *testing.T) {
	raw := []byte(`{"profile":{"last_used":"Profile 9"}}`)
	s, err := ParseChromeLocalState(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got := s.DisplayNameFor("Profile 9"); got != "Profile 9" {
		t.Errorf("expected fallback to dir name, got %q", got)
	}
}

func TestParseChromeLocalStateInvalidJSON(t *testing.T) {
	if _, err := ParseChromeLocalState([]byte("not json")); err == nil {
		t.Fatal("expected json error")
	}
}

func TestParseChromeLocalStateEmptyDoc(t *testing.T) {
	s, err := ParseChromeLocalState([]byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if s.LastUsed != "" || len(s.Profiles) != 0 {
		t.Errorf("expected empty state, got %+v", s)
	}
}
