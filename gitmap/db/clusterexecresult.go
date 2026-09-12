package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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

const (
	sqlInsertClusterExecResult = `
		INSERT INTO ClusterExecResult (
			ClusterRunId, NodeId, SubCommand, CommandText, ResultStatus,
			ExitCode, Stdout, Stderr, StartedAt, FinishedAt, DurationMs, ErrorMessage
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	sqlUpdateClusterExecResult = `
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
	sqlSelectClusterExecResultsByRunId = `
		SELECT 
			ClusterExecResultId, ClusterRunId, NodeId, SubCommand, CommandText,
			ResultStatus, ExitCode, Stdout, Stderr, StartedAt, FinishedAt, DurationMs, ErrorMessage
		FROM ClusterExecResult
		WHERE ClusterRunId = ?
		ORDER BY ClusterExecResultId ASC
	`
)

func capString(s *string, maxLen int) *string {
	if s == nil {
		return nil
	}

	hasExceeded := len(*s) > maxLen
	if hasExceeded {
		capped := (*s)[:maxLen]

		return &capped
	}

	return s
}

func InsertClusterExecResult(
	ctx context.Context,
	db *sql.DB,
	result ClusterExecResult,
) (int64, *apperror.AppError) {
	result.Stdout = capString(result.Stdout, 64*1024)
	result.Stderr = capString(result.Stderr, 16*1024)

	res, err := db.ExecContext(ctx, sqlInsertClusterExecResult,
		result.ClusterRunId, result.NodeId, result.SubCommand, result.CommandText,
		result.ResultStatus, result.ExitCode, result.Stdout, result.Stderr,
		result.StartedAt, result.FinishedAt, result.DurationMs, result.ErrorMessage,
	)
	if err != nil {
		return 0, apperror.WrapSimple(err, "InsertClusterExecResult.Exec")
	}

	id, errId := res.LastInsertId()
	if errId != nil {
		return 0, apperror.WrapSimple(errId, "InsertClusterExecResult.LastInsertId")
	}

	return id, nil
}

func UpdateClusterExecResult(ctx context.Context, db *sql.DB, result ClusterExecResult) *apperror.AppError {
	_, err := db.ExecContext(ctx, sqlUpdateClusterExecResult,
		result.CommandText, result.ResultStatus, result.ExitCode, result.Stdout,
		result.Stderr, result.StartedAt, result.FinishedAt, result.DurationMs,
		result.ErrorMessage, result.ClusterExecResultId,
	)
	if err != nil {
		return apperror.WrapSimple(err, "UpdateClusterExecResult.Exec")
	}

	return nil
}

func SelectClusterExecResultsByRunId(
	ctx context.Context,
	db *sql.DB,
	runId int64,
) ([]ClusterExecResult, *apperror.AppError) {
	rows, err := db.QueryContext(ctx, sqlSelectClusterExecResultsByRunId, runId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "SelectClusterExecResultsByRunId.Query")
	}

	defer rows.Close()

	return scanClusterExecResultRows(rows)
}

func scanClusterExecResultRows(rows *sql.Rows) ([]ClusterExecResult, *apperror.AppError) {
	var results []ClusterExecResult
	for rows.Next() {
		var res ClusterExecResult
		err := rows.Scan(
			&res.ClusterExecResultId, &res.ClusterRunId, &res.NodeId, &res.SubCommand,
			&res.CommandText, &res.ResultStatus, &res.ExitCode, &res.Stdout,
			&res.Stderr, &res.StartedAt, &res.FinishedAt, &res.DurationMs, &res.ErrorMessage,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scanClusterExecResultRows.Scan")
		}

		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "scanClusterExecResultRows.Rows")
	}

	return results, nil
}
