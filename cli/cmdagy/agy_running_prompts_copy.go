// Package cmdagy — agy_running_prompts_copy.go implements prompt batch and item cloning logic.
package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ExecutePromptCopy clones a prompt batch or individual prompt item into a destination name or namespace.
func ExecutePromptCopy(opts PromptCopyOptions) (*PromptCopyResult, *appfault.AppError) {
	cleanSource := strings.TrimSpace(opts.SourceID)
	isSourceEmpty := len(cleanSource) == 0
	if isSourceEmpty {
		return nil, appfault.NewValidation("ExecutePromptCopy", "E1000", "sourceId is required")
	}

	db, err := store.OpenBackupPromptsSplitDB(opts.TargetFile)
	isDbErr := err != nil
	if isDbErr {
		return nil, appfault.WrapSimple(err, "open backup split db")
	}

	defer db.Close()

	summary, bErr := db.GetBackupBatch(cleanSource)
	isBatchFound := bErr == nil && summary != nil
	if isBatchFound {
		return copyBatchSummary(db, *summary, opts)
	}

	item, iErr := findSingleItemByID(db.Conn(), cleanSource)
	isItemFound := iErr == nil && item != nil
	if isItemFound {
		return copySingleItem(db, *item, cleanSource, opts)
	}

	return nil, appfault.NewNotFound("ExecutePromptCopy", "E404", "prompt batch or item not found: "+cleanSource)
}

// CopyPromptBatch clones an existing prompt batch into a new destination batch.
func CopyPromptBatch(db *store.BackupPromptsSplitDB, opts PromptCopyOptions) (*PromptCopyResult, *appfault.AppError) {
	summary, err := db.GetBackupBatch(opts.SourceID)
	isGetErr := err != nil
	if isGetErr {
		return nil, appfault.WrapSimple(err, "get backup batch")
	}

	return copyBatchSummary(db, *summary, opts)
}

// CopyPromptItem clones a single prompt item into a new destination batch.
func CopyPromptItem(db *store.BackupPromptsSplitDB, opts PromptCopyOptions) (*PromptCopyResult, *appfault.AppError) {
	item, err := findSingleItemByID(db.Conn(), opts.SourceID)
	isFindErr := err != nil
	if isFindErr {
		return nil, err
	}

	return copySingleItem(db, *item, opts.SourceID, opts)
}

func copyBatchSummary(db *store.BackupPromptsSplitDB, summary store.PromptBackupSummary, opts PromptCopyOptions) (*PromptCopyResult, *appfault.AppError) {
	destBatchID := resolveDestinationBatchID(opts.DestinationID, opts.SourceID)
	clonedItems := clonePromptRecords(summary.Items, opts.DestinationName, opts.Namespace)
	rCount, eCount := countPromptStatuses(clonedItems)
	now := time.Now().UTC().Format(time.RFC3339)
	projectNames := extractCleanProjectNames(clonedItems)

	newSummary := store.PromptBackupSummary{
		BatchId:       destBatchID,
		TotalPrompts:  len(clonedItems),
		RunningCount:  rCount,
		EnqueuedCount: eCount,
		ProjectNames:  projectNames,
		DatabasePath:  db.Path(),
		CreatedAt:     now,
		Items:         clonedItems,
	}

	isDryRun := opts.IsDryRun
	if isDryRun {
		return buildCopyDryRunResult(opts, destBatchID, len(clonedItems)), nil
	}

	insertErr := db.InsertBackupBatch(newSummary)
	isInsertErr := insertErr != nil
	if isInsertErr {
		return nil, appfault.WrapSimple(insertErr, "insert cloned backup batch")
	}

	return buildCopySuccessResult(opts, destBatchID, len(clonedItems)), nil
}

func copySingleItem(db *store.BackupPromptsSplitDB, item store.RunningPromptRecord, itemID string, opts PromptCopyOptions) (*PromptCopyResult, *appfault.AppError) {
	destBatchID := resolveDestinationBatchID(opts.DestinationID, "single")
	cloned := cloneSingleRecord(item, opts.DestinationName, opts.Namespace)
	clonedList := []store.RunningPromptRecord{cloned}
	rCount, eCount := countPromptStatuses(clonedList)
	now := time.Now().UTC().Format(time.RFC3339)
	projectNames := extractCleanProjectNames(clonedList)

	newSummary := store.PromptBackupSummary{
		BatchId:       destBatchID,
		TotalPrompts:  1,
		RunningCount:  rCount,
		EnqueuedCount: eCount,
		ProjectNames:  projectNames,
		DatabasePath:  db.Path(),
		CreatedAt:     now,
		Items:         clonedList,
	}

	isDryRun := opts.IsDryRun
	if isDryRun {
		return buildCopyDryRunResult(opts, destBatchID, 1), nil
	}

	insertErr := db.InsertBackupBatch(newSummary)
	isInsertErr := insertErr != nil
	if isInsertErr {
		return nil, appfault.WrapSimple(insertErr, "insert cloned item batch")
	}

	return buildCopySuccessResult(opts, destBatchID, 1), nil
}

