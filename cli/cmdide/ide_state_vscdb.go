package cmdide

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"
)

const (
	TableItemTable             = "ItemTable"
	KeyRecentlyOpenedPathsList = "history.recentlyOpenedPathsList"
	DriverSQLite               = "sqlite"
	QueryInitItemTable         = "CREATE TABLE IF NOT EXISTS " + TableItemTable + " (key TEXT PRIMARY KEY, value TEXT);"
	QuerySelectRecentPaths     = "SELECT value FROM " + TableItemTable + " WHERE key = '" + KeyRecentlyOpenedPathsList + "'"
	QueryUpsertRecentPaths     = "INSERT OR REPLACE INTO " + TableItemTable + " (key, value) VALUES ('" + KeyRecentlyOpenedPathsList + "', ?)"
)

// VSCodeRecentEntry represents a single item in history.recentlyOpenedPathsList.
type VSCodeRecentEntry struct {
	FolderURI string `json:"folderUri,omitempty"`
	FileURI   string `json:"fileUri,omitempty"`
	Label     string `json:"label,omitempty"`
}

// VSCodeRecentList encapsulates the entries array stored in state.vscdb.
type VSCodeRecentList struct {
	Entries []VSCodeRecentEntry `json:"entries"`
}

// ResolveStateVscdbPaths discovers state.vscdb candidate paths for VS Code and Cursor.
func ResolveStateVscdbPaths() []string {
	if runtime.GOOS == "windows" {
		return resolveWindowsVscdbPaths(os.Getenv("APPDATA"))
	}

	home, errHome := os.UserHomeDir()
	if errHome != nil {
		return nil
	}

	if runtime.GOOS == "darwin" {
		return resolveDarwinVscdbPaths(home)
	}

	return resolveUnixVscdbPaths(home)
}

func resolveWindowsVscdbPaths(appData string) []string {
	if appData == "" {
		return nil
	}

	return []string{
		filepath.Join(appData, "Code", "User", "globalStorage", "state.vscdb"),
		filepath.Join(appData, "Cursor", "User", "globalStorage", "state.vscdb"),
	}
}

func resolveDarwinVscdbPaths(home string) []string {
	appSupport := filepath.Join(home, "Library", "Application Support")

	return []string{
		filepath.Join(appSupport, "Code", "User", "globalStorage", "state.vscdb"),
		filepath.Join(appSupport, "Cursor", "User", "globalStorage", "state.vscdb"),
	}
}

func resolveUnixVscdbPaths(home string) []string {
	configDir := filepath.Join(home, ".config")

	return []string{
		filepath.Join(configDir, "Code", "User", "globalStorage", "state.vscdb"),
		filepath.Join(configDir, "Cursor", "User", "globalStorage", "state.vscdb"),
	}
}

func ensureVscdbParentDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, "ensureVscdbParentDir: mkdir failed")
	}

	return nil
}

func openVscdb(dbPath string) (*sql.DB, error) {
	if errDir := ensureVscdbParentDir(dbPath); errDir != nil {
		return nil, errDir
	}

	conn, errOpen := sql.Open(DriverSQLite, dbPath)
	if errOpen != nil {
		return nil, apperror.WrapSimple(errOpen, "openVscdb: sql.Open failed")
	}

	return conn, nil
}

func initItemTable(conn *sql.DB) error {
	if _, err := conn.Exec(QueryInitItemTable); err != nil {
		return apperror.WrapSimple(err, "initItemTable: create table failed")
	}

	return nil
}

func readRecentRawJSON(conn *sql.DB) string {
	var raw string
	row := conn.QueryRow(QuerySelectRecentPaths)
	if row.Scan(&raw) != nil {
		return ""
	}

	return raw
}

func loadRecentEntries(conn *sql.DB) []VSCodeRecentEntry {
	raw := readRecentRawJSON(conn)
	if raw == "" {
		return nil
	}

	var data VSCodeRecentList
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil
	}

	return data.Entries
}

