package cmdstatus

import (
	"testing"

	statusenums "github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus/enums"
)

func TestRenderStatusHelp(t *testing.T) {
	// Verify that RenderStatusHelp runs without panic
	RenderStatusHelp()
}

func TestBuildStatusColumns(t *testing.T) {
	cols := buildStatusColumns()
	expected := statusenums.AllStatusTableColumns()
	if len(cols) != len(expected) {
		t.Fatalf("expected %d columns, got %d", len(expected), len(cols))
	}

	for i, col := range cols {
		if col.Title != expected[i].String() {
			t.Errorf("column %d title mismatch: expected %q, got %q", i, expected[i].String(), col.Title)
		}
	}
}

func TestRenderStatusTableNil(t *testing.T) {
	RenderStatusTable(nil)
}

func TestRenderStatusTableWithRows(t *testing.T) {
	ctx := newStatusTableContext()
	ctx.addRow(statusRow{
		Missing:   false,
		RepoName:  "test-repo",
		Branch:    "main",
		StateIcon: "✔",
		SyncText:  "-",
		StashText: "-",
		FilesText: "-",
	})
	ctx.addRow(statusRow{
		Missing:  true,
		RepoName: "missing-repo",
	})
	RenderStatusTable(ctx)
}
