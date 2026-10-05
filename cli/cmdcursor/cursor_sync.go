package cmdcursor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ProjectManagerEntry mirrors the alefragnani.project-manager entry schema.
type ProjectManagerEntry struct {
	Name     string   `json:"name"`
	RootPath string   `json:"rootPath"`
	Paths    []string `json:"paths"`
	Tags     []string `json:"tags"`
	Enabled  bool     `json:"enabled"`
}

func resolvePlatformConfigRoot(home string) string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "Cursor")
	}
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		xdg = filepath.Join(home, ".config")
	}
	return filepath.Join(xdg, "Cursor")
}

func resolveWindowsAppData(home string) string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return filepath.Join(home, "AppData", "Roaming")
	}
	return appData
}

func getCursorUserDataRoot() (string, error) {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(resolveWindowsAppData(home), "Cursor"), nil
	}
	return resolvePlatformConfigRoot(home), nil
}

// GetCursorProjectsJSONPath resolves the path to Cursor's projects.json file.
func GetCursorProjectsJSONPath() (string, error) {
	root, err := getCursorUserDataRoot()
	if err != nil {
		return "", err
	}
	extDir := filepath.Join(root, "User", "globalStorage", "alefragnani.project-manager")
	return filepath.Join(extDir, "projects.json"), nil
}

// GetCursorBackupProjectsJSONPath resolves the path to Cursor's backup projects.json.
func GetCursorBackupProjectsJSONPath() (string, error) {
	root, err := getCursorUserDataRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "User", "projects.json"), nil
}

func loadProjectsFromDisk(path string) ([]ProjectManagerEntry, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []ProjectManagerEntry{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read cursor projects.json")
	}
	var entries []ProjectManagerEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return []ProjectManagerEntry{}, nil
	}
	return entries, nil
}

func saveProjectsToDisk(path string, entries []ProjectManagerEntry) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir cursor project manager dir")
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal cursor projects")
	}
	return os.WriteFile(path, data, 0644)
}

func getSecondaryProjectsJSONPath() string {
	if runtime.GOOS != "windows" {
		return ""
	}
	username := os.Getenv("USERNAME")
	sysDrive := os.Getenv("SystemDrive")
	if sysDrive == "" {
		sysDrive = "C:"
	}
	return filepath.Join(sysDrive, "Users", username, "AppData", "Roaming", "Cursor", "User", "globalStorage", "alefragnani.project-manager", "projects.json")
}

func saveBackupProjects(jsonPath string, merged []ProjectManagerEntry) {
	if backupPath, bErr := GetCursorBackupProjectsJSONPath(); bErr == nil && backupPath != "" {
		_ = saveProjectsToDisk(backupPath, merged)
	}
	if secPath := getSecondaryProjectsJSONPath(); secPath != "" && secPath != jsonPath {
		_ = saveProjectsToDisk(secPath, merged)
	}
}

func isGitRepoDir(dirPath string) bool {
	info, err := os.Stat(filepath.Join(dirPath, ".git"))
	return err == nil && info.IsDir()
}

func scanDirectoryRepoPaths(baseDir string) ([]string, error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read directory for repos")
	}
	var paths []string
	for _, e := range entries {
		fullPath := filepath.Join(baseDir, e.Name())
		if e.IsDir() && isGitRepoDir(fullPath) {
			paths = append(paths, fullPath)
		}
	}
	return paths, nil
}

func appendRepoPath(paths []string, r model.ScanRecord) []string {
	p := r.AbsolutePath
	if p == "" {
		p = r.RelativePath
	}
	if p != "" {
		return append(paths, p)
	}
	return paths
}

func fetchGitMapRepoPaths() ([]string, error) {
	s, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open gitmap store")
	}
	defer s.Close()
	repos, repoErr := s.ListRepos()
	if repoErr != nil {
		return nil, apperror.WrapSimple(repoErr, "list repos")
	}
	paths := make([]string, 0, len(repos))
	for _, r := range repos {
		paths = appendRepoPath(paths, r)
	}
	return paths, nil
}

