package cmdfixreleasetags

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestDeleteExecutor_FullFourTierSuccess(t *testing.T) {
	tempDir := t.TempDir()
	relDir := filepath.Join(tempDir, constants.GitMapDir, constants.ReleaseDirName)
	_ = os.MkdirAll(relDir, 0o755)
	sidecarPath := filepath.Join(relDir, "v6.521.0.json")
	_ = os.WriteFile(sidecarPath, []byte("{}"), 0o644)

	mockExec := newTestMockExecutor()
	mockExec.responses["gh release delete v6.521.0 -y"] = []byte("deleted")
	mockExec.responses["git push origin :refs/tags/v6.521.0"] = []byte("deleted remote")
	mockExec.responses["git tag -d v6.521.0"] = []byte("Deleted tag 'v6.521.0'")

	record := ReleaseTagAuditRecord{
		Tag:                   "v6.521.0",
		AuditReason:           ReasonDraftRelease,
		HasGitHubRelease:      true,
		HasRemoteTag:          true,
		HasLocalTag:           true,
		IsEligibleForDeletion: true,
	}

	res, err := ExecuteTagDeletion(tempDir, record, false, mockExec)
	if err != nil {
		t.Fatalf("ExecuteTagDeletion returned error: %v", err)
	}

	if res.IsFailed() {
		t.Fatalf("expected successful deletion, got isCompleted=%v isSuccess=%v err=%s", res.IsCompleted, res.IsSuccess, res.ErrorMessage)
	}

	if len(res.Steps) != 4 {
		t.Fatalf("expected 4 steps, got %d", len(res.Steps))
	}

	for _, s := range res.Steps {
		if s.Status != StepSuccess {
			t.Errorf("step %s expected StepSuccess, got %s", s.StepName, s.Status)
		}
	}

	if _, statErr := os.Stat(sidecarPath); !os.IsNotExist(statErr) {
		t.Errorf("expected sidecar file to be purged")
	}
}

func TestDeleteExecutor_OrphanTagSkipsGHRelease(t *testing.T) {
	mockExec := newTestMockExecutor()
	mockExec.responses["git push origin :refs/tags/v6.520.0"] = []byte("deleted remote")
	mockExec.responses["git tag -d v6.520.0"] = []byte("Deleted tag 'v6.520.0'")

	record := ReleaseTagAuditRecord{
		Tag:                   "v6.520.0",
		AuditReason:           ReasonOrphanTag,
		HasGitHubRelease:      false,
		HasRemoteTag:          true,
		HasLocalTag:           true,
		IsEligibleForDeletion: true,
	}

	res, err := ExecuteTagDeletion(".", record, false, mockExec)
	if err != nil {
		t.Fatalf("ExecuteTagDeletion returned error: %v", err)
	}

	if res.Steps[0].Status != StepSkipped {
		t.Errorf("expected Step 1 skipped for orphan tag, got %s", res.Steps[0].Status)
	}
	if res.Steps[1].Status != StepSuccess {
		t.Errorf("expected Step 2 success, got %s", res.Steps[1].Status)
	}
	if res.Steps[2].Status != StepSuccess {
		t.Errorf("expected Step 3 success, got %s", res.Steps[2].Status)
	}
}

func TestDeleteExecutor_ErrorIsolation(t *testing.T) {
	tempDir := t.TempDir()
	mockExec := newTestMockExecutor()
	mockExec.responses["gh release delete v6.524.0 -y"] = []byte("deleted")
	mockExec.errors["git push origin :refs/tags/v6.524.0"] = fmt.Errorf("remote rejected")
	mockExec.responses["git tag -d v6.524.0"] = []byte("Deleted tag 'v6.524.0'")

	record := ReleaseTagAuditRecord{
		Tag:                   "v6.524.0",
		AuditReason:           ReasonCICDFailed,
		HasGitHubRelease:      true,
		HasRemoteTag:          true,
		HasLocalTag:           true,
		IsEligibleForDeletion: true,
	}

	res, err := ExecuteTagDeletion(tempDir, record, false, mockExec)
	if err != nil {
		t.Fatalf("expected error isolation, got err: %v", err)
	}

	if res.IsSuccess {
		t.Errorf("expected isSuccess=false due to remote failure")
	}

	if res.Steps[1].Status != StepFailed {
		t.Errorf("expected Step 2 failed, got %s", res.Steps[1].Status)
	}
	if res.Steps[2].Status != StepSuccess {
		t.Errorf("expected Step 3 to proceed and succeed despite Step 2 failure, got %s", res.Steps[2].Status)
	}
	if res.Steps[3].Status != StepSuccess {
		t.Errorf("expected Step 4 to succeed, got %s", res.Steps[3].Status)
	}
}

func TestDeleteExecutor_DryRun(t *testing.T) {
	tempDir := t.TempDir()
	relDir := filepath.Join(tempDir, constants.GitMapDir, constants.ReleaseDirName)
	_ = os.MkdirAll(relDir, 0o755)
	sidecarPath := filepath.Join(relDir, "v6.522.0.json")
	_ = os.WriteFile(sidecarPath, []byte("{}"), 0o644)

	mockExec := newTestMockExecutor()

	record := ReleaseTagAuditRecord{
		Tag:                   "v6.522.0",
		AuditReason:           ReasonMissingAssets,
		HasGitHubRelease:      true,
		HasRemoteTag:          true,
		HasLocalTag:           true,
		IsEligibleForDeletion: true,
	}

	res, err := ExecuteTagDeletion(tempDir, record, true, mockExec)
	if err != nil {
		t.Fatalf("ExecuteTagDeletion failed in dry-run: %v", err)
	}

	if len(mockExec.calls) != 0 {
		t.Errorf("expected zero subprocess calls in dry-run mode, got %d", len(mockExec.calls))
	}

	if _, statErr := os.Stat(sidecarPath); os.IsNotExist(statErr) {
		t.Errorf("sidecar file should not be removed in dry-run mode")
	}

	if res.IsFailed() {
		t.Errorf("expected isSuccess=true for dry-run preview")
	}
}

func TestExecuteDeletion_PlanBatch(t *testing.T) {
	mockExec := newTestMockExecutor()
	mockExec.responses["git tag -d v6.520.0"] = []byte("Deleted tag")

	plan := &DeletionPlan{
		RepoPath:        ".",
		IsDryRun:        false,
		CommandExecutor: mockExec,
		Records: []ReleaseTagAuditRecord{
			{
				Tag:                   "v6.520.0",
				AuditReason:           ReasonOrphanTag,
				HasLocalTag:           true,
				IsEligibleForDeletion: true,
			},
		},
	}

	result, err := ExecuteDeletion(".", plan)
	if err != nil {
		t.Fatalf("ExecuteDeletion returned error: %v", err)
	}

	if result.Tag != "v6.520.0" {
		t.Errorf("expected result tag v6.520.0, got %s", result.Tag)
	}
}
