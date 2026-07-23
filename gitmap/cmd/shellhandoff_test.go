package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestWriteShellHandoffNoEnvIsNoop verifies that without the wrapper env
// var, no file is written and no error occurs.
func TestWriteShellHandoffNoEnvIsNoop(t *testing.T) {
	t.Setenv(constants.EnvGitmapHandoffFile, "")
	WriteShellHandoff("/tmp/somewhere") // must not panic
}

// TestWriteShellHandoffWritesPath verifies that with the env var set,
// the target path is written verbatim to the sentinel file.
func TestWriteShellHandoffWritesPath(t *testing.T) {
	dir := t.TempDir()
	sentinel := filepath.Join(dir, "handoff.txt")
	t.Setenv(constants.EnvGitmapHandoffFile, sentinel)

	target := filepath.Join(dir, "target-folder")
	WriteShellHandoff(target)

	got, err := os.ReadFile(sentinel)
	if err != nil {
		t.Fatalf("expected sentinel file to exist: %v", err)
	}
	if string(got) != target {
		t.Fatalf("expected %q in sentinel; got %q", target, string(got))
	}
}

// TestWriteShellHandoffEmptyPathIsNoop ensures empty path does not
// truncate the sentinel file (preserves whatever an earlier call wrote).
func TestWriteShellHandoffEmptyPathIsNoop(t *testing.T) {
	dir := t.TempDir()
	sentinel := filepath.Join(dir, "handoff.txt")
	if err := os.WriteFile(sentinel, []byte("preserved"), 0o600); err != nil {
		t.Fatalf("seed sentinel: %v", err)
	}
	t.Setenv(constants.EnvGitmapHandoffFile, sentinel)

	WriteShellHandoff("")

	got, _ := os.ReadFile(sentinel)
	if string(got) != "preserved" {
		t.Fatalf("expected sentinel preserved; got %q", string(got))
	}
}