func filterExistingRecentEntry(entries []VSCodeRecentEntry, folderURI string) []VSCodeRecentEntry {
	filtered := make([]VSCodeRecentEntry, 0, len(entries))
	for _, e := range entries {
		if !strings.EqualFold(e.FolderURI, folderURI) && !strings.EqualFold(e.FileURI, folderURI) {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

func prependRecentEntry(entries []VSCodeRecentEntry, folderURI, label string) []VSCodeRecentEntry {
	filtered := filterExistingRecentEntry(entries, folderURI)
	newEntry := VSCodeRecentEntry{
		FolderURI: folderURI,
		Label:     label,
	}

	return append([]VSCodeRecentEntry{newEntry}, filtered...)
}

func saveRecentEntries(conn *sql.DB, entries []VSCodeRecentEntry) error {
	payload := VSCodeRecentList{Entries: entries}
	bytes, errJSON := json.Marshal(payload)
	if errJSON != nil {
		return apperror.WrapSimple(errJSON, "saveRecentEntries: json marshal")
	}

	if _, errExec := conn.Exec(QueryUpsertRecentPaths, string(bytes)); errExec != nil {
		return apperror.WrapSimple(errExec, "saveRecentEntries: sql exec failed")
	}

	return nil
}

// UpdateStateVscdb adds or updates a folder entry in a specific state.vscdb database.
func UpdateStateVscdb(dbPath, folderURI, label string) error {
	conn, errOpen := openVscdb(dbPath)
	if errOpen != nil {
		return errOpen
	}

	defer conn.Close()

	if errInit := initItemTable(conn); errInit != nil {
		return errInit
	}

	entries := loadRecentEntries(conn)
	updated := prependRecentEntry(entries, folderURI, label)

	return saveRecentEntries(conn, updated)
}

// AddRecentWorkspaceToAllStateVscdb updates all discovered VS Code and Cursor state.vscdb databases.
func AddRecentWorkspaceToAllStateVscdb(folderURI, label string) []string {
	paths := ResolveStateVscdbPaths()
	var registered []string

	for _, p := range paths {
		if err := UpdateStateVscdb(p, folderURI, label); err == nil {
			editor := resolveEditorNameFromPath(p)
			registered = append(registered, editor)
		}
	}

	return registered
}

func resolveEditorNameFromPath(p string) string {
	lower := strings.ToLower(filepath.ToSlash(p))
	if strings.Contains(lower, string(TargetCursor)) {
		return string(TargetCursor)
	}

	return string(TargetVSCode)
}

// UpdateStateVscdbForVSCode updates history.recentlyOpenedPathsList in all VS Code state.vscdb databases.
func UpdateStateVscdbForVSCode(folderURI, label string) error {
	paths := ResolveStateVscdbPaths()
	var lastErr error
	for _, p := range paths {
		if !strings.Contains(strings.ToLower(filepath.ToSlash(p)), "code") {
			continue
		}

		if err := UpdateStateVscdb(p, folderURI, label); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// UpdateStateVscdbForCursor updates history.recentlyOpenedPathsList in all Cursor state.vscdb databases.
func UpdateStateVscdbForCursor(folderURI, label string) error {
	paths := ResolveStateVscdbPaths()
	var lastErr error
	for _, p := range paths {
		if !strings.Contains(strings.ToLower(filepath.ToSlash(p)), string(TargetCursor)) {
			continue
		}

		if err := UpdateStateVscdb(p, folderURI, label); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// PinnedProjectEntry represents an item in .gemini/config/pinned_projects.json.
type PinnedProjectEntry struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Branch   string `json:"branch,omitempty"`
	PinnedAt string `json:"pinnedAt"`
}

// PinnedProjectsConfig encapsulates the pinned_projects.json schema.
type PinnedProjectsConfig struct {
	Version   string               `json:"version"`
	UpdatedAt string               `json:"updatedAt"`
	Projects  []PinnedProjectEntry `json:"projects"`
}

// SyncAntigravityPinned updates .gemini/config/pinned_projects.json in pure Go.
func SyncAntigravityPinned(repoPath, repoName string) error {
	home, errHome := os.UserHomeDir()
	if errHome != nil {
		return errHome
	}

	pinnedPath := filepath.Join(home, ".gemini", "config", "pinned_projects.json")
	if errDir := os.MkdirAll(filepath.Dir(pinnedPath), 0755); errDir != nil {
		return errDir
	}

	return appendPinnedStoreEntry(pinnedPath, repoPath, repoName)
}

func appendPinnedStoreEntry(pinnedPath, repoPath, repoName string) error {
	nowStr := time.Now().UTC().Format(time.RFC3339Nano)
	store := loadOrCreatePinnedStore(pinnedPath, nowStr)
	entry := PinnedProjectEntry{
		ID:       uuidutil.GenerateUUID(),
		Name:     repoName,
		Path:     repoPath,
		Branch:   "main",
		PinnedAt: nowStr,
	}

	store.Projects = upsertPinnedEntry(store.Projects, entry)
	store.UpdatedAt = nowStr
	bytes, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pinnedPath, bytes, 0644)
}

func loadOrCreatePinnedStore(pinnedPath, nowStr string) PinnedProjectsConfig {
	data, err := os.ReadFile(pinnedPath)
	if err != nil || len(data) == 0 {
		return defaultPinnedStore(nowStr)
	}

	var store PinnedProjectsConfig
	if errJSON := json.Unmarshal(data, &store); errJSON != nil {
		return defaultPinnedStore(nowStr)
	}

	return store
}

func defaultPinnedStore(nowStr string) PinnedProjectsConfig {
	return PinnedProjectsConfig{
		Version:   "1.0.0",
		UpdatedAt: nowStr,
		Projects:  make([]PinnedProjectEntry, 0),
	}
}

func upsertPinnedEntry(list []PinnedProjectEntry, entry PinnedProjectEntry) []PinnedProjectEntry {
	idx := -1
	for i, p := range list {
		if strings.EqualFold(p.Path, entry.Path) || strings.EqualFold(p.Name, entry.Name) {
			idx = i
			break
		}
	}

	if idx >= 0 {
		list[idx] = entry

		return list
	}

	return append(list, entry)
}
