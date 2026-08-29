package db

import (
	"context"
	"database/sql"
	"time"
)

type RepoCGVersion struct {
	RepoAlias   string
	Version     string
	InstalledAt time.Time
}

func InsertOrUpdateRepoCGVersion(ctx context.Context, db *sql.DB, version RepoCGVersion) error {
	query := `
		INSERT INTO RepoCGVersion (
			RepoAlias, Version, InstalledAt
		) VALUES (?, ?, ?)
		ON CONFLICT(RepoAlias) DO UPDATE SET
			Version = excluded.Version,
			InstalledAt = excluded.InstalledAt
	`
	_, err := db.ExecContext(ctx, query,
		version.RepoAlias,
		version.Version,
		version.InstalledAt,
	)
	return err
}

func GetRepoCGVersion(ctx context.Context, db *sql.DB, repoAlias string) (*RepoCGVersion, error) {
	query := `SELECT RepoAlias, Version, InstalledAt FROM RepoCGVersion WHERE RepoAlias = ?`
	row := db.QueryRowContext(ctx, query, repoAlias)

	var v RepoCGVersion
	err := row.Scan(&v.RepoAlias, &v.Version, &v.InstalledAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}
