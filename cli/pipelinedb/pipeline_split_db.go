package pipelinedb

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

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
	slug := store.SanitizeSlug(repo)
	if slug == "default" {
		return "pipeline-default"
	}

	return slug
}

// PipelineDbDir returns the dedicated directory where pipeline split DBs live.
func PipelineDbDir() string {
	return filepath.Dir(store.ResolveSplitDbDir(store.SectionPipeline, "default", ""))
}

// PipelineDBDir is an alias to PipelineDbDir.
var PipelineDBDir = PipelineDbDir

// RepoPipelineDir returns the dedicated repository directory inside pipeline data root.
func RepoPipelineDir(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)
	return store.ResolveSplitDbDir(store.SectionPipeline, slug, "")
}

// RepoScopedPipelineDbDir returns the dedicated repository pipeline directory.
func RepoScopedPipelineDbDir(repoSlug string) string {
	return RepoPipelineDir(repoSlug)
}

// ResolvePipelineDbPath resolves the CLI-anchored SQLite database path for a repository.
func ResolvePipelineDbPath(repoSlug string) string {
	slug := SanitizeRepoSlug(repoSlug)
	primary := store.ResolveSplitDbPath(store.SectionPipeline, slug, "")
	candidate := findPipelineDbCandidate(slug)
	if candidate == "" {
		return primary
	}

	primarySize := getFileSize(primary)
	candidateSize := getFileSize(candidate)

	if candidateSize > primarySize {
		return candidate
	}

	return primary
}

func findPipelineDbCandidate(slug string) string {
	baseDir := PipelineDbDir()
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return ""
	}

	var bestPath string
	var bestSize int64
	targetSuffix := "-" + strings.ToLower(slug)
	shortSlug := strings.TrimPrefix(strings.ToLower(slug), "alimtvnetwork-")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		nameLower := strings.ToLower(entry.Name())
		isMatch := strings.HasSuffix(nameLower, targetSuffix) ||
			(shortSlug != strings.ToLower(slug) && nameLower == shortSlug)
		if !isMatch {
			continue
		}
		dbFile := filepath.Join(baseDir, entry.Name(), store.DbFileName)
		sz := getFileSize(dbFile)
		if sz > bestSize {
			bestPath = filepath.ToSlash(dbFile)
			bestSize = sz
		}
	}

	return bestPath
}

// PipelineDbPath returns the full SQLite database file path for a repository.
func PipelineDbPath(repoSlug string) string {
	return ResolvePipelineDbPath(repoSlug)
}

// PipelineDBPath is an alias to PipelineDbPath.
var PipelineDBPath = PipelineDbPath
