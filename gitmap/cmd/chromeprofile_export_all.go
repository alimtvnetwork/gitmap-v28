// Package cmd — chromeprofile_export_all.go: multi-profile and smart format
// exporter for Chrome profiles (JSON, SQLite DB, YAML, ZIP).
package cmd

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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
		return 0, fmt.Errorf("marshal all profiles YAML: %w", err)
	}
	if err := os.WriteFile(outPath, raw, constants.FilePermission); err != nil {
		return 0, fmt.Errorf("write %s: %w", outPath, err)
	}
	printAllProfilesExportRows(all.Profiles, names, "YAML snapshot")
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
	);`
	_, err := db.Exec(schema)
	return err
}

func populateProfilesInSQLite(db *sql.DB, names []string) error {
	maxW := maxChromeProfileLabelWidth(names)
	for _, name := range names {
		exp, hasExp := loadSingleChromeProfileExport(name)
		if !hasExp {
			continue
		}
		if err := insertProfileToSQLite(db, name, exp); err != nil {
			return err
		}
		label := formatChromeProfileLabel(name, exp.Preferences)
		pad := calculateLabelPadding(maxW, label)
		fmt.Printf("  \033[1;92m✓\033[0m %s%s → SQLite tables populated\n", label, strings.Repeat(" ", pad))
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

	_ = addManifestToZip(zw, names)
	maxW := maxChromeProfileLabelWidth(names)
	for _, name := range names {
		writeProfileToZipWithProgress(zw, name, maxW)
	}
	_ = zw.Close()
	return getFileSize(outPath), nil
}

func writeProfileToZipWithProgress(zw *zip.Writer, name string, maxW int) {
	srcPath, hasDir := resolveChromeProfileDir(name)
	if !hasDir {
		return
	}
	_ = addProfileToZip(zw, srcPath, name)
	label := formatChromeProfileLabel(name, nil)
	pad := calculateLabelPadding(maxW, label)
	fmt.Printf("  \033[1;92m✓\033[0m %s%s → packaged in zip\n", label, strings.Repeat(" ", pad))
}

func addManifestToZip(zw *zip.Writer, names []string) error {
	m := chromeProfileManifest{
		GitMapVersion: constants.Version,
		SchemaVersion: chromeExportSchemaVersion,
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		ProfileCount:  len(names),
		Profiles:      make([]chromeManifestProfile, 0, len(names)),
	}
	for _, name := range names {
		displayName, email := resolveProfileNameAndEmail(name, nil)
		srcPath, _ := resolveChromeProfileDir(name)
		extCount := len(listExtensionIDs(filepath.Join(srcPath, "Extensions")))
		m.Profiles = append(m.Profiles, chromeManifestProfile{
			Name:           name,
			DisplayName:    displayName,
			Email:          email,
			ExtensionCount: extCount,
		})
	}
	raw, err := json.MarshalIndent(m, "", constants.JSONIndent)
	if err != nil {
		return err
	}
	w, err := zw.Create("manifest.json")
	if err != nil {
		return err
	}
	_, err = w.Write(raw)
	return err
}

func addProfileToZip(zw *zip.Writer, srcPath, name string) error {
	tmpJSON := filepath.Join(os.TempDir(), name+"-snapshot.json")
	if _, err := writeChromeExport(srcPath, name, tmpJSON); err == nil {
		defer os.Remove(tmpJSON)
		_ = copyFileToZipPath(zw, tmpJSON, filepath.ToSlash(filepath.Join(name, name+".json")))
	}
	_ = copyFileToZipPath(zw, filepath.Join(srcPath, "Bookmarks"), filepath.ToSlash(filepath.Join(name, "Bookmarks")))
	_ = copyFileToZipPath(zw, filepath.Join(srcPath, "Preferences"), filepath.ToSlash(filepath.Join(name, "Preferences")))
	for _, dbName := range constants.ChromeProfileSQLiteEntries {
		_ = copyFileToZipPath(zw, filepath.Join(srcPath, dbName), filepath.ToSlash(filepath.Join(name, dbName)))
	}
	return nil
}

func copyFileToZipPath(zw *zip.Writer, srcPath, zipPath string) error {
	info, err := os.Stat(srcPath)
	if err != nil || info.IsDir() {
		return nil
	}
	f, err := os.Open(srcPath)
	if err != nil {
		return nil
	}
	defer f.Close()
	w, err := zw.Create(zipPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err
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
	case displayName != "" && displayName != dirName && displayName != email && email != "":
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
