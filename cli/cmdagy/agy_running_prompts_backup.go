// Package cmdagy — agy_running_prompts_backup.go manages prompt backup, restoration, discovery, and cleanup.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	reMarkdownMedia = regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	reRawUUID       = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

func getAllProjects() ([]AgyProject, error) {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return nil, apperror.WrapSimple(err, "projects dir path")
	}
	return loadAllAgyProjects(dirPath)
}

// CollectActiveAndQueuedPrompts discovers running and pending prompts in parallel across all non-restricted workspaces.
func CollectActiveAndQueuedPrompts(maxWords int, isFull bool) ([]store.RunningPromptRecord, error) {
	projects, _ := getAllProjects()
	eligible := filterNonRestrictedProjects(projects)
	return collectPromptsParallel(eligible, maxWords, isFull), nil
}

func filterNonRestrictedProjects(projects []AgyProject) []AgyProject {
	var out []AgyProject
	for _, p := range projects {
		path := strings.TrimSpace(p.GetPath())
		if len(path) > 0 && !IsRestrictedSystemOrHomeDir(path) {
			out = append(out, p)
		}
	}
	return out
}

func collectPromptsParallel(projects []AgyProject, maxWords int, isFull bool) []store.RunningPromptRecord {
	var wg sync.WaitGroup
	var mu sync.Mutex
	seen := make(map[string]bool)
	var records []store.RunningPromptRecord
	for _, proj := range projects {
		wg.Add(1)
		go collectSingleProjectWorker(proj, maxWords, isFull, &wg, &mu, seen, &records)
	}
	wg.Wait()
	return records
}

func collectSingleProjectWorker(p AgyProject, maxWords int, isFull bool, wg *sync.WaitGroup, mu *sync.Mutex, seen map[string]bool, records *[]store.RunningPromptRecord) {
	defer wg.Done()
	localSeen := make(map[string]bool)
	var local []store.RunningPromptRecord
	projName := resolveHumanProjectName(p.Name, p.GetPath())
	local = collectProjectPrompts(local, localSeen, projName, p.GetPath(), p.ID, maxWords, isFull)
	mergeProjectPromptRecords(mu, seen, records, local)
}

func mergeProjectPromptRecords(mu *sync.Mutex, seen map[string]bool, dest *[]store.RunningPromptRecord, local []store.RunningPromptRecord) {
	mu.Lock()
	defer mu.Unlock()
	for _, rec := range local {
		key := fmt.Sprintf("%s:%s", rec.ProjectPath, rec.Prompt)
		if !seen[key] {
			seen[key] = true
			*dest = append(*dest, rec)
		}
	}
}

func resolveHumanProjectName(name, projectPath string) string {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) > 0 && !reRawUUID.MatchString(trimmed) {
		return trimmed
	}
	return filepath.Base(filepath.Clean(projectPath))
}

func collectProjectPrompts(records []store.RunningPromptRecord, seen map[string]bool, name, path, id string, maxWords int, isFull bool) []store.RunningPromptRecord {
	if len(path) == 0 || IsRestrictedSystemOrHomeDir(path) {
		return records
	}
	records = appendQueuePrompts(records, seen, name, path, id, maxWords, isFull)
	records = appendActiveTextPrompt(records, seen, name, path, id, maxWords, isFull)
	return records
}

func appendQueuePrompts(records []store.RunningPromptRecord, seen map[string]bool, name, path, id string, maxWords int, isFull bool) []store.RunningPromptRecord {
	q, hasQueue := loadWorkspaceQueueFile(path)
	if !hasQueue {
		return records
	}
	records = processActiveQueueItem(records, seen, q.Active, name, path, id, maxWords, isFull)
	for _, item := range q.Queued {
		records = processSingleQueueEntry(records, seen, item, name, path, id, maxWords, isFull)
	}
	return records
}

func loadWorkspaceQueueFile(path string) (AgyPromptQueueFile, bool) {
	data, err := os.ReadFile(resolveQueueFilePath(path))
	if err != nil {
		return AgyPromptQueueFile{}, false
	}
	var q AgyPromptQueueFile
	return q, json.Unmarshal(data, &q) == nil
}

func resolveQueueFilePath(dir string) string {
	nested := filepath.Join(dir, ".ai-memory", "temp", "agy-prompt-queue.json")
	if _, err := os.Stat(nested); err == nil {
		return nested
	}
	return filepath.Join(dir, "agy-prompt-queue.json")
}

