package cmdpipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePipelineDirForRepo_ForwardSlashes(t *testing.T) {
	dir := resolvePipelineDirForRepo("alimtvnetwork/gitmap-v28")
	if strings.Contains(dir, "\\") {
		t.Fatalf("expected all forward slashes in repo pipeline dir, got: %s", dir)
	}
	if !strings.Contains(dir, "alimtvnetwork-gitmap-v28") {
		t.Fatalf("expected repo folder name in path, got: %s", dir)
	}
}

func TestResolvePipelineErrorReportPathForRepo_ForwardSlashes(t *testing.T) {
	reportPath := resolvePipelineErrorReportPathForRepo("alimtvnetwork/gitmap-v28")
	if strings.Contains(reportPath, "\\") {
		t.Fatalf("expected all forward slashes in error report path, got: %s", reportPath)
	}
	if !strings.HasSuffix(reportPath, "alimtvnetwork-gitmap-v28/pipeline_errors.log") {
		t.Fatalf("expected path to end with alimtvnetwork-gitmap-v28/pipeline_errors.log, got: %s", reportPath)
	}
}

func TestGetCachedPipelineLogPathForRepo_ForwardSlashes(t *testing.T) {
	logPath := getCachedPipelineLogPathForRepo("alimtvnetwork/gitmap-v28", 999888)
	if strings.Contains(logPath, "\\") {
		t.Fatalf("expected all forward slashes in cached log path, got: %s", logPath)
	}
	if !strings.HasSuffix(logPath, "alimtvnetwork-gitmap-v28/999888.log") {
		t.Fatalf("expected path to end with alimtvnetwork-gitmap-v28/999888.log, got: %s", logPath)
	}
}

func TestExtractTargetRepo(t *testing.T) {
	defaultRepo := extractTargetRepo([]string{"-y"})
	if len(defaultRepo) == 0 {
		t.Fatalf("expected non-empty default repo slug")
	}

	customRepo := extractTargetRepo([]string{"my-org/my-repo", "-y"})
	if customRepo != "my-org/my-repo" {
		t.Fatalf("expected my-org/my-repo, got: %s", customRepo)
	}
}

func TestPurgeRepoPipelineFolder(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "123.log")
	jsonFile := filepath.Join(tmpDir, "123.json")
	errLog := filepath.Join(tmpDir, "pipeline_errors.log")
	_ = os.WriteFile(logFile, []byte("log data"), 0644)
	_ = os.WriteFile(jsonFile, []byte("{}"), 0644)
	_ = os.WriteFile(errLog, []byte("error report"), 0644)

	purged, reclaimed := purgeRepoPipelineFolder(tmpDir)
	if purged != 3 {
		t.Fatalf("expected 3 purged files, got %d", purged)
	}
	if reclaimed <= 0 {
		t.Fatalf("expected positive reclaimed bytes, got %d", reclaimed)
	}
}

func TestPurgeRepoPipelineFolder_IncludesSqlDb(t *testing.T) {
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, "sql.db")
	walFile := filepath.Join(tmpDir, "sql.db-wal")
	shmFile := filepath.Join(tmpDir, "sql.db-shm")
	_ = os.WriteFile(dbFile, make([]byte, 5000), 0644)
	_ = os.WriteFile(walFile, make([]byte, 2000), 0644)
	_ = os.WriteFile(shmFile, make([]byte, 1000), 0644)

	purged, reclaimed := purgeRepoPipelineFolder(tmpDir)
	if purged != 3 {
		t.Fatalf("expected 3 purged db files, got %d", purged)
	}
	if reclaimed < 8000 {
		t.Fatalf("expected at least 8000 reclaimed bytes, got %d", reclaimed)
	}
}

func TestFormatRelativeDbPath_AppData(t *testing.T) {
	appDataDb := "C:/Users/Administrator/AppData/Local/gitmap-cli/data/pipeline/alimtvnetwork-antigravity-manager/sql.db"
	formatted := FormatRelativeDbPath(appDataDb)
	if strings.Contains(formatted, ".gitmap") {
		t.Fatalf("expected real path without fake .gitmap prefix, got: %s", formatted)
	}
	if !strings.HasSuffix(formatted, "alimtvnetwork-antigravity-manager/sql.db") {
		t.Fatalf("expected path to end with slug/sql.db, got: %s", formatted)
	}
}

func TestCollectRepoPipelineDirs(t *testing.T) {
	dirs := collectRepoPipelineDirs("alimtvnetwork/gitmap-v28")
	if len(dirs) == 0 {
		t.Fatalf("expected at least one pipeline directory to be found")
	}
}
