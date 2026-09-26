//go:build tempe2e

package e2e_test

import (
	"os"
	"os/exec"
	"testing"
)

// TestCommitPullLocalWorkflow_TempE2E verifies the complete multi-repo replay workflow
// against local test instances.
//
// Strictly isolated:
//   - Build tag: tempe2e (excluded from go test ./... by default)
//   - Env guard: RUN_TEMP_E2E=1 required
func TestCommitPullLocalWorkflow_TempE2E(t *testing.T) {
	if os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "../../../scripts/run-e2e-commit-pull.ps1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Temporary E2E commit-pull test failed: %v", err)
	}
}
