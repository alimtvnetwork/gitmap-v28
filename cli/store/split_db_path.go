package store

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

const (
	DbFileName = "sql.db"

	SectionAutomation   = "automation"
	SectionPipeline     = "pipeline"
	SectionInstallation = "installation"
	SectionStartup      = "startup"
	SectionSites        = "sites"
	SectionSchedule     = "schedule"
)

// SanitizeSlug converts an arbitrary name or repo identifier into a clean filesystem slug.
func SanitizeSlug(raw string) string {
	clean := cleanRepoPrefix(raw)
	lower := strings.ToLower(clean)
	slug := lazyregex.SlugSanitizeRegex.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "default"
	}

	return slug
}

func cleanRepoPrefix(raw string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(raw), ".git")
	for _, token := range []string{"github.com/", "github.com:", "gitlab.com/", "gitlab.com:"} {
		if idx := strings.Index(clean, token); idx != -1 {
			return clean[idx+len(token):]
		}
	}

	return clean
}

// ResolveSplitDbDir resolves the directory containing the split SQLite database.
func ResolveSplitDbDir(section, slug, repoRoot string) string {
	cleanSec := SanitizeSlug(section)
	cleanSlug := SanitizeSlug(slug)
	baseDir := resolveBaseDataDir(repoRoot)
	dir := filepath.Join(baseDir, cleanSec, cleanSlug)
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(dir)
}

func resolveBaseDataDir(repoRoot string) string {
	root := resolveTargetRoot(repoRoot)
	if !isDirExisting(filepath.Join(root, ".gitmap")) {
		return BinaryDataDir()
	}
	if root == "." {
		return filepath.Join(".gitmap", "data")
	}

	return filepath.Join(root, ".gitmap", "data")
}

func resolveTargetRoot(repoRoot string) string {
	if repoRoot != "" {
		return repoRoot
	}

	return findWorkingRepoRoot()
}

func findWorkingRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if isDirExisting(filepath.Join(dir, ".gitmap")) || isDirExisting(filepath.Join(dir, ".git")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "."
}

// ResolveSplitDbPath resolves the canonical split SQLite database file path.
func ResolveSplitDbPath(section, slug, repoRoot string) string {
	cleanSec := SanitizeSlug(section)
	cleanSlug := SanitizeSlug(slug)
	dir := ResolveSplitDbDir(cleanSec, cleanSlug, repoRoot)
	targetPath := filepath.ToSlash(filepath.Join(dir, DbFileName))
	if isFileExisting(targetPath) {
		return targetPath
	}

	migrateLegacySplitDb(cleanSec, cleanSlug, repoRoot, targetPath)

	return targetPath
}

// ResolveTasksRootDbPath resolves the tasks root SQLite database file path (.gitmap/data/tasks/sql.db).
func ResolveTasksRootDbPath(repoRoot string) string {
	baseDir := resolveBaseDataDir(repoRoot)
	dir := filepath.Join(baseDir, "tasks")
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(filepath.Join(dir, DbFileName))
}

// ResolveSectionTasksDbPath resolves a section-scoped tasks database file path (.gitmap/data/<section>/<section>-tasks.db).
func ResolveSectionTasksDbPath(section, repoRoot string) string {
	cleanSec := SanitizeSlug(section)
	baseDir := resolveBaseDataDir(repoRoot)
	dir := filepath.Join(baseDir, cleanSec)
	_ = os.MkdirAll(dir, 0755)

	fileName := cleanSec + "-tasks.db"
	return filepath.ToSlash(filepath.Join(dir, fileName))
}

// ResolveAiInstructionDbPath resolves the AI instruction SQLite database file path (.gitmap/data/ai-instruction/sql.db).
func ResolveAiInstructionDbPath(repoRoot string) string {
	baseDir := resolveBaseDataDir(repoRoot)
	dir := filepath.Join(baseDir, "ai-instruction")
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(filepath.Join(dir, DbFileName))
}

// ResolveSearchDbPath resolves the search SQLite database file path (.gitmap/data/search/sql.db).
func ResolveSearchDbPath(repoRoot string) string {
	baseDir := resolveBaseDataDir(repoRoot)
	dir := filepath.Join(baseDir, "search")
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(filepath.Join(dir, DbFileName))
}

func migrateLegacySplitDb(section, slug, repoRoot, targetPath string) {
	candidates := collectLegacyCandidates(section, slug, repoRoot)
	for _, cand := range candidates {
		if !isFileExisting(cand) {
			continue
		}
		if tryMigrateFile(cand, targetPath) {
			return
		}
	}
}

func tryMigrateFile(source, target string) bool {
	if err := os.Rename(source, target); err == nil {
		return true
	}

	return copyAndRemoveFile(source, target)
}

func copyAndRemoveFile(source, target string) bool {
	srcFile, err := os.Open(source)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	dstFile, err := os.Create(target)
	if err != nil {
		return false
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return false
	}
	_ = srcFile.Close()
	_ = os.Remove(source)

	return true
}

func collectLegacyCandidates(section, slug, repoRoot string) []string {
	var candidates []string
	root := resolveTargetRoot(repoRoot)
	candidates = append(candidates, collectRootLegacyCandidates(root, section, slug)...)
	candidates = append(candidates, collectBinaryLegacyCandidates(section, slug)...)

	return candidates
}

func collectRootLegacyCandidates(root, section, slug string) []string {
	gitmapData := filepath.Join(root, ".gitmap", "data")

	return []string{
		filepath.Join(gitmapData, slug, section, DbFileName),
		filepath.Join(gitmapData, slug, section+".db"),
		filepath.Join(gitmapData, section, slug+".db"),
		filepath.Join(gitmapData, section, DbFileName),
		filepath.Join(gitmapData, section+".db"),
		filepath.Join(root, ".gitmap", section, slug, "pipeline.db"),
		filepath.Join(root, ".gitmap", section, slug, DbFileName),
	}
}

func collectBinaryLegacyCandidates(section, slug string) []string {
	binDir := BinaryDataDir()

	return []string{
		filepath.Join(binDir, section, slug, "pipeline.db"),
		filepath.Join(binDir, section, slug+".db"),
		filepath.Join(binDir, "schedules", slug+".db"),
		filepath.Join(binDir, "pipeline", "pipeline_"+slug+".db"),
		filepath.Join(binDir, "pipeline_db", "pipeline_"+slug+".db"),
		filepath.Join(binDir, section+".db"),
		filepath.Join(binDir, section, section+".db"),
	}
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
