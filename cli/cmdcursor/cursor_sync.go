package cmdcursor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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

func getCursorUserDataRoot() (string, error) {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Cursor"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Cursor"), nil
	default:
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg == "" {
			xdg = filepath.Join(home, ".config")
		}
		return filepath.Join(xdg, "Cursor"), nil
	}
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
		p := r.AbsolutePath
		if p == "" {
			p = r.RelativePath
		}
		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths, nil
}

func mergeReposIntoCursorProjects(existing []ProjectManagerEntry, repoPaths []string) ([]ProjectManagerEntry, int) {
	byPath := make(map[string]bool)
	for _, e := range existing {
		byPath[filepath.Clean(e.RootPath)] = true
	}
	addedCount := 0
	for _, p := range repoPaths {
		clean := filepath.Clean(p)
		if !byPath[clean] {
			existing = append(existing, ProjectManagerEntry{
				Name:     filepath.Base(clean),
				RootPath: clean,
				Paths:    []string{},
				Tags:     []string{"gitmap"},
				Enabled:  true,
			})
			byPath[clean] = true
			addedCount++
		}
	}
	return existing, addedCount
}

// RunCursorSync synchronizes GitMap repositories into Cursor Project Manager.
func RunCursorSync(_ []string) error {
	jsonPath, err := GetCursorProjectsJSONPath()
	if err != nil {
		return err
	}
	existing, _ := loadProjectsFromDisk(jsonPath)
	repoPaths, fetchErr := fetchGitMapRepoPaths()
	if fetchErr != nil {
		return fetchErr
	}
	merged, added := mergeReposIntoCursorProjects(existing, repoPaths)
	if saveErr := saveProjectsToDisk(jsonPath, merged); saveErr != nil {
		return saveErr
	}
	fmt.Printf("%s✔ Synchronized Cursor Project Manager:%s %s\n", constants.ColorGreen, constants.ColorReset, jsonPath)
	fmt.Printf("    • Added:     %d\n", added)
	fmt.Printf("    • Total:     %d\n", len(merged))
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

