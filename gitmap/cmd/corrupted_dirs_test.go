package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsCorruptedDirName(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantIsBad  bool
		wantReason string
	}{
		{"literal tilde", "~", true, "literal tilde directory"},
		{"ansi escape prefix", "\x1b[36mtest", true, "contains ANSI escape sequence"},
		{"ansi escape octal", "\033[36mtest", true, "contains ANSI escape sequence"},
		{"embedded newline", "foo\nbar", true, "contains embedded newline or carriage return"},
		{"embedded carriage return", "foo\rbar", true, "contains embedded newline or carriage return"},
		{"quick installer lowercase", "my quick installer dir", true, "contains installer banner prompt text"},
		{"gitmap installer", "gitmap installer dir", true, "contains installer banner prompt text"},
		{"default prompt text", "Default: /home/user", true, "contains installer prompt text"},
		{"choose install folder", "Choose install folder", true, "contains installer prompt text"},
		{"clean directory name", "normal_folder", false, ""},
		{"gitmap standard folder", ".gitmap", false, ""},
		{"local bin folder", "bin", false, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isBad, reason := isCorruptedDirName(tc.input)
			if isBad != tc.wantIsBad {
				t.Errorf("got isBad %v, want %v for %q", isBad, tc.wantIsBad, tc.input)
			}
			if tc.wantIsBad && reason == "" {
				t.Errorf("expected non-empty reason for %q", tc.input)
			}
		})
	}
}

func TestIsProtectedPath(t *testing.T) {
	home := "/home/testuser"
	tests := []struct {
		path          string
		wantProtected bool
	}{
		{"/", true},
		{".", true},
		{"/tmp", true},
		{"/usr", true},
		{"/usr/local", true},
		{"/usr/local/bin", true},
		{home, true},
		{"/home/testuser/corrupted_folder", false},
		{"/home/testuser/.local/bin/corrupted", false},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			got := isProtectedPath(tc.path, home)
			if got != tc.wantProtected {
				t.Errorf("isProtectedPath(%q) = %v, want %v", tc.path, got, tc.wantProtected)
			}
		})
	}
}

func TestIsInsidePath(t *testing.T) {
	base := filepath.Join("test", "dir")
	sub := filepath.Join(base, "child")
	outside := filepath.Join("test", "other")

	if !isInsidePath(base, base) {
		t.Errorf("expected isInsidePath to be true for identical path")
	}
	if !isInsidePath(sub, base) {
		t.Errorf("expected isInsidePath to be true for child path")
	}
	if isInsidePath(outside, base) {
		t.Errorf("expected isInsidePath to be false for outside path")
	}
}

func TestRecoverAndCleanCorruptedDir(t *testing.T) {
	tempRoot := t.TempDir()
	corruptedDir := filepath.Join(tempRoot, "gitmap quick installer folder")
	if err := os.MkdirAll(corruptedDir, 0755); err != nil {
		t.Fatalf("failed to create corrupted dir: %v", err)
	}

	dummyBin := filepath.Join(corruptedDir, "gitmap")
	if err := os.WriteFile(dummyBin, []byte("#!/bin/sh\necho test\n"), 0755); err != nil {
		t.Fatalf("failed to write dummy binary: %v", err)
	}

	targetDir := filepath.Join(tempRoot, "recovered_bin")
	info := CorruptedDirInfo{
		Path:     corruptedDir,
		Name:     "gitmap quick installer folder",
		HasFiles: true,
	}

	rec, err := RecoverCorruptedDirAssets(info, targetDir)
	if err != nil {
		t.Fatalf("unexpected recovery error: %v", err)
	}
	if len(rec) != 1 {
		t.Errorf("expected 1 recovered file, got %d", len(rec))
	}

	recoveredBin := filepath.Join(targetDir, "gitmap")
	if _, err := os.Stat(recoveredBin); os.IsNotExist(err) {
		t.Errorf("expected recovered file to exist at %s", recoveredBin)
	}

	var cleanRes CleanResult
	opts := CleanOptions{IsDryRun: false, IsForce: true}
	err = cleanSingleDir(info, opts, targetDir, tempRoot, &cleanRes)
	if err != nil {
		t.Fatalf("unexpected clean error: %v", err)
	}

	if _, err := os.Stat(corruptedDir); !os.IsNotExist(err) {
		t.Errorf("expected corrupted dir to be deleted")
	}
}
