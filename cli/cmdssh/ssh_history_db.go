package cmdssh

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

func copyLegacyFile(src, dst string) {
	s, errOpen := os.Open(src)
	if errOpen != nil {
		return
	}
	defer s.Close()

	d, errCreate := os.Create(dst)
	if errCreate != nil {
		return
	}
	defer d.Close()

	_, _ = io.Copy(d, s)
}

func migrateLegacySSHHistoryFile(targetPath string) {
	legacyPath := store.ResolveSplitDbPath("history", "task", "")
	fi, errStat := os.Stat(legacyPath)
	if errStat == nil && !fi.IsDir() {
		copyLegacyFile(legacyPath, targetPath)
	}
}

func resolveSSHHistoryDBPath() string {
	targetPath := store.ResolveSectionTasksDbPath("ssh", "")
	fi, errStat := os.Stat(targetPath)
	if errStat != nil || fi.IsDir() {
		_ = os.MkdirAll(filepath.Dir(targetPath), 0755)
		migrateLegacySSHHistoryFile(targetPath)
	}
	return targetPath
}

func openSSHHistoryDB() (*sql.DB, error) {
	_ = resolveSSHHistoryDBPath()
	secDB, err := store.OpenSectionTasksDB("ssh", "")
	if err != nil {
		return nil, apperror.WrapSimple(err, "openSSHHistoryDB.OpenSectionTasksDB")
	}

	conn := secDB.Conn()
	migrateLegacySSHHistoryTables(conn)

	return conn, nil
}

// SnapshotSSHNodesBeforeRemoval stores nodes prior to deletion for undo support.
func SnapshotSSHNodesBeforeRemoval(ctx context.Context, action, target string, nodes []store.SSHHost) (string, error) {
	if len(nodes) == 0 {
		return "", nil
	}

	db, err := openSSHHistoryDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	return insertHistorySnapshot(ctx, db, action, target, nodes)
}

func insertHistorySnapshot(ctx context.Context, db *sql.DB, action, target string, nodes []store.SSHHost) (string, error) {
	payload, errMarshal := json.Marshal(nodes)
	if errMarshal != nil {
		return "", apperror.WrapSimple(errMarshal, "insertHistorySnapshot.Marshal")
	}

	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	createdAt := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO TaskHistory (TaskId, Section, Action, Target, ForwardPayload, InversePayload, CreatedAt) VALUES (?, 'ssh', ?, ?, ?, ?, ?)`
	res := store.ExecWrapper(db, query, taskID, action, target, target, string(payload), createdAt)
	if res.IsFailure {
		return "", apperror.WrapSimple(res.Error, "insertHistorySnapshot.Exec")
	}

	recordSSHTaskToRootTasksDB(taskID, action, target, target, string(payload))
	return taskID, nil
}

func recordSSHTaskToRootTasksDB(taskID, action, target, fwd, inv string) {
	rdb, errOpen := store.OpenTasksRootSplitDB()
	if errOpen != nil {
		return
	}
	defer rdb.Close()

	sqlQuery := `INSERT OR IGNORE INTO TaskHistory 
		(TaskId, Section, Action, Target, ForwardPayload, InversePayload, Status) 
		VALUES (?, 'ssh', ?, ?, ?, ?, 'completed')`
	res := store.ExecWrapper(rdb.Conn(), sqlQuery, taskID, action, target, fwd, inv)
	_, _ = res.Destruct() // lint-allow: ignore-db-error reason="optional cross-split sync"
}
