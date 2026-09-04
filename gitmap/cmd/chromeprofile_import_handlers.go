// Package cmd — chromeprofile_import_handlers.go: multi-format snapshot importers
// (SQLite DB, multi-profile JSON, YAML, ZIP).
package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

func importChromeSnapshot(srcFile, targetName string) error {
	return importChromeSnapshotWithOptions(srcFile, targetName, 0)
}

func importChromeSnapshotWithOptions(srcFile, targetName string, limit int) error {
	lower := strings.ToLower(srcFile)
	if isSQLiteSnapshot(lower) {
		return importChromeSQLite(srcFile, targetName, limit)
	}
	if strings.HasSuffix(lower, constants.ExtZIP) {
		return importChromeArchiveWithOptions(srcFile, targetName, limit)
	}
	if strings.HasSuffix(lower, constants.ExtYAML) || strings.HasSuffix(lower, constants.ExtYML) {
		return importChromeYAML(srcFile, targetName, limit)
	}
	if strings.HasSuffix(lower, constants.ExtCSV) {
		return importChromeCSV(srcFile, targetName)
	}
	return importChromeJSON(srcFile, targetName, limit)
}

func isSQLiteSnapshot(lower string) bool {
	return strings.HasSuffix(lower, constants.ExtDB) ||
		strings.HasSuffix(lower, constants.ExtSQLite) ||
		strings.HasSuffix(lower, ".sqlite3")
}

func importChromeCSV(srcFile, targetName string) error {
	exp, err := readChromeExportCSV(srcFile)
	if err != nil {
		return err
	}
	return applyImportedProfile(exp, srcFile, targetName)
}

func importChromeJSON(srcFile, targetName string, limit int) error {
	raw, err := os.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", srcFile, err)
	}
	var all chromeAllProfilesExport
	if json.Unmarshal(raw, &all) == nil && len(all.Profiles) > 0 {
		checkSnapshotVersion(all.GitMapVersion)
		return dispatchImportAllProfiles(all.Profiles, srcFile, targetName, limit)
	}
	var exp chromeExport
	if err := json.Unmarshal(raw, &exp); err != nil {
		return fmt.Errorf("parse %s: %w", srcFile, err)
	}
	checkSnapshotVersion(exp.GitMapVersion)
	return applyImportedProfile(&exp, srcFile, targetName)
}

func importChromeYAML(srcFile, targetName string, limit int) error {
	raw, err := os.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", srcFile, err)
	}
	var all chromeAllProfilesExport
	if yaml.Unmarshal(raw, &all) == nil && len(all.Profiles) > 0 {
		checkSnapshotVersion(all.GitMapVersion)
		return dispatchImportAllProfiles(all.Profiles, srcFile, targetName, limit)
	}
	var exp chromeExport
	if err := yaml.Unmarshal(raw, &exp); err != nil {
		return fmt.Errorf("parse %s: %w", srcFile, err)
	}
	checkSnapshotVersion(exp.GitMapVersion)
	return applyImportedProfile(&exp, srcFile, targetName)
}

func dispatchImportAllProfiles(profiles []chromeExport, srcFile, targetName string, limit int) error {
	if targetName != "" {
		return importSingleProfileFromCollection(profiles, srcFile, targetName)
	}
	if limit > 0 && limit < len(profiles) {
		profiles = profiles[:limit]
	}
	return applyAllProfilesExport(profiles, srcFile)
}

func importSingleProfileFromCollection(profiles []chromeExport, srcFile, targetName string) error {
	for _, p := range profiles {
		if p.Name == targetName || p.DisplayName == targetName {
			return applyImportedProfile(&p, srcFile, targetName)
		}
	}
	return fmt.Errorf("profile %q not found in snapshot %s", targetName, srcFile)
}

