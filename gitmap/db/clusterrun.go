package db

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type ClusterRun struct {
	ClusterRunId   int64
	RunRef         string
	CommandKind    CommandKindType
	RawCommand     string
	TargetSelector string
	ExceptClause   *string
	StartedAt      time.Time
	FinishedAt     *time.Time
	TotalNodes     *int
	SucceededNodes *int
	FailedNodes    *int
	SkippedNodes   *int
}

const (
	sqlInsertClusterRun = `
		INSERT INTO ClusterRun (
			RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause,
			StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	sqlUpdateClusterRun = `
		UPDATE ClusterRun
		SET FinishedAt = ?, TotalNodes = ?, SucceededNodes = ?, FailedNodes = ?, SkippedNodes = ?
		WHERE ClusterRunId = ?
	`
	sqlSelectClusterRunByRef = `
		SELECT 
			ClusterRunId, RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause,
			StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes
		FROM ClusterRun
		WHERE RunRef = ?
	`
	sqlSelectListClusterRuns = `
		SELECT 
			ClusterRunId, RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause,
			StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes
		FROM ClusterRun
		ORDER BY StartedAt DESC
	`
)

func InsertClusterRun(ctx context.Context, db *sql.DB, run ClusterRun) (int64, *apperror.AppError) {
	res, err := db.ExecContext(ctx, sqlInsertClusterRun,
		run.RunRef, run.CommandKind, run.RawCommand, run.TargetSelector, run.ExceptClause,
		run.StartedAt, run.FinishedAt, run.TotalNodes, run.SucceededNodes, run.FailedNodes, run.SkippedNodes,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, "InsertClusterRun.Exec")
	}
	id, errId := res.LastInsertId()
	if errId != nil {
		return 0, apperror.WrapSimple(errId, "InsertClusterRun.LastInsertId")
	}

	return id, nil
}

func UpdateClusterRun(
	ctx context.Context,
	db *sql.DB,
	id int64,
	finishedAt *time.Time,
	totalNodes, succeededNodes, failedNodes, skippedNodes *int,
) *apperror.AppError {
	_, err := db.ExecContext(ctx, sqlUpdateClusterRun, finishedAt, totalNodes, succeededNodes, failedNodes, skippedNodes, id)
	if err != nil {
		return apperror.WrapSimple(err, "UpdateClusterRun.Exec")
	}

	return nil
}

func SelectClusterRun(ctx context.Context, db *sql.DB, runRef string) (ClusterRun, *apperror.AppError) {
	row := db.QueryRowContext(ctx, sqlSelectClusterRunByRef, runRef)
	var run ClusterRun
	err := row.Scan(
		&run.ClusterRunId, &run.RunRef, &run.CommandKind, &run.RawCommand, &run.TargetSelector,
		&run.ExceptClause, &run.StartedAt, &run.FinishedAt, &run.TotalNodes, &run.SucceededNodes,
		&run.FailedNodes, &run.SkippedNodes,
	)
	if err == sql.ErrNoRows {
		return ClusterRun{}, apperror.NewWithDetails(
			"SelectClusterRun", "E4004", "ClusterRun not found for RunRef: "+runRef, "db",
			apperror.ErrorTypeNotFound, apperror.SeverityError, map[string]any{"runRef": runRef},
		)
	}
	if err != nil {
		return ClusterRun{}, apperror.WrapSimple(err, "SelectClusterRun.Scan")
	}

	return run, nil
}

func ListClusterRuns(ctx context.Context, db *sql.DB, limit int) ([]ClusterRun, *apperror.AppError) {
	query := sqlSelectListClusterRuns
	if limit > 0 {
		query += " LIMIT " + strconv.Itoa(limit)
	}
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "ListClusterRuns.Query")
	}
	defer rows.Close()

	return scanClusterRunRows(rows)
}

func scanClusterRunRows(rows *sql.Rows) ([]ClusterRun, *apperror.AppError) {
	var runs []ClusterRun
	for rows.Next() {
		var run ClusterRun
		err := rows.Scan(
			&run.ClusterRunId, &run.RunRef, &run.CommandKind, &run.RawCommand, &run.TargetSelector,
			&run.ExceptClause, &run.StartedAt, &run.FinishedAt, &run.TotalNodes, &run.SucceededNodes,
			&run.FailedNodes, &run.SkippedNodes,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scanClusterRunRows.Scan")
		}
		runs = append(runs, run)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "scanClusterRunRows.Rows")
	}

	return runs, nil
}
