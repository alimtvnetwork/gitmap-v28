package cmdsequence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
)

func queryAndPrintSequenceGet(ctx context.Context, repoDB *sql.DB, targetDir string) error {
	query := "SELECT Filename, SequenceNumber, BaseName FROM FileSequence WHERE Directory = ? ORDER BY SequenceNumber ASC"
	rows, err := repoDB.QueryContext(ctx, query, targetDir)
	if err != nil {
		return err
	}

	defer rows.Close()

	files, err := scanSequenceRows(rows, targetDir)
	if err != nil {
		return err
	}

	payload := buildSequencePayload(targetDir, files, len(files))
	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(data))

	return nil
}

func scanSequenceRows(rows *sql.Rows, targetDir string) ([]SequenceItem, error) {
	var files []SequenceItem
	for rows.Next() {
		var fn, bn string
		var seq int
		if err := rows.Scan(&fn, &seq, &bn); err != nil {
			return nil, apperror.WrapSimple(err, "scan sequence row")
		}

		files = append(files, SequenceItem{
			Sequence: seq,
			Filename: fn,
			BaseName: bn,
			Path:     filepath.Join(targetDir, fn),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "iterate sequence rows")
	}

	return files, nil
}

func handleSequenceHistory(args []string) error {
	ctx := context.Background()
	mainDB, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return err
	}

	defer mainDB.Close()
	defer repoDB.Close()

	targetDir := resolveSequenceDir(args)
	targetDir = strings.ReplaceAll(targetDir, "\\", "/")

	return queryAndPrintSequenceHistory(ctx, repoDB, targetDir)
}

func queryAndPrintSequenceHistory(ctx context.Context, repoDB *sql.DB, targetDir string) error {
	query := "SELECT SequenceHistoryId, Directory, OperationsJson, CreatedAt FROM SequenceHistory WHERE Directory = ? ORDER BY CreatedAt DESC LIMIT 20"
	rows, err := repoDB.QueryContext(ctx, query, targetDir)
	if err != nil {
		return err
	}

	defer rows.Close()

	return printSequenceHistoryRows(rows, targetDir)
}

func printSequenceHistoryRows(rows *sql.Rows, targetDir string) error {
	fmt.Printf("\nSequence History for %s:\n", targetDir)
	for rows.Next() {
		var id int64
		var dir, opsJson string
		var created int64
		if err := rows.Scan(&id, &dir, &opsJson, &created); err != nil {
			return apperror.WrapSimple(err, "scan sequence history")
		}

		t := time.Unix(created, 0).UTC().Format(time.RFC3339)
		fmt.Printf("  [%s] ID #%d:\n    %s\n", t, id, opsJson)
	}

	if err := rows.Err(); err != nil {
		return apperror.WrapSimple(err, "iterate sequence history")
	}

	fmt.Println()

	return nil
}

func saveSequenceToRepoDB(payload *SequencePayload) error {
	ctx := context.Background()
	mainDB, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "get repo db for sequence")
	}

	defer mainDB.Close()
	defer repoDB.Close()

	wrapper, appErr := dbengine.WrapDB(repoDB, dbengine.DBSQLite)
	if appErr != nil {
		return appErr
	}

	if txErr := wrapper.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return executeSaveSequenceTx(ctx, tx, payload)
	}); txErr != nil {
		return txErr
	}

	return nil
}

func executeSaveSequenceTx(ctx context.Context, tx *dbengine.TxWrapper, payload *SequencePayload) *apperror.AppError {
	now := time.Now().Unix()
	if execRes := tx.Exec(ctx, "DELETE FROM FileSequence WHERE Directory = ?", payload.Directory); execRes.IsFailure() {
		return apperror.WrapSimple(execRes.AppError(), "delete existing file sequence")
	}

	return insertSequenceFilesTx(ctx, tx, payload, now)
}

func insertSequenceFilesTx(ctx context.Context, tx *dbengine.TxWrapper, payload *SequencePayload, now int64) *apperror.AppError {
	stmtRes := tx.Prepare(ctx, "INSERT OR REPLACE INTO FileSequence (Directory, Filename, SequenceNumber, BaseName, UpdatedAt) VALUES (?, ?, ?, ?, ?)")
	if stmtRes.IsFailure() {
		return apperror.WrapSimple(stmtRes.AppError(), "prepare insert file sequence")
	}

	stmt := stmtRes.Value
	defer stmt.Close()

	for _, f := range payload.Files {
		if _, execErr := stmt.ExecContext(ctx, payload.Directory, f.Filename, f.Sequence, f.BaseName, now); execErr != nil {
			return apperror.WrapSimple(execErr, "exec insert file sequence item")
		}
	}

	return nil
}

func recordSequenceHistoryInDB(report SequenceFixReport) error {
	ctx := context.Background()
	mainDB, repoDB, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "get repo db for history")
	}

	defer mainDB.Close()
	defer repoDB.Close()

	return insertSequenceHistoryTx(ctx, repoDB, report)
}

func insertSequenceHistoryTx(ctx context.Context, repoDB *sql.DB, report SequenceFixReport) error {
	opsJson, err := json.Marshal(report.Operations)
	if err != nil {
		return apperror.WrapSimple(err, "marshal sequence operations")
	}

	now := time.Now().Unix()
	query := "INSERT INTO SequenceHistory (Directory, OperationsJson, CreatedAt) VALUES (?, ?, ?)"
	_, err = repoDB.ExecContext(ctx, query, report.Directory, string(opsJson), now)
	if err != nil {
		return apperror.WrapSimple(err, "insert sequence history")
	}

	return nil
}