func processActiveQueueItem(records []store.RunningPromptRecord, seen map[string]bool, item *AgyPromptQueueEntry, name, path, id string, maxWords int, isFull bool) []store.RunningPromptRecord {
	if item == nil || len(strings.TrimSpace(item.Prompt)) == 0 {
		return records
	}
	return processSingleQueueEntry(records, seen, *item, name, path, id, maxWords, isFull)
}

func processSingleQueueEntry(records []store.RunningPromptRecord, seen map[string]bool, item AgyPromptQueueEntry, name, path, id string, maxWords int, isFull bool) []store.RunningPromptRecord {
	clean := strings.TrimSpace(item.Prompt)
	key := fmt.Sprintf("%s:%s", path, clean)
	if len(clean) == 0 || seen[key] {
		return records
	}
	seen[key] = true
	status := resolveQueueEntryStatus(item.Status)
	rec := buildRunningPromptRecord(name, path, id, clean, item.CreatedAt, status, maxWords, isFull)
	return append(records, rec)
}

func resolveQueueEntryStatus(status string) store.PromptStatusType {
	if status == "dispatched" || status == "running" {
		return store.PromptStatusRunning
	}
	return store.PromptStatusEnqueued
}

func buildRunningPromptRecord(name, path, id, text, createdAt string, status store.PromptStatusType, maxWords int, isFull bool) store.RunningPromptRecord {
	snippet, wordCount := TruncateWords(text, maxWords)
	humanName := resolveHumanProjectName(name, path)
	return store.RunningPromptRecord{
		ProjectName: humanName, ProjectPath: path, ProjectId: id,
		Prompt: text, Snippet: resolveSnippet(text, snippet, isFull),
		Status: status, WordCount: wordCount,
		MediaPaths: extractPromptMediaPaths(text, path),
		CreatedAt:  createdAt,
	}
}

func extractPromptMediaPaths(text, projectPath string) []string {
	seen := make(map[string]bool)
	var paths []string
	paths = appendMarkdownMediaPaths(paths, seen, text)
	paths = appendTokenMediaPaths(paths, seen, text)
	paths = appendProjectScreenshotPaths(paths, seen, projectPath)
	return paths
}

func appendMarkdownMediaPaths(paths []string, seen map[string]bool, text string) []string {
	matches := reMarkdownMedia.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		if len(m) > 1 {
			paths = addUniqueMediaPath(paths, seen, m[1])
		}
	}
	return paths
}

func appendTokenMediaPaths(paths []string, seen map[string]bool, text string) []string {
	for _, tok := range strings.Fields(text) {
		clean := strings.Trim(tok, "\"'`()[]{},;:")
		if isSupportedMediaFile(clean) {
			paths = addUniqueMediaPath(paths, seen, clean)
		}
	}
	return paths
}

func appendProjectScreenshotPaths(paths []string, seen map[string]bool, projectPath string) []string {
	if len(strings.TrimSpace(projectPath)) == 0 {
		return paths
	}
	screenshotsDir := filepath.Join(projectPath, "assets", "screenshots")
	entries, err := os.ReadDir(screenshotsDir)
	if err != nil {
		return paths
	}
	for _, e := range entries {
		paths = appendScreenshotEntry(paths, seen, screenshotsDir, e)
	}
	return paths
}

func appendScreenshotEntry(paths []string, seen map[string]bool, dir string, entry os.DirEntry) []string {
	if entry.IsDir() || !isSupportedMediaFile(entry.Name()) {
		return paths
	}
	return addUniqueMediaPath(paths, seen, filepath.Join(dir, entry.Name()))
}

func isSupportedMediaFile(candidate string) bool {
	low := strings.ToLower(strings.TrimSpace(candidate))
	return strings.HasSuffix(low, ".png") || strings.HasSuffix(low, ".jpg") ||
		strings.HasSuffix(low, ".jpeg") || strings.HasSuffix(low, ".webp")
}

func addUniqueMediaPath(paths []string, seen map[string]bool, candidate string) []string {
	clean := strings.TrimSpace(candidate)
	if len(clean) == 0 || seen[clean] {
		return paths
	}
	seen[clean] = true
	return append(paths, clean)
}

func resolveSnippet(prompt, snippet string, isFull bool) string {
	if isFull {
		return prompt
	}
	return snippet
}

func appendActiveTextPrompt(records []store.RunningPromptRecord, seen map[string]bool, name, path, id string, maxWords int, isFull bool) []store.RunningPromptRecord {
	data, err := os.ReadFile(resolveActiveTextPath(path))
	clean := strings.TrimSpace(string(data))
	key := fmt.Sprintf("%s:%s", path, clean)
	if err != nil || len(clean) == 0 || seen[key] {
		return records
	}
	seen[key] = true
	now := time.Now().UTC().Format(time.RFC3339)
	rec := buildRunningPromptRecord(name, path, id, clean, now, store.PromptStatusRunning, maxWords, isFull)
	return append(records, rec)
}

