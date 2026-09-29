// Package cmdagy — agy_wpr_split_db.go coordinates Split-DB snapshots and deduplication for watch-prompts-running.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SnapshotAllRunningPrompts captures all running and queued prompts into repo-scoped Split-DBs.
func SnapshotAllRunningPrompts() ([]store.WatchPromptsSummary, error) {
	items, err := CollectActiveAndQueuedPrompts(0, true)
	if err != nil {
		return nil, err
	}

	return persistPromptRecordsToSplitDB(items)
}

func persistPromptRecordsToSplitDB(items []store.RunningPromptRecord) ([]store.WatchPromptsSummary, error) {
	var summaries []store.WatchPromptsSummary
	now := time.Now().UTC().Format(time.RFC3339)

	for _, item := range items {
		summary, saveErr := persistSinglePromptRecord(item, now)
		if saveErr == nil {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

func persistSinglePromptRecord(item store.RunningPromptRecord, now string) (store.WatchPromptsSummary, error) {
	slug := resolveRepoSlug(item.ProjectPath, item.ProjectName)
	db, err := store.OpenWatchPromptsSplitDB(slug)
	if err != nil {
		return store.WatchPromptsSummary{}, err
	}
	defer db.Close()

	_ = db.PrunePreviousPromptRecords(slug)
	rec := buildWatchPromptRecord(item, slug, now)
	if err := db.SavePromptRecord(rec); err != nil {
		return store.WatchPromptsSummary{}, err
	}

	_ = db.InsertWatchLog(slug, "snapshot", "Backed up recent prompt into split-db", "success")

	return buildWatchSummaryFromRecord(rec, db.Path()), nil
}

func resolveRepoSlug(projectPath, projectName string) string {
	if len(strings.TrimSpace(projectPath)) > 0 {
		return store.SanitizeSlug(filepath.Base(filepath.Clean(projectPath)))
	}

	return store.SanitizeSlug(projectName)
}

func buildWatchPromptRecord(item store.RunningPromptRecord, slug, now string) store.WatchPromptRecord {
	mediaPaths := extractPromptMediaPaths(item.Prompt, item.ProjectPath)
	recordId := fmt.Sprintf("wpr-%x", time.Now().UnixNano()%0xffffffff)

	return store.WatchPromptRecord{
		RecordId:       recordId,
		RepoSlug:       slug,
		ProjectName:    item.ProjectName,
		ProjectPath:    item.ProjectPath,
		ProjectId:      item.ProjectId,
		SequenceId:     item.ProjectId,
		ConversationId: item.ProjectId,
		PromptText:     item.Prompt,
		PromptStatus:   string(item.Status),
		WordCount:      item.WordCount,
		MediaPaths:     mediaPaths,
		IsActive:       true,
		UpdatedAt:      now,
		CreatedAt:      item.CreatedAt,
	}
}

func buildWatchSummaryFromRecord(rec store.WatchPromptRecord, dbPath string) store.WatchPromptsSummary {
	return store.WatchPromptsSummary{
		RepoSlug:     rec.RepoSlug,
		ProjectName:  rec.ProjectName,
		ProjectPath:  rec.ProjectPath,
		ProjectId:    rec.ProjectId,
		SequenceId:   rec.SequenceId,
		PromptStatus: rec.PromptStatus,
		WordCount:    rec.WordCount,
		MediaCount:   len(rec.MediaPaths),
		MediaPaths:   rec.MediaPaths,
		IsWatching:   true,
		DatabasePath: dbPath,
		UpdatedAt:    rec.UpdatedAt,
	}
}

// LoadWatchPromptsSummary loads the latest backed-up prompt record for a repo slug.
func LoadWatchPromptsSummary(slug string) (store.WatchPromptsSummary, bool) {
	db, err := store.OpenWatchPromptsSplitDB(slug)
	if err != nil {
		return store.WatchPromptsSummary{}, false
	}
	defer db.Close()

	records, err := db.ListPromptRecords(slug)
	if err != nil || len(records) == 0 {
		return store.WatchPromptsSummary{}, false
	}

	return buildWatchSummaryFromRecord(records[0], db.Path()), true
}

// ImportWatchPromptsSplitDB imports an external database file into repo-scoped Split-DB location.
func ImportWatchPromptsSplitDB(sourceFile, targetSlug string) error {
	cleanSlug := store.SanitizeSlug(targetSlug)
	targetPath := store.ResolveWatchPromptsDbPath(cleanSlug)
	_ = os.MkdirAll(filepath.Dir(targetPath), 0755)

	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return apperror.WrapSimple(err, "read source split db for import")
	}

	writeErr := os.WriteFile(targetPath, data, 0644)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "write imported watch-prompts split db")
	}

	return nil
}
