// Package cmdtask — task_audit_test.go tests audit task logging and badge formatting.
package cmdtask

import (
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestFormatBadges(t *testing.T) {
	macroBadge := formatSectionBadge("macro")
	if !strings.Contains(macroBadge, "MACRO") {
		t.Errorf("expected MACRO badge, got %q", macroBadge)
	}

	sshBadge := formatSectionBadge("ssh")
	if !strings.Contains(sshBadge, "SSH") {
		t.Errorf("expected SSH badge, got %q", sshBadge)
	}

	okBadge := formatStatusBadge("completed")
	if !strings.Contains(okBadge, "COMPLETED") {
		t.Errorf("expected COMPLETED badge, got %q", okBadge)
	}

	failBadge := formatStatusBadge("failed")
	if !strings.Contains(failBadge, "FAILED") {
		t.Errorf("expected FAILED badge, got %q", failBadge)
	}
}

func TestRecordTaskAudit_Roundtrip(t *testing.T) {
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		t.Fatalf("OpenTasksRootSplitDB failed: %v", err)
	}
	defer tasksDB.Close()

	testAction := "test-audit-action-" + time.Now().Format("150405.000")
	RecordTaskAudit("macro", testAction, "test-target", "payload-data", "completed")

	records, listErr := tasksDB.ListTaskHistory("macro", 10, 0)
	if listErr != nil {
		t.Fatalf("ListTaskHistory failed: %v", listErr)
	}

	var found bool
	for _, r := range records {
		if r.Action != testAction {
			continue
		}
		found = true
		if r.Target != "test-target" || r.Section != "macro" || r.Status != "completed" {
			t.Errorf("record mismatch: %+v", r)
		}
		break
	}

	if !found {
		t.Errorf("expected to find audit record for action %q", testAction)
	}
}
