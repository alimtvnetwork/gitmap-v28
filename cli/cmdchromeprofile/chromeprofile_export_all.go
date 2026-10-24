// Package cmd — chromeprofile_export_all.go: multi-profile and smart format
// exporter for Chrome profiles (JSON, SQLite DB, YAML, ZIP).
package cmdchromeprofile

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

type chromeAllProfilesExport struct {
	SchemaVersion int            `json:"schemaVersion" yaml:"schemaVersion"`
	GitMapVersion string         `json:"gitmapVersion,omitempty" yaml:"gitmapVersion,omitempty"`
	ExportedAt    string         `json:"exportedAt" yaml:"exportedAt"`
	ProfileCount  int            `json:"profileCount" yaml:"profileCount"`
	Profiles      []chromeExport `json:"profiles" yaml:"profiles"`
}

type chromeProfileManifest struct {
	GitMapVersion string                  `json:"gitmapVersion" yaml:"gitmapVersion"`
	SchemaVersion int                     `json:"schemaVersion" yaml:"schemaVersion"`
	ExportedAt    string                  `json:"exportedAt" yaml:"exportedAt"`
	ProfileCount  int                     `json:"profileCount" yaml:"profileCount"`
	Profiles      []chromeManifestProfile `json:"profiles" yaml:"profiles"`
}

type chromeManifestProfile struct {
	Name           string `json:"name" yaml:"name"`
	DisplayName    string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Email          string `json:"email,omitempty" yaml:"email,omitempty"`
	ExtensionCount int    `json:"extensionCount" yaml:"extensionCount"`
}

func stringDisplayWidth(s string) int {
	return runewidth.StringWidth(s)
}

func maxChromeProfileLabelWidth(names []string) int {
	maxW := 0
	for _, name := range names {
		label := formatChromeProfileLabel(name, nil)
		w := stringDisplayWidth(label)
		if w > maxW {
			maxW = w
		}
	}

	return maxW
}

func calculateLabelPadding(maxW int, label string) int {
	pad := maxW - stringDisplayWidth(label)
	if pad < 0 {
		return 0
	}

	return pad
}

func inferExportFormatFromPath(path, explicit string) string {
	if explicit != "" {
		return explicit
	}

	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, constants.ExtDB) || strings.HasSuffix(lower, constants.ExtSQLite) || strings.HasSuffix(lower, ".sqlite3") {
		return constants.OutputSQLite
	}

	if strings.HasSuffix(lower, constants.ExtYAML) || strings.HasSuffix(lower, constants.ExtYML) {
		return constants.OutputYAML
	}

	if strings.HasSuffix(lower, constants.ExtZIP) {
		return constants.OutputZIP
	}

	if strings.HasSuffix(lower, constants.ExtJSON) {
		return constants.OutputJSON
	}

	return constants.OutputZIP
}

func buildAllProfilesExport(names []string) chromeAllProfilesExport {
	all := chromeAllProfilesExport{
		SchemaVersion: chromeExportSchemaVersion,
		GitMapVersion: constants.Version,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		ProfileCount:  len(names),
		Profiles:      make([]chromeExport, 0, len(names)),
	}

	for _, name := range names {
		if exp, hasProfile := loadSingleChromeProfileExport(name); hasProfile {
			all.Profiles = append(all.Profiles, exp)
		}
	}

	return all
}

func loadSingleChromeProfileExport(name string) (chromeExport, bool) {
	srcPath, hasDir := resolveChromeProfileDir(name)
	if !hasDir {
		return chromeExport{}, false
	}

	exp := buildExportFromDisk(name, srcPath)

	return exp, true
}

func buildExportFromDisk(name, srcPath string) chromeExport {
	return buildExportSnapshot(srcPath, name)
}

func writeAllChromeProfilesJSON(names []string, outPath string) (int, error) {
	all := buildAllProfilesExport(names)
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil && filepath.Dir(outPath) != "." {
		return 0, err
	}

	raw, err := json.MarshalIndent(all, "", constants.JSONIndent)
	if err != nil {
		return 0, apperror.WrapSimple(err, "marshal all profiles JSON")
	}

	if err := os.WriteFile(outPath, raw, constants.FilePermission); err != nil {
		return 0, apperror.WrapSimple(err, "write "+outPath)
	}

	printAllProfilesExportRows(all.Profiles, names, "JSON snapshot")

	return len(raw), nil
}

