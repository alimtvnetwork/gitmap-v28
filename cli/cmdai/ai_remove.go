// Package cmdai — ai_remove.go: zero-loss file removal with temp vault staging and SHA-256 recording.
package cmdai

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	removeCmd = &cobra.Command{
		Use:     "rm [files...]",
		Aliases: []string{"remove", "delete"},
		Short:   "Safely remove files with zero-loss temp vault staging and SHA-256 recording",
		RunE:    runRemoveCmd,
	}

	removeTaskFlag   string
	removeReasonFlag string
	removeDryRunFlag bool
	removeJsonFlag   bool
)

// RemovalManifest encapsulates metadata about staged files in the vault.
type RemovalManifest struct {
	TaskId    string                `json:"taskId"`
	UpdatedAt string                `json:"updatedAt"`
	Reason    string                `json:"reason,omitempty"`
	Files     []RemovalManifestFile `json:"files"`
}

// RemovalManifestFile records a single file staged in the removal vault.
type RemovalManifestFile struct {
	RelPath      string `json:"relPath"`
	SizeBytes    int64  `json:"sizeBytes"`
	BeforeSha256 string `json:"beforeSha256"`
	Timestamp    string `json:"timestamp"`
	BackupPath   string `json:"backupPath"`
}

// AiRemoveResponse defines JSON output envelope for file removals.
type AiRemoveResponse struct {
	Status     string                `json:"status"`
	TaskId     string                `json:"taskId"`
	TotalFiles int                   `json:"totalFiles"`
	Files      []RemovalManifestFile `json:"files"`
	VaultDir   string                `json:"vaultDir"`
}

func initRemoveFlags() {
	initRemoveCmdFlags(removeCmd)
}

func initRemoveCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&removeTaskFlag, "task", "t", "", "Task session identifier")
	cmd.Flags().StringVarP(&removeReasonFlag, "reason", "r", "", "Reasoning explanation for removal")
	cmd.Flags().BoolVarP(&removeDryRunFlag, "dry-run", "n", false, "Preview files to remove without deletion")
	cmd.Flags().BoolVar(&removeJsonFlag, "json", false, "Output results in JSON format")
}

func init() {
	initRemoveFlags()
	AiCmd.RemoveCommand(AiAnalysisCmd)
	AiAnalysisCmd.AddCommand(removeCmd)

	bridgeCmd := &cobra.Command{
		Use:                "analysis",
		Aliases:            []string{"ai-analysis", "trace", "session"},
		Short:              "AI task analysis session recording and granular reasoning tracer",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return DispatchAiAnalysis(args)
		},
	}
	AiCmd.AddCommand(bridgeCmd)

	aiRmCmd := &cobra.Command{
		Use:     "rm [files...]",
		Aliases: []string{"remove", "delete"},
		Short:   "Safely remove files with zero-loss temp vault staging and SHA-256 recording",
		RunE:    runRemoveCmd,
	}
	initRemoveCmdFlags(aiRmCmd)
	AiCmd.AddCommand(aiRmCmd)
}

func runRemoveCmd(cmd *cobra.Command, args []string) error {
	taskId, filePaths, err := parseRemoveArgs(args)
	if err != nil {
		return err
	}

	reason := resolveRemovalReason(removeReasonFlag)
	runErr := RunAiRemove(taskId, filePaths, reason, removeDryRunFlag)
	if runErr != nil {
		return runErr
	}

	return renderRemoveSuccess(taskId, filePaths, removeDryRunFlag, removeJsonFlag)
}

func resolveRemovalReason(raw string) string {
	clean := strings.TrimSpace(raw)
	hasReason := clean != ""
	if hasReason {
		return clean
	}

	return "obsolete file removed by autonomous ai agent"
}

func parseRemoveArgs(args []string) (string, []string, error) {
	hasTaskFlag := strings.TrimSpace(removeTaskFlag) != ""
	if hasTaskFlag {
		return parseFlaggedArgs(args)
	}

	return parsePositionalRemoveArgs(args)
}

