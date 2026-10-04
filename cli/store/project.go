// Package store — project.go manages DetectedProject CRUD operations.
package store

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// UpsertDetectedProject inserts or updates a detected project record.
func (db *DB) UpsertDetectedProject(p model.DetectedProject) error {
	validated, err := db.validateProjectForeignKeys(p)
	if err != nil {
		return err
	}

	_, err = ExecWrapper(db.conn, constants.SQLUpsertDetectedProject,
		validated.RepoID, validated.ProjectTypeID, validated.ProjectName,
		validated.AbsolutePath, validated.RepoPath, validated.RelativePath, validated.PrimaryIndicator).Destruct()

	return err
}

func (db *DB) validateProjectForeignKeys(p model.DetectedProject) (model.DetectedProject, error) {
	repoID, err := db.resolveProjectRepoID(p)
	if err != nil {
		return p, err
	}

	p.RepoID = repoID
	if err := db.ensureProjectTypeExists(p.ProjectTypeID); err != nil {
		return p, err
	}

	return p, nil
}

func (db *DB) resolveProjectRepoID(p model.DetectedProject) (int64, error) {
	if db.RepoExists(p.RepoID) {
		return p.RepoID, nil
	}

	if id, err := db.SelectRepoIDByPath(p.RepoPath); err == nil && id > 0 {
		return id, nil
	}

	if id, err := db.SelectRepoIDByPath(p.AbsolutePath); err == nil && id > 0 {
		return id, nil
	}

	return 0, fmt.Errorf("cannot upsert detected project: RepoId %d not found in Repo table for %s", p.RepoID, p.AbsolutePath)
}

func (db *DB) ensureProjectTypeExists(typeID int64) error {
	if typeID <= 0 {
		return fmt.Errorf("cannot upsert detected project: invalid ProjectTypeID %d", typeID)
	}

	var exists bool
	_ = db.conn.QueryRow("SELECT EXISTS(SELECT 1 FROM ProjectType WHERE ProjectTypeId = ?)", typeID).Scan(&exists)
	if !exists {
		_ = db.SeedProjectTypes()
		_ = db.conn.QueryRow("SELECT EXISTS(SELECT 1 FROM ProjectType WHERE ProjectTypeId = ?)", typeID).Scan(&exists)
	}

	if !exists {
		return fmt.Errorf("cannot upsert detected project: ProjectTypeId %d not found in ProjectType table", typeID)
	}

	return nil
}

// SelectDetectedProjectID returns the persisted ID for a project identity tuple.
func (db *DB) SelectDetectedProjectID(
	repoID,
	projectTypeID int64,
	relativePath string,
) (int64, error) {
	var id int64
	err := db.conn.QueryRow(constants.SQLSelectDetectedProjectID,
		repoID, projectTypeID, relativePath).Scan(&id)

	return id, err
}

// SelectProjectsByTypeKey returns all detected projects of a given type.
func (db *DB) SelectProjectsByTypeKey(key string) ([]model.DetectedProject, error) {
	rows, err := QueryWrapper(db.conn, constants.SQLSelectProjectsByTypeKey, key).Destruct()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrProjectQuery, err)
	}

	defer rows.Close()

	return scanProjectRows(rows)
}

// CountProjectsByTypeKey returns the count of projects for a given type.
func (db *DB) CountProjectsByTypeKey(key string) (int, error) {
	var count int
	err := db.conn.QueryRow(constants.SQLCountProjectsByTypeKey, key).Scan(&count)

	return count, err
}

// DeleteStaleProjects removes projects not in the given ID list for a repo.
func (db *DB) DeleteStaleProjects(repoID int64, keepIDs []int64) (int64, error) {
	if len(keepIDs) == 0 {
		return 0, nil
	}

	placeholders := buildPlaceholders(len(keepIDs))
	query := fmt.Sprintf(constants.SQLDeleteStaleProjects, placeholders)
	args := buildStaleArgsInt64(repoID, keepIDs)
	result, err := ExecWrapper(db.conn, query, args...).Destruct()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// scanProjectRows scans SQL rows into DetectedProject slices.
func scanProjectRows(rows interface {
	Next() bool
	Scan(...interface{}) error
}) ([]model.DetectedProject, error) {
	var projects []model.DetectedProject
	for rows.Next() {
		var p model.DetectedProject
		err := rows.Scan(&p.ID, &p.RepoID, &p.ProjectType, &p.ProjectName,
			&p.AbsolutePath, &p.RepoPath, &p.RelativePath,
			&p.PrimaryIndicator, &p.DetectedAt, &p.RepoName)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

// buildPlaceholders creates a comma-separated list of ? placeholders.
func buildPlaceholders(count int) string {
	p := make([]string, count)
	for i := range p {
		p[i] = "?"
	}

	return strings.Join(p, ", ")
}

// buildStaleArgsInt64 creates the argument slice for stale cleanup queries with int64 IDs.
func buildStaleArgsInt64(parentID int64, keepIDs []int64) []interface{} {
	args := make([]interface{}, 0, len(keepIDs)+1)
	args = append(args, parentID)
	for _, id := range keepIDs {
		args = append(args, id)
	}

	return args
}
