package cmdpipeline

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
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

	if isDirectoryTarget(trimmed) {
		return resolveDirectoryTarget(trimmed)
	}

	return resolveDBOrExplicitSlug(trimmed)
}

func isDirectoryTarget(target string) bool {
	info, err := os.Stat(target)

	return err == nil && info.IsDir()
}

func resolveDirectoryTarget(dirPath string) (string, string) {
	abs, absErr := filepath.Abs(dirPath)
	if absErr != nil {
		return resolveDirectorySlug(dirPath), dirPath
	}

	return resolveDirectorySlug(dirPath), abs
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

	return constants.DefaultGitHubOwner + "/" + target, ""
}

func resolveAliasKey(low string) string {
	switch low {
	case "rs", "repo-secrets":
		return "repo-secrets"
	case "rc", "repo-cache", "repo-storage":
		return "repo-cache"
	case "gm", "gitmap":
		return "gitmap"
	case "cg", "coding-guidelines":
		return "coding-guidelines"
	default:
		return ""
	}
}

func resolveSpecialAlias(target string) (string, string) {
	key := resolveAliasKey(strings.ToLower(target))
	if key == "" {
		return "", ""
	}

	slug, absPath, isFound := queryRepoSlugAndPathFromDB(key)
	if !isFound {
		return "", ""
	}

	return slug, absPath
}

func resolveLocalOrCwdSlug(target string) string {
	dir := "."
	if target != "" {
		dir = target
	}

	origin := gitutil.GetRemoteOriginURL(dir)
	if origin != "" {
		return parseSlugFromGitURL(origin)
	}

	return resolveCurrentRepoSlug()
}

// DirectorySlugResult wraps resolved directory slug provenance.
type DirectorySlugResult struct {
	Slug      string
	IsFound   bool
	IsGuessed bool
	Err       error
}

func resolveDirectorySlug(dirPath string) string {
	res := resolveDirectorySlugWithProvenance(dirPath)

	return res.Slug
}

func resolveOriginDirectorySlug(dirPath string) (DirectorySlugResult, bool) {
	origin := gitutil.GetRemoteOriginURL(dirPath)
	if origin != "" {
		return DirectorySlugResult{
			Slug:    parseSlugFromGitURL(origin),
			IsFound: true,
		}, true
	}

	return DirectorySlugResult{}, false
}

func resolveDBDirectorySlug(dirPath string) (DirectorySlugResult, bool) {
	base := filepath.Base(dirPath)
	if slug, isFound := queryRepoSlugFromDB(base); isFound {
		return DirectorySlugResult{
			Slug:    slug,
			IsFound: true,
		}, true
	}

	return DirectorySlugResult{}, false
}

func resolveDirectorySlugWithProvenance(dirPath string) DirectorySlugResult {
	if res, hasOrigin := resolveOriginDirectorySlug(dirPath); hasOrigin {
		return res
	}

	if res, hasDB := resolveDBDirectorySlug(dirPath); hasDB {
		return res
	}

	base := filepath.Base(dirPath)

	return DirectorySlugResult{
		Slug:      constants.DefaultGitHubOwner + "/" + base,
		IsFound:   false,
		IsGuessed: true,
	}
}

func isPipelineGitURL(s string) bool {
	low := strings.ToLower(s)

	return strings.HasPrefix(low, string(constants.GitPrefixHTTPS)) ||
		strings.HasPrefix(low, string(constants.GitPrefixHTTP)) ||
		strings.HasPrefix(low, string(constants.GitPrefixSSH)) ||
		strings.HasPrefix(low, string(constants.GitPrefixSSHScheme)) ||
		strings.HasSuffix(low, constants.GitDotGitSuffix) ||
		strings.Contains(low, constants.GitHostTokenGitHubSlash) ||
		strings.Contains(low, constants.GitHostTokenGitHubColon)
}

func queryRepoSlugFromDB(identifier string) (string, bool) {
	slug, _, isFound := queryRepoSlugAndPathFromDB(identifier)

	return slug, isFound
}

func queryRepoSlugAndPathFromDB(identifier string) (string, string, bool) {
	sdb, err := store.OpenDefault()
	if err != nil {
		return "", "", false
	}
	defer sdb.Close()

	conn := sdb.Conn()
	if slug, absPath, isFound := queryExactAlias(conn, identifier); isFound {
		return slug, absPath, true
	}

	if slug, absPath, isFound := queryExactRepo(conn, identifier); isFound {
		return slug, absPath, true
	}

	return queryPrefixRepo(conn, identifier)
}

