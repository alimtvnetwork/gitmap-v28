package store

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

const (
	sqlSelectReposByTag = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE Notes LIKE ? COLLATE NOCASE ORDER BY Slug ASC"

	sqlSelectUnpushedRepos = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE (Notes LIKE '%unpushed%' OR Notes LIKE '%ahead%') OR RepoName IN (SELECT RepoSlug FROM HeatmapLocalTracking WHERE IsDirty = 1) ORDER BY Slug ASC"
)

// ListReposByTag returns all repositories matching a specific tag stored in metadata or notes.
func (db *DB) ListReposByTag(tag string) ([]model.ScanRecord, error) {
	cleanTag := strings.TrimSpace(tag)
	hasEmptyTag := cleanTag == ""
	if hasEmptyTag {
		return db.ListRepos()
	}

	pattern := "%" + cleanTag + "%"
	rows, err := QueryWrapper(db.conn, sqlSelectReposByTag, pattern).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}

	defer rows.Close()

	return scanRows(rows)
}

// ListReposUnpushed returns repositories flagged with unpushed commits in metadata or local tracking.
func (db *DB) ListReposUnpushed() ([]model.ScanRecord, error) {
	rows, err := QueryWrapper(db.conn, sqlSelectUnpushedRepos).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}

	defer rows.Close()

	return scanRows(rows)
}
