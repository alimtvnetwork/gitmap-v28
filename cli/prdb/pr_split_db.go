package prdb

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PrSplitDb encapsulates an isolated SQLite database connection for PR data.
type PrSplitDb struct {
	conn     *sql.DB
	RepoSlug string
	Path     string
}

// Conn returns the underlying SQLite connection.
func (p *PrSplitDb) Conn() *sql.DB {
	return p.conn
}

// DbPath returns the filesystem path of the database.
func (p *PrSplitDb) DbPath() string {
	return p.Path
}

// Slug returns the repository slug for this database.
func (p *PrSplitDb) Slug() string {
	return p.RepoSlug
}

// Close terminates the SQLite connection.
func (p *PrSplitDb) Close() *apperror.AppError {
	if p.conn == nil {
		return nil
	}

	if err := p.conn.Close(); err != nil {
		return apperror.WrapSimple(err, "prdb.close")
	}

	return nil
}

func cleanRepoURLPrefix(repo string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(repo), ".git")
	for _, token := range []string{"github.com/", "github.com:", "gitlab.com/", "gitlab.com:"} {
		if idx := strings.Index(clean, token); idx != -1 {
			return clean[idx+len(token):]
		}
	}

	return clean
}

// SanitizeRepoSlug converts a repository slug or URL into a valid filesystem slug.
func SanitizeRepoSlug(repo string) string {
	lower := strings.ToLower(cleanRepoURLPrefix(repo))
	lower = strings.ReplaceAll(lower, "_", "-")
	slug := lazyregex.SlugSanitizeRegex.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "pr-default"
	}

	return slug
}

// ResolvePrDbDir returns the directory where the PR database lives for a repository.
func ResolvePrDbDir(repoSlug, repoRoot string) string {
	slug := SanitizeRepoSlug(repoSlug)
	if repoRoot != "" {
		return filepath.ToSlash(filepath.Join(repoRoot, ".gitmap", "data", "pr", slug))
	}

	return filepath.ToSlash(filepath.Join(store.BinaryDataDir(), "pr", slug))
}

// ResolvePrDbPath returns the canonical path to sql.db for the PR subsystem.
func ResolvePrDbPath(repoSlug, repoRoot string) string {
	dir := ResolvePrDbDir(repoSlug, repoRoot)
	target := filepath.ToSlash(filepath.Join(dir, "sql.db"))

	return migrateOrFallbackPrDb(dir, repoSlug, target)
}

func isFileExisting(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func migrateOrFallbackPrDb(dir, slug, targetPath string) string {
	candidates := []string{
		filepath.Join(dir, "pr.db"),
		filepath.Join(dir, slug+".db"),
		filepath.Join(store.BinaryDataDir(), "pr", "pr_"+slug+".db"),
	}
	for _, cand := range candidates {
		if isFileExisting(cand) && !isFileExisting(targetPath) {
			_ = os.Rename(cand, targetPath)

			return targetPath
		}
	}

	return targetPath
}
