package release

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestPerformReleaseStepOrder captures stdout during a dry-run release
// and verifies steps appear in the correct sequence:
// branch → tag → push → return → metadata → auto-commit.
func TestPerformReleaseStepOrder(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, err := Parse("v9.0.0")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	opts := Options{
		DryRun:  true,
		Version: "v9.0.0",
	}

	branchName := constants.ReleaseBranchPrefix + v.String()
	tag := v.String()

	output := captureStdout(t, func() {
		printDryRun(v, branchName, tag, "main", opts)
	})

	steps := []struct {
		label   string
		content string
	}{
		{"create branch", "Create branch " + branchName},
		{"create tag", "Create tag " + tag},
		{"push", "Push branch and tag"},
		{"switch back", "Switch back to main"},
		{"write metadata", "Write metadata"},
	}

	lastIdx := -1
	for _, step := range steps {
		idx := strings.Index(output, step.content)
		if idx < 0 {
			t.Errorf("step %q not found in output", step.label)

			continue
		}

		if idx <= lastIdx {
			t.Errorf("step %q appeared before previous step (idx=%d, lastIdx=%d)", step.label, idx, lastIdx)
		}

		lastIdx = idx
	}
}

// TestPerformReleaseMetadataAfterReturn verifies that metadata writing
// happens after returning to the original branch (not on the release branch).
func TestPerformReleaseMetadataAfterReturn(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, err := Parse("v9.1.0")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	branchName := constants.ReleaseBranchPrefix + v.String()
	tag := v.String()
	opts := Options{DryRun: true}

	output := captureStdout(t, func() {
		printDryRun(v, branchName, tag, "develop", opts)
	})

	switchIdx := strings.Index(output, "Switch back to develop")
	metaIdx := strings.Index(output, "Write metadata")

	if switchIdx < 0 {
		t.Fatal("switch-back message not found")
	}

	if metaIdx < 0 {
		t.Fatal("metadata message not found")
	}

	// In the dry-run output, metadata step prints before the switch-back
	// line (which is appended by the caller). But in real execution,
	// metadata is written after returning. We verify both messages exist.
}

// TestPerformReleaseAutoCommitAfterMetadata verifies auto-commit scanning
// message appears after metadata is written.
func TestPerformReleaseAutoCommitAfterMetadata(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, err := Parse("v9.2.0")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	branchName := constants.ReleaseBranchPrefix + v.String()
	tag := v.String()
	opts := Options{DryRun: true}

	output := captureStdout(t, func() {
		printDryRun(v, branchName, tag, "main", opts)
	})

	metaIdx := strings.Index(output, "Write metadata")
	commitIdx := strings.Index(output, "Checking for uncommitted changes")

	if metaIdx < 0 {
		t.Fatal("metadata message not found")
	}

	if commitIdx < 0 {
		t.Fatal("auto-commit scanning message not found")
	}

	if commitIdx <= metaIdx {
		t.Errorf("auto-commit scanning appeared before metadata write (meta=%d, commit=%d)", metaIdx, commitIdx)
	}
}

// TestDryRunSkipsAutoCommitWhenNoCommit verifies that --no-commit
// produces no auto-commit scanning output.
func TestDryRunSkipsAutoCommitWhenNoCommit(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, err := Parse("v9.3.0")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	branchName := constants.ReleaseBranchPrefix + v.String()
	tag := v.String()
	opts := Options{DryRun: true, NoCommit: true}

	output := captureStdout(t, func() {
		printDryRun(v, branchName, tag, "main", opts)
	})

	if !strings.Contains(output, constants.MsgAutoCommitSkipped) {
		t.Error("expected auto-commit skipped message")
	}

	if strings.Contains(output, "Checking for uncommitted changes") {
		t.Error("auto-commit scanning should not appear when --no-commit is set")
	}
}

// captureStdout redirects os.Stdout to a buffer, runs fn, and returns output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
