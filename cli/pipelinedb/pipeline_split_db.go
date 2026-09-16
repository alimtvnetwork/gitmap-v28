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

// SanitizeRepoSlug converts a repository slug into a valid safe filesystem name.
func SanitizeRepoSlug(repo string) string {
	lower := strings.ToLower(strings.TrimSpace(repo))
	if strings.Contains(lower, "://") || strings.Contains(lower, "@") {
		canonical := gitutil.CanonicalRepoID(lower)
		if idx := strings.Index(canonical, "/"); idx >= 0 {
			lower = canonical[idx+1:]
		}
	}
	slug := lazyregex.SlugSanitizeRegex.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "pipeline-default"
	}

	return slug
}

// PipelineDbDir returns the dedicated directory where pipeline split DBs live.
func PipelineDbDir() string {
<<<<<<< HEAD
	binDir := filepath.Dir(store.BinaryDataDir())
	altDir := filepath.Join(binDir, "pipeline")
	if isDirExisting(altDir) {
		return altDir
	}

=======
>>>>>>> ee457fe2ea6ccd129001691694fc9d207018a5a6
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

<<<<<<< HEAD
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
=======
func isRepoMatchingSlug(repoRoot, repoSlug string) bool {
	remote, err := gitutil.RemoteURL(repoRoot)
	if err == nil && strings.Contains(strings.ToLower(remote), strings.ToLower(repoSlug)) {
		return true
	}
	base := filepath.Base(repoRoot)
	if strings.EqualFold(base, filepath.Base(repoSlug)) {
		return true
	}
	trimmed := stripVersionSuffix(filepath.Base(repoSlug))

	return strings.EqualFold(base, trimmed)
}

func findRepoRootInStore(repoSlug string) string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	records, err := db.FindBySlug(repoSlug)
	if err != nil || len(records) == 0 {
		return ""
	}

	for _, rec := range records {
		if isDirExisting(rec.AbsolutePath) {
			return rec.AbsolutePath
		}
	}

	return ""
}

func resolveTargetRepoRoot(repoSlug string) string {
	root, err := gitutil.RepoRoot(".")
	if err == nil && root != "" && isRepoMatchingSlug(root, repoSlug) {
		return root
	}
	if cand := findCandidateRepoRoot(repoSlug); cand != "" {
		return cand
	}
	if storeRoot := findRepoRootInStore(repoSlug); storeRoot != "" {
		return storeRoot
	}
	if err == nil && root != "" {
		return root
	}

	return ""
}

// ResolvePipelineDbPath resolves the dedicated CLI-scoped SQLite database path for a repository.
func ResolvePipelineDbPath(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)
	primary := filepath.Join(PipelineDbDir(), "pipeline_"+slug+".db")
	if isFileExisting(primary) {
		return primary
	}
	legacyOldDir := filepath.Join(store.BinaryDataDir(), "pipeline_db", "pipeline_"+slug+".db")
	if isFileExisting(legacyOldDir) {
		return legacyOldDir
	}

	return primary
>>>>>>> ee457fe2ea6ccd129001691694fc9d207018a5a6
}

// PipelineDbPath returns the full SQLite database file path for a repository.
func PipelineDbPath(repoSlug string) string {
	return ResolvePipelineDbPath(repoSlug)
}

// PipelineDBPath is an alias to PipelineDbPath.
var PipelineDBPath = PipelineDbPath
