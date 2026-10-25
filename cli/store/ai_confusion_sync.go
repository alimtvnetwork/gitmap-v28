package store

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"
)

// AiConfusionSyncResult describes the outcome of a 6-stage safe synchronization.
type AiConfusionSyncResult struct {
	TotalExported int    `json:"totalExported"`
	TotalReplayed int    `json:"totalReplayed"`
	DbPath        string `json:"dbPath"`
	IsCommitted   bool   `json:"isCommitted"`
	IsPushed      bool   `json:"isPushed"`
	Message       string `json:"message"`
}

// SyncAiConfusion performs the safe 6-stage synchronization cycle on the AI confusion database.
func SyncAiConfusion(repoCacheDir string) (*AiConfusionSyncResult, error) {
	root := resolveSyncRepoCacheRoot(repoCacheDir)
	dbPath := ResolveAiConfusionDbPath(root)
	tempPath := generateTempSyncDbPath()
	defer os.Remove(tempPath)

	exported, err := stage1ExportToTempDb(dbPath, tempPath)
	if err != nil {
		return nil, err
	}

	relPath := resolveRelDbPath(root, dbPath)
	executePullSequence(root, relPath)

	return completeSyncPipeline(dbPath, root, tempPath, relPath, exported)
}

func executePullSequence(root, relPath string) {
	_ = stage2DiscardBinaryChanges(root, relPath)
	_ = stage3RebasePull(root)
}

func completeSyncPipeline(dbPath, root, tempPath, relPath string, exported int) (*AiConfusionSyncResult, error) {
	replayed, repErr := stage4ReplayRecords(dbPath, tempPath)
	if repErr != nil {
		return nil, repErr
	}

	isCommitted := stage5CommitSync(root, relPath)
	isPushed := stage6PushSync(root)

	return buildSyncResult(dbPath, exported, replayed, isCommitted, isPushed), nil
}

func stage1ExportToTempDb(srcDbPath, tempDbPath string) (int, error) {
	hasSrc := hasFile(srcDbPath)
	if hasSrc {
		return copyRecordsToTempDb(srcDbPath, tempDbPath)
	}

	return 0, nil
}

func copyRecordsToTempDb(srcDbPath, tempDbPath string) (int, error) {
	srcStore, err := OpenAiConfusionStore(srcDbPath)
	if err != nil {
		return 0, err
	}
	defer srcStore.Close()

	records, listErr := srcStore.ListAiConfusions(AiConfusionFilter{})
	if listErr != nil {
		return 0, listErr
	}

	return insertRecordsIntoTempDb(tempDbPath, records)
}

func insertRecordsIntoTempDb(tempDbPath string, records []AiConfusionRecord) (int, error) {
	tempStore, err := OpenAiConfusionStore(tempDbPath)
	if err != nil {
		return 0, err
	}
	defer tempStore.Close()

	for _, rec := range records {
		recCopy := rec
		if insErr := tempStore.InsertAiConfusion(&recCopy); insErr != nil {
			return 0, insErr
		}
	}

	return len(records), nil
}

func stage2DiscardBinaryChanges(repoRoot, relPath string) error {
	hasGit := isGitDirectory(repoRoot)
	if hasGit {
		_ = runGitCmd(repoRoot, "checkout", "--", relPath)
	}

	return nil
}

func stage3RebasePull(repoRoot string) error {
	hasGit := isGitDirectory(repoRoot)
	if hasGit {
		_ = runGitCmd(repoRoot, "pull", "--rebase")
	}

	return nil
}

const sqlUpsertAiConfusion = `
INSERT INTO AiConfusion (
    ConfusionUuid, AgentSlug, FileNames, MdFiles, Spec, Task, Subtask,
    HeatmapInfo, LineNumbers, CodeSnippet, Reason, InitialIntention,
    ReasoningThinking, ConfusionQuestion, IsCompleted, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(ConfusionUuid) DO UPDATE SET
    Reason = CASE WHEN excluded.Reason != '' THEN excluded.Reason ELSE AiConfusion.Reason END,
    IsCompleted = CASE WHEN excluded.IsCompleted = 1 THEN 1 ELSE AiConfusion.IsCompleted END,
    UpdatedAt = excluded.UpdatedAt;`

func stage4ReplayRecords(targetDbPath, tempDbPath string) (int, error) {
	tempStore, err := OpenAiConfusionStore(tempDbPath)
	if err != nil {
		return 0, err
	}
	defer tempStore.Close()

	records, listErr := tempStore.ListAiConfusions(AiConfusionFilter{})
	if listErr != nil {
		return 0, listErr
	}

	return executeReplayMerge(targetDbPath, records)
}

func executeReplayMerge(targetDbPath string, records []AiConfusionRecord) (int, error) {
	destStore, err := OpenAiConfusionStore(targetDbPath)
	if err != nil {
		return 0, err
	}
	defer destStore.Close()

	for _, rec := range records {
		args := destStore.buildInsertArgs(&rec)
		_, execErr := destStore.db.Exec(sqlUpsertAiConfusion, args...)
		if execErr != nil {
			return 0, apperror.WrapSimple(execErr, "replay ai confusion record")
		}
	}

	return len(records), nil
}

func stage5CommitSync(repoRoot, relPath string) bool {
	hasGit := isGitDirectory(repoRoot)
	if hasGit {
		_ = runGitCmd(repoRoot, "add", "--", relPath)
		commitMsg := "AI_CONFUSION - sync and replay local agent episodes"
		err := runGitCmd(repoRoot, "commit", "-m", commitMsg)

		return err == nil
	}

	return false
}

func stage6PushSync(repoRoot string) bool {
	hasGit := isGitDirectory(repoRoot)
	if hasGit {
		err := runGitCmd(repoRoot, "push")

		return err == nil
	}

	return false
}

func runGitCmd(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	return cmd.Run()
}

func isGitDirectory(dir string) bool {
	gitMarker := filepath.Join(dir, ".git")
	_, err := os.Stat(gitMarker)

	return err == nil
}

func hasFile(filePath string) bool {
	info, err := os.Stat(filePath)
	hasStat := err == nil && !info.IsDir()

	return hasStat
}

func generateTempSyncDbPath() string {
	fileName := fmt.Sprintf("ai_confusion_sync_%d_%s.db", time.Now().UnixNano(), uuidutil.GenerateUUID()[:8])

	return filepath.Join(os.TempDir(), fileName)
}

func resolveSyncRepoCacheRoot(repoCacheDir string) string {
	hasDir := len(strings.TrimSpace(repoCacheDir)) > 0
	if hasDir {
		return repoCacheDir
	}

	return ResolveRepoCacheRoot()
}

func resolveRelDbPath(root, dbPath string) string {
	rel, err := filepath.Rel(root, dbPath)
	hasRel := err == nil
	if hasRel {
		return filepath.ToSlash(rel)
	}

	return filepath.Base(dbPath)
}

func buildSyncResult(dbPath string, exported, replayed int, isCommitted, isPushed bool) *AiConfusionSyncResult {
	msg := fmt.Sprintf("Synced %d exported and %d replayed records (committed=%v, pushed=%v)",
		exported, replayed, isCommitted, isPushed)

	return &AiConfusionSyncResult{
		TotalExported: exported,
		TotalReplayed: replayed,
		DbPath:        dbPath,
		IsCommitted:   isCommitted,
		IsPushed:      isPushed,
		Message:       msg,
	}
}
