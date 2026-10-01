package cmdcache

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const maxCacheFileSize = 200 * 1024 // 200 KB
const maxJsonFileSize = 150 * 1024  // 150 KB

// CreateCache indexes target paths into the Split-DB cache.
func CreateCache(args []string) *appfault.AppError {
	opts := parseCreateOptions(args)
	repoRoot := findRepoRoot()
	queueId, tasksDB := enqueueCacheTask("create", strings.Join(opts.Targets, ","))
	defer closeTasksDB(tasksDB)

	rootDB, err := store.OpenRootCacheDB(repoRoot)
	if err != nil {
		logAndFailCacheTask(tasksDB, queueId, err, "open_root_cache_db")
		return err
	}
	defer rootDB.Close()

	indexed, totalBytes, procErr := processTargets(rootDB, opts, repoRoot)
	if procErr != nil {
		logAndFailCacheTask(tasksDB, queueId, procErr, "process_targets")
		return procErr
	}
	updateRepoStats(rootDB, repoRoot, indexed, totalBytes)
	completeCacheTask(tasksDB, queueId)
	fmt.Printf("\n%s✓ Split-DB cache created: indexed %d file(s) (<= 200KB).%s\n",
		constants.ColorGreen, indexed, constants.ColorReset)
	return nil
}

// ReconcileCache checks mtime vs filesystem and updates outdated cache entries.
func ReconcileCache(repoRoot string) *appfault.AppError {
	queueId, tasksDB := enqueueCacheTask("reconcile", repoRoot)
	defer closeTasksDB(tasksDB)

	rootDB, err := store.OpenRootCacheDB(repoRoot)
	if err != nil {
		logAndFailCacheTask(tasksDB, queueId, err, "open_root_cache_db")
		return err
	}
	defer rootDB.Close()

	updated, removed := reconcileExistingFiles(rootDB, repoRoot)
	opts := CacheCreateOptions{Targets: []string{"."}, IsKeep: false}
	indexed, totalBytes, _ := processTargets(rootDB, opts, repoRoot)
	updateRepoStats(rootDB, repoRoot, indexed, totalBytes)
	completeCacheTask(tasksDB, queueId)

	fmt.Printf("\n%s✓ Split-DB cache reconciled: updated %d, removed %d, total %d file(s).%s\n",
		constants.ColorGreen, updated, removed, indexed, constants.ColorReset)
	return nil
}

func parseCreateOptions(args []string) CacheCreateOptions {
	opts := CacheCreateOptions{IsKeep: false}
	for _, a := range args {
		trimmed := strings.TrimSpace(a)
		if trimmed == "--keep" || trimmed == "-k" {
			opts.IsKeep = true
		} else if !strings.HasPrefix(trimmed, "-") {
			parts := extractTargetParts(trimmed)
			opts.Targets = append(opts.Targets, parts...)
		}
	}
	if len(opts.Targets) == 0 {
		opts.Targets = []string{"."}
	}
	return opts
}

func extractTargetParts(raw string) []string {
	parts := strings.Split(raw, ",")
	var cleaned []string
	for _, p := range parts {
		trimmed := strings.Trim(strings.TrimSpace(p), `"'`)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func enqueueCacheTask(action, target string) (string, *store.TasksSplitDB) {
	queueId := fmt.Sprintf("cache-%s-%d", action, time.Now().UnixNano())
	tasksDB, err := store.OpenTasksRootSplitDB()
	if err != nil {
		return "", nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	stmt := "INSERT INTO TaskQueue (QueueId, Section, Action, Target, Status, CreatedAt, UpdatedAt) VALUES (?, 'cache', ?, ?, 'pending', ?, ?)"
	store.ExecWrapper(tasksDB.Conn(), stmt, queueId, action, target, now, now)
	return queueId, tasksDB
}

func completeCacheTask(tasksDB *store.TasksSplitDB, queueId string) {
	if tasksDB == nil {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(tasksDB.Conn(), "UPDATE TaskQueue SET Status = 'completed', UpdatedAt = ? WHERE QueueId = ?", now, queueId)
}

func logAndFailCacheTask(tasksDB *store.TasksSplitDB, queueId string, err error, action string) {
	if tasksDB != nil {
		now := time.Now().UTC().Format(time.RFC3339)
		store.ExecWrapper(tasksDB.Conn(), "UPDATE TaskQueue SET Status = 'failed', UpdatedAt = ? WHERE QueueId = ?", now, queueId)
	}
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	store.LogInternalError("cache", "E1030", msg, action, "cache_create.go")
}

func closeTasksDB(tasksDB *store.TasksSplitDB) {
	if tasksDB != nil {
		_ = tasksDB.Close()
	}
}

func processTargets(rootDB *sql.DB, opts CacheCreateOptions, repoRoot string) (int, int64, *appfault.AppError) {
	totalIndexed := 0
	var totalBytes int64
	folderStats := make(map[string]struct {
		count int
		bytes int64
	})

	for _, target := range opts.Targets {
		processSingleTarget(rootDB, target, repoRoot, opts.IsKeep, &totalIndexed, &totalBytes, folderStats)
	}
	for folderSlug, stat := range folderStats {
		_ = store.UpdateFolderTree(rootDB, folderSlug, folderSlug, stat.count, stat.bytes)
	}
	return totalIndexed, totalBytes, nil
}

func processSingleTarget(rootDB *sql.DB, target, repoRoot string, isKeep bool, totalIndexed *int, totalBytes *int64, folderStats map[string]struct {
	count int
	bytes int64
}) {
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return
	}
	info, statErr := os.Stat(absTarget)
	if statErr != nil {
		return
	}
	if info.IsDir() {
		walkAndIndexDir(rootDB, absTarget, repoRoot, isKeep, totalIndexed, totalBytes, folderStats)
	} else {
		indexFile(rootDB, absTarget, repoRoot, isKeep, totalIndexed, totalBytes, folderStats)
	}
}

func walkAndIndexDir(rootDB *sql.DB, targetDir, repoRoot string, isKeep bool, totalIndexed *int, totalBytes *int64, folderStats map[string]struct {
	count int
	bytes int64
}) {
	_ = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			if isSkippedDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		indexFile(rootDB, path, repoRoot, isKeep, totalIndexed, totalBytes, folderStats)
		return nil
	})
}

func isSkippedDir(name string) bool {
	lower := strings.ToLower(name)
	return lower == ".git" || lower == "node_modules" || lower == ".vscode" ||
		lower == ".idea" || lower == ".gitmap" || lower == "vendor"
}

func isFileEligible(info fs.FileInfo, absPath string, isKeep bool) bool {
	if info.IsDir() {
		return false
	}
	if isKeep {
		return !isBinaryFile(absPath)
	}
	if info.Size() > maxCacheFileSize || isLargeJson(absPath, info.Size()) {
		return false
	}
	return !isBinaryFile(absPath)
}

func isLargeJson(absPath string, size int64) bool {
	if !strings.HasSuffix(strings.ToLower(absPath), ".json") {
		return false
	}
	return size > maxJsonFileSize
}

func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, readErr := f.Read(buf)
	if readErr != nil && n == 0 {
		return false
	}
	return bytes.Contains(buf[:n], []byte{0})
}

