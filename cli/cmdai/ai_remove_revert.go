// Package cmdai — ai_remove_revert.go: cryptographic revert and restoration engine for safely staged removals.
package cmdai

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fspath"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	revertCmd = newRevertCommand("revert", []string{"restore"})
	undoCmd   = newRevertCommand("undo", []string{"revert", "restore"})
	rmUndoCmd = newRevertCommand("undo", []string{"revert", "restore"})
	aiUndoCmd = newRevertCommand("undo", []string{"revert", "restore"})

	revertTaskFlag string
	revertJsonFlag bool
)

// AiRevertResponse defines JSON output envelope for revert operations.
type AiRevertResponse struct {
	Status     string             `json:"status"`
	TaskId     string             `json:"taskId"`
	TotalFiles int                `json:"totalFiles"`
	Files      []AiRevertFileInfo `json:"files"`
}

// AiRevertFileInfo describes a single file restored to working tree.
type AiRevertFileInfo struct {
	RelPath      string `json:"relPath"`
	RestoredPath string `json:"restoredPath"`
	Sha256       string `json:"sha256"`
	Verified     bool   `json:"verified"`
}

func newRevertCommand(use string, aliases []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     use + " [task_id] [files...]",
		Aliases: aliases,
		Short:   "Restore files previously removed during an AI analysis task",
		RunE:    runRevertCmd,
	}
	cmd.Flags().StringVarP(&revertTaskFlag, "task", "t", "", "Task session identifier")
	cmd.Flags().BoolVar(&revertJsonFlag, "json", false, "Output results in JSON format")

	return cmd
}

func init() {
	AiAnalysisCmd.AddCommand(revertCmd)
	AiAnalysisCmd.AddCommand(undoCmd)
	AiCmd.AddCommand(aiUndoCmd)
	removeCmd.AddCommand(rmUndoCmd)
}

func runRevertCmd(cmd *cobra.Command, args []string) error {
	taskId, filePaths, err := parseRevertArgs(args)
	if err != nil {
		return err
	}

	runErr := RunAiRevert(taskId, filePaths)
	if runErr != nil {
		return runErr
	}

	return renderRevertSuccess(taskId, filePaths, revertJsonFlag)
}

func parseRevertArgs(args []string) (string, []string, error) {
	hasTaskFlag := strings.TrimSpace(revertTaskFlag) != ""
	if hasTaskFlag {
		return strings.TrimSpace(revertTaskFlag), args, nil
	}

	hasArgs := len(args) > 0
	if hasArgs {
		return args[0], args[1:], nil
	}

	return resolveActiveRevertTask()
}

func resolveActiveRevertTask() (string, []string, error) {
	repoRoot := getEffectiveRepoRoot()
	db, err := store.OpenAiAnalysisSplitDB("", repoRoot)
	if err != nil {
		return "", nil, apperror.NewValidationError("task ID required via --task or as first argument")
	}
	defer db.Close()

	activeId, activeErr := resolveActiveTaskId(db)
	if activeErr == nil && activeId != "" {
		return activeId, nil, nil
	}

	return "", nil, apperror.NewValidationError("task ID required via --task or as first argument")
}

// RunAiRevert validates SHA-256 and restores removed files from the OS temp vault.
func RunAiRevert(taskId string, filePaths []string) error {
	cleanTaskId := strings.TrimSpace(taskId)
	hasTaskId := cleanTaskId != ""
	if !hasTaskId {
		return apperror.NewValidationError("task ID cannot be empty")
	}

	vaultDir, checkErr := checkVaultExists(cleanTaskId)
	if checkErr != nil {
		return checkErr
	}

	repoRoot := getEffectiveRepoRoot()
	db, _ := store.OpenAiAnalysisSplitDB("", repoRoot)
	if db != nil {
		defer db.Close()
	}

	files, collectErr := collectFilesToRevert(db, cleanTaskId, vaultDir, filePaths)
	if collectErr != nil {
		return collectErr
	}

	return revertFiles(db, repoRoot, cleanTaskId, vaultDir, files)
}

