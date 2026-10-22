// Package store — zipgroupjson.go persists zip group data to .gitmap/zip-groups.json.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// ZipGroupJSON represents the JSON structure for .gitmap/zip-groups.json.
type ZipGroupJSON struct {
	Groups []ZipGroupJSONEntry `json:"groups"`
}

// ZipGroupJSONEntry represents a single zip group with its items.
type ZipGroupJSONEntry struct {
	Name        string               `json:"name"`
	ArchiveName string               `json:"archiveName,omitempty"`
	Items       []model.ZipGroupItem `json:"items"`
}

// WriteZipGroupsJSON persists all zip groups to .gitmap/zip-groups.json
// in the given repo root directory.
func (db *DB) WriteZipGroupsJSON(repoRoot string) error {
	groups, err := db.ListZipGroups()
	if err != nil {
		return err
	}

	var entries []ZipGroupJSONEntry

	for _, g := range groups {
		items, _ := db.ListZipGroupItems(g.Name)
		entries = append(entries, ZipGroupJSONEntry{
			Name:        g.Name,
			ArchiveName: g.ArchiveName,
			Items:       items,
		})
	}

	data := ZipGroupJSON{Groups: entries}

	dir := filepath.Join(repoRoot, constants.ZGJSONDir)

	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("failed to write zip-groups.json at %s (operation: write)", dir))
	}

	jsonPath := filepath.Join(dir, constants.ZGJSONFile)

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("failed to write zip-groups.json at %s (operation: write)", jsonPath))
	}

	err = os.WriteFile(jsonPath, jsonBytes, 0o644)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("failed to write zip-groups.json at %s (operation: write)", jsonPath))
	}

	return nil
}