func indexFile(rootDB *sql.DB, absPath, repoRoot string, isKeep bool, totalIndexed *int, totalBytes *int64, folderStats map[string]struct {
	count int
	bytes int64
}) {
	info, err := os.Stat(absPath)
	if err != nil || !isFileEligible(info, absPath, isKeep) {
		return
	}
	relPath, _ := filepath.Rel(repoRoot, absPath)
	cleanRel := filepath.ToSlash(relPath)
	slug := resolveFileSlug(cleanRel)

	rec := store.CacheFileRecord{
		RelativePath: cleanRel,
		AbsolutePath: absPath,
		FileSize:     info.Size(),
		ModifiedTime: info.ModTime().Unix(),
		FolderSlug:   slug,
		IsKeep:       isKeep,
	}
	if insertErr := store.InsertCacheFile(rootDB, rec); insertErr != nil {
		return
	}
	storeFileLines(slug, cleanRel, absPath, repoRoot)
	*totalIndexed++
	*totalBytes += info.Size()
	cur := folderStats[slug]
	cur.count++
	cur.bytes += info.Size()
	folderStats[slug] = cur
}

func storeFileLines(slug, relPath, absPath, repoRoot string) {
	contentBytes, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	slugDB, openErr := store.OpenSlugCacheDB(slug, repoRoot)
	if openErr != nil {
		return
	}
	defer slugDB.Close()

	lines := strings.Split(string(contentBytes), "\n")
	_ = store.InsertCachedLines(slugDB, relPath, lines)
}

func reconcileExistingFiles(rootDB *sql.DB, repoRoot string) (int, int) {
	files, err := store.ListCacheFiles(rootDB)
	if err != nil {
		return 0, 0
	}
	updated, removed := 0, 0
	for _, rec := range files {
		info, statErr := os.Stat(rec.AbsolutePath)
		if statErr != nil {
			if !rec.IsKeep {
				_ = store.DeleteCacheFile(rootDB, rec.RelativePath)
				removeLinesFromSlug(rec.FolderSlug, rec.RelativePath, repoRoot)
				removed++
			}
			continue
		}
		if info.ModTime().Unix() != rec.ModifiedTime {
			reIndexChangedFile(rootDB, rec, info, repoRoot)
			updated++
		}
	}
	return updated, removed
}

func removeLinesFromSlug(slug, relPath, repoRoot string) {
	slugDB, err := store.OpenSlugCacheDB(slug, repoRoot)
	if err == nil {
		defer slugDB.Close()
		_ = store.DeleteCachedLines(slugDB, relPath)
	}
}

func reIndexChangedFile(rootDB *sql.DB, rec store.CacheFileRecord, info fs.FileInfo, repoRoot string) {
	rec.ModifiedTime = info.ModTime().Unix()
	rec.FileSize = info.Size()
	_ = store.InsertCacheFile(rootDB, rec)
	storeFileLines(rec.FolderSlug, rec.RelativePath, rec.AbsolutePath, repoRoot)
}

func updateRepoStats(rootDB *sql.DB, repoRoot string, totalFiles int, totalBytes int64) {
	now := time.Now().UTC().Format(time.RFC3339)
	slug := store.ResolveRepoSlug(repoRoot)
	_ = store.SetRepoMetadata(rootDB, "repo_slug", slug)
	_ = store.SetRepoMetadata(rootDB, "repo_path", repoRoot)
	_ = store.SetRepoMetadata(rootDB, "last_indexed_at", now)
	_ = store.SetRepoMetadata(rootDB, "total_files", fmt.Sprintf("%d", totalFiles))
	_ = store.SetRepoMetadata(rootDB, "total_bytes", fmt.Sprintf("%d", totalBytes))
}

func resolveFileSlug(relPath string) string {
	clean := filepath.ToSlash(relPath)
	clean = strings.TrimPrefix(clean, "./")
	parts := strings.Split(clean, "/")
	if len(parts) > 1 && parts[0] != "" {
		return store.SanitizeSlug(parts[0])
	}
	return "root"
}

func findRepoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}
