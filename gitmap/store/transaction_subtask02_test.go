package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func setupSubtask02DB(t *testing.T) *DB {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_subtask02.db")

	db, err := OpenAt(dbPath)
	if err != nil {
		t.Fatalf("OpenAt failed: %v", err)
	}

	if err := db.Migrate(); err != nil {
		db.Close()
		t.Fatalf("Migrate failed: %v", err)
	}

	return db
}

func seedVisibilityRun(t *testing.T, db *DB) int64 {
	t.Helper()
	run := model.MakeAllVisibilityRunRecord{
		CommandKind:      "mapub",
		TargetVisibility: "public",
		Provider:         "github",
		Owner:            "testorg",
		TargetRaw:        "all",
		PatternList:      "*",
		StartedAt:        "2026-09-12T10:00:00Z",
	}

	id, err := db.InsertMakeAllVisibilityRun(run)
	if err != nil {
		t.Fatalf("InsertMakeAllVisibilityRun failed: %v", err)
	}

	return id
}

func TestMakeAllVisibilityPendingAndExclude(t *testing.T) {
	db := setupSubtask02DB(t)
	defer db.Close()

	runID := seedVisibilityRun(t, db)
	rows := []model.MakeAllVisibilityResultRecord{
		{RepoName: "repo-a", MatchedPattern: "*", StartedAt: "2026-09-12T10:00:00Z"},
		{RepoName: "repo-b", MatchedPattern: "*", StartedAt: "2026-09-12T10:00:00Z"},
	}

	ids, err := db.InsertMakeAllVisibilityPendingResults(runID, rows)
	if err != nil || len(ids) != 2 {
		t.Fatalf("InsertMakeAllVisibilityPendingResults failed: %v", err)
	}

	if err := db.MarkMakeAllVisibilityResultsExcluded(ids, "2026-09-12T10:01:00Z"); err != nil {
		t.Fatalf("MarkMakeAllVisibilityResultsExcluded failed: %v", err)
	}
}

func TestOwnerRepoNameIndex_Tx(t *testing.T) {
	db := setupSubtask02DB(t)
	defer db.Close()

	if err := db.EnsureOwnerRepoNameIndex(); err != nil {
		t.Fatalf("EnsureOwnerRepoNameIndex failed: %v", err)
	}

	names := []string{"macro-ahk", "macro-ahk-v1", "macro-ahk-v2"}
	now := time.Now().UTC()
	if err := db.UpsertOwnerRepoNameIndex("github", "testowner", names, now); err != nil {
		t.Fatalf("UpsertOwnerRepoNameIndex failed: %v", err)
	}

	repoName, ver, hasVersion := db.LookupHighestVersion("github", "testowner", "macro-ahk")
	if hasVersion == false || ver != 2 || repoName != "macro-ahk-v2" {
		t.Fatalf("LookupHighestVersion failed: ver=%d, name=%s", ver, repoName)
	}
}

func TestPendingTaskComplete_Tx(t *testing.T) {
	db := setupSubtask02DB(t)
	defer db.Close()

	id, err := db.InsertPendingTask(1, "/path/target", "/path/work", "git", "status")
	if err != nil {
		t.Fatalf("InsertPendingTask failed: %v", err)
	}

	if err := db.CompleteTask(id); err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	completed, err := db.ListCompletedTasks()
	if err != nil || len(completed) != 1 {
		t.Fatalf("expected 1 completed task, got: %v", completed)
	}
}

func TestChromeProfileDelete_Tx(t *testing.T) {
	db := setupSubtask02DB(t)
	defer db.Close()

	id, err := db.UpsertChromeProfile("TestProf", "/src/path", true)
	if err != nil {
		t.Fatalf("UpsertChromeProfile failed: %v", err)
	}

	if err := db.InsertChromeProfileExport(id, "json", "/out/test.json", 100); err != nil {
		t.Fatalf("InsertChromeProfileExport failed: %v", err)
	}

	paths, err := db.DeleteChromeProfile("TestProf")
	if err != nil || len(paths) != 1 {
		t.Fatalf("DeleteChromeProfile failed: %v, paths=%v", err, paths)
	}
}

func seedPruneTxns(t *testing.T, db *DB) {
	t.Helper()
	for i := 0; i < 3; i++ {
		rec := model.TransactionRecord{Kind: "rollback", Argv: "args", Cwd: "/cwd"}
		if _, err := db.InsertTransaction(rec); err != nil {
			t.Fatalf("InsertTransaction failed: %v", err)
		}
	}
}

func TestTransactionPrune_Tx(t *testing.T) {
	db := setupSubtask02DB(t)
	defer db.Close()

	seedPruneTxns(t, db)
	prunedIDs, err := db.PruneOldestTransactions(1)
	if err != nil {
		t.Fatalf("PruneOldestTransactions failed: %v", err)
	}

	if len(prunedIDs) != 2 {
		t.Fatalf("expected 2 pruned txns, got %d", len(prunedIDs))
	}
}

func TestInsertSSHHost_TxWrapper(t *testing.T) {
	db := setupSubtask02DB(t)
	defer db.Close()

	if err := RegisterSSHHostMigration(db.conn, 1, false); err != nil {
		t.Fatalf("RegisterSSHHostMigration failed: %v", err)
	}

	wrap, appErr := dbengine.WrapDb(db.conn, dbengine.DbSQLite)
	if appErr != nil {
		t.Fatalf("WrapDb failed: %v", appErr)
	}

	ctx := context.Background()
	now := time.Now().UTC()
	h1 := SSHHost{ID: "h-1", Alias: "srv1", IP: "10.0.0.1", Username: "u1", CreatedAt: now}
	h2 := SSHHost{ID: "h-2", Alias: "srv2", IP: "10.0.0.2", Username: "u2", CreatedAt: now}

	appErr = wrap.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		if err := InsertSSHHost(ctx, h1, tx); err != nil {
			return apperror.WrapSimple(err, "insert h1")
		}

		if err := InsertSSHHostTx(ctx, h2, tx); err != nil {
			return apperror.WrapSimple(err, "insert h2")
		}

		return nil
	})
	if appErr != nil {
		t.Fatalf("WithTransaction failed: %v", appErr)
	}

	hosts, err := ListHosts(ctx, db.conn)
	if err != nil || len(hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d, err: %v", len(hosts), err)
	}
}