func parseFlaggedArgs(args []string) (string, []string, error) {
	hasArgs := len(args) > 0
	if !hasArgs {
		return "", nil, apperror.NewValidationError("no files specified for removal")
	}

	return strings.TrimSpace(removeTaskFlag), args, nil
}

func parsePositionalRemoveArgs(args []string) (string, []string, error) {
	hasArgs := len(args) > 0
	if !hasArgs {
		return "", nil, apperror.NewValidationError("task ID and file paths are required")
	}

	isMulti := len(args) > 1
	if isMulti {
		return args[0], args[1:], nil
	}

	return resolveActiveOrFallbackTask(args)
}

func resolveActiveOrFallbackTask(args []string) (string, []string, error) {
	db, err := store.OpenAiAnalysisSplitDB("", getEffectiveRepoRoot())
	if err != nil {
		return "", nil, apperror.NewValidationError("task ID required via --task or as first argument")
	}
	defer db.Close()

	activeId, activeErr := resolveActiveTaskId(db)
	if activeErr == nil && activeId != "" {
		return activeId, args, nil
	}

	return "", nil, apperror.NewValidationError("task ID required via --task or as first argument")
}

func resolveActiveTaskId(db *store.AiAnalysisSplitDB) (string, error) {
	tasks, err := db.ListTasks(true, 1)
	hasActive := err == nil && len(tasks) > 0
	if hasActive {
		return tasks[0].TaskUuid, nil
	}

	return "", apperror.NewValidationError("no active AI task found")
}

// RunAiRemove safely stages target files into temp vault and deletes them from working tree.
func RunAiRemove(taskId string, filePaths []string, reason string, isDryRun bool) error {
	hasFiles := len(filePaths) > 0
	if !hasFiles {
		return apperror.NewValidationError("no file paths provided for removal")
	}

	repoRoot := getEffectiveRepoRoot()
	db, err := store.OpenAiAnalysisSplitDB("", repoRoot)
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.RunAiRemove.open_db")
	}
	defer db.Close()

	taskUuid, taskErr := resolveOrEnsureTaskUuid(db, taskId)
	if taskErr != nil {
		return taskErr
	}

	task, ensureErr := ensureTaskRecord(db, taskUuid, reason)
	if ensureErr != nil {
		return ensureErr
	}

	return executeFileRemovals(db, task, repoRoot, taskUuid, reason, filePaths, isDryRun)
}

func resolveOrEnsureTaskUuid(db *store.AiAnalysisSplitDB, explicit string) (string, error) {
	clean := strings.TrimSpace(explicit)
	hasClean := clean != ""
	if hasClean {
		return clean, nil
	}

	return resolveActiveTaskId(db)
}

func ensureTaskRecord(db *store.AiAnalysisSplitDB, taskUuid, reason string) (*store.AiTask, error) {
	existing, err := db.GetTask(taskUuid)
	isFound := err == nil && existing != nil
	if isFound {
		return existing, nil
	}

	newTask := &store.AiTask{
		TaskUuid:  taskUuid,
		Goal:      "Zero-loss file removal: " + reason,
		Category:  "cleanup",
		Reasoning: reason,
		Status:    "in_progress",
		IsActive:  true,
	}
	id, createErr := db.CreateTask(newTask)
	if createErr != nil {
		return nil, apperror.WrapSimple(createErr, "cmdai.ensureTaskRecord")
	}
	newTask.AiTaskId = id

	return newTask, nil
}

func executeFileRemovals(db *store.AiAnalysisSplitDB, task *store.AiTask, repoRoot, taskId, reason string, files []string, isDryRun bool) error {
	vaultDir := fspath.RepoTempDir("removed", taskId)
	manifestPath := filepath.Join(vaultDir, "manifest.json")
	manifest := loadOrCreateManifest(manifestPath, taskId, reason)

	for _, p := range files {
		entry, err := processSingleRemoval(db, task, repoRoot, vaultDir, p, reason, isDryRun)
		if err != nil {
			return err
		}
		if !isDryRun && entry != nil {
			manifest.Files = appendOrUpdateManifestFile(manifest.Files, *entry)
		}
	}

	if !isDryRun {
		return saveManifest(manifestPath, manifest)
	}

	return nil
}

