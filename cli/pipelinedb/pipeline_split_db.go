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
		return altDir
	}

	dir := filepath.Join(store.BinaryDataDir(), "pipeline")
	_ = os.MkdirAll(dir, 0755)

	return dir
}

// PipelineDBDir is an alias to PipelineDbDir.
var PipelineDBDir = PipelineDbDir

// RepoScopedPipelineDbDir returns the pipeline db directory co-located with the CLI installation.
func RepoScopedPipelineDbDir(_ string) string {
	return PipelineDbDir()
}

func isFileExisting(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func isDirExisting(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func migrateOrFallbackPipelineDb(dir, slug, targetPrefixed string) string {
	legacyDir := filepath.Join(store.BinaryDataDir(), "pipeline_db")
	legacyFile := filepath.Join(legacyDir, "pipeline_"+slug+".db")
	if isFileExisting(legacyFile) {
		if err := os.Rename(legacyFile, targetPrefixed); err == nil {
			return targetPrefixed
		}

		return legacyFile
	}

	return targetPrefixed
}

// ResolvePipelineDbPath resolves the CLI-anchored SQLite database path for a repository.
func ResolvePipelineDbPath(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)
	dir := PipelineDbDir()

	direct := filepath.Join(dir, slug+".db")
	if isFileExisting(direct) {
		return direct
	}

	prefixed := filepath.Join(dir, "pipeline_"+slug+".db")
	if isFileExisting(prefixed) {
		return prefixed
	}

	return migrateOrFallbackPipelineDb(dir, slug, prefixed)
}

// PipelineDbPath returns the full SQLite database file path for a repository.
func PipelineDbPath(repoSlug string) string {
	return ResolvePipelineDbPath(repoSlug)
}

// PipelineDBPath is an alias to PipelineDbPath.
var PipelineDBPath = PipelineDbPath
