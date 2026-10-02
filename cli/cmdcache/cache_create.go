package cmdcache

import (
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

// CreateCache indexes target paths into the Split-DB cache.
func CreateCache(args []string) *appfault.AppError {
	opts := parseCreateOptions(args)
	repoRoot := findRepoRoot()
	queueId, tasksDb := enqueueCacheTask("create", strings.Join(opts.Targets, ","))
	defer closeTasksDB(tasksDb)

	return executeCreateCache(opts, repoRoot, queueId, tasksDb)
}

func executeCreateCache(opts CacheCreateOptions, repoRoot, queueId string, tasksDb *store.TasksSplitDB) *appfault.AppError {
	rootDb, err := store.OpenRootCacheDB(repoRoot)
	hasErr := err != nil
	if hasErr {
		logAndFailCacheTask(tasksDb, queueId, err, "open_root_cache_db")
		return err
	}
	defer rootDb.Close()

	return runIndexAndFinalize(rootDb, opts, repoRoot, queueId, tasksDb)
}

func runIndexAndFinalize(rootDb *sql.DB, opts CacheCreateOptions, repoRoot, queueId string, tasksDb *store.TasksSplitDB) *appfault.AppError {
	indexed, totalBytes, procErr := processTargets(rootDb, opts, repoRoot)
	hasProcErr := procErr != nil
	if hasProcErr {
		logAndFailCacheTask(tasksDb, queueId, procErr, "process_targets")
		return procErr
	}
	updateRepoStats(rootDb, repoRoot, indexed, totalBytes)
	completeCacheTask(tasksDb, queueId)
	fmt.Printf("\n%s✓ Split-DB cache created: indexed %d file(s) (<= 200KB).%s\n",
		constants.ColorGreen, indexed, constants.ColorReset)
	return nil
}

// ReconcileCache checks mtime vs filesystem and updates outdated cache entries.
func ReconcileCache(repoRoot string) *appfault.AppError {
	queueId, tasksDb := enqueueCacheTask("reconcile", repoRoot)
	defer closeTasksDB(tasksDb)

	rootDb, err := store.OpenRootCacheDB(repoRoot)
	hasErr := err != nil
	if hasErr {
		logAndFailCacheTask(tasksDb, queueId, err, "open_root_cache_db")
		return err
	}
	defer rootDb.Close()

	return executeReconcile(rootDb, repoRoot, queueId, tasksDb)
}

func executeReconcile(rootDb *sql.DB, repoRoot, queueId string, tasksDb *store.TasksSplitDB) *appfault.AppError {
	updated, removed := reconcileExistingFiles(rootDb, repoRoot)
	opts := CacheCreateOptions{Targets: []string{"."}, IsKeep: false}
	indexed, totalBytes, _ := processTargets(rootDb, opts, repoRoot)
	updateRepoStats(rootDb, repoRoot, indexed, totalBytes)
	completeCacheTask(tasksDb, queueId)

	fmt.Printf("\n%s✓ Split-DB cache reconciled: updated %d, removed %d, total %d file(s).%s\n",
		constants.ColorGreen, updated, removed, indexed, constants.ColorReset)
	return nil
}

func enqueueCacheTask(action, target string) (string, *store.TasksSplitDB) {
	queueId := fmt.Sprintf("cache-%s-%d", action, time.Now().UnixNano())
	tasksDb, err := store.OpenTasksRootSplitDB()
	hasErr := err != nil
	if hasErr {
		return "", nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	stmt := "INSERT INTO TaskQueue (QueueId, Section, Action, Target, Status, CreatedAt, UpdatedAt) VALUES (?, 'cache', ?, ?, 'pending', ?, ?)"
	store.ExecWrapper(tasksDb.Conn(), stmt, queueId, action, target, now, now)
	return queueId, tasksDb
}

func completeCacheTask(tasksDb *store.TasksSplitDB, queueId string) {
	hasTasksDb := tasksDb != nil
	if !hasTasksDb {
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	store.ExecWrapper(tasksDb.Conn(), "UPDATE TaskQueue SET Status = 'completed', UpdatedAt = ? WHERE QueueId = ?", now, queueId)
}

func logAndFailCacheTask(tasksDb *store.TasksSplitDB, queueId string, err error, action string) {
	hasTasksDb := tasksDb != nil
	if hasTasksDb {
		now := time.Now().UTC().Format(time.RFC3339)
		store.ExecWrapper(tasksDb.Conn(), "UPDATE TaskQueue SET Status = 'failed', UpdatedAt = ? WHERE QueueId = ?", now, queueId)
	}
	msg := ""
	hasErr := err != nil
	if hasErr {
		msg = err.Error()
	}
	store.LogInternalError("cache", "E1030", msg, action, "cache_create.go")
}

func closeTasksDB(tasksDb *store.TasksSplitDB) {
	hasTasksDb := tasksDb != nil
	if hasTasksDb {
		_ = tasksDb.Close()
	}
}

func processTargets(rootDb *sql.DB, opts CacheCreateOptions, repoRoot string) (int, int64, *appfault.AppError) {
	totalIndexed := 0
	var totalBytes int64
	folderStats := make(map[string]FolderStats)

	for _, target := range opts.Targets {
		processSingleTarget(rootDb, target, repoRoot, opts.IsKeep, &totalIndexed, &totalBytes, folderStats)
	}
	updateFolderTreeStats(rootDb, folderStats)
	return totalIndexed, totalBytes, nil
}

func updateFolderTreeStats(rootDb *sql.DB, folderStats map[string]FolderStats) {
	for folderSlug, stat := range folderStats {
		_ = store.UpdateFolderTree(rootDb, folderSlug, folderSlug, stat.Count, stat.Bytes)
	}
}

func processSingleTarget(rootDb *sql.DB, target, repoRoot string, isKeep bool, totalIndexed *int, totalBytes *int64, folderStats map[string]FolderStats) {
	absTarget, err := filepath.Abs(target)
	hasErr := err != nil
	if hasErr {
		return
	}
	info, statErr := os.Stat(absTarget)
	hasStatErr := statErr != nil
	if hasStatErr {
		return
	}
	dispatchTargetIndex(rootDb, absTarget, repoRoot, isKeep, info.IsDir(), totalIndexed, totalBytes, folderStats)
}

func dispatchTargetIndex(rootDb *sql.DB, absTarget, repoRoot string, isKeep, isDir bool, totalIndexed *int, totalBytes *int64, folderStats map[string]FolderStats) {
	if isDir {
		walkAndIndexDir(rootDb, absTarget, repoRoot, isKeep, totalIndexed, totalBytes, folderStats)
		return
	}
	indexFile(rootDb, absTarget, repoRoot, isKeep, totalIndexed, totalBytes, folderStats)
}

func walkAndIndexDir(rootDb *sql.DB, targetDir, repoRoot string, isKeep bool, totalIndexed *int, totalBytes *int64, folderStats map[string]FolderStats) {
	_ = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, walkErr error) error {
		hasWalkErr := walkErr != nil
		if hasWalkErr {
			return nil
		}
		if d.IsDir() {
			return handleDirWalk(d.Name())
		}
		indexFile(rootDb, path, repoRoot, isKeep, totalIndexed, totalBytes, folderStats)
		return nil
	})
}

func handleDirWalk(name string) error {
	if isSkippedDir(name) {
		return filepath.SkipDir
	}
	return nil
}

func indexFile(rootDb *sql.DB, absPath, repoRoot string, isKeep bool, totalIndexed *int, totalBytes *int64, folderStats map[string]FolderStats) {
	info, err := os.Stat(absPath)
	hasErr := err != nil
	if hasErr {
		return
	}
	isEligible := isFileEligible(info, absPath, isKeep)
	if !isEligible {
		return
	}
	recordAndStoreFile(rootDb, absPath, repoRoot, isKeep, info, totalIndexed, totalBytes, folderStats)
}

func recordAndStoreFile(rootDb *sql.DB, absPath, repoRoot string, isKeep bool, info fs.FileInfo, totalIndexed *int, totalBytes *int64, folderStats map[string]FolderStats) {
	cleanRel := resolveCleanRelPath(repoRoot, absPath)
	slug := resolveFileSlug(cleanRel)
	rec := buildCacheFileRecord(cleanRel, absPath, slug, isKeep, info)
	if insertErr := store.InsertCacheFile(rootDb, rec); insertErr != nil {
		return
	}
	storeFileLines(slug, cleanRel, absPath, repoRoot)
	accumulateStats(slug, info.Size(), totalIndexed, totalBytes, folderStats)
}

func buildCacheFileRecord(cleanRel, absPath, slug string, isKeep bool, info fs.FileInfo) store.CacheFileRecord {
	return store.CacheFileRecord{
		RelativePath: cleanRel,
		AbsolutePath: absPath,
		FileSize:     info.Size(),
		ModifiedTime: info.ModTime().Unix(),
		FolderSlug:   slug,
		IsKeep:       isKeep,
	}
}

func resolveCleanRelPath(repoRoot, absPath string) string {
	relPath, _ := filepath.Rel(repoRoot, absPath)
	return filepath.ToSlash(relPath)
}

func accumulateStats(slug string, size int64, totalIndexed *int, totalBytes *int64, folderStats map[string]FolderStats) {
	*totalIndexed++
	*totalBytes += size
	cur := folderStats[slug]
	cur.Count++
	cur.Bytes += size
	folderStats[slug] = cur
}

func storeFileLines(slug, relPath, absPath, repoRoot string) {
	contentBytes, err := os.ReadFile(absPath)
	hasErr := err != nil
	if hasErr {
		return
	}
	slugDb, openErr := store.OpenSlugCacheDB(slug, repoRoot)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return
	}
	defer slugDb.Close()

	lines := strings.Split(string(contentBytes), "\n")
	_ = store.InsertCachedLines(slugDb, relPath, lines)
}