// RunAiUndo is an alias for RunAiRevert.
func RunAiUndo(taskId string, filePaths []string) error {
	return RunAiRevert(taskId, filePaths)
}

func checkVaultExists(taskId string) (string, error) {
	vaultDir := filepath.Join(os.TempDir(), fspath.RepoTempSubdir, "removed", taskId)
	info, statErr := os.Stat(vaultDir)
	isMissing := statErr != nil || !info.IsDir()
	if isMissing {
		return "", fmt.Errorf("E1077: Cannot revert task '%s' — OS temporary backup was purged", taskId)
	}

	entries, readErr := os.ReadDir(vaultDir)
	isEmpty := readErr != nil || len(entries) == 0
	if isEmpty {
		return "", fmt.Errorf("E1077: Cannot revert task '%s' — OS temporary backup was purged", taskId)
	}

	return vaultDir, nil
}

func queryDbFilesIfPresent(db *store.AiAnalysisSplitDB, taskId string, filter []string) []store.AiTaskFile {
	if db == nil {
		return nil
	}

	return queryRemovedFilesFromDb(db, taskId, filter)
}

func collectFilesToRevert(db *store.AiAnalysisSplitDB, taskId, vaultDir string, filter []string) ([]store.AiTaskFile, error) {
	dbFiles := queryDbFilesIfPresent(db, taskId, filter)
	hasDbFiles := len(dbFiles) > 0
	if hasDbFiles {
		return dbFiles, nil
	}

	manifestFiles := queryRemovedFilesFromManifest(vaultDir, filter)
	hasManifest := len(manifestFiles) > 0
	if hasManifest {
		return manifestFiles, nil
	}

	return nil, apperror.NewNotFoundError("no removed files found to revert for task: " + taskId)
}

func queryRemovedFilesFromDb(db *store.AiAnalysisSplitDB, taskId string, filter []string) []store.AiTaskFile {
	task, err := db.GetTask(taskId)
	isFound := err == nil && task != nil
	if !isFound {
		return nil
	}

	query := `SELECT 
		AiTaskFileId, AiTaskId, RelPath, AbsPath, Action, BeforeSha256, BackupPath
	FROM AiTaskFile WHERE AiTaskId = ? AND IsRemoved = 1 ORDER BY AiTaskFileId ASC`

	rows, qErr := db.Conn().Query(query, task.AiTaskId)
	if qErr != nil {
		return nil
	}
	defer rows.Close()

	return scanAndFilterRemovedRows(rows, filter)
}

func scanAndFilterRemovedRows(rows *sql.Rows, filter []string) []store.AiTaskFile {
	var results []store.AiTaskFile
	for rows.Next() {
		var f store.AiTaskFile
		err := rows.Scan(&f.AiTaskFileId, &f.AiTaskId, &f.RelPath, &f.AbsPath, &f.Action, &f.BeforeSha256, &f.BackupPath)
		if err != nil {
			continue
		}
		f.IsRemoved = true
		if isTargetFilterMatch(f.RelPath, filter) {
			results = append(results, f)
		}
	}

	return results
}

func queryRemovedFilesFromManifest(vaultDir string, filter []string) []store.AiTaskFile {
	manifestPath := filepath.Join(vaultDir, "manifest.json")
	manifest := loadOrCreateManifest(manifestPath, "", "")
	hasFiles := len(manifest.Files) > 0
	if !hasFiles {
		return nil
	}

	var results []store.AiTaskFile
	for _, f := range manifest.Files {
		if isTargetFilterMatch(f.RelPath, filter) {
			results = append(results, store.AiTaskFile{
				RelPath:      f.RelPath,
				BeforeSha256: f.BeforeSha256,
				BackupPath:   f.BackupPath,
				IsRemoved:    true,
				Action:       "remove",
			})
		}
	}

	return results
}

