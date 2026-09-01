// Package cmd — chromeprofile_export_all.go: multi-profile and smart format
// exporter for Chrome profiles (JSON, SQLite DB, YAML, ZIP).
package cmd

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

type chromeAllProfilesExport struct {
	SchemaVersion int            `json:"schemaVersion" yaml:"schemaVersion"`
	ExportedAt    string         `json:"exportedAt" yaml:"exportedAt"`
	ProfileCount  int            `json:"profileCount" yaml:"profileCount"`
	Profiles      []chromeExport `json:"profiles" yaml:"profiles"`
}

func inferExportFormatFromPath(path, explicit string) string {
	if explicit != "" && explicit != constants.OutputJSON {
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
	return constants.OutputJSON
}

func buildAllProfilesExport(names []string) chromeAllProfilesExport {
	all := chromeAllProfilesExport{
		SchemaVersion: chromeExportSchemaVersion,
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
	exp := chromeExport{
		SchemaVersion: chromeExportSchemaVersion,
		Name:          name,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		Bookmarks:     readOptionalJSON(filepath.Join(srcPath, "Bookmarks")),
		Preferences:   readOptionalJSON(filepath.Join(srcPath, "Preferences")),
		ExtensionIDs:  listExtensionIDs(filepath.Join(srcPath, "Extensions")),
	}
	return exp, true
}

func writeAllChromeProfilesJSON(names []string, outPath string) (int, error) {
	all := buildAllProfilesExport(names)
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil && filepath.Dir(outPath) != "." {
		return 0, err
	}
	raw, err := json.MarshalIndent(all, "", constants.JSONIndent)
	if err != nil {
		return 0, fmt.Errorf("marshal all profiles JSON: %w", err)
	}
	if err := os.WriteFile(outPath, raw, constants.FilePermission); err != nil {
		return 0, fmt.Errorf("write %s: %w", outPath, err)
	}
	for _, p := range all.Profiles {
		label := formatChromeProfileLabel(p.Name, p.Preferences)
		fmt.Printf("  \033[1;92m✓\033[0m %s → JSON snapshot\n", label)
	}
	return len(raw), nil
}

func writeAllChromeProfilesYAML(names []string, outPath string) (int, error) {
	all := buildAllProfilesExport(names)
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil && filepath.Dir(outPath) != "." {
		return 0, err
	}
	raw, err := yaml.Marshal(all)
	if err != nil {
		return 0, fmt.Errorf("marshal all profiles YAML: %w", err)
	}
	if err := os.WriteFile(outPath, raw, constants.FilePermission); err != nil {
		return 0, fmt.Errorf("write %s: %w", outPath, err)
	}
	for _, p := range all.Profiles {
		label := formatChromeProfileLabel(p.Name, p.Preferences)
		fmt.Printf("  \033[1;92m✓\033[0m %s → YAML snapshot\n", label)
	}
	return len(raw), nil
}

func writeAllChromeProfilesSQLite(names []string, outPath string) (int, error) {
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil && filepath.Dir(outPath) != "." {
		return 0, err
	}
	db, err := sql.Open("sqlite", outPath)
	if err != nil {
		return 0, fmt.Errorf("open sqlite db %s: %w", outPath, err)
	}
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
	);`
	_, err := db.Exec(schema)
	return err
}

func populateProfilesInSQLite(db *sql.DB, names []string) error {
	for _, name := range names {
		exp, hasExp := loadSingleChromeProfileExport(name)
		if !hasExp {
			continue
		}
		if err := insertProfileToSQLite(db, name, exp); err != nil {
			return err
		}
		label := formatChromeProfileLabel(name, exp.Preferences)
		fmt.Printf("  \033[1;92m✓\033[0m %s → SQLite tables populated\n", label)
	}
	return nil
}

func insertProfileToSQLite(db *sql.DB, name string, exp chromeExport) error {
	displayName := chromeProfileDisplayName(name)
	_, err := db.Exec(
		"INSERT OR REPLACE INTO chrome_profiles (name, display_name, exported_at, extension_count) VALUES (?, ?, ?, ?)",
		name, displayName, exp.ExportedAt, len(exp.ExtensionIDs),
	)
	if err != nil {
		return err
	}
	if err := insertOptionalJSONToSQLite(db, "chrome_preferences", "preferences_json", name, exp.Preferences); err != nil {
		return err
	}
	if err := insertOptionalJSONToSQLite(db, "chrome_bookmarks", "bookmarks_json", name, exp.Bookmarks); err != nil {
		return err
	}
	return insertExtensionsAndBlobsToSQLite(db, name, exp.ExtensionIDs)
}

func insertOptionalJSONToSQLite(db *sql.DB, table, col, name string, raw json.RawMessage) error {
	if len(raw) == 0 {
		return nil
	}
	query := fmt.Sprintf("INSERT OR REPLACE INTO %s (profile_name, %s) VALUES (?, ?)", table, col)
	_, err := db.Exec(query, name, string(raw))
	return err
}

func insertExtensionsAndBlobsToSQLite(db *sql.DB, name string, extIDs []string) error {
	for _, id := range extIDs {
		_, _ = db.Exec("INSERT OR REPLACE INTO chrome_extensions (profile_name, extension_id) VALUES (?, ?)", name, id)
	}
	srcPath, _ := resolveChromeProfileDir(name)
	for _, blobFile := range constants.ChromeProfileSQLiteEntries {
		filePath := filepath.Join(srcPath, blobFile)
		if bytes, readErr := os.ReadFile(filePath); readErr == nil && len(bytes) > 0 {
			_, _ = db.Exec("INSERT OR REPLACE INTO chrome_blobs (profile_name, file_name, payload) VALUES (?, ?, ?)", name, blobFile, bytes)
		}
	}
	return nil
}

func writeAllChromeProfilesZIP(names []string, outPath string) (int, error) {
	if err := os.MkdirAll(filepath.Dir(outPath), constants.DirPermission); err != nil && filepath.Dir(outPath) != "." {
		return 0, err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	for _, name := range names {
		srcPath, hasDir := resolveChromeProfileDir(name)
		if hasDir {
			_ = addProfileToZip(zw, srcPath, name)
			label := formatChromeProfileLabel(name, nil)
			fmt.Printf("  \033[1;92m✓\033[0m %s → packaged in zip\n", label)
		}
	}
	_ = zw.Close()
	return getFileSize(outPath), nil
}

func addProfileToZip(zw *zip.Writer, srcPath, name string) error {
	_ = exportChromeProfileJSON(zw, srcPath, name)
	for _, dbName := range constants.ChromeProfileSQLiteEntries {
		_ = copySQLiteEntryToZip(zw, srcPath, dbName)
	}
	return nil
}

func getFileSize(path string) int {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return int(info.Size())
}

func formatChromeProfileLabel(dirName string, prefsRaw json.RawMessage) string {
	displayName, email := resolveProfileNameAndEmail(dirName, prefsRaw)
	switch {
	case displayName != "" && displayName != dirName && email != "":
		return fmt.Sprintf("%s • %s (%s)", dirName, displayName, email)
	case email != "":
		return fmt.Sprintf("%s (%s)", dirName, email)
	case displayName != "" && displayName != dirName:
		return fmt.Sprintf("%s (%s)", dirName, displayName)
	default:
		return dirName
	}
}

func resolveProfileNameAndEmail(dirName string, prefsRaw json.RawMessage) (string, string) {
	displayName, email := fetchLocalStateProfileInfo(dirName)
	if email == "" && len(prefsRaw) > 0 {
		email = extractEmailFromPreferences(prefsRaw)
	}
	if email == "" {
		email = extractEmailFromProfileDisk(dirName)
	}
	return displayName, email
}

func fetchLocalStateProfileInfo(dirName string) (string, string) {
	state := readChromeLocalState()
	if state == nil {
		return "", ""
	}
	info, ok := state.Profile.InfoCache[dirName]
	if !ok {
		return "", ""
	}
	return info.Name, info.UserName
}

func extractEmailFromProfileDisk(dirName string) string {
	srcPath, hasDir := resolveChromeProfileDir(dirName)
	if !hasDir {
		return ""
	}
	prefBytes, err := os.ReadFile(filepath.Join(srcPath, "Preferences"))
	if err != nil || len(prefBytes) == 0 {
		return ""
	}
	return extractEmailFromPreferences(prefBytes)
}

func extractEmailFromPreferences(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var pref struct {
		AccountInfo []struct {
			Email string `json:"email"`
		} `json:"account_info"`
		Signin struct {
			UserName string `json:"user_name"`
		} `json:"signin"`
		Google struct {
			Services struct {
				Username string `json:"username"`
			} `json:"services"`
		} `json:"google"`
	}
	if err := json.Unmarshal(raw, &pref); err != nil {
		return ""
	}
	if len(pref.AccountInfo) > 0 && pref.AccountInfo[0].Email != "" {
		return pref.AccountInfo[0].Email
	}
	if pref.Signin.UserName != "" {
		return pref.Signin.UserName
	}
	return pref.Google.Services.Username
}