func resolveActiveTextPath(dir string) string {
	nested := filepath.Join(dir, ".ai-memory", "temp", "active-agy-pipeline-fix-prompt.txt")
	if _, err := os.Stat(nested); err == nil {
		return nested
	}
	return filepath.Join(dir, "active-agy-pipeline-fix-prompt.txt")
}

func countPromptStatuses(items []store.RunningPromptRecord) (int, int) {
	running, enqueued := 0, 0
	for _, it := range items {
		if it.Status == store.PromptStatusRunning {
			running++
		}
		if it.Status == store.PromptStatusEnqueued {
			enqueued++
		}
	}
	return running, enqueued
}

// RunRunningPromptsBackup snapshots all running and queued prompts into the backup database.
func RunRunningPromptsBackup(customFile string, isJSON, isSSH bool) error {
	if isSSH {
		return AggregateSSHRunningPromptsBackup(customFile, isJSON)
	}
	summary, err := RunRunningPromptsBackupSummary(customFile)
	if err != nil {
		return err
	}
	return outputBackupResult(summary, isJSON)
}

// RunRunningPromptsBackupSummary snapshots running prompts and returns the PromptBackupSummary.
func RunRunningPromptsBackupSummary(customFile string) (store.PromptBackupSummary, error) {
	db, err := store.OpenBackupPromptsSplitDB(customFile)
	if err != nil {
		return store.PromptBackupSummary{}, err
	}
	defer db.Close()
	_, _ = db.PruneOldRestoredEntries(24*time.Hour, false)
	items, cErr := CollectActiveAndQueuedPrompts(0, true)
	if cErr != nil {
		return store.PromptBackupSummary{}, cErr
	}
	return saveBackupBatchSummary(db, items)
}

func executeBackupBatchSave(db *store.BackupPromptsSplitDB, items []store.RunningPromptRecord, isJSON bool) error {
	summary, err := saveBackupBatchSummary(db, items)
	if err != nil {
		return err
	}
	return outputBackupResult(summary, isJSON)
}

func saveBackupBatchSummary(db *store.BackupPromptsSplitDB, items []store.RunningPromptRecord) (store.PromptBackupSummary, error) {
	batchID := fmt.Sprintf("b-%x", time.Now().UnixNano()%0xffffffff)
	rCount, eCount := countPromptStatuses(items)
	now := time.Now().UTC().Format(time.RFC3339)
	projNames := extractCleanProjectNames(items)
	summary := store.PromptBackupSummary{
		BatchId: batchID, TotalPrompts: len(items),
		RunningCount: rCount, EnqueuedCount: eCount,
		ProjectNames: projNames,
		DatabasePath: db.Path(), CreatedAt: now, Items: items,
	}
	if err := db.InsertBackupBatch(summary); err != nil {
		return store.PromptBackupSummary{}, err
	}
	return summary, nil
}

func extractCleanProjectNames(items []store.RunningPromptRecord) []string {
	seen := make(map[string]bool)
	var names []string
	for _, it := range items {
		clean := resolveHumanProjectName(it.ProjectName, it.ProjectPath)
		if len(clean) > 0 && !seen[clean] {
			seen[clean] = true
			names = append(names, clean)
		}
	}
	return names
}

func outputBackupResult(summary store.PromptBackupSummary, isJSON bool) error {
	if isJSON {
		return printJSON(summary)
	}
	projLabel := formatBackupProjectsLabel(summary.ProjectNames)
	fmt.Printf("%s✔ Backed up %d running prompts across %d project(s) [%s] (Batch: %s, Running: %d, Enqueued: %d)%s\n",
		constants.ColorGreen, summary.TotalPrompts, len(summary.ProjectNames), projLabel,
		summary.BatchId, summary.RunningCount, summary.EnqueuedCount, constants.ColorReset)
	return nil
}

func formatBackupProjectsLabel(names []string) string {
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}

// RunRunningPromptsBackupLs lists backup batches recorded in the database.
func RunRunningPromptsBackupLs(customFile string, isJSON, isSSH bool) error {
	if isSSH {
		return AggregateSSHRunningPromptsBackupLs(customFile, isJSON)
	}
	db, err := store.OpenBackupPromptsSplitDB(customFile)
	if err != nil {
		return err
	}
	defer db.Close()
	_, _ = db.PruneOldRestoredEntries(24*time.Hour, false)
	return displayBackupBatches(db, isJSON)
}

