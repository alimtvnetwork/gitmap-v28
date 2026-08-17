package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }
func ptrTime(t time.Time) *time.Time { return &t }

func setupTestDB(t *testing.T) *sql.DB {
	// Enable foreign keys in sqlite
	conn, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	ctx := context.Background()
	err = ApplyMigrations(ctx, conn)
	if err != nil {
		t.Fatalf("ApplyMigrations failed: %v", err)
	}

	return conn
}

func TestClusterRunAndExecResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	ctx := context.Background()

	// 1. Insert ClusterRun
	now := time.Now().UTC()
	run := ClusterRun{
		RunRef:         "RUN-20260818-001",
		CommandKind:    CommandKindPsCommand,
		RawCommand:     "Get-Process",
		TargetSelector: "@all",
		ExceptClause:   ptrString("node3"),
		StartedAt:      now,
	}

	runId, err := InsertClusterRun(ctx, db, run)
	if err != nil {
		t.Fatalf("InsertClusterRun failed: %v", err)
	}
	if runId == 0 {
		t.Fatalf("Expected non-zero runId")
	}

	// 2. Select ClusterRun
	retrieved, err := SelectClusterRun(ctx, db, run.RunRef)
	if err != nil {
		t.Fatalf("SelectClusterRun failed: %v", err)
	}
	if retrieved.ClusterRunId != runId {
		t.Errorf("Expected runId %d, got %d", runId, retrieved.ClusterRunId)
	}
	if retrieved.CommandKind != run.CommandKind {
		t.Errorf("Expected CommandKind %v, got %v", run.CommandKind, retrieved.CommandKind)
	}

	// 3. Update ClusterRun
	finishedAt := now.Add(5 * time.Second)
	err = UpdateClusterRun(ctx, db, runId, ptrTime(finishedAt), ptrInt(5), ptrInt(4), ptrInt(1), ptrInt(0))
	if err != nil {
		t.Fatalf("UpdateClusterRun failed: %v", err)
	}

	retrieved, err = SelectClusterRun(ctx, db, run.RunRef)
	if err != nil {
		t.Fatalf("SelectClusterRun after update failed: %v", err)
	}
	if retrieved.TotalNodes == nil || *retrieved.TotalNodes != 5 {
		t.Errorf("Expected TotalNodes to be 5, got %v", retrieved.TotalNodes)
	}

	// 4. Insert ClusterNode for FK constraints
	_, err = db.ExecContext(ctx, `
		INSERT INTO ClusterNode (NodeId, Alias, DisplayId, IPAddress, NodeRole, OS, Status)
		VALUES ('node1', 'node-1', 1, '10.0.0.1', 'worker', 'windows', 'online')
	`)
	if err != nil {
		t.Fatalf("Failed to insert mock ClusterNode: %v", err)
	}

	// 5. Insert 5 ClusterExecResults
	for i := 1; i <= 5; i++ {
		res := ClusterExecResult{
			ClusterRunId: runId,
			NodeId:       "node1",
			SubCommand:   "Get-Process",
			CommandText:  ptrString("Get-Process"),
			ResultStatus: ResultStatusSucceeded,
			ExitCode:     ptrInt(0),
			Stdout:       ptrString("processes..."),
			StartedAt:    ptrTime(now),
			FinishedAt:   ptrTime(now.Add(1 * time.Second)),
			DurationMs:   ptrInt(1000),
		}
		_, err := InsertClusterExecResult(ctx, db, res)
		if err != nil {
			t.Fatalf("InsertClusterExecResult %d failed: %v", i, err)
		}
	}

	// 6. Query by RunId and verify counts
	results, err := SelectClusterExecResultsByRunId(ctx, db, runId)
	if err != nil {
		t.Fatalf("SelectClusterExecResultsByRunId failed: %v", err)
	}
	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	// 7. Verify FK cascade on run delete
	_, err = db.ExecContext(ctx, "DELETE FROM ClusterRun WHERE ClusterRunId = ?", runId)
	if err != nil {
		t.Fatalf("Failed to delete ClusterRun: %v", err)
	}

	// Check if results are deleted
	resultsAfterDelete, err := SelectClusterExecResultsByRunId(ctx, db, runId)
	if err != nil {
		t.Fatalf("SelectClusterExecResultsByRunId after delete failed: %v", err)
	}
	if len(resultsAfterDelete) != 0 {
		t.Errorf("Expected 0 results after run deletion due to cascade, got %d", len(resultsAfterDelete))
	}
}