func reconcileExistingFiles(rootDb *sql.DB, repoRoot string) (int, int) {
	files, err := store.ListCacheFiles(rootDb)
	hasErr := err != nil
	if hasErr {
		return 0, 0
	}
	return processReconcileFiles(rootDb, files, repoRoot)
}

func processReconcileFiles(rootDb *sql.DB, files []store.CacheFileRecord, repoRoot string) (int, int) {
	updated, removed := 0, 0
	for _, rec := range files {
		isUpdated, isRemoved := reconcileSingleFile(rootDb, rec, repoRoot)
		if isUpdated {
			updated++
		}
		if isRemoved {
			removed++
		}
	}
	return updated, removed
}

func reconcileSingleFile(rootDb *sql.DB, rec store.CacheFileRecord, repoRoot string) (bool, bool) {
	info, statErr := os.Stat(rec.AbsolutePath)
	hasStatErr := statErr != nil
	if hasStatErr {
		return false, handleMissingFile(rootDb, rec, repoRoot)
	}
	isModified := info.ModTime().Unix() != rec.ModifiedTime
	if isModified {
		reIndexChangedFile(rootDb, rec, info, repoRoot)
		return true, false
	}
	return false, false
}

func handleMissingFile(rootDb *sql.DB, rec store.CacheFileRecord, repoRoot string) bool {
	if rec.IsKeep {
		return false
	}
	_ = store.DeleteCacheFile(rootDb, rec.RelativePath)
	removeLinesFromSlug(rec.FolderSlug, rec.RelativePath, repoRoot)
	return true
}

