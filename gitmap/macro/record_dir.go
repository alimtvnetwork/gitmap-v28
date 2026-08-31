package macro

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// DirTracker maintains current and previous working directories during interactive recording and execution.
type DirTracker struct {
	CurrentDir string
	PrevDir    string
}

// NewDirTracker initializes a new directory tracker.
func NewDirTracker(initialDir string) *DirTracker {
	return &DirTracker{
		CurrentDir: initialDir,
		PrevDir:    initialDir,
	}
}

// ProcessCd checks if a command changes directory and updates tracked paths.
func (dt *DirTracker) ProcessCd(cmdText string) bool {
	target, isStandardCd, isGitmapCd := parseDirectoryChange(cmdText)
	if !isStandardCd && !isGitmapCd {
		return false
	}
	resolved := dt.resolveTarget(target, isGitmapCd)
	if resolved == "" {
		return false
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.IsDir() {
		return false
	}
	dt.PrevDir = dt.CurrentDir
	dt.CurrentDir = resolved
	fmt.Printf("  ➜ 📁 Directory: %s%s%s\n", constants.ColorGreen, dt.CurrentDir, constants.ColorReset)
	return true
}

func (dt *DirTracker) resolveTarget(target string, isGitmapCd bool) string {
	if isGitmapCd {
		return dt.resolveGitmapTarget(target)
	}
	if target == "-" {
		return dt.PrevDir
	}
	if filepath.IsAbs(target) {
		return filepath.Clean(target)
	}
	return filepath.Clean(filepath.Join(dt.CurrentDir, target))
}

func (dt *DirTracker) resolveGitmapTarget(target string) string {
	repoPath, hasRepo := resolveGitmapCD(target, dt.CurrentDir)
	if hasRepo {
		return repoPath
	}
	return ""
}

func parseDirectoryChange(cmdText string) (string, bool, bool) {
	trimmed := strings.TrimSpace(cmdText)
	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return "", false, false
	}
	if target, isCd := parseStandardCd(trimmed, parts); isCd {
		return target, true, false
	}
	if target, isGitmapCd := parseGitmapCd(trimmed, parts); isGitmapCd {
		return target, false, true
	}
	return "", false, false
}

func parseStandardCd(trimmed string, parts []string) (string, bool) {
	lower0 := strings.ToLower(parts[0])
	isCdCmd := (lower0 == "cd" || lower0 == "chdir") && len(parts) >= 2
	if !isCdCmd {
		return "", false
	}
	target := strings.TrimSpace(trimmed[len(parts[0]):])
	return target, true
}

func parseGitmapCd(trimmed string, parts []string) (string, bool) {
	lower0 := strings.ToLower(parts[0])
	isGitmapPrefix := lower0 == "gitmap" || lower0 == "gitmap.exe" || lower0 == "gitmap-v28"
	if !isGitmapPrefix || len(parts) < 3 {
		return "", false
	}
	lower1 := strings.ToLower(parts[1])
	if lower1 != "cd" {
		return "", false
	}
	target := strings.TrimSpace(trimmed[len(parts[0])+1+len(parts[1]):])
	return target, true
}

func resolveGitmapCD(repoName, currentDir string) (string, bool) {
	candidate := filepath.Join(currentDir, repoName)
	info, err := os.Stat(candidate)
	if err == nil && info.IsDir() {
		return filepath.Clean(candidate), true
	}
	if path, hasPath := queryDBForRepo(repoName); hasPath {
		return path, true
	}
	return "", false
}

func queryDBForRepo(repoName string) (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	dbPath := filepath.Join(home, constants.GitMapDir, "gitmap.db")
	if info, statErr := os.Stat(dbPath); statErr != nil || info.Size() == 0 {
		return "", false
	}
	return queryOpenDBForRepo(repoName)
}

func queryOpenDBForRepo(repoName string) (string, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return "", false
	}
	defer db.Close()
	cleanName := strings.TrimRight(repoName, "/\\")
	repos, err := db.FindBySlug(strings.ToLower(cleanName))
	if err == nil && len(repos) > 0 {
		return repos[0].AbsolutePath, true
	}
	return findMatchingRepoInList(db, cleanName)
}

func findMatchingRepoInList(db *store.DB, cleanName string) (string, bool) {
	all, listErr := db.ListRepos()
	if listErr != nil {
		return "", false
	}
	for _, r := range all {
		isNameMatch := strings.EqualFold(r.RepoName, cleanName)
		isBaseMatch := strings.EqualFold(filepath.Base(r.AbsolutePath), cleanName)
		if isNameMatch || isBaseMatch {
			return r.AbsolutePath, true
		}
	}
	return "", false
}
