// Package cmdagy — agy_running_prompts_backup.go manages prompt backup, restoration, discovery, and cleanup.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func getAllProjects() ([]AgyProject, error) {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return nil, apperror.WrapSimple(err, "projects dir path")
	}
	return loadAllAgyProjects(dirPath)
}

// CollectActiveAndQueuedPrompts discovers running and pending prompts across all configured workspaces.
func CollectActiveAndQueuedPrompts(maxWords int, isFull bool) ([]store.RunningPromptRecord, error) {
	projects, _ := getAllProjects()
	seen := make(map[string]bool)
	var records []store.RunningPromptRecord
	for _, p := range projects {
		records = collectProjectPrompts(records, seen, p.Name, p.GetPath(), p.ID, maxWords, isFull)
	}
	cwd, _ := os.Getwd()
	if len(cwd) > 0 {
		records = collectProjectPrompts(records, seen, filepath.Base(cwd), cwd, "cwd", maxWords, isFull)
	}
	return records, nil
}

func collectProjectPrompts(records []store.RunningPromptRecord, seen map[string]bool, name, path, id string, maxWords int, isFull bool) []store.RunningPromptRecord {
	if len(path) == 0 {
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
	return store.RunningPromptRecord{
		ProjectName: name, ProjectPath: path, ProjectId: id,
		Prompt: text, Snippet: resolveSnippet(text, snippet, isFull),
		Status: status, WordCount: wordCount, CreatedAt: createdAt,
	}
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
	db, err := store.OpenBackupPromptsSplitDB(customFile)
	if err != nil {
		return err
	}
	defer db.Close()
	_, _ = db.PruneOldRestoredEntries(24*time.Hour, false)
	items, cErr := CollectActiveAndQueuedPrompts(0, true)
	if cErr != nil {
		return cErr
	}
	return executeBackupBatchSave(db, items, isJSON)
}

func executeBackupBatchSave(db *store.BackupPromptsSplitDB, items []store.RunningPromptRecord, isJSON bool) error {
	batchID := fmt.Sprintf("b-%x", time.Now().UnixNano()%0xffffffff)
	rCount, eCount := countPromptStatuses(items)
	now := time.Now().UTC().Format(time.RFC3339)
	summary := store.PromptBackupSummary{
		BatchId: batchID, TotalPrompts: len(items),
		RunningCount: rCount, EnqueuedCount: eCount,
		DatabasePath: db.Path(), CreatedAt: now, Items: items,
	}
	if err := db.InsertBackupBatch(summary); err != nil {
		return err
	}
	return outputBackupResult(summary, isJSON)
}

func outputBackupResult(summary store.PromptBackupSummary, isJSON bool) error {
	if isJSON {
		return printJSON(summary)
	}
	fmt.Printf("%s✔ Backed up %d running prompts (Batch: %s, Running: %d, Enqueued: %d)%s\n",
		constants.ColorGreen, summary.TotalPrompts, summary.BatchId,
		summary.RunningCount, summary.EnqueuedCount, constants.ColorReset)
	return nil
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
	db, err := store.OpenBackupPromptsSplitDB(opts.TargetFile)
	if err != nil {
		return err
	}
	defer db.Close()
	_, _ = db.PruneOldRestoredEntries(24*time.Hour, false)
	batches, bErr := db.ListBackupBatches()
	if bErr != nil {
		return bErr
	}
	return restoreFromLatestBatch(db, batches, opts)
}

func restoreFromLatestBatch(db *store.BackupPromptsSplitDB, batches []store.PromptBackupBatchRecord, opts store.RestoreOptions) error {
	if len(batches) == 0 {
		return apperror.NewSimple("no backup batches found to restore", "E404")
	}
	return executeRestoreBatch(db, batches[0], opts)
}

func executeRestoreBatch(db *store.BackupPromptsSplitDB, latest store.PromptBackupBatchRecord, opts store.RestoreOptions) error {
	summary, err := db.GetBackupBatch(latest.BatchID)
	if err != nil {
		return err
	}
	restoredCount := restoreItemsToWorkspaces(summary.Items)
	if mErr := db.MarkBatchRestored(latest.BatchID, opts.IsKeep, restoredCount); mErr != nil {
		return mErr
	}
	return outputRestoreResult(latest.BatchID, restoredCount, opts)
}

func outputRestoreResult(batchID string, count int, opts store.RestoreOptions) error {
	if opts.IsJSON {
		res := map[string]interface{}{
			"batchId": batchID, "restoredCount": count, "isKept": opts.IsKeep,
		}
		return printJSON(res)
	}
	fmt.Printf("%s✔ Restored %d prompts from batch %s (keep: %v)%s\n",
		constants.ColorGreen, count, batchID, opts.IsKeep, constants.ColorReset)
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
