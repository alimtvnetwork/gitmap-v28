// Package store — pendingtaskscan.go contains row scanners for pending tasks.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// scanPendingTaskRows reads all rows into PendingTaskRecord slices.
func scanPendingTaskRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.PendingTaskRecord, error) {
	var results []model.PendingTaskRecord

	for rows.Next() {
		var r model.PendingTaskRecord

		err := rows.Scan(&r.ID, &r.TaskTypeId, &r.TaskTypeName, &r.TargetPath,
			&r.WorkingDirectory, &r.SourceCommand, &r.CommandArgs,
			&r.FailureReason, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrPendingTaskQuery, err)
		}

		results = append(results, r)
	}

	return results, nil
}

// scanCompletedTaskRows reads all rows into CompletedTaskRecord slices.
func scanCompletedTaskRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.CompletedTaskRecord, error) {
	var results []model.CompletedTaskRecord

	for rows.Next() {
		var r model.CompletedTaskRecord

		err := rows.Scan(&r.ID, &r.OriginalTaskId, &r.TaskTypeId, &r.TaskTypeName,
			&r.TargetPath, &r.WorkingDirectory, &r.SourceCommand, &r.CommandArgs,
			&r.CompletedAt, &r.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrPendingTaskQuery, err)
		}

		results = append(results, r)
	}

	return results, nil
}