func printAllProfilesExportRows(profiles []chromeExport, names []string, destName string) {
	maxW := maxChromeProfileLabelWidth(names)
	for _, p := range profiles {
		label := formatChromeProfileLabel(p.Name, p.Preferences)
		pad := calculateLabelPadding(maxW, label)
		fmt.Printf("  \033[1;92m✓\033[0m %s%s → %s\n", label, strings.Repeat(" ", pad), destName)
	}
}

func writeAllChromeProfilesYAML(names []string, outPath string) (int, error) {
	all := buildAllProfilesExport(names)
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil && filepath.Dir(outPath) != "." {
		return 0, err
	}

	raw, err := yaml.Marshal(all)
	if err != nil {
		return 0, apperror.WrapSimple(err, "marshal all profiles YAML")
	}

	if err := os.WriteFile(outPath, raw, constants.FilePermission); err != nil {
		return 0, apperror.WrapSimple(err, "write "+outPath)
	}

	printAllProfilesExportRows(all.Profiles, names, "YAML snapshot")

	return len(raw), nil
}

func writeAllChromeProfilesSQLite(names []string, outPath string) (int, error) {
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil && filepath.Dir(outPath) != "." {
		return 0, err
	}

	openRes := store.OpenSQLiteDB(outPath)
	if openRes.IsFailure() {
		return 0, apperror.WrapSimple(openRes.AppError(), "open sqlite db "+outPath)
	}
	db := openRes.Value

	defer db.Close()

	if err := initChromeSQLiteTables(db); err != nil {
		return 0, err
	}

	if err := populateProfilesInSQLite(db, names); err != nil {
		return 0, err
	}

	return getFileSize(outPath), nil
}

func initChromeSQLiteTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS gitmap_metadata (
		key TEXT PRIMARY KEY,
		value TEXT
	);
	CREATE TABLE IF NOT EXISTS chrome_profiles (
		name TEXT PRIMARY KEY,
		display_name TEXT,
		exported_at TEXT,
		extension_count INTEGER
	);
	CREATE TABLE IF NOT EXISTS chrome_preferences (
		profile_name TEXT PRIMARY KEY,
		preferences_json TEXT
	);
	CREATE TABLE IF NOT EXISTS chrome_bookmarks (
		profile_name TEXT PRIMARY KEY,
		bookmarks_json TEXT
	);
	CREATE TABLE IF NOT EXISTS chrome_extensions (
		profile_name TEXT,
		extension_id TEXT,
		PRIMARY KEY (profile_name, extension_id)
	);
	CREATE TABLE IF NOT EXISTS chrome_blobs (
		profile_name TEXT,
		file_name TEXT,
		payload BLOB,
		PRIMARY KEY (profile_name, file_name)
	);
	CREATE TABLE IF NOT EXISTS chrome_tokens (
		profile_name TEXT,
		service TEXT,
		account_id TEXT,
		raw_base64 TEXT,
		double_base64 TEXT,
		PRIMARY KEY (profile_name, service)
	);`
	_, err := db.Exec(schema)

	return err
}

func populateProfilesInSQLite(db *sql.DB, names []string) error {
	wrapper, appErr := dbengine.WrapDB(db, dbengine.DBSQLite)
	if appErr != nil {
		return appErr
	}

	if txErr := wrapper.WithTransaction(context.Background(), func(tx *dbengine.TxWrapper) *apperror.AppError {
		return insertAllProfilesTx(tx, names)
	}); txErr != nil {
		return txErr
	}

	return nil
}

func insertAllProfilesTx(tx *dbengine.TxWrapper, names []string) *apperror.AppError {
	maxW := maxChromeProfileLabelWidth(names)
	for _, name := range names {
		exp, hasExp := loadSingleChromeProfileExport(name)
		if !hasExp {
			continue
		}

		if err := insertProfileToSQLite(tx, name, exp); err != nil {
			return err
		}

		printProfilePopulatedRow(name, exp.Preferences, maxW)
	}

	return nil
}

func printProfilePopulatedRow(name string, prefs json.RawMessage, maxW int) {
	label := formatChromeProfileLabel(name, prefs)
	pad := calculateLabelPadding(maxW, label)
	fmt.Printf("  \033[1;92m✓\033[0m %s%s → SQLite tables populated\n", label, strings.Repeat(" ", pad))
}

func insertProfileToSQLite(tx *dbengine.TxWrapper, name string, exp chromeExport) *apperror.AppError {
	if err := insertProfileMetaToSQLite(tx, name, exp); err != nil {
		return err
	}

	if err := insertTokensToSQLite(tx, name, exp.TokenVault); err != nil {
		return err
	}

	return insertExtensionsAndBlobsToSQLite(tx, name, exp.ExtensionIDs)
}