func isTargetFilterMatch(relPath string, filter []string) bool {
	isEmpty := len(filter) == 0
	if isEmpty {
		return true
	}

	norm := filepath.ToSlash(filepath.Clean(relPath))
	for _, target := range filter {
		cleanTarget := filepath.ToSlash(filepath.Clean(target))
		isMatch := norm == cleanTarget || strings.HasSuffix(norm, "/"+cleanTarget) || filepath.Base(norm) == cleanTarget
		if isMatch {
			return true
		}
	}

	return false
}

func revertFiles(db *store.AiAnalysisSplitDB, repoRoot, taskId, vaultDir string, files []store.AiTaskFile) error {
	for _, f := range files {
		err := restoreSingleFile(db, repoRoot, taskId, vaultDir, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func restoreSingleFile(db *store.AiAnalysisSplitDB, repoRoot, taskId, vaultDir string, f store.AiTaskFile) error {
	backupPath := resolveVaultBackupPath(vaultDir, f)
	data, valErr := validateBackupData(taskId, backupPath, f.BeforeSha256)
	if valErr != nil {
		return valErr
	}

	targetPath := filepath.Join(repoRoot, filepath.FromSlash(f.RelPath))
	writeErr := writeRestoredBytes(targetPath, data)
	if writeErr != nil {
		return writeErr
	}

	return updateDbRevertRecord(db, f.AiTaskFileId, f.AiTaskId, f.RelPath)
}

func isValidBackupFile(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)

	return err == nil && !info.IsDir()
}

func resolveVaultBackupPath(vaultDir string, f store.AiTaskFile) string {
	if isValidBackupFile(f.BackupPath) {
		return f.BackupPath
	}

	return filepath.Join(vaultDir, filepath.FromSlash(f.RelPath))
}

func validateBackupData(taskId, backupPath, expectedSha string) ([]byte, error) {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, fmt.Errorf("E1077: Cannot revert task '%s' — OS temporary backup was purged", taskId)
	}

	actualSha := computeSha256(data)
	hasExpected := expectedSha != ""
	isMatch := !hasExpected || actualSha == expectedSha
	if !isMatch {
		return nil, apperror.NewValidationError(fmt.Sprintf("checksum mismatch for staged file: expected %s, got %s", expectedSha, actualSha))
	}

	return data, nil
}

func writeRestoredBytes(destPath string, data []byte) error {
	dirErr := os.MkdirAll(filepath.Dir(destPath), 0755)
	if dirErr != nil {
		return apperror.WrapSimple(dirErr, "cmdai.writeRestoredBytes.mkdir")
	}

	writeErr := os.WriteFile(destPath, data, 0644)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "cmdai.writeRestoredBytes.write")
	}

	return nil
}

func updateDbRevertByFileId(db *store.AiAnalysisSplitDB, now string, fileId int64) error {
	query := `UPDATE AiTaskFile SET IsRemoved = 0, Action = 'revert', UpdatedAt = ? WHERE AiTaskFileId = ?`
	res := store.ExecWrapper(db.Conn(), query, now, fileId)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.updateDbRevertRecord.id")
	}

	return nil
}

func updateDbRevertRecord(db *store.AiAnalysisSplitDB, fileId, taskId int64, relPath string) error {
	if db == nil {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	hasFileId := fileId > 0
	if hasFileId {
		return updateDbRevertByFileId(db, now, fileId)
	}

	query := `UPDATE AiTaskFile SET IsRemoved = 0, Action = 'revert', UpdatedAt = ? WHERE RelPath = ?`
	res := store.ExecWrapper(db.Conn(), query, now, relPath)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.updateDbRevertRecord.rel")
	}

	return nil
}

func renderRevertSuccess(taskId string, files []string, isJson bool) error {
	if isJson {
		resp := AiRevertResponse{
			Status:     "restored",
			TaskId:     taskId,
			TotalFiles: len(files),
		}
		return printJsonOutput(resp)
	}

	fmt.Printf("Successfully restored removed file(s) for task %s\n", taskId)

	return nil
}
