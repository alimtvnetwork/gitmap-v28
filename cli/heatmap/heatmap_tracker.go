package heatmap

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TrackFileTouch records a file access operation into the root heatmap database.
func TrackFileTouch(filePath string, touchType store.FileTouchType) *apperror.AppError {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimSpace(filePath)))
	isValidPath := len(cleanPath) > 0 && cleanPath != "."
	if !isValidPath {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.cwd")
	}

	return recordFileTouchAtRoot(cwd, cleanPath, touchType)
}

func recordFileTouchAtRoot(cwd, cleanPath string, touchType store.FileTouchType) *apperror.AppError {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.openDB")
	}

	defer hdb.Close()

	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
		ActionType: string(touchType),
		FilePath:   cleanPath,
	})
	if logErr != nil {
		return apperror.WrapSimple(logErr, "heatmap.logAction")
	}

	return nil
}

// TrackBatchFileTouches records multiple file touch records in a single database transaction.
func TrackBatchFileTouches(cwd string, touches []store.FileTouchCount) *apperror.AppError {
	if len(touches) == 0 {
		return nil
	}

	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.openDB")
	}

	defer hdb.Close()

	return insertBatchTouchesTx(hdb.Conn(), touches)
}

func insertBatchTouchesTx(conn *sql.DB, touches []store.FileTouchCount) *apperror.AppError {
	tx, err := conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.batchTx")
	}

	defer tx.Rollback()

	if execErr := executeTouchInserts(tx, touches); execErr != nil {
		return execErr
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return apperror.WrapSimple(commitErr, "heatmap.commit")
	}

	return nil
}

func executeTouchInserts(tx *sql.Tx, touches []store.FileTouchCount) *apperror.AppError {
	stmt, err := tx.Prepare(`INSERT INTO HeatmapActionLog (TaskSlug, ActionType, FilePath, ExecutionDepth, ExecutionOrder) VALUES ('batch-touch', 'BATCH_TOUCH', ?, 0, 1)`)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.prepareBatch")
	}

	defer stmt.Close()

	for _, t := range touches {
		if _, err := stmt.Exec(t.FilePath); err != nil {
			return apperror.WrapSimple(err, "heatmap.execTouch")
		}
	}

	return nil
}
