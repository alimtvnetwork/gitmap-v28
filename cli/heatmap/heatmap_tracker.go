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
func TrackFileTouch(filePath string, touchType store.FileTouchType) error {
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

func recordFileTouchAtRoot(cwd, cleanPath string, touchType store.FileTouchType) error {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return err
	}

	defer hdb.Close()

	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
		ActionType: string(touchType),
		FilePath:   cleanPath,
	})

	return logErr
}

// TrackBatchFileTouches records multiple file touch records in a single database transaction.
func TrackBatchFileTouches(cwd string, touches []store.FileTouchCount) error {
	if len(touches) == 0 {
		return nil
	}

	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return err
	}

	defer hdb.Close()

	return insertBatchTouchesTx(hdb.Conn(), touches)
}

func insertBatchTouchesTx(conn *sql.DB, touches []store.FileTouchCount) error {
	tx, err := conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.batchTx")
	}

	defer tx.Rollback()

	if execErr := executeTouchInserts(tx, touches); execErr != nil {
		return execErr
	}

	return tx.Commit()
}

func executeTouchInserts(tx *sql.Tx, touches []store.FileTouchCount) error {
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

