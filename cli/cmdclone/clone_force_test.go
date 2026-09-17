package cmdclone

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCloneFlags_ForceAndRecloneParsing(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "clean_flag",
			args:     []string{"--clean", "my-repo"},
			expected: true,
		},
		{
			name:     "force_long_flag",
			args:     []string{"--force", "my-repo"},
			expected: true,
		},
		{
			name:     "force_short_flag",
			args:     []string{"-f", "my-repo"},
			expected: true,
		},
		{
			name:     "reclone_flag",
			args:     []string{"--reclone", "my-repo"},
			expected: true,
		},
		{
			name:     "no_force_flag",
			args:     []string{"my-repo"},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cf := parseCloneFlags(tc.args)
			if cf.Clean != tc.expected {
				t.Fatalf("expected Clean=%v, got %v for args %v", tc.expected, cf.Clean, tc.args)
			}
		})
	}
}

func TestCleanDirectCloneTarget_Hermetic(t *testing.T) {
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "repo-to-clean")

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("failed to create target dir: %v", err)
	}

	dummyFile := filepath.Join(targetDir, "README.md")
	if err := os.WriteFile(dummyFile, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to write dummy file: %v", err)
	}

	// No-op when isClean = false
	cleanDirectCloneTarget(targetDir, false)
	if _, err := os.Stat(targetDir); err != nil {
		t.Fatalf("expected targetDir to remain when isClean=false")
	}

	// Deleted when isClean = true
	cleanDirectCloneTarget(targetDir, true)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Fatalf("expected targetDir to be removed when isClean=true")
	}
}
