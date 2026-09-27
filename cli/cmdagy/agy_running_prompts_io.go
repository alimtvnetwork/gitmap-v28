// Package cmdagy — agy_running_prompts_io.go handles importing, exporting, and workspace queue restoration.
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

func isJSONPath(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".json")
}

func resolveIOFilePath(path, defaultName string) string {
	if len(strings.TrimSpace(path)) > 0 {
		return path
	}
	return defaultName
}

func restoreItemsToWorkspaces(items []store.RunningPromptRecord) int {
	count := 0
	for _, it := range items {
		if enqueueSingleRestoredItem(it) {
			count++
		}
	}
	return count
}

func enqueueSingleRestoredItem(item store.RunningPromptRecord) bool {
	if len(item.ProjectPath) == 0 || len(item.Prompt) == 0 {
		return false
	}
	qPath := filepath.Join(item.ProjectPath, ".ai-memory", "temp", "agy-prompt-queue.json")
	_ = os.MkdirAll(filepath.Dir(qPath), 0755)
	q := loadQueueOrCreate(qPath)
	q.Queued = append(q.Queued, makeRestoredQueueEntry(q, item))
	q.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return saveQueueFile(qPath, q)
}

func makeRestoredQueueEntry(q AgyPromptQueueFile, item store.RunningPromptRecord) AgyPromptQueueEntry {
	return AgyPromptQueueEntry{
		ID:        computeNextQueueID(q),
		Type:      "restored_prompt",
		Title:     fmt.Sprintf("Restored: %s", item.ProjectName),
		Prompt:    item.Prompt,
		Status:    "queued",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func loadQueueOrCreate(qPath string) AgyPromptQueueFile {
	data, err := os.ReadFile(qPath)
	if err != nil {
		return AgyPromptQueueFile{Queued: make([]AgyPromptQueueEntry, 0)}
	}
	var q AgyPromptQueueFile
	if unmarshalErr := json.Unmarshal(data, &q); unmarshalErr != nil {
		return AgyPromptQueueFile{Queued: make([]AgyPromptQueueEntry, 0)}
	}
	return q
}

func saveQueueFile(qPath string, q AgyPromptQueueFile) bool {
	data, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return false
	}
	return os.WriteFile(qPath, data, 0644) == nil
}

// RunRunningPromptsExport exports running prompts into a SQLite database or JSON file.
func RunRunningPromptsExport(targetFile string, maxWords int) error {
	filePath := resolveIOFilePath(targetFile, "gitmap-running-prompts.db")
	isFull := maxWords <= 0
	items, err := CollectActiveAndQueuedPrompts(maxWords, isFull)
	if err != nil {
		return err
	}
	if isJSONPath(filePath) {
		return exportToJSONFile(filePath, items)
	}
	return exportToSQLiteFile(filePath, items)
}

func exportToJSONFile(filePath string, items []store.RunningPromptRecord) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal export json")
	}
	if writeErr := os.WriteFile(filePath, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write export json")
	}
	fmt.Printf("%s✔ Exported %d running prompts to %s%s\n", constants.ColorGreen, len(items), filePath, constants.ColorReset)
	return nil
}

func exportToSQLiteFile(filePath string, items []store.RunningPromptRecord) error {
	db, err := store.OpenBackupPromptsSplitDB(filePath)
	if err != nil {
		return err
	}
	defer db.Close()
	return executeBackupBatchSave(db, items, false)
}

// RunRunningPromptsImport imports prompts from a SQLite database or JSON file into project queues.
func RunRunningPromptsImport(targetFile string, maxWords int) error {
	filePath := resolveIOFilePath(targetFile, "gitmap-running-prompts.db")
	items, err := readImportItems(filePath)
	if err != nil {
		return err
	}
	if maxWords > 0 {
		items = applyImportWordCount(items, maxWords)
	}
	count := restoreItemsToWorkspaces(items)
	fmt.Printf("%s✔ Imported and enqueued %d prompts from %s%s\n", constants.ColorGreen, count, filePath, constants.ColorReset)
	return nil
}

func readImportItems(filePath string) ([]store.RunningPromptRecord, error) {
	if isJSONPath(filePath) {
		return importFromJSONFile(filePath)
	}
	return importFromSQLiteFile(filePath)
}

func applyImportWordCount(items []store.RunningPromptRecord, maxWords int) []store.RunningPromptRecord {
	for i := range items {
		truncated, wc := TruncateWords(items[i].Prompt, maxWords)
		items[i].Prompt = truncated
		items[i].WordCount = wc
	}
	return items
}

func importFromJSONFile(filePath string) ([]store.RunningPromptRecord, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read import json file")
	}
	var items []store.RunningPromptRecord
	if unmarshalErr := json.Unmarshal(data, &items); unmarshalErr != nil {
		return nil, apperror.WrapSimple(unmarshalErr, "unmarshal import json")
	}
	return items, nil
}

func importFromSQLiteFile(filePath string) ([]store.RunningPromptRecord, error) {
	db, err := store.OpenBackupPromptsSplitDB(filePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return readLatestBatchItems(db, filePath)
}

func readLatestBatchItems(db *store.BackupPromptsSplitDB, path string) ([]store.RunningPromptRecord, error) {
	batches, err := db.ListBackupBatches()
	if err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return nil, apperror.NewSimple("no backup batches found in "+path, "E404")
	}
	summary, sErr := db.GetBackupBatch(batches[0].BatchID)
	if sErr != nil {
		return nil, sErr
	}
	return summary.Items, nil
}
