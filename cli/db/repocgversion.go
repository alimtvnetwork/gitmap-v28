package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type RepoCGVersion struct {
	RepoAlias   string
	Version     string
	InstalledAt time.Time
}

const (
	sqlUpsertRepoCGVersion = `
		INSERT INTO RepoCGVersion (
			RepoAlias, Version, InstalledAt
		) VALUES (?, ?, ?)
		ON CONFLICT(RepoAlias) DO UPDATE SET
			Version = excluded.Version,
			InstalledAt = excluded.InstalledAt
	`
	sqlSelectRepoCGVersion = `SELECT RepoAlias, Version, InstalledAt FROM RepoCGVersion WHERE RepoAlias = ?`
)

func InsertOrUpdateRepoCGVersion(ctx context.Context, db *sql.DB, version RepoCGVersion) *apperror.AppError {
	_, err := db.ExecContext(ctx, sqlUpsertRepoCGVersion,
		version.RepoAlias,
		version.Version,
		version.InstalledAt,
	)
	if err != nil {
		return apperror.WrapSimple(err, "InsertOrUpdateRepoCGVersion.Exec")
	}

	return nil
}

func GetRepoCGVersion(ctx context.Context, db *sql.DB, repoAlias string) (*RepoCGVersion, *apperror.AppError) {
	row := db.QueryRowContext(ctx, sqlSelectRepoCGVersion, repoAlias)

	var v RepoCGVersion
	err := row.Scan(&v.RepoAlias, &v.Version, &v.InstalledAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, apperror.WrapSimple(err, "GetRepoCGVersion.Scan")
	}

	return &v, nil
}
