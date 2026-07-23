// Package store manages the SQLite database for gitmap.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// ZipGroupWithCount extends ZipGroup with an item count.
type ZipGroupWithCount struct {
	model.ZipGroup
	ItemCount int
}

// CreateZipGroup inserts a new zip group.
func (db *DB) CreateZipGroup(name, archiveName string) (model.ZipGroup, error) {
	_, err := db.conn.Exec(constants.SQLInsertZipGroup, name, archiveName)
	if err != nil {
		return model.ZipGroup{}, fmt.Errorf(constants.ErrZGCreate, err)
	}

	return db.FindZipGroupByName(name)
}

// FindZipGroupByName retrieves a zip group by its name.
func (db *DB) FindZipGroupByName(name string) (model.ZipGroup, error) {
	row := db.conn.QueryRow(constants.SQLSelectZipGroupByName, name)

	var g model.ZipGroup

	err := row.Scan(&g.ID, &g.Name, &g.ArchiveName, &g.CreatedAt)
	if err != nil {
		return model.ZipGroup{}, fmt.Errorf(constants.ErrZGNotFound, name)
	}

	return g, nil
}

// ListZipGroups returns all zip groups ordered by name.
func (db *DB) ListZipGroups() ([]model.ZipGroup, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllZipGroups)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrZGQuery, err)
	}
	defer rows.Close()

	var groups []model.ZipGroup

	for rows.Next() {
		var g model.ZipGroup

		err := rows.Scan(&g.ID, &g.Name, &g.ArchiveName, &g.CreatedAt)
		if err != nil {
			continue
		}

		groups = append(groups, g)
	}

	return groups, nil
}

// ListZipGroupsWithCount returns all zip groups with item counts.
func (db *DB) ListZipGroupsWithCount() ([]ZipGroupWithCount, error) {
	rows, err := db.conn.Query(constants.SQLSelectAllZipGroupsWithCount)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrZGQuery, err)
	}
	defer rows.Close()

	var groups []ZipGroupWithCount

	for rows.Next() {
		var g ZipGroupWithCount

		err := rows.Scan(&g.ID, &g.Name, &g.ArchiveName, &g.CreatedAt, &g.ItemCount)
		if err != nil {
			continue
		}

		groups = append(groups, g)
	}

	return groups, nil
}

// DeleteZipGroup removes a zip group by name (items cascade).
func (db *DB) DeleteZipGroup(name string) error {
	result, err := db.conn.Exec(constants.SQLDeleteZipGroup, name)
	if err != nil {
		return fmt.Errorf(constants.ErrZGDelete, err)
	}

	return checkDeleted(result, name)
}

// UpdateZipGroupArchive sets a custom archive name for a group.
func (db *DB) UpdateZipGroupArchive(name, archiveName string) error {
	_, err := db.conn.Exec(constants.SQLUpdateZipGroupArchive, archiveName, name)
	if err != nil {
		return fmt.Errorf(constants.ErrZGCreate, err)
	}

	return nil
}

// AddZipGroupItem adds a file or folder to a zip group with full path metadata.
func (db *DB) AddZipGroupItem(groupName, repoPath, relativePath, fullPath string, isFolder bool) error {
	g, err := db.FindZipGroupByName(groupName)
	if err != nil {
		return err
	}

	folderFlag := 0
	if isFolder {
		folderFlag = 1
	}

	_, err = db.conn.Exec(constants.SQLInsertZipGroupItem, g.ID, repoPath, relativePath, fullPath, folderFlag)
	if err != nil {
		return fmt.Errorf(constants.ErrZGAddItem, err)
	}

	return nil
}

// RemoveZipGroupItem removes a path from a zip group.
func (db *DB) RemoveZipGroupItem(groupName, fullPath string) error {
	g, err := db.FindZipGroupByName(groupName)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(constants.SQLDeleteZipGroupItem, g.ID, fullPath)
	if err != nil {
		return fmt.Errorf(constants.ErrZGRemoveItem, err)
	}

	return nil
}

// ListZipGroupItems returns all items in a zip group.
func (db *DB) ListZipGroupItems(groupName string) ([]model.ZipGroupItem, error) {
	g, err := db.FindZipGroupByName(groupName)
	if err != nil {
		return nil, err
	}

	rows, err := db.conn.Query(constants.SQLSelectZipGroupItems, g.ID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrZGQuery, err)
	}
	defer rows.Close()

	var items []model.ZipGroupItem

	for rows.Next() {
		var item model.ZipGroupItem
		var folderFlag int

		err := rows.Scan(&item.GroupID, &item.RepoPath, &item.RelativePath, &item.FullPath, &folderFlag)
		if err != nil {
			continue
		}

		item.IsFolder = folderFlag == 1
		item.Path = item.FullPath
		items = append(items, item)
	}

	return items, nil
}

// CountZipGroupItems returns the number of items in a zip group.
func (db *DB) CountZipGroupItems(groupName string) (int, error) {
	g, err := db.FindZipGroupByName(groupName)
	if err != nil {
		return 0, err
	}

	var count int

	err = db.conn.QueryRow(constants.SQLCountZipGroupItems, g.ID).Scan(&count)

	return count, err
}

// ZipGroupExists returns true if a zip group with the given name exists.
func (db *DB) ZipGroupExists(name string) bool {
	_, err := db.FindZipGroupByName(name)

	return err == nil
}
