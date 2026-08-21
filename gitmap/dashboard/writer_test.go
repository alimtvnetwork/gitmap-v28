package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func sampleDashboardData(recent bool) model.DashboardData {
	return model.DashboardData{
		Meta: model.DashboardMeta{
			RepoName:    "gitmap-test",
			GeneratedAt: "2026-08-21T00:00:00Z",
			Recent:      recent,
		},
		Commits: []model.CommitInfo{},
	}
}

func TestWriteJSONWithRecent(t *testing.T) {
	tmpDir := t.TempDir()
	data := sampleDashboardData(true)
	path, err := WriteJSON(tmpDir, data)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read JSON output: %v", err)
	}
	verifyJSONRecent(t, bytes)
}

func verifyJSONRecent(t *testing.T, bytes []byte) {
	var readData model.DashboardData
	if err := json.Unmarshal(bytes, &readData); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if !readData.Meta.Recent {
		t.Errorf("expected readData.Meta.Recent to be true")
	}
}

func TestWriteHTML(t *testing.T) {
	tmpDir := t.TempDir()
	data := sampleDashboardData(true)
	path, err := WriteHTML(tmpDir, data)
	if err != nil {
		t.Fatalf("WriteHTML failed: %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read HTML output: %v", err)
	}
	if !strings.Contains(string(content), `"recent":true`) {
		t.Errorf("expected HTML to contain recent:true in JSON payload")
	}
}

func TestSummaryAndFormatSize(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	s := Summary(testFile)
	if !strings.Contains(s, "11 B") {
		t.Errorf("expected summary to contain size, got %s", s)
	}
}
