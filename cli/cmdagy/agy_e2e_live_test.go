//go:build e2e

package cmdagy

import (
	"testing"
)

func TestE2E_SequencePersistence(t *testing.T) {
	skipIfInCI(t)
	testID := "test-e2e-project-id"
	testName := "test-e2e-project"
	testPath := "d:\\work\\test-e2e"

	seq1, err1 := GetOrAssignProjectSequence(testID, testName, testPath)
	if err1 != nil {
		t.Fatalf("failed to assign sequence: %v", err1)
	}
	if seq1 <= 0 {
		t.Fatalf("expected positive sequence number, got %d", seq1)
	}

	seq2, err2 := GetOrAssignProjectSequence(testID, testName, testPath)
	if err2 != nil {
		t.Fatalf("failed to retrieve sequence: %v", err2)
	}
	if seq1 != seq2 {
		t.Fatalf("sequence mismatch: first=%d, second=%d", seq1, seq2)
	}

	rec, getErr := GetProjectBySequence(seq1)
	if getErr != nil || rec == nil {
		t.Fatalf("failed to get project by sequence: %v", getErr)
	}
	if rec.ProjectID != testID {
		t.Fatalf("expected project ID %s, got %s", testID, rec.ProjectID)
	}
}

func TestE2E_ConversationSummariesDiscovery(t *testing.T) {
	skipIfInCI(t)
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		t.Logf("conversation_summaries.db not found on this machine: %v", err)
		return
	}
	t.Logf("conversation_summaries.db found at: %s", dbPath)
}

func TestE2E_DiagnosticsString(t *testing.T) {
	skipIfInCI(t)
	diag := DiagnoseAntigravityIDEAndCLI()
	if len(diag) == 0 {
		t.Fatalf("expected non-empty diagnostics output")
	}
	t.Logf("Diagnostics:\n%s", diag)
}