func applyAllProfilesExport(profiles []chromeExport, srcFile string) error {
	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name
	}
	maxW := maxChromeProfileLabelWidth(names)
	for _, p := range profiles {
		dstPath := chromeProfilePath(p.Name)
		if err := applyChromeExport(&p, dstPath); err != nil {
			fmt.Fprintf(os.Stderr, "  \033[1;91m✗ %s:\033[0m %v\n", p.Name, err)
			continue
		}
		label := formatChromeProfileLabel(p.Name, p.Preferences)
		pad := calculateLabelPadding(maxW, label)
		fmt.Printf("  \033[1;92m✓\033[0m %s%s → imported\n", label, strings.Repeat(" ", pad))
	}
	return nil
}

func importChromeSQLite(dbPath, targetName string, limit int) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open sqlite %s: %w", dbPath, err)
	}
	defer db.Close()

	var dbVer string
	_ = db.QueryRow("SELECT value FROM gitmap_metadata WHERE key = 'gitmap_version'").Scan(&dbVer)
	checkSnapshotVersion(dbVer)

	rows, err := db.Query("SELECT name, display_name FROM chrome_profiles")
	if err != nil {
		return fmt.Errorf("query chrome_profiles from %s: %w", dbPath, err)
	}
	defer rows.Close()

	return restoreProfilesFromSQLite(db, rows, dbPath, targetName, limit)
}

func restoreProfilesFromSQLite(db *sql.DB, rows *sql.Rows, dbPath, targetName string, limit int) error {
	restoredCount := 0
	for rows.Next() {
		if limit > 0 && restoredCount >= limit {
			break
		}
		var name, displayName string
		if err := rows.Scan(&name, &displayName); err != nil {
			continue
		}
		if targetName != "" && name != targetName && displayName != targetName {
			continue
		}
		if err := restoreSingleProfileFromSQLite(db, name, dbPath); err != nil {
			fmt.Fprintf(os.Stderr, "  \033[1;91m✗ %s:\033[0m %v\n", name, err)
			continue
		}
		restoredCount++
	}
	return nil
}

func restoreSingleProfileFromSQLite(db *sql.DB, name, dbPath string) error {
	dstPath := chromeProfilePath(name)
	exp := &chromeExport{SchemaVersion: chromeExportSchemaVersion, Name: name}

	var prefs, bms sql.NullString
	_ = db.QueryRow("SELECT preferences_json FROM chrome_preferences WHERE profile_name = ?", name).Scan(&prefs)
	if prefs.Valid && len(prefs.String) > 0 {
		exp.Preferences = json.RawMessage(prefs.String)
	}
	_ = db.QueryRow("SELECT bookmarks_json FROM chrome_bookmarks WHERE profile_name = ?", name).Scan(&bms)
	if bms.Valid && len(bms.String) > 0 {
		exp.Bookmarks = json.RawMessage(bms.String)
	}
	loadExtensionsFromSQLite(db, name, exp)
	if err := applyChromeExport(exp, dstPath); err != nil {
		return err
	}
	restoreBlobsFromSQLite(db, name, dstPath)
	fmt.Printf(constants.MsgChromeProfileImportOk, dbPath, name)
	return nil
}

func loadExtensionsFromSQLite(db *sql.DB, name string, exp *chromeExport) {
	extRows, err := db.Query("SELECT extension_id FROM chrome_extensions WHERE profile_name = ?", name)
	if err != nil {
		return
	}
	defer extRows.Close()
	for extRows.Next() {
		var extID string
		if extRows.Scan(&extID) == nil {
			exp.ExtensionIDs = append(exp.ExtensionIDs, extID)
		}
	}
}

func restoreBlobsFromSQLite(db *sql.DB, name, dstPath string) {
	blobRows, err := db.Query("SELECT file_name, payload FROM chrome_blobs WHERE profile_name = ?", name)
	if err != nil {
		return
	}
	defer blobRows.Close()
	for blobRows.Next() {
		var fileName string
		var payload []byte
		if blobRows.Scan(&fileName, &payload) == nil && len(payload) > 0 {
			_ = os.WriteFile(filepath.Join(dstPath, fileName), payload, constants.FilePermission)
		}
	}
}
