package cmdai

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export AI analysis sessions and line telemetry to JSON, ZIP, or DB snapshot",
		RunE:  runExportCmd,
	}

	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import AI analysis sessions from an exported JSON or DB snapshot",
		RunE:  runImportCmd,
	}

	exportFlags AiExportOptions
	importFlags AiImportOptions
)

func runExportCmd(cmd *cobra.Command, args []string) error {
	return ExecuteAiExport(exportFlags)
}

func runImportCmd(cmd *cobra.Command, args []string) error {
	return ExecuteAiImport(importFlags)
}

// ExecuteAiExport handles multi-format export of analysis sessions.
func ExecuteAiExport(opts AiExportOptions) error {
	format := strings.ToLower(string(opts.Format))
	if format == "" {
		format = string(ExportFormatJson)
	}

	isDbFormat := format == string(ExportFormatDb)
	if isDbFormat {
		return exportDatabaseSnapshot(opts)
	}

	db, err := OpenAiAnalysisSplitDB(opts.RepoRoot)
	if err != nil {
		return err
	}
	defer db.Close()

	bundle, bundleErr := buildExportBundle(db, opts)
	if bundleErr != nil {
		return bundleErr
	}

	isZipFormat := format == string(ExportFormatZip)
	if isZipFormat {
		return exportZipBundle(bundle, opts.Output)
	}

	return exportJsonBundle(bundle, opts.Output)
}

func buildExportBundle(db *sql.DB, opts AiExportOptions) (AiExportBundle, error) {
	tasks, lines, err := collectTasksAndLines(db, opts.TaskId)
	if err != nil {
		return AiExportBundle{}, err
	}

	slug := "default"
	hasRepo := opts.RepoRoot != ""
	if hasRepo {
		slug = store.SanitizeSlug(filepath.Base(opts.RepoRoot))
	}

	return AiExportBundle{
		ExportedAt: time.Now().UTC(),
		RepoSlug:   slug,
		Tasks:      tasks,
		Lines:      lines,
	}, nil
}

func collectTasksAndLines(db *sql.DB, taskID string) ([]AiAnalysisTask, []AiAnalysisLine, error) {
	cleanID := strings.TrimSpace(taskID)
	hasTaskID := cleanID != ""
	if hasTaskID {
		return collectSingleTaskAndLines(db, cleanID)
	}

	return collectAllTasksAndLines(db)
}

func collectSingleTaskAndLines(db *sql.DB, taskID string) ([]AiAnalysisTask, []AiAnalysisLine, error) {
	task, err := GetAiAnalysisTask(db, taskID)
	if err != nil {
		return nil, nil, err
	}

	lines, linesErr := GetAiAnalysisLinesForTask(db, taskID)
	if linesErr != nil {
		return nil, nil, linesErr
	}

	return []AiAnalysisTask{*task}, lines, nil
}

func collectAllTasksAndLines(db *sql.DB) ([]AiAnalysisTask, []AiAnalysisLine, error) {
	tasks, err := ListAiAnalysisTasks(db, 1000, 0)
	if err != nil {
		return nil, nil, err
	}

	var allLines []AiAnalysisLine
	for _, t := range tasks {
		lines, linesErr := GetAiAnalysisLinesForTask(db, t.TaskId)
		if linesErr != nil {
			return nil, nil, linesErr
		}
		allLines = append(allLines, lines...)
	}

	return tasks, allLines, nil
}

func exportJsonBundle(bundle AiExportBundle, outputPath string) error {
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.exportJsonBundle.marshal")
	}

	hasOutput := strings.TrimSpace(outputPath) != ""
	if !hasOutput {
		fmt.Println(string(data))
		return nil
	}

	writeErr := os.WriteFile(outputPath, data, 0644)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "cmdai.exportJsonBundle.write")
	}

	fmt.Printf("Exported %d tasks and %d lines to %q\n", len(bundle.Tasks), len(bundle.Lines), outputPath)
	return nil
}

func exportZipBundle(bundle AiExportBundle, outputPath string) error {
	target := outputPath
	if strings.TrimSpace(target) == "" {
		target = "ai_analysis_export.zip"
	}

	zipFile, err := os.Create(target)
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.exportZipBundle.create")
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	if err := writeZipManifest(w, bundle); err != nil {
		return err
	}
	if err := writeZipRecords(w, bundle); err != nil {
		return err
	}

	fmt.Printf("Exported ZIP bundle with %d tasks and %d lines to %q\n", len(bundle.Tasks), len(bundle.Lines), target)
	return nil
}

func writeZipManifest(w *zip.Writer, bundle AiExportBundle) error {
	f, err := w.Create("manifest.json")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.writeZipManifest.create")
	}

	manifest := map[string]any{
		"exportedAt": bundle.ExportedAt.Format(time.RFC3339),
		"repoSlug":   bundle.RepoSlug,
		"taskCount":  len(bundle.Tasks),
		"lineCount":  len(bundle.Lines),
	}
	return json.NewEncoder(f).Encode(manifest)
}

