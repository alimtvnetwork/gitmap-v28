package cmdinstall

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestResolveToolStatusFromDB(t *testing.T) {
	installed := map[string]string{"node": "20.11.0"}
	status, version := resolveToolStatus("node", installed)
	if status != constants.StatusInstalled {
		t.Fatalf("expected installed glyph for DB hit, got %q", status)
	}

	if version != "20.11.0" {
		t.Fatalf("expected DB version 20.11.0, got %q", version)
	}
}

func TestResolveToolStatusGitmap(t *testing.T) {
	status, version := resolveToolStatus(constants.ToolGitmap, map[string]string{})
	if status != constants.StatusInstalled {
		t.Fatalf("expected installed status for gitmap, got %q", status)
	}

	if version != constants.Version {
		t.Fatalf("expected constants.Version %q, got %q", constants.Version, version)
	}
}

func TestParseVersionFromOutput(t *testing.T) {
	cases := map[string]string{
		"go version go1.24.1 windows/amd64": "1.24.1",
		"v22.14.0\n":                        "v22.14.0",
		"git version 2.47.0.windows.1":      "2.47.0",
	}

	for input, expected := range cases {
		got := parseVersionFromOutput(input)
		if got != expected {
			t.Fatalf("input %q: expected %q, got %q", input, expected, got)
		}
	}
}

func TestResolveToolStatusUnknownTool(t *testing.T) {
	status, version := resolveToolStatus("definitely-not-a-real-binary-xyz", map[string]string{})
	if status != constants.StatusNotInstalled {
		t.Fatalf("expected not-installed glyph for missing tool, got %q", status)
	}

	if version != "—" {
		t.Fatalf("expected dash version for missing tool, got %q", version)
	}
}

func TestPickDisplayVersion(t *testing.T) {
	if pickDisplayVersion(store.InstalledTool{VersionString: "1.2.3"}) != "1.2.3" {
		t.Fatal("expected 1.2.3")
	}

	if pickDisplayVersion(store.InstalledTool{VersionString: ""}) != "—" {
		t.Fatal("expected dash for empty")
	}

	if pickDisplayVersion(store.InstalledTool{VersionString: "0.0.0"}) != "—" {
		t.Fatal("expected dash for zeros")
	}
}

func TestSortedCategoryNamesCoreFirst(t *testing.T) {
	got := sortedCategoryNames()
	if len(got) == 0 {
		t.Skip("no categories defined; nothing to assert")
	}

	if got[0] != constants.ToolCategoryCore {
		t.Fatalf("ToolCategoryCore must sort first; got order %v", got)
	}
}

func TestResolvePackageManagerOverride(t *testing.T) {
	got := resolvePackageManager("brew", "")
	if got != "brew" {
		t.Fatalf("override must win; got %q", got)
	}
}

func TestIsInstallLogsCommand(t *testing.T) {
	if !isInstallLogsCommand([]string{"logs"}) || !isInstallLogsCommand([]string{"--logs"}) {
		t.Fatal("expected logs and --logs to be recognized")
	}

	if !isInstallLogsCommand([]string{"log"}) || !isInstallLogsCommand([]string{"-logs"}) {
		t.Fatal("expected log and -logs to be recognized")
	}

	if isInstallLogsCommand([]string{"node"}) || isInstallLogsCommand([]string{}) {
		t.Fatal("expected node or empty args to not be recognized as logs")
	}
}

func TestExtractInstallLogsArgs(t *testing.T) {
	got1 := extractInstallLogsArgs([]string{"logs", "--failed"})
	if len(got1) != 1 || got1[0] != "--failed" {
		t.Fatalf("expected [--failed], got %v", got1)
	}

	got2 := extractInstallLogsArgs([]string{"--logs", "--tool", "go"})
	if len(got2) != 2 || got2[0] != "--tool" || got2[1] != "go" {
		t.Fatalf("expected [--tool go], got %v", got2)
	}

	if len(extractInstallLogsArgs([]string{"--logs"})) != 0 {
		t.Fatal("expected empty args when only --logs given")
	}
}

func TestParseInstallLogsFlags(t *testing.T) {
	opts, err := parseInstallLogsFlags([]string{"--failed", "--tool", "rust", "--limit", "10"})
	if err != nil || !opts.Failed || opts.Tool != "rust" || opts.Limit != 10 {
		t.Fatalf("unexpected parsed flags: %+v, err=%v", opts, err)
	}

	opts2, err := parseInstallLogsFlags([]string{"python"})
	if err != nil || opts2.Tool != "python" || opts2.Limit != 50 || opts2.Failed {
		t.Fatalf("unexpected parsed fallback tool: %+v, err=%v", opts2, err)
	}
}

func TestFormatLogHelpers(t *testing.T) {
	if formatLogDuration(0) != "0ms" || formatLogDuration(500) != "500ms" {
		t.Fatal("unexpected formatLogDuration ms values")
	}

	if formatLogDuration(1500) != "1.5s" || formatLogDuration(65000) != "1m5s" {
		t.Fatal("unexpected formatLogDuration seconds/minutes values")
	}

	if formatLogStatus(true) != "success" || formatLogStatus(false) != "failed" {
		t.Fatal("unexpected formatLogStatus")
	}

	if formatLogDisplayVal("") != "-" || formatLogDisplayVal("1.0") != "1.0" {
		t.Fatal("unexpected formatLogDisplayVal")
	}
}

func setupTestInstallDB(t *testing.T) (*store.InstallationSplitDB, func()) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "installation.db")
	db, err := store.OpenInstallationSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test split db: %v", err)
	}

	seedTestLogs(db)

	return db, func() { _ = db.Close() }
}

func seedTestLogs(db *store.InstallationSplitDB) {
	_ = db.RecordExecution("go", "install", "1.22.0", "winget", 1200, true, 0, "", "", "winget install", "", "")
	_ = db.RecordExecution("go", "install", "1.21.0", "choco", 800, false, 1, "", "err", "choco install", "", "")
	_ = db.RecordExecution("rust", "install", "1.75.0", "winget", 2500, true, 0, "", "", "winget install", "", "")
}

func TestExecuteInstallLogs(t *testing.T) {
	db, cleanup := setupTestInstallDB(t)
	defer cleanup()

	var buf bytes.Buffer
	if err := executeInstallLogs(&buf, db, []string{}); err != nil {
		t.Fatalf("executeInstallLogs failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "TOOL") || !strings.Contains(out, "go") || !strings.Contains(out, "rust") {
		t.Fatalf("expected table headers and records, got: %s", out)
	}
}

func TestExecuteInstallLogsFiltering(t *testing.T) {
	db, cleanup := setupTestInstallDB(t)
	defer cleanup()

	var buf bytes.Buffer
	if err := executeInstallLogs(&buf, db, []string{"--failed"}); err != nil {
		t.Fatalf("executeInstallLogs failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "failed") || strings.Contains(out, "rust") {
		t.Fatalf("expected only failed log, got: %s", out)
	}
}