func displayBackupBatches(db *store.BackupPromptsSplitDB, isJSON bool) error {
	batches, err := db.ListBackupBatches()
	if err != nil {
		return err
	}
	if isJSON {
		return printJSON(batches)
	}
	dbPath, dbSize, _ := db.GetStorageInfo()
	RenderBackupBatchesTable(dbPath, dbSize, batches)
	return nil
}

func printJSON(v interface{}) error {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
	return nil
}

// RunRunningPromptsRestore restores prompts from the latest backup batch.
func RunRunningPromptsRestore(opts store.RestoreOptions) error {
	if opts.IsSSH {
		return AggregateSSHRunningPromptsRestore(opts)
	}
	summary, restoredCount, err := RunRunningPromptsRestoreSummary(opts)
	if err != nil {
		return err
	}
	return outputRestoreResult(summary.BatchId, restoredCount, summary.ProjectNames, opts)
}

// RunRunningPromptsRestoreSummary restores prompts from the latest batch and returns the summary and count.
func RunRunningPromptsRestoreSummary(opts store.RestoreOptions) (store.PromptBackupSummary, int, error) {
	db, err := store.OpenBackupPromptsSplitDB(opts.TargetFile)
	if err != nil {
		return store.PromptBackupSummary{}, 0, err
	}
	defer db.Close()
	_, _ = db.PruneOldRestoredEntries(24*time.Hour, false)
	batches, bErr := db.ListBackupBatches()
	if bErr != nil {
		return store.PromptBackupSummary{}, 0, bErr
	}
	if len(batches) == 0 {
		return store.PromptBackupSummary{}, 0, apperror.NewSimple("no backup batches found to restore", "E404")
	}
	return executeRestoreBatchSummary(db, batches[0], opts)
}

func executeRestoreBatchSummary(db *store.BackupPromptsSplitDB, latest store.PromptBackupBatchRecord, opts store.RestoreOptions) (store.PromptBackupSummary, int, error) {
	summary, err := db.GetBackupBatch(latest.BatchID)
	if err != nil {
		return store.PromptBackupSummary{}, 0, err
	}
	summary.ProjectNames = extractCleanProjectNames(summary.Items)
	restoredCount := restoreItemsToWorkspaces(summary.Items)
	if mErr := db.MarkBatchRestored(latest.BatchID, opts.IsKeep, restoredCount); mErr != nil {
		return store.PromptBackupSummary{}, 0, mErr
	}
	return *summary, restoredCount, nil
}

func outputRestoreResult(batchID string, count int, projectNames []string, opts store.RestoreOptions) error {
	if opts.IsJSON {
		res := map[string]interface{}{
			"batchId": batchID, "restoredCount": count,
			"projectCount": len(projectNames), "projectNames": projectNames,
			"isKept": opts.IsKeep,
		}
		return printJSON(res)
	}
	projLabel := formatBackupProjectsLabel(projectNames)
	fmt.Printf("%s✔ Restored %d prompts across %d project(s) [%s] from batch %s (keep: %v)%s\n",
		constants.ColorGreen, count, len(projectNames), projLabel, batchID, opts.IsKeep, constants.ColorReset)
	return nil
}

// RunRunningPromptsClean prunes expired restored entries or forces deletion of all backups.
func RunRunningPromptsClean(customFile string, isForce bool) error {
	db, err := store.OpenBackupPromptsSplitDB(customFile)
	if err != nil {
		return err
	}
	defer db.Close()
	count, pErr := db.PruneOldRestoredEntries(24*time.Hour, isForce)
	if pErr != nil {
		return pErr
	}
	if count == 0 {
		fmt.Println("No old backup data to remove.")
	}
	return nil
}

// RunRunningPromptsLs lists current active and enqueued prompts with truncation and limit options.
func RunRunningPromptsLs(limit, wordCount int, isFull, isJSON, isSSH bool) error {
	items, err := CollectActiveAndQueuedPrompts(wordCount, isFull)
	if err != nil {
		return err
	}
	if isSSH {
		return AggregateSSHRunningPromptsLs(items, limit, wordCount, isFull, isJSON)
	}
	items = applyPromptsLimit(items, limit)
	if isJSON {
		return printJSON(items)
	}
	RenderRunningPromptsTable(items, isFull)
	return nil
}

func applyPromptsLimit(items []store.RunningPromptRecord, limit int) []store.RunningPromptRecord {
	if limit > 0 && len(items) > limit {
		return items[:limit]
	}
	return items
}
