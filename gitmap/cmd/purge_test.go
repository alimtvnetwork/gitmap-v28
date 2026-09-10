package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePurgeArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantPat     string
		wantRestore bool
		wantConfirm bool
	}{
		{
			name:        "restore flag",
			args:        []string{"--restore"},
			wantPat:     "",
			wantRestore: true,
			wantConfirm: false,
		},
		{
			name:        "confirm and pattern positional",
			args:        []string{"-y", "*.db"},
			wantPat:     "*.db",
			wantRestore: false,
			wantConfirm: true,
		},
		{
			name:        "confirm long and path flag",
			args:        []string{"--confirm", "--path", "secrets/*.key"},
			wantPat:     "secrets/*.key",
			wantRestore: false,
			wantConfirm: true,
		},
		{
			name:        "pattern first",
			args:        []string{"build/*", "-y"},
			wantPat:     "build/*",
			wantRestore: false,
			wantConfirm: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pat, isRestore, isAutoConfirm := parsePurgeArgs(tc.args)
			if pat != tc.wantPat {
				t.Errorf("got pat %q, want %q", pat, tc.wantPat)
			}
			if isRestore != tc.wantRestore {
				t.Errorf("got isRestore %v, want %v", isRestore, tc.wantRestore)
			}
			if isAutoConfirm != tc.wantConfirm {
				t.Errorf("got isAutoConfirm %v, want %v", isAutoConfirm, tc.wantConfirm)
			}
		})
	}
}

func TestPurgeCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.txt")
	dst := filepath.Join(tmpDir, "dest.txt")

	if err := os.WriteFile(src, []byte("purge content"), 0644); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}

	if err := copyPurgeFile(src, dst); err != nil {
		t.Fatalf("copyPurgeFile failed: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read dest: %v", err)
	}
	if string(data) != "purge content" {
		t.Errorf("got content %q, want %q", string(data), "purge content")
	}

	if err := copyPurgeFile(filepath.Join(tmpDir, "nonexistent"), dst); err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}
}

func TestPurgePatternNormalization(t *testing.T) {
	winPath := "dir\\subdir\\*.secret"
	expected := "dir/subdir/*.secret"
	if norm := filepath.ToSlash(winPath); norm != expected {
		t.Errorf("got %q, want %q", norm, expected)
	}
}