func queryExactAlias(conn *sql.DB, identifier string) (string, string, bool) {
	var slug, absPath, httpsUrl, sshUrl string
	q := "SELECT r.Slug, r.AbsolutePath, COALESCE(r.HttpsUrl, ''), COALESCE(r.SshUrl, '') FROM Alias a JOIN Repo r ON a.RepoId = r.RepoId WHERE a.Alias = ? LIMIT 1"
	if err := conn.QueryRow(q, identifier).Scan(&slug, &absPath, &httpsUrl, &sshUrl); err == nil && len(slug) > 0 {
		return canonicalizeRepoSlug(slug, httpsUrl, sshUrl), absPath, true
	}

	return "", "", false
}

func queryRepoByPattern(conn *sql.DB, name, slashPat, backslashPat string) (string, string, bool) {
	var slug, absPath, httpsUrl, sshUrl string
	q := "SELECT Slug, AbsolutePath, COALESCE(HttpsUrl, ''), COALESCE(SshUrl, '') FROM Repo WHERE Slug = ? OR RepoName = ? OR AbsolutePath = ? OR AbsolutePath LIKE ? OR AbsolutePath LIKE ? LIMIT 1"
	if err := conn.QueryRow(q, name, name, name, slashPat, backslashPat).Scan(&slug, &absPath, &httpsUrl, &sshUrl); err == nil && len(slug) > 0 {
		return canonicalizeRepoSlug(slug, httpsUrl, sshUrl), absPath, true
	}

	return "", "", false
}

func queryExactRepo(conn *sql.DB, identifier string) (string, string, bool) {
	slashPat := "%/" + filepath.ToSlash(identifier)
	backslashPat := "%\\" + filepath.FromSlash(identifier)
	if slug, absPath, isFound := queryRepoByPattern(conn, identifier, slashPat, backslashPat); isFound {
		return slug, absPath, true
	}

	if !strings.Contains(identifier, "/") {
		return "", "", false
	}

	shortName := identifier[strings.LastIndex(identifier, "/")+1:]

	return queryRepoByPattern(conn, shortName, slashPat, backslashPat)
}

func scanFirstRepoRow(rows *sql.Rows) (string, string, string, string, bool) {
	var slug, absPath, httpsUrl, sshUrl string
	hasRow := false
	for rows.Next() {
		var s, a, h, sc string
		if errScan := rows.Scan(&s, &a, &h, &sc); errScan != nil {
			return "", "", "", "", false
		}

		if !hasRow {
			slug, absPath, httpsUrl, sshUrl = s, a, h, sc
			hasRow = true
		}
	}

	return slug, absPath, httpsUrl, sshUrl, hasRow
}

func queryPrefixRepo(conn *sql.DB, identifier string) (string, string, bool) {
	q := "SELECT Slug, AbsolutePath, COALESCE(HttpsUrl, ''), COALESCE(SshUrl, '') FROM Repo WHERE Slug LIKE ? OR RepoName LIKE ? ORDER BY Slug DESC LIMIT 5"
	rows, err := conn.Query(q, identifier+"%", identifier+"%")
	if err != nil {
		return "", "", false
	}
	defer rows.Close()

	slug, absPath, httpsUrl, sshUrl, hasRow := scanFirstRepoRow(rows)
	if !hasRow || slug == "" {
		return "", "", false
	}

	return canonicalizeRepoSlug(slug, httpsUrl, sshUrl), absPath, true
}

func slugFromRemoteURL(url string) string {
	if len(url) == 0 {
		return ""
	}

	parsed := parseSlugFromGitURL(url)
	if strings.Contains(parsed, "/") {
		return parsed
	}

	return ""
}

func canonicalizeRepoSlug(slug, httpsUrl, sshUrl string) string {
	if strings.Contains(slug, "/") {
		return slug
	}

	if parsed := slugFromRemoteURL(httpsUrl); parsed != "" {
		return parsed
	}

	if parsed := slugFromRemoteURL(sshUrl); parsed != "" {
		return parsed
	}

	if len(slug) > 0 {
		return constants.DefaultGitHubOwner + "/" + slug
	}

	return ""
}
