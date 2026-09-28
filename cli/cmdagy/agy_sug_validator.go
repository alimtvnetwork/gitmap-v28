package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ValidateSUGTarget checks if target is a local path, git URL, GitMap DB repo, or active AGY project.
func ValidateSUGTarget(raw string) (string, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", false
	}

	// 1. Local filesystem directory
	if info, err := os.Stat(trimmed); err == nil && info.IsDir() {
		return resolveLocalDirPath(trimmed)
	}


	// 2. Git remote URL
	if isGitRemoteURL(trimmed) {
		return trimmed, true
	}

	// 3. GitMap DB indexed repository
	if path, isFound := matchGitMapDBRepo(trimmed); isFound {
		return path, true
	}

	// 4. Antigravity project
	if path, isFound := matchAGYProject(trimmed); isFound {
		return path, true
	}

	return "", false
}

func resolveLocalDirPath(trimmed string) (string, bool) {
	abs, err := filepath.Abs(trimmed)
	if err != nil {
		return trimmed, true
	}
	return abs, true
}

func isGitRemoteURL(target string) bool {

	low := strings.ToLower(target)
	return strings.HasPrefix(low, "http://") ||
		strings.HasPrefix(low, "https://") ||
		strings.HasPrefix(low, "git@") ||
		strings.HasPrefix(low, "ssh://") ||
		strings.HasSuffix(low, ".git")
}

func matchGitMapDBRepo(target string) (string, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return "", false
	}
	defer db.Close()

	repos, errList := db.ListRepos()
	if errList != nil {
		return "", false
	}

	for _, r := range repos {
		if strings.EqualFold(r.Slug, target) || strings.EqualFold(r.RepoName, target) {
			return r.AbsolutePath, true
		}
		if strings.EqualFold(filepath.Base(r.AbsolutePath), target) {
			return r.AbsolutePath, true
		}
	}
	return "", false
}

func matchAGYProject(target string) (string, bool) {
	projects, err := getAllProjects()
	if err != nil {
		return "", false
	}
	return findMatchingAGYProject(projects, target)
}

func findMatchingAGYProject(projects []AgyProject, target string) (string, bool) {
	for _, p := range projects {
		path := p.GetPath()
		if isAGYProjectNameMatch(p, target) {
			return resolveAGYProjectPathOrName(p, path), true
		}
		if path != "" && strings.EqualFold(filepath.Base(path), target) {
			return path, true
		}
	}
	return "", false
}

func isAGYProjectNameMatch(p AgyProject, target string) bool {
	return strings.EqualFold(p.Name, target) || strings.EqualFold(p.ID, target)
}

func resolveAGYProjectPathOrName(p AgyProject, path string) string {
	if path != "" {
		return path
	}
	return p.Name
}


// PrintSUGTargetValidationDiagnostic prints actionable error messages when target validation fails.
func PrintSUGTargetValidationDiagnostic(target string) {
	fmt.Printf("\n  %s✖ Target not found:%s %q\n", constants.ColorRed, constants.ColorReset, target)
	fmt.Printf("  Target must be a local folder, GitMap repository name, AGY project, or Git URL.\n\n")
	fmt.Printf("  %sTo discover valid targets, use:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap list%s                 (list all indexed local repositories)\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy ls%s               (list all Antigravity projects)\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap agy running-projects%s (list projects running active prompts)\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %sgitmap scan <directory>%s     (index local repositories into GitMap)\n\n", constants.ColorYellow, constants.ColorReset)
}
