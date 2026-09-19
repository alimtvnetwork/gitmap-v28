package pipelinedb

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// PipelineSplitDb encapsulates an isolated SQLite database connection for a single repository's pipeline data.
type PipelineSplitDb struct {
	conn     *sql.DB
	RepoSlug string
	Path     string
}

// PipelineSplitDB is an alias to PipelineSplitDb for backward compatibility.
type PipelineSplitDB = PipelineSplitDb

func cleanRepoURLPrefix(repo string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(repo), ".git")
	for _, token := range []string{"github.com/", "github.com:", "gitlab.com/", "gitlab.com:"} {
		if idx := strings.Index(clean, token); idx != -1 {
			return clean[idx+len(token):]
		}
	}

	return clean
}

// SanitizeRepoSlug converts a repository slug into a valid safe filesystem name.
func SanitizeRepoSlug(repo string) string {
	lower := strings.ToLower(cleanRepoURLPrefix(repo))
	slug := lazyregex.SlugSanitizeRegex.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "pipeline-default"
	}

	return slug
}

// PipelineDbDir returns the dedicated directory where pipeline split DBs live.
func PipelineDbDir() string {
	binDir := filepath.Dir(store.BinaryDataDir())
	altDir := filepath.Join(binDir, "pipeline")
	if isDirExisting(altDir) {
		return filepath.ToSlash(altDir)
	}

	dir := filepath.Join(store.BinaryDataDir(), "pipeline")
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(dir)
}

// PipelineDBDir is an alias to PipelineDbDir.
var PipelineDBDir = PipelineDbDir

// RepoPipelineDir returns the dedicated repository directory inside pipeline data root.
func RepoPipelineDir(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)
	repoDir := filepath.Join(PipelineDbDir(), slug)
	_ = os.MkdirAll(repoDir, 0755)

	return filepath.ToSlash(repoDir)
}

// RepoScopedPipelineDbDir returns the dedicated repository pipeline directory.
func RepoScopedPipelineDbDir(repoSlug string) string {
	return RepoPipelineDir(repoSlug)
}

// ResolvePipelineDbPath resolves the CLI-anchored SQLite database path for a repository.
func ResolvePipelineDbPath(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)
	repoDir := RepoPipelineDir(slug)
	target := filepath.Join(repoDir, "pipeline.db")
	if isFileExisting(target) {
		return filepath.ToSlash(target)
	}

	if sqlDb := filepath.Join(repoDir, "sql.db"); isFileExisting(sqlDb) {
		return filepath.ToSlash(sqlDb)
	}

	return filepath.ToSlash(migrateOrFallbackPipelineDb(PipelineDbDir(), slug, target))
}

// PipelineDbPath returns the full SQLite database file path for a repository.
func PipelineDbPath(repoSlug string) string {
	return ResolvePipelineDbPath(repoSlug)
}

// PipelineDBPath is an alias to PipelineDbPath.
var PipelineDBPath = PipelineDbPath
