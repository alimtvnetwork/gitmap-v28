package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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

func InsertClusterRun(ctx context.Context, db *sql.DB, run ClusterRun) (int64, error) {
	query := `
		INSERT INTO ClusterRun (
			RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause,
			StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := db.ExecContext(ctx, query,
		run.RunRef,
		run.CommandKind,
		run.RawCommand,
		run.TargetSelector,
		run.ExceptClause,
		run.StartedAt,
		run.FinishedAt,
		run.TotalNodes,
		run.SucceededNodes,
		run.FailedNodes,
		run.SkippedNodes,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert ClusterRun: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func UpdateClusterRun(ctx context.Context, db *sql.DB, id int64, finishedAt *time.Time, totalNodes, succeededNodes, failedNodes, skippedNodes *int) error {
	query := `
		UPDATE ClusterRun
		SET FinishedAt = ?, TotalNodes = ?, SucceededNodes = ?, FailedNodes = ?, SkippedNodes = ?
		WHERE ClusterRunId = ?
	`
	_, err := db.ExecContext(ctx, query, finishedAt, totalNodes, succeededNodes, failedNodes, skippedNodes, id)
	if err != nil {
		return fmt.Errorf("failed to update ClusterRun %d: %w", id, err)
	}

	return nil
}

func SelectClusterRun(ctx context.Context, db *sql.DB, runRef string) (ClusterRun, error) {
	query := `
		SELECT 
			ClusterRunId, RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause,
			StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes
		FROM ClusterRun
		WHERE RunRef = ?
	`
	row := db.QueryRowContext(ctx, query, runRef)

	var run ClusterRun
	err := row.Scan(
		&run.ClusterRunId,
		&run.RunRef,
		&run.CommandKind,
		&run.RawCommand,
		&run.TargetSelector,
		&run.ExceptClause,
		&run.StartedAt,
		&run.FinishedAt,
		&run.TotalNodes,
		&run.SucceededNodes,
		&run.FailedNodes,
		&run.SkippedNodes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ClusterRun{}, fmt.Errorf("ClusterRun not found for RunRef %s: %w", runRef, err)
		}
		return ClusterRun{}, fmt.Errorf("failed to scan ClusterRun: %w", err)
	}

	return run, nil
}

func ListClusterRuns(ctx context.Context, db *sql.DB, limit int) ([]ClusterRun, error) {
	query := `
		SELECT 
			ClusterRunId, RunRef, CommandKind, RawCommand, TargetSelector, ExceptClause,
			StartedAt, FinishedAt, TotalNodes, SucceededNodes, FailedNodes, SkippedNodes
		FROM ClusterRun
		ORDER BY StartedAt DESC
	`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query ClusterRuns: %w", err)
	}
	defer rows.Close()

	var runs []ClusterRun
	for rows.Next() {
		var run ClusterRun
		err := rows.Scan(
			&run.ClusterRunId,
			&run.RunRef,
			&run.CommandKind,
			&run.RawCommand,
			&run.TargetSelector,
			&run.ExceptClause,
			&run.StartedAt,
			&run.FinishedAt,
			&run.TotalNodes,
			&run.SucceededNodes,
			&run.FailedNodes,
			&run.SkippedNodes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ClusterRun: %w", err)
		}
		runs = append(runs, run)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over ClusterRuns: %w", err)
	}

	return runs, nil
}