func processSingleRemoval(db *store.AiAnalysisSplitDB, task *store.AiTask, repoRoot, vaultDir, inputPath, reason string, isDryRun bool) (*RemovalManifestFile, error) {
	relPath, absPath, pathErr := resolveFileRelAndAbs(repoRoot, inputPath)
	if pathErr != nil {
		return nil, pathErr
	}

	info, statErr := os.Stat(absPath)
	if statErr != nil {
		return nil, apperror.WrapSimple(statErr, "cmdai.processSingleRemoval.stat: "+absPath)
	}
	if info.IsDir() {
		return nil, apperror.NewValidationError("target is a directory, only files can be removed: " + relPath)
	}

	if isDryRun {
		return dryRunSingleFile(absPath, relPath, info.Size())
	}

	return stageAndRemoveFile(db, task, absPath, relPath, vaultDir, reason, info.Mode(), info.Size())
}

func dryRunSingleFile(absPath, relPath string, size int64) (*RemovalManifestFile, error) {
	data, readErr := os.ReadFile(absPath)
	if readErr != nil {
		return nil, apperror.WrapSimple(readErr, "cmdai.dryRunSingleFile.read")
	}

	shaHex := computeSha256(data)
	fmt.Printf("[DRY RUN] Would stage and remove: %s (SHA-256: %s, %d bytes)\n", relPath, shaHex, size)

	return &RemovalManifestFile{
		RelPath:      relPath,
		SizeBytes:    size,
		BeforeSha256: shaHex,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func stageAndRemoveFile(db *store.AiAnalysisSplitDB, task *store.AiTask, absPath, relPath, vaultDir, reason string, mode os.FileMode, size int64) (*RemovalManifestFile, error) {
	data, readErr := os.ReadFile(absPath)
	if readErr != nil {
		return nil, apperror.WrapSimple(readErr, "cmdai.stageAndRemoveFile.read")
	}

	shaHex := computeSha256(data)
	backupPath, backupErr := backupFileToVault(vaultDir, relPath, data, mode, shaHex)
	if backupErr != nil {
		return nil, backupErr
	}

	lineCount := countLines(data)
	dbErr := recordFileInDatabase(db, task.AiTaskId, relPath, absPath, shaHex, backupPath, reason, lineCount)
	if dbErr != nil {
		return nil, dbErr
	}

	rmErr := os.Remove(absPath)
	if rmErr != nil {
		return nil, apperror.WrapSimple(rmErr, "cmdai.stageAndRemoveFile.unlink: "+absPath)
	}

	return &RemovalManifestFile{
		RelPath:      relPath,
		SizeBytes:    size,
		BeforeSha256: shaHex,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		BackupPath:   backupPath,
	}, nil
}

func backupFileToVault(vaultDir, relPath string, data []byte, mode os.FileMode, expectedSha string) (string, error) {
	backupPath := filepath.Join(vaultDir, filepath.FromSlash(relPath))
	dirErr := os.MkdirAll(filepath.Dir(backupPath), 0755)
	if dirErr != nil {
		return "", apperror.WrapSimple(dirErr, "cmdai.backupFileToVault.mkdir")
	}

	writeErr := os.WriteFile(backupPath, data, mode)
	if writeErr != nil {
		return "", apperror.WrapSimple(writeErr, "cmdai.backupFileToVault.write")
	}

	verifyData, verifyErr := os.ReadFile(backupPath)
	if verifyErr != nil {
		return "", apperror.WrapSimple(verifyErr, "cmdai.backupFileToVault.verify_read")
	}

	verifySha := computeSha256(verifyData)
	isMatch := verifySha == expectedSha
	if !isMatch {
		return "", apperror.NewValidationError("checksum mismatch in staged backup: " + backupPath)
	}

	return backupPath, nil
}

func recordFileInDatabase(db *store.AiAnalysisSplitDB, taskId int64, relPath, absPath, sha256Hex, backupPath, reason string, lines int) error {
	record := &store.AiTaskFile{
		AiTaskId:     taskId,
		RelPath:      relPath,
		AbsPath:      absPath,
		Action:       "remove",
		BeforeSha256: sha256Hex,
		BackupPath:   backupPath,
		LineCount:    lines,
		Reasoning:    reason,
		IsModified:   false,
		IsRemoved:    true,
	}

	_, err := db.RecordFile(record)
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.recordFileInDatabase")
	}

	return nil
}

func resolveFileRelAndAbs(repoRoot, inputPath string) (string, string, error) {
	var absPath string
	isAbs := filepath.IsAbs(inputPath)
	if isAbs {
		absPath = filepath.Clean(inputPath)
	} else {
		absPath = filepath.Clean(filepath.Join(repoRoot, inputPath))
	}

	rel, err := filepath.Rel(repoRoot, absPath)
	isOutside := err != nil || strings.HasPrefix(rel, "..")
	if isOutside {
		return "", "", apperror.NewValidationError("target file is outside repository root: " + inputPath)
	}

	return filepath.ToSlash(rel), absPath, nil
}

func getEffectiveRepoRoot() string {
	root, err := ResolveRepoRoot()
	hasRoot := err == nil && root != ""
	if hasRoot {
		return root
	}

	cwd, cwdErr := os.Getwd()
	hasCwd := cwdErr == nil && cwd != ""
	if hasCwd {
		return cwd
	}

	return "."
}

func computeSha256(data []byte) string {
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])
}

