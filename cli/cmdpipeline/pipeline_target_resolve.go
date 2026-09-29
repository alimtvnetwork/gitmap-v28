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
	slug, _ := ResolvePipelineTargetAndPath(target)
	return slug
}

// ResolvePipelineTargetAndPath resolves a path, alias, URL, or identifier to (slug, absPath).
func ResolvePipelineTargetAndPath(target string) (string, string) {
	trimmed := strings.TrimSpace(target)
	if trimmed == "" || trimmed == "." {
		return resolveLocalOrCwdSlug(trimmed), ""
	}
	if isPipelineGitURL(trimmed) {
		return parseSlugFromGitURL(trimmed), ""
	}
	if info, err := os.Stat(trimmed); err == nil && info.IsDir() {
		abs, _ := filepath.Abs(trimmed)
		return resolveDirectorySlug(trimmed), abs
	}
	return resolveDBOrExplicitSlug(trimmed)
}

func resolveDBOrExplicitSlug(target string) (string, string) {
	if slug, absPath := resolveSpecialAlias(target); slug != "" {
		return slug, absPath
	}
	if slug, absPath, isFound := queryRepoSlugAndPathFromDB(target); isFound {
		return slug, absPath
	}
	if strings.Contains(target, "/") {
		return target, ""
	}
	return "", ""
}

func resolveSpecialAlias(target string) (string, string) {
	low := strings.ToLower(target)
	switch low {
	case "rs", "repo-secrets":
		if slug, absPath, isFound := queryRepoSlugAndPathFromDB("repo-secrets"); isFound {
			return slug, absPath
		}
	case "rc", "repo-cache", "repo-storage":
		if slug, absPath, isFound := queryRepoSlugAndPathFromDB("repo-cache"); isFound {
			return slug, absPath
		}
	case "gm", "gitmap":
		if slug, absPath, isFound := queryRepoSlugAndPathFromDB("gitmap"); isFound {
			return slug, absPath
		}
	case "cg", "coding-guidelines":
		if slug, absPath, isFound := queryRepoSlugAndPathFromDB("coding-guidelines"); isFound {
			return slug, absPath
		}
	}
	return "", ""
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
	slug, _, ok := queryRepoSlugAndPathFromDB(identifier)
	return slug, ok
}

func queryRepoSlugAndPathFromDB(identifier string) (string, string, bool) {
	sdb, err := store.OpenDefault()
	if err != nil {
		return "", "", false
	}
	defer sdb.Close()
	conn := sdb.Conn()

	if slug, absPath, ok := queryExactAlias(conn, identifier); ok {
		return slug, absPath, true
	}
	if slug, absPath, ok := queryExactRepo(conn, identifier); ok {
		return slug, absPath, true
	}
	return queryPrefixRepo(conn, identifier)
}

func queryExactAlias(conn *sql.DB, identifier string) (string, string, bool) {
	var slug, absPath string
	q := "SELECT r.Slug, r.AbsolutePath FROM Alias a JOIN Repo r ON a.RepoId = r.RepoId WHERE a.Alias = ? LIMIT 1"
	if err := conn.QueryRow(q, identifier).Scan(&slug, &absPath); err == nil && len(slug) > 0 {
		return slug, absPath, true
	}
	return "", "", false
}

func queryExactRepo(conn *sql.DB, identifier string) (string, string, bool) {
	var slug, absPath string
	q := "SELECT Slug, AbsolutePath FROM Repo WHERE Slug = ? OR RepoName = ? OR AbsolutePath = ? OR AbsolutePath LIKE ? OR AbsolutePath LIKE ? LIMIT 1"
	slashPattern := "%/" + filepath.ToSlash(identifier)
	backslashPattern := "%\\" + filepath.FromSlash(identifier)
	if err := conn.QueryRow(q, identifier, identifier, identifier, slashPattern, backslashPattern).Scan(&slug, &absPath); err == nil && len(slug) > 0 {
		return slug, absPath, true
	}
	return "", "", false
}

func queryPrefixRepo(conn *sql.DB, identifier string) (string, string, bool) {
	q := "SELECT Slug, AbsolutePath FROM Repo WHERE Slug LIKE ? OR RepoName LIKE ? ORDER BY Slug DESC LIMIT 5"
	rows, err := conn.Query(q, identifier+"%", identifier+"%")
	if err != nil {
		return "", "", false
	}
	defer rows.Close()
	var slug, absPath, s, a string
	count := 0
	for rows.Next() {
		if errScan := rows.Scan(&s, &a); errScan != nil {
			return "", "", false
		}
		if count == 0 {
			slug, absPath = s, a
		}
		count++
	}
	if count > 0 && slug != "" {
		return slug, absPath, true
	}
	return "", "", false
}
