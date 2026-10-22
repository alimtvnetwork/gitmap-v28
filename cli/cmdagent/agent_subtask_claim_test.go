package cmdagent

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

func TestBuildClaimResultMap_Nil(t *testing.T) {
	out := buildClaimResultMap(nil, "worker-01", "some/path.db")
	if out["status"] != "empty" {
		t.Fatalf("expected empty status for nil subtask, got %+v", out)
	}
}

func TestBuildClaimResultMap_WithFiles(t *testing.T) {
	sub := &types.Subtask{
		SubtaskId:      "sub-12345",
		ParentTaskId:   "task-abcde",
		TaskCode:       "Subtask-01",
		Title:          "Refactor agent claim",
		OwnedFilesJson: `["read:02-spec/spec.md", "edit:cli/main.go", "created:cli/new.go"]`,
		Status:         "PENDING",
	}

	out := buildClaimResultMap(sub, "worker-02", ".ai-memory/temp-agents/task-abcde/agent-task.db")

	if out["status"] != "claimed" {
		t.Fatalf("expected status claimed, got %+v", out["status"])
	}
	if out["assignedAgentRole"] != "worker-02" {
		t.Fatalf("expected role worker-02, got %+v", out["assignedAgentRole"])
	}
	if out["databasePath"] == nil || out["databasePath"] == "" {
		t.Fatalf("expected non-empty databasePath, got %+v", out["databasePath"])
	}

	files, ok := out["files"].(fsutil.ClassifiedFiles)
	if !ok {
		t.Fatalf("expected files to be fsutil.ClassifiedFiles, got %T", out["files"])
	}

	if len(files.Read) != 1 || files.Read[0] != "02-spec/spec.md" {
		t.Fatalf("unexpected read files: %+v", files.Read)
	}
	if len(files.Edit) != 1 || files.Edit[0] != "cli/main.go" {
		t.Fatalf("unexpected edit files: %+v", files.Edit)
	}
	if len(files.Created) != 1 || files.Created[0] != "cli/new.go" {
		t.Fatalf("unexpected created files: %+v", files.Created)
	}
	if files.Total != 3 {
		t.Fatalf("expected total 3 files, got %d", files.Total)
	}
}