func countLines(data []byte) int {
	isEmpty := len(data) == 0
	if isEmpty {
		return 0
	}

	var count int
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}

	return count + 1
}

func tryLoadExistingManifest(manifestPath string) *RemovalManifest {
	data, err := os.ReadFile(manifestPath)
	if err != nil || len(data) == 0 {
		return nil
	}

	var m RemovalManifest
	if jsonErr := json.Unmarshal(data, &m); jsonErr != nil {
		return nil
	}

	m.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	return &m
}

func loadOrCreateManifest(manifestPath, taskId, reason string) *RemovalManifest {
	if existing := tryLoadExistingManifest(manifestPath); existing != nil {
		return existing
	}

	return &RemovalManifest{
		TaskId:    taskId,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Reason:    reason,
		Files:     []RemovalManifestFile{},
	}
}

func appendOrUpdateManifestFile(files []RemovalManifestFile, entry RemovalManifestFile) []RemovalManifestFile {
	for i, f := range files {
		isSame := f.RelPath == entry.RelPath
		if isSame {
			files[i] = entry
			return files
		}
	}

	return append(files, entry)
}

func saveManifest(manifestPath string, manifest *RemovalManifest) error {
	manifest.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.saveManifest.marshal")
	}

	writeErr := os.WriteFile(manifestPath, data, 0644)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "cmdai.saveManifest.write")
	}

	return nil
}

func renderRemoveSuccess(taskId string, files []string, isDryRun, isJson bool) error {
	status := "staged_and_removed"
	if isDryRun {
		status = "dry_run"
	}

	if isJson {
		resp := AiRemoveResponse{
			Status:     status,
			TaskId:     taskId,
			TotalFiles: len(files),
			VaultDir:   fspath.RepoTempDir("removed", taskId),
		}
		return printJsonOutput(resp)
	}

	if !isDryRun {
		fmt.Printf("Safely staged and removed %d file(s) for task %s (vault: %s)\n",
			len(files), taskId, fspath.RepoTempDir("removed", taskId))
	}

	return nil
}
