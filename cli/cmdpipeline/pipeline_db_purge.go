package cmdpipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func purgeAllGlobalDirs(dirs []string) (int, int64) {
	seen := make(map[string]bool)
	var totalFiles int
	var totalBytes int64
	for _, d := range dirs {
		clean := filepath.Clean(d)
		if len(d) > 0 && !seen[clean] && isDirExisting(clean) {
			seen[clean] = true
			f, b := purgeSubtreeAndFiles(clean)
			totalFiles += f
			totalBytes += b
		}
	}

	return totalFiles, totalBytes
}

func purgeSubtreeAndFiles(dir string) (int, int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}

	var count int
	var bytes int64
	for _, e := range entries {
		c, b := purgeEntryRecursive(dir, e)
		count += c
		bytes += b
	}

	return count, bytes
}

func purgeEntryRecursive(parent string, e os.DirEntry) (int, int64) {
	target := filepath.Join(parent, e.Name())
	if e.IsDir() {
		c, b := purgeRepoPipelineFolder(target)
		_ = os.Remove(target)

		return c, b
	}
	if isPurgeablePipelineFile(e.Name()) {
		return purgeSinglePipelineFile(target)
	}

	return 0, 0
}

func purgeAllCandidateDirs(dirs []string, activeDir string) (int, int64) {
	var totalFiles int
	var totalBytes int64
	for _, d := range dirs {
		f, b := purgeRepoPipelineFolder(d)
		totalFiles += f
		totalBytes += b
		if d != activeDir {
			_ = os.Remove(d)
		}
	}

	return totalFiles, totalBytes
}

func reinitFreshPipelineDb(repo string) {
	freshDb, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err == nil {
		_ = freshDb.Close()
	}
}

func collectRepoPipelineDirs(repo string) []string {
	seen := make(map[string]bool)
	var dirs []string
	slug := pipelinedb.SanitizeRepoSlug(repo)
	underscoreSlug := strings.ReplaceAll(repo, "/", "_")
	candidates := buildDirCandidates(repo, slug, underscoreSlug)

	for _, d := range candidates {
		clean := filepath.Clean(d)
		if isDirExisting(clean) && !seen[clean] {
			seen[clean] = true
			dirs = append(dirs, clean)
		}
	}

	return dirs
}

func resolveGlobalAppDataPipelineDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if len(localAppData) > 0 {
		return filepath.Join(localAppData, "gitmap-cli", "data", "pipeline")
	}

	home, err := os.UserHomeDir()
	if err == nil && len(home) > 0 {
		return filepath.Join(home, "AppData", "Local", "gitmap-cli", "data", "pipeline")
	}

	return ""
}

func buildDirCandidates(repo, slug, underscoreSlug string) []string {
	appData := store.BinaryDataDir()
	candidates := []string{
		pipelinedb.RepoPipelineDir(repo),
		filepath.Join(appData, "pipeline", slug),
		filepath.Join(appData, "pipeline", underscoreSlug),
	}
	candidates = appendGlobalAppDataCandidates(candidates, slug, underscoreSlug)

	return appendRepoRootCandidates(candidates, slug, underscoreSlug)
}

func appendGlobalAppDataCandidates(list []string, slug, underscoreSlug string) []string {
	globalDir := resolveGlobalAppDataPipelineDir()
	if len(globalDir) > 0 {
		list = append(list,
			filepath.Join(globalDir, slug),
			filepath.Join(globalDir, underscoreSlug),
		)
	}

	return list
}

func appendRepoRootCandidates(list []string, slug, underscoreSlug string) []string {
	repoRoot := resolveRepoRootDir()
	if len(repoRoot) > 0 && repoRoot != "." {
		list = append(list,
			filepath.Join(repoRoot, ".gitmap", "data", "pipeline", slug),
			filepath.Join(repoRoot, ".gitmap", "data", "pipeline", underscoreSlug),
		)
	}

	return list
}

func purgeParentPipelineOrphans(repo string) (int, int64) {
	slug := pipelinedb.SanitizeRepoSlug(repo)
	underscoreSlug := strings.ReplaceAll(repo, "/", "_")
	dirs := []string{
		filepath.Join(store.BinaryDataDir(), "pipeline"),
		resolveGlobalAppDataPipelineDir(),
	}

	return purgeOrphanFilesAcrossDirs(dirs, slug, underscoreSlug)
}

func purgeOrphanFilesAcrossDirs(dirs []string, slug, underscoreSlug string) (int, int64) {
	var totalFiles int
	var totalBytes int64
	for _, d := range dirs {
		if len(d) == 0 || !isDirExisting(d) {
			continue
		}
		f, b := purgeOrphansInDir(d, slug, underscoreSlug)
		totalFiles += f
		totalBytes += b
	}

	return totalFiles, totalBytes
}

func purgeOrphansInDir(dir, slug, underscoreSlug string) (int, int64) {
	targets := []string{
		filepath.Join(dir, fmt.Sprintf("pipeline_%s.db", slug)),
		filepath.Join(dir, fmt.Sprintf("%s.db", slug)),
		filepath.Join(dir, fmt.Sprintf("pipeline_%s.db", underscoreSlug)),
		filepath.Join(dir, fmt.Sprintf("%s.db", underscoreSlug)),
	}
	var count int
	var bytes int64
	for _, t := range targets {
		c, b := purgeSinglePipelineFile(t)
		count += c
		bytes += b
	}

	return count, bytes
}

func isDirExisting(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func purgeRepoPipelineFolder(dir string) (int, int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}

	return iteratePurgeEntries(dir, entries)
}

func iteratePurgeEntries(dir string, entries []os.DirEntry) (int, int64) {
	var count int
	var totalBytes int64
	for _, e := range entries {
		c, b := purgeEntryIfEligible(dir, e)
		count += c
		totalBytes += b
	}

	return count, totalBytes
}

func purgeEntryIfEligible(dir string, e os.DirEntry) (int, int64) {
	target := filepath.Join(dir, e.Name())
	if e.IsDir() {
		c, b := purgeRepoPipelineFolder(target)
		_ = os.Remove(target)

		return c, b
	}

	return purgeSinglePipelineFile(target)
}

func purgeSinglePipelineFile(filePath string) (int, int64) {
	var size int64
	if fi, sErr := os.Stat(filePath); sErr == nil {
		size = fi.Size()
	}
	if rErr := os.Remove(filePath); rErr == nil {
		return 1, size
	}

	return 0, 0
}

func isPurgeablePipelineFile(name string) bool {
	lower := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lower, ".log"):
		return true
	case strings.HasSuffix(lower, ".json"):
		return true
	case strings.HasSuffix(lower, ".db-journal") || strings.HasSuffix(lower, ".db-wal") || strings.HasSuffix(lower, ".db-shm"):
		return true
	case lower == "sql.db" || strings.HasSuffix(lower, ".db"):
		return true
	default:
		return false
	}
}