func createProjectEntry(cleanPath string) ProjectManagerEntry {
	return ProjectManagerEntry{
		Name:     filepath.Base(cleanPath),
		RootPath: cleanPath,
		Paths:    []string{},
		Tags:     []string{"gitmap"},
		Enabled:  true,
	}
}

func indexExistingPaths(existing []ProjectManagerEntry) map[string]bool {
	byPath := make(map[string]bool, len(existing))
	for _, e := range existing {
		byPath[filepath.Clean(e.RootPath)] = true
	}
	return byPath
}

func mergeReposIntoCursorProjects(existing []ProjectManagerEntry, repoPaths []string) ([]ProjectManagerEntry, int) {
	byPath := indexExistingPaths(existing)
	addedCount := 0
	for _, p := range repoPaths {
		clean := filepath.Clean(p)
		if !byPath[clean] {
			existing = append(existing, createProjectEntry(clean))
			byPath[clean] = true
			addedCount++
		}
	}
	return existing, addedCount
}

// SyncCursorProjects merges target repository paths into Cursor Project Manager.
func SyncCursorProjects(repoPaths []string) (int, int, error) {
	jsonPath, err := GetCursorProjectsJSONPath()
	if err != nil {
		return 0, 0, err
	}
	existing, _ := loadProjectsFromDisk(jsonPath)
	merged, added := mergeReposIntoCursorProjects(existing, repoPaths)
	if saveErr := saveProjectsToDisk(jsonPath, merged); saveErr != nil {
		return 0, 0, saveErr
	}
	saveBackupProjects(jsonPath, merged)
	return added, len(merged), nil
}

func delegateRemoteSync(targetNode string) error {
	fmt.Printf("\n%s● Delegating Cursor Project Sync to remote node:%s %s\n", constants.ColorCyan, constants.ColorReset, targetNode)
	cmd := "python3 repo-secrets/05-scripts/setup-cursor-ubuntu.py --node " + targetNode
	if err := cmdssh.RunSSHExec([]string{targetNode, cmd}); err != nil {
		fmt.Printf("%s⚠ Note: Remote sync notice: %v%s\n", constants.ColorYellow, err, constants.ColorReset)
		return nil
	}
	fmt.Printf("%s✔ Remote Cursor project sync completed on %s.%s\n", constants.ColorGreen, targetNode, constants.ColorReset)
	return nil
}

func resolveSyncPaths(args []string) ([]string, string, error) {
	targetNode := ""
	for i := 0; i < len(args); i++ {
		if (args[i] == "--node" || args[i] == "-n") && i+1 < len(args) {
			targetNode, i = args[i+1], i+1
		} else if isTargetDirectory(args[i]) {
			paths, err := scanDirectoryRepoPaths(args[i])
			return paths, "", err
		}
	}
	if targetNode != "" {
		return nil, targetNode, nil
	}
	paths, err := fetchGitMapRepoPaths()
	return paths, "", err
}

// RunCursorSync synchronizes GitMap repositories into Cursor Project Manager.
func RunCursorSync(args []string) error {
	paths, targetNode, err := resolveSyncPaths(args)
	if targetNode != "" {
		return delegateRemoteSync(targetNode)
	}
	if err != nil {
		return err
	}
	added, total, syncErr := SyncCursorProjects(paths)
	if syncErr != nil {
		return syncErr
	}
	fmt.Printf("%s✔ Synchronized Cursor Project Manager:%s added %d, total %d\n", constants.ColorGreen, constants.ColorReset, added, total)
	return nil
}

// RunCursorListProjects lists all projects currently configured in Cursor PM.
func RunCursorListProjects(_ []string) error {
	jsonPath, err := GetCursorProjectsJSONPath()
	if err != nil {
		return err
	}
	entries, _ := loadProjectsFromDisk(jsonPath)
	fmt.Printf("\n%s● Cursor Project Manager (%d projects):%s\n", constants.ColorCyan, len(entries), constants.ColorReset)
	for _, e := range entries {
		fmt.Printf("  • %-24s %s\n", e.Name, e.RootPath)
	}
	fmt.Println()
	return nil
}
