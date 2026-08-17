package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ClusterExecResult struct {
	ClusterExecResultId int64
	ClusterRunId        int64
	NodeId              string
	SubCommand          string
	CommandText         *string
	ResultStatus        ResultStatusType
	ExitCode            *int
	Stdout              *string
	Stderr              *string
	StartedAt           *time.Time
	FinishedAt          *time.Time
	DurationMs          *int
	ErrorMessage        *string
}

func InsertClusterExecResult(ctx context.Context, db *sql.DB, result ClusterExecResult) (int64, error) {
	query := `
		INSERT INTO ClusterExecResult (
			ClusterRunId, NodeId, SubCommand, CommandText, ResultStatus,
			ExitCode, Stdout, Stderr, StartedAt, FinishedAt, DurationMs, ErrorMessage
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := db.ExecContext(ctx, query,
		result.ClusterRunId,
		result.NodeId,
		result.SubCommand,
		result.CommandText,
		result.ResultStatus,
		result.ExitCode,
		result.Stdout,
		result.Stderr,
		result.StartedAt,
		result.FinishedAt,
		result.DurationMs,
		result.ErrorMessage,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert ClusterExecResult: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id for ClusterExecResult: %w", err)
	}

	return id, nil
}

func UpdateClusterExecResult(ctx context.Context, db *sql.DB, result ClusterExecResult) error {
	query := `
		UPDATE ClusterExecResult
		SET 
			CommandText = ?,
			ResultStatus = ?,
			ExitCode = ?,
			Stdout = ?,
			Stderr = ?,
			StartedAt = ?,
			FinishedAt = ?,
			DurationMs = ?,
			ErrorMessage = ?
		WHERE ClusterExecResultId = ?
	`
	_, err := db.ExecContext(ctx, query,
		result.CommandText,
		result.ResultStatus,
		result.ExitCode,
		result.Stdout,
		result.Stderr,
		result.StartedAt,
		result.FinishedAt,
		result.DurationMs,
		result.ErrorMessage,
		result.ClusterExecResultId,
	)
	if err != nil {
		return fmt.Errorf("failed to update ClusterExecResult %d: %w", result.ClusterExecResultId, err)
	}

	return nil
}

func SelectClusterExecResultsByRunId(ctx context.Context, db *sql.DB, runId int64) ([]ClusterExecResult, error) {
	query := `
		SELECT 
			ClusterExecResultId, ClusterRunId, NodeId, SubCommand, CommandText,
			ResultStatus, ExitCode, Stdout, Stderr, StartedAt, FinishedAt, DurationMs, ErrorMessage
		FROM ClusterExecResult
		WHERE ClusterRunId = ?
		ORDER BY ClusterExecResultId ASC
	`
	rows, err := db.QueryContext(ctx, query, runId)
	if err != nil {
		return nil, fmt.Errorf("failed to query ClusterExecResults for run %d: %w", runId, err)
	}
	defer rows.Close()

	var results []ClusterExecResult
	for rows.Next() {
		var res ClusterExecResult
		err := rows.Scan(
			&res.ClusterExecResultId,
			&res.ClusterRunId,
			&res.NodeId,
			&res.SubCommand,
			&res.CommandText,
			&res.ResultStatus,
			&res.ExitCode,
			&res.Stdout,
			&res.Stderr,
			&res.StartedAt,
			&res.FinishedAt,
			&res.DurationMs,
			&res.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ClusterExecResult: %w", err)
		}
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over ClusterExecResults: %w", err)
	}

	return results, nil
}
