package cmdpipeline

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// ResolvePipelineTarget resolves a path, alias, URL, or identifier to a canonical repository slug.
func ResolvePipelineTarget(target string) string {
	trimmed := strings.TrimSpace(target)
	if trimmed == "" || trimmed == "." {
		return resolveLocalOrCwdSlug(trimmed)
	}

	// 1. Git remote URL
	if isPipelineGitURL(trimmed) {
		return parseSlugFromGitURL(trimmed)
	}

	// 2. Local directory
	if info, err := os.Stat(trimmed); err == nil && info.IsDir() {
		return resolveDirectorySlug(trimmed)
	}

	// 3. Database repository lookup by slug or alias
	if slug, isFound := queryRepoSlugFromDB(trimmed); isFound {
		return slug
	}

	// 4. Fallback to slug directly or with default owner if missing slash
	if strings.Contains(trimmed, "/") {
		return trimmed
	}
	return "alimtvnetwork/" + trimmed
}

func resolveLocalOrCwdSlug(target string) string {
	dir := "."
	if target != "" {
		dir = target
	}
	out, err := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url").Output()
	if err == nil && len(out) > 0 {
		return parseSlugFromGitURL(strings.TrimSpace(string(out)))
	}
	return resolveCurrentRepoSlug()
}

func resolveDirectorySlug(dirPath string) string {
	out, err := exec.Command("git", "-C", dirPath, "config", "--get", "remote.origin.url").Output()
	if err == nil && len(out) > 0 {
		return parseSlugFromGitURL(strings.TrimSpace(string(out)))
	}
	base := filepath.Base(dirPath)
	if slug, isFound := queryRepoSlugFromDB(base); isFound {
		return slug
	}
	return "alimtvnetwork/" + base
}

func isPipelineGitURL(s string) bool {
	low := strings.ToLower(s)
	return strings.HasPrefix(low, "http://") ||
		strings.HasPrefix(low, "https://") ||
		strings.HasPrefix(low, "git@") ||
		strings.HasPrefix(low, "ssh://") ||
		strings.HasSuffix(low, ".git") ||
		strings.Contains(low, "github.com/") ||
		strings.Contains(low, "github.com:")
}

func queryRepoSlugFromDB(identifier string) (string, bool) {
	dbPath := store.DefaultDBPath()
	info, err := os.Stat(dbPath)
	if err != nil || info.IsDir() {
		return "", false
	}

	conn, err := sql.Open("sqlite", "file:"+filepath.ToSlash(dbPath)+"?mode=ro")
	if err != nil {
		return "", false
	}
	defer conn.Close()

	var slug string
	query := "SELECT Slug FROM Repo WHERE Slug = ? OR RepoName = ? OR AbsolutePath LIKE ? LIMIT 1"
	likePattern := "%" + filepath.ToSlash(identifier)
	row := conn.QueryRow(query, identifier, identifier, likePattern)
	if scanErr := row.Scan(&slug); scanErr == nil && len(slug) > 0 {
		return slug, true
	}
	return "", false
}
