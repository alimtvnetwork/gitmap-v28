package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ReleaseAcrossRepos is a Release row joined with the owning Repo's slug
// and absolute path. Returned by ListReleasesAcrossRepos for the
// `gitmap releases --all-repos` view.
type ReleaseAcrossRepos struct {
	ReleaseID    int64
	RepoID       int64
	RepoSlug     string
	RepoPath     string
	Version      string
	Tag          string
	Branch       string
	CommitSha    string
	Source       string
	IsDraft      bool
	IsLatest     bool
	IsPreRelease bool
	CreatedAt    string
}

// ListReleasesAcrossRepos returns every Release row in the DB joined with
// its owning Repo, ordered by CreatedAt DESC. This is the multi-repo batch
// view that exercises the IdxRelease_RepoId index added in v17.
func (db *DB) ListReleasesAcrossRepos() ([]ReleaseAcrossRepos, error) {
	if !db.tableExists("Release") || !db.tableExists("Repo") {
		return nil, nil
	}
	if !db.columnExists("Release", "RepoId") {
		return nil, nil
	}

	rows, err := db.conn.Query(constants.SQLSelectAllReleasesAcrossRepos)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBReleaseQuery, err)
	}
	defer rows.Close()

	return scanAcrossRepoRows(rows)
}

// scanAcrossRepoRows materializes the join result into typed records.
func scanAcrossRepoRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]ReleaseAcrossRepos, error) {
	var out []ReleaseAcrossRepos
	for rows.Next() {
		rec, err := scanOneAcrossRepoRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, rec)
	}

	return out, nil
}

// scanOneAcrossRepoRow scans a single joined Release+Repo row.
func scanOneAcrossRepoRow(row interface{ Scan(dest ...any) error }) (ReleaseAcrossRepos, error) {
	var r ReleaseAcrossRepos
	var draft, latest, pre int
	err := row.Scan(&r.ReleaseID, &r.RepoID, &r.RepoSlug, &r.RepoPath,
		&r.Version, &r.Tag, &r.Branch, &r.CommitSha, &r.Source,
		&draft, &latest, &pre, &r.CreatedAt)
	if err != nil {
		return ReleaseAcrossRepos{}, err
	}
	r.IsDraft = draft == 1
	r.IsLatest = latest == 1
	r.IsPreRelease = pre == 1

	return r, nil
}