func removeLinesFromSlug(slug, relPath, repoRoot string) {
	slugDb, err := store.OpenSlugCacheDB(slug, repoRoot)
	hasErr := err != nil
	if !hasErr {
		defer slugDb.Close()
		_ = store.DeleteCachedLines(slugDb, relPath)
	}
}

func reIndexChangedFile(rootDb *sql.DB, rec store.CacheFileRecord, info fs.FileInfo, repoRoot string) {
	rec.ModifiedTime = info.ModTime().Unix()
	rec.FileSize = info.Size()
	_ = store.InsertCacheFile(rootDb, rec)
	storeFileLines(rec.FolderSlug, rec.RelativePath, rec.AbsolutePath, repoRoot)
}

func updateRepoStats(rootDb *sql.DB, repoRoot string, totalFiles int, totalBytes int64) {
	now := time.Now().UTC().Format(time.RFC3339)
	slug := store.ResolveRepoSlug(repoRoot)
	_ = store.SetRepoMetadata(rootDb, "repo_slug", slug)
	_ = store.SetRepoMetadata(rootDb, "repo_path", repoRoot)
	_ = store.SetRepoMetadata(rootDb, "last_indexed_at", now)
	_ = store.SetRepoMetadata(rootDb, "total_files", fmt.Sprintf("%d", totalFiles))
	_ = store.SetRepoMetadata(rootDb, "total_bytes", fmt.Sprintf("%d", totalBytes))
}
