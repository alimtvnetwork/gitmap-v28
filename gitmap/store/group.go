package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// CreateGroup inserts a new group with the given name, description, and color.
func (db *DB) CreateGroup(name, description, color string) (model.Group, error) {
	_, err := db.conn.Exec(constants.SQLInsertGroup, name, description, color)
	if err != nil {
		return model.Group{}, fmt.Errorf(constants.ErrDBGroupCreate, err)
	}

	return db.findGroupByName(name)
}

// ListGroups returns all groups ordered by name.
func (db *DB) ListGroups() ([]model.Group, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllGroups)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBGroupQuery, err)
	}
	defer rows.Close()

	return scanGroupRows(rows)
}

// FindGroupByName returns the group matching the given name.
func (db *DB) findGroupByName(name string) (model.Group, error) {
	row := db.conn.QueryRow(constants.SQLSelectGroupByName, name)
	var g model.Group
	err := row.Scan(&g.ID, &g.Name, &g.Description, &g.Color, &g.CreatedAt)
	if err != nil {
		return model.Group{}, fmt.Errorf(constants.ErrDBGroupNone, name)
	}

	return g, nil
}

// AddRepoToGroup links a repo to a group (silent no-op if already linked).
func (db *DB) AddRepoToGroup(groupName string, repoID int64) error {
	group, err := db.findGroupByName(groupName)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(constants.SQLInsertGroupRepo, group.ID, repoID)
	if err != nil {
		return fmt.Errorf(constants.ErrDBGroupAdd, err)
	}

	return nil
}

// RemoveRepoFromGroup unlinks a repo from a group.
func (db *DB) RemoveRepoFromGroup(groupName string, repoID int64) error {
	group, err := db.findGroupByName(groupName)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(constants.SQLDeleteGroupRepo, group.ID, repoID)
	if err != nil {
		return fmt.Errorf(constants.ErrDBGroupRemove, err)
	}

	return nil
}

// ShowGroup returns all repos belonging to a group.
func (db *DB) ShowGroup(name string) ([]model.ScanRecord, error) {
	group, err := db.findGroupByName(name)
	if err != nil {
		return nil, err
	}

	rows, err := db.conn.Query(constants.SQLSelectGroupRepos, group.ID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrDBQuery, err)
	}
	defer rows.Close()

	return scanRows(rows)
}

// DeleteGroup removes a group by name (repos are not deleted).
func (db *DB) DeleteGroup(name string) error {
	result, err := db.conn.Exec(constants.SQLDeleteGroup, name)
	if err != nil {
		return fmt.Errorf(constants.ErrDBGroupDelete, err)
	}

	return checkDeleted(result, name)
}

// CountGroupRepos returns the number of repos in a group.
func (db *DB) CountGroupRepos(name string) (int, error) {
	group, err := db.findGroupByName(name)
	if err != nil {
		return 0, err
	}

	var count int
	err = db.conn.QueryRow(constants.SQLCountGroupRepos, group.ID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrDBGroupQuery, err)
	}

	return count, nil
}

// scanGroupRows reads Group values from query result rows.
func scanGroupRows(rows interface {
	Next() bool
	Scan(dest ...any) error
}) ([]model.Group, error) {
	var results []model.Group

	for rows.Next() {
		var g model.Group
		err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.Color, &g.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, g)
	}

	return results, nil
}

// checkDeleted verifies at least one row was affected by a delete.
func checkDeleted(result interface{ RowsAffected() (int64, error) }, name string) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(constants.ErrDBGroupDelete, err)
	}

	if affected == 0 {
		return fmt.Errorf(constants.ErrDBGroupNone, name)
	}

	return nil
}