func buildCopyDryRunResult(opts PromptCopyOptions, destBatchID string, count int) *PromptCopyResult {
	return &PromptCopyResult{
		SourceID:        opts.SourceID,
		DestinationID:   destBatchID,
		CopiedCount:     count,
		DestinationName: opts.DestinationName,
		Namespace:       opts.Namespace,
		IsSuccess:       true,
		Message:         fmt.Sprintf("dry-run: cloned prompt %s -> %s (%d items)", opts.SourceID, destBatchID, count),
	}
}

func buildCopySuccessResult(opts PromptCopyOptions, destBatchID string, count int) *PromptCopyResult {
	return &PromptCopyResult{
		SourceID:        opts.SourceID,
		DestinationID:   destBatchID,
		CopiedCount:     count,
		DestinationName: opts.DestinationName,
		Namespace:       opts.Namespace,
		IsSuccess:       true,
		Message:         fmt.Sprintf("cloned prompt %s -> %s (%d items)", opts.SourceID, destBatchID, count),
	}
}

func resolveDestinationBatchID(destID, sourceID string) string {
	cleanDest := strings.TrimSpace(destID)
	isDestProvided := len(cleanDest) > 0
	if isDestProvided {
		return cleanDest
	}

	cleanSrc := strings.TrimPrefix(sourceID, "copy-")

	return fmt.Sprintf("copy-%x-%s", time.Now().UnixNano()%0xffffff, cleanSrc)
}

func clonePromptRecords(items []store.RunningPromptRecord, destName, namespace string) []store.RunningPromptRecord {
	cloned := make([]store.RunningPromptRecord, len(items))
	for i, it := range items {
		cloned[i] = cloneSingleRecord(it, destName, namespace)
	}

	return cloned
}

func cloneSingleRecord(it store.RunningPromptRecord, destName, namespace string) store.RunningPromptRecord {
	copyItem := it
	cleanDestName := strings.TrimSpace(destName)
	isDestNameProvided := len(cleanDestName) > 0
	if isDestNameProvided {
		copyItem.ProjectName = cleanDestName
	}

	cleanNamespace := strings.TrimSpace(namespace)
	isNamespaceProvided := len(cleanNamespace) > 0
	if isNamespaceProvided {
		copyItem.ProjectName = fmt.Sprintf("%s/%s", cleanNamespace, copyItem.ProjectName)
		copyItem.ProjectPath = fmt.Sprintf("%s/%s", cleanNamespace, copyItem.ProjectPath)
	}

	copyItem.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	return copyItem
}

func findSingleItemByID(conn *sql.DB, itemID string) (*store.RunningPromptRecord, *appfault.AppError) {
	isConnNil := conn == nil
	if isConnNil {
		return nil, appfault.NewValidation("findSingleItemByID", "E1000", "database connection is nil")
	}

	q := `SELECT ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt
	FROM PromptBackupItem WHERE ItemId = ? LIMIT 1`

	var rec store.RunningPromptRecord
	var statusStr, mediaJSON string
	err := conn.QueryRow(q, itemID).Scan(
		&rec.ProjectName, &rec.ProjectPath, &rec.ProjectId,
		&rec.ConversationId, &rec.SequenceId, &rec.Prompt,
		&statusStr, &rec.WordCount, &mediaJSON, &rec.CreatedAt,
	)

	isNotFound := err == sql.ErrNoRows
	if isNotFound {
		return nil, appfault.NewNotFound("findSingleItemByID", "E404", "item not found: "+itemID)
	}

	isQueryErr := err != nil
	if isQueryErr {
		return nil, appfault.WrapSimple(err, "query prompt backup item")
	}

	rec.Status = store.PromptStatusType(statusStr)
	rec.MediaPaths = parseMediaPathsJSON(mediaJSON)

	return &rec, nil
}

func parseMediaPathsJSON(raw string) []string {
	isRawEmpty := len(strings.TrimSpace(raw)) == 0
	if isRawEmpty {
		return nil
	}

	var paths []string
	_ = json.Unmarshal([]byte(raw), &paths)

	return paths
}

// RunRunningPromptsCopy executes the prompt copy CLI action with formatted terminal or JSON output.
func RunRunningPromptsCopy(opts PromptCopyOptions) *appfault.AppError {
	res, err := ExecutePromptCopy(opts)
	isCopyErr := err != nil
	if isCopyErr {
		return err
	}

	isJSON := opts.IsJSON
	if isJSON {
		_ = printJSON(res)
		return nil
	}

	fmt.Printf("%s✔ %s%s\n", constants.ColorGreen, res.Message, constants.ColorReset)

	return nil
}