func writeZipRecords(w *zip.Writer, bundle AiExportBundle) error {
	f, err := w.Create("records.json")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.writeZipRecords.create")
	}

	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.writeZipRecords.marshal")
	}

	_, writeErr := f.Write(data)
	return writeErr
}

func exportDatabaseSnapshot(opts AiExportOptions) error {
	target := opts.Output
	if strings.TrimSpace(target) == "" {
		target = "ai_analysis_snapshot.db"
	}

	srcPath := store.ResolveAiAnalysisSplitDbPath("", opts.RepoRoot)
	db, err := OpenAiAnalysisSplitDBAt(srcPath)
	if err != nil {
		return err
	}

	_ = store.ExecWrapper(db, "PRAGMA wal_checkpoint(TRUNCATE);")
	_ = db.Close()

	return copyDatabaseFile(srcPath, target)
}

func copyDatabaseFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.copyDatabaseFile.open")
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.copyDatabaseFile.create")
	}
	defer out.Close()

	_, copyErr := io.Copy(out, in)
	if copyErr != nil {
		return apperror.WrapSimple(copyErr, "cmdai.copyDatabaseFile.copy")
	}

	fmt.Printf("Exported atomic SQLite snapshot to %q\n", dst)
	return nil
}

// ExecuteAiImport handles ingestion of analysis data into Split-DB.
func ExecuteAiImport(opts AiImportOptions) error {
	isEmptyInput := strings.TrimSpace(opts.InputPath) == ""
	if isEmptyInput {
		return apperror.NewValidation("cmd.ai.import", "E3011", "flag --input is required")
	}

	data, readErr := os.ReadFile(opts.InputPath)
	if readErr != nil {
		return apperror.WrapSimple(readErr, "cmdai.ExecuteAiImport.read")
	}

	var bundle AiExportBundle
	unmarshalErr := json.Unmarshal(data, &bundle)
	if unmarshalErr != nil {
		return apperror.WrapSimple(unmarshalErr, "cmdai.ExecuteAiImport.unmarshal")
	}

	db, err := OpenAiAnalysisSplitDB(opts.RepoRoot)
	if err != nil {
		return err
	}
	defer db.Close()

	return importBundleIntoDb(db, bundle)
}

func importBundleIntoDb(db *sql.DB, bundle AiExportBundle) error {
	tx, err := db.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.importBundleIntoDb.begin")
	}
	defer tx.Rollback()

	for _, task := range bundle.Tasks {
		if err := upsertTaskInTx(tx, task); err != nil {
			return err
		}
	}

	for _, line := range bundle.Lines {
		if err := executeInsertLine(tx, line); err != nil {
			return err
		}
		_ = executeCounterUpdate(tx, line.TaskId)
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return apperror.WrapSimple(commitErr, "cmdai.importBundleIntoDb.commit")
	}

	fmt.Printf("Successfully imported %d tasks and %d lines into AI analysis Split-DB\n", len(bundle.Tasks), len(bundle.Lines))
	return nil
}

func upsertTaskInTx(tx *sql.Tx, task AiAnalysisTask) error {
	startedAt := formatTimestamp(task.StartedAt)
	createdAt := formatTimestamp(task.CreatedAt)
	var completedAt *string
	hasCompleted := task.CompletedAt != nil
	if hasCompleted {
		val := formatTimestamp(*task.CompletedAt)
		completedAt = &val
	}

	query := `INSERT OR REPLACE INTO AiAnalysisTask (
		task_id, repo_path, task_description, model_name, status,
		started_at, completed_at, total_files_read, total_files_modified,
		total_files_deleted, reasoning_summary, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res := store.ExecWrapper(tx, query,
		task.TaskId, task.RepoPath, task.TaskDescription, task.ModelName, string(task.Status),
		startedAt, completedAt, task.TotalFilesRead, task.TotalFilesModified,
		task.TotalFilesDeleted, task.ReasoningSummary, createdAt,
	)
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.upsertTaskInTx")
	}

	return nil
}

func initExportCommands() {
	exportCmd.Flags().StringVarP(&exportFlags.TaskId, "task", "t", "", "Filter export by specific task ID")
	exportCmd.Flags().StringVarP((*string)(&exportFlags.Format), "format", "f", "json", "Export format: json, zip, db")
	exportCmd.Flags().StringVarP(&exportFlags.Output, "output", "o", "", "Destination output file path")

	importCmd.Flags().StringVarP(&importFlags.InputPath, "input", "i", "", "Path to imported JSON dataset")

	AiCmd.AddCommand(exportCmd)
	AiCmd.AddCommand(importCmd)
}
