package heatmap

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestRcHeatmapSingleShotJsonInput(t *testing.T) {
	raw := `{"task":"test-313","action":"READ","file":"cli/main.go","expected":"verify entrypoint"}`
	rec, err := ParseJsonActionPayload(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing json payload: %v", err)
	}

	if rec.TaskSlug != "test-313" {
		t.Errorf("expected TaskSlug test-313, got %s", rec.TaskSlug)
	}

	if rec.ActionType != "READ" {
		t.Errorf("expected ActionType READ, got %s", rec.ActionType)
	}

	if rec.FilePath != "cli/main.go" {
		t.Errorf("expected FilePath cli/main.go, got %s", rec.FilePath)
	}
}

func TestRcHeatmapModelAndToolTelemetry(t *testing.T) {
	os.Setenv("AI_MODEL", "gemini-2.5-pro")
	os.Setenv("ANTIGRAVITY_AGENT_ID", "agent-worker-03")
	defer os.Unsetenv("AI_MODEL")
	defer os.Unsetenv("ANTIGRAVITY_AGENT_ID")

	modelName := store.DetectCurrentModelName()
	toolName := store.DetectCurrentToolName()

	if modelName != "gemini-2.5-pro" {
		t.Errorf("expected gemini-2.5-pro, got %s", modelName)
	}

	if toolName != "Antigravity IDE" {
		t.Errorf("expected Antigravity IDE, got %s", toolName)
	}
}

func TestRcHeatmapDispatch_Canonical(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("GITMAP_DATA_DIR", tmpDir)
	defer os.Unsetenv("GITMAP_DATA_DIR")

	var buf bytes.Buffer
	args := []string{"log", "--task", "test-dispatch", "--action", "ANALYZE", "--file", "cli/main.go", "--json"}
	err := DispatchRcHeatmap(context.Background(), args, &buf)
	if err != nil {
		t.Fatalf("unexpected dispatch error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"status": "success"`) && !strings.Contains(out, `"status":"success"`) {
		t.Errorf("expected json success status in output, got: %s", out)
	}
}
