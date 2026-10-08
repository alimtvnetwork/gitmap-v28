// Package cmdai — ai_transfer.go: Cross-system export and import for AI analysis Split-DB datasets.
package cmdai

import (
	"archive/zip"
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

// AiTransferDataset packages analysis tasks, files, and lines for portable transfer.
type AiTransferDataset struct {
	Version    string             `json:"version"`
	ExportedAt string             `json:"exportedAt"`
	TaskCount  int                `json:"taskCount"`
	FileCount  int                `json:"fileCount"`
	LineCount  int                `json:"lineCount"`
	Tasks      []store.AiTask     `json:"tasks"`
	Files      []store.AiTaskFile `json:"files"`
	Lines      []store.AiTaskLine `json:"lines"`
}

// TransferExportOptions holds parsed CLI flags for the export command.
type TransferExportOptions struct {
	Format string
	Output string
	TaskId string
}

// TransferImportOptions holds parsed CLI flags for the import command.
type TransferImportOptions struct {
	FilePath  string
	IsMerge   bool
	IsReplace bool
}

var (
	exportOptions TransferExportOptions
	importOptions TransferImportOptions

	aiAnalysisExportCmd = &cobra.Command{
		Use:     "export",
		Aliases: []string{"exp"},
		Short:   "Export AI analysis sessions to JSON, ZIP, or SQLite snapshot",
		RunE:    runAnalysisExportCommand,
	}

	aiAnalysisImportCmd = &cobra.Command{
		Use:     "import",
		Aliases: []string{"imp"},
		Short:   "Import AI analysis sessions into local Split-DB without ID collisions",
		RunE:    runAnalysisImportCommand,
	}
)

func runAnalysisExportCommand(cmd *cobra.Command, args []string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.export.open_db")
	}
	defer db.Close()

	return ExecuteAiAnalysisExport(db, exportOptions)
}

func runAnalysisImportCommand(cmd *cobra.Command, args []string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.import.open_db")
	}
	defer db.Close()

	return ExecuteAiAnalysisImport(db, importOptions)
}

// ExecuteAiAnalysisExport serializes and saves the dataset according to target format.
func ExecuteAiAnalysisExport(db *store.AiAnalysisSplitDB, opts TransferExportOptions) error {
	cleanFormat := strings.ToLower(strings.TrimSpace(opts.Format))
	isSqlite := cleanFormat == "sqlite" || cleanFormat == "db"
	if isSqlite {
		return exportAsSqliteSnapshot(db, opts.Output)
	}

	dataset, err := collectTransferDataset(db, opts.TaskId)
	if err != nil {
		return err
	}

	isZip := cleanFormat == "zip"
	if isZip {
		return exportAsZipBundle(dataset, opts.Output)
	}

	return exportAsJsonBundle(dataset, opts.Output)
}

func collectTransferDataset(db *store.AiAnalysisSplitDB, taskId string) (AiTransferDataset, error) {
	tasks, err := queryExportTasks(db, taskId)
	if err != nil {
		return AiTransferDataset{}, err
	}

	var allFiles []store.AiTaskFile
	var allLines []store.AiTaskLine

	for _, t := range tasks {
		files, lines, fetchErr := collectTaskFilesAndLines(db, t.AiTaskId)
		if fetchErr != nil {
			return AiTransferDataset{}, fetchErr
		}
		allFiles = append(allFiles, files...)
		allLines = append(allLines, lines...)
	}

	return AiTransferDataset{
		Version:    "1.0.0",
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		TaskCount:  len(tasks),
		FileCount:  len(allFiles),
		LineCount:  len(allLines),
		Tasks:      tasks,
		Files:      allFiles,
		Lines:      allLines,
	}, nil
}

func queryExportTasks(db *store.AiAnalysisSplitDB, taskId string) ([]store.AiTask, error) {
	trimmed := strings.TrimSpace(taskId)
	if trimmed == "" {
		return db.ListTasks(false, 10000)
	}

	task, err := db.GetTask(trimmed)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.export.get_task")
	}

	return []store.AiTask{*task}, nil
}

func collectTaskFilesAndLines(db *store.AiAnalysisSplitDB, taskId int64) ([]store.AiTaskFile, []store.AiTaskLine, error) {
	files, err := db.GetTaskFiles(taskId)
	if err != nil {
		return nil, nil, apperror.WrapSimple(err, "cmdai.export.get_files")
	}

	var allLines []store.AiTaskLine
	sanitizedFiles := make([]store.AiTaskFile, 0, len(files))

	for _, f := range files {
		cleanFile := sanitizeFileForExport(f)
		sanitizedFiles = append(sanitizedFiles, cleanFile)

		lines, lineErr := db.GetTaskLines(f.AiTaskFileId)
		if lineErr != nil {
			return nil, nil, apperror.WrapSimple(lineErr, "cmdai.export.get_lines")
		}
		allLines = append(allLines, lines...)
	}

	return sanitizedFiles, allLines, nil
}

func sanitizeFileForExport(f store.AiTaskFile) store.AiTaskFile {
	f.RelPath = filepath.ToSlash(filepath.Clean(f.RelPath))
	f.AbsPath = ""
	f.BackupPath = ""

	return f
}

func exportAsJsonBundle(dataset AiTransferDataset, outPath string) error {
	dest := resolveExportPath(outPath, "ai-analysis-export.json")
	data, err := json.MarshalIndent(dataset, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.export.marshal_json")
	}

	if err := os.WriteFile(dest, data, 0644); err != nil {
		return apperror.WrapSimple(err, "cmdai.export.write_json")
	}

	fmt.Printf("Exported %d tasks, %d files, and %d lines to JSON: %s\n", dataset.TaskCount, dataset.FileCount, dataset.LineCount, dest)

	return nil
}

func exportAsZipBundle(dataset AiTransferDataset, outPath string) error {
	dest := resolveExportPath(outPath, "ai-analysis-export.zip")
	zipFile, err := os.Create(dest)
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.export.create_zip")
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	if err := writeZipDataset(w, dataset); err != nil {
		return err
	}

	fmt.Printf("Exported %d tasks, %d files, and %d lines to ZIP: %s\n", dataset.TaskCount, dataset.FileCount, dataset.LineCount, dest)

	return nil
}

func writeZipDataset(w *zip.Writer, dataset AiTransferDataset) error {
	entry, err := w.Create("dataset.json")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.export.create_zip_entry")
	}

	encoder := json.NewEncoder(entry)
	encoder.SetIndent("", "  ")

	return encoder.Encode(dataset)
}

func exportAsSqliteSnapshot(db *store.AiAnalysisSplitDB, outPath string) error {
	dest := resolveExportPath(outPath, "ai-analysis-export.db")
	srcPath := db.Path

	_ = store.ExecWrapper(db.Conn(), "PRAGMA wal_checkpoint(TRUNCATE);")

	in, err := os.Open(srcPath)
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.export.open_sqlite_source")
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.export.create_sqlite_target")
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return apperror.WrapSimple(err, "cmdai.export.copy_sqlite")
	}

	fmt.Printf("Exported atomic SQLite snapshot to %s\n", dest)

	return nil
}

func resolveExportPath(outPath, defaultName string) string {
	hasPath := strings.TrimSpace(outPath) != ""
	if hasPath {
		return strings.TrimSpace(outPath)
	}

	return defaultName
}

// ExecuteAiAnalysisImport imports tasks, files, and lines from an external bundle.
func ExecuteAiAnalysisImport(db *store.AiAnalysisSplitDB, opts TransferImportOptions) error {
	isEmpty := strings.TrimSpace(opts.FilePath) == ""
	if isEmpty {
		return apperror.NewValidation("cmdai.import", "E3014", "flag --file is required")
	}

	dataset, err := loadImportDataset(opts.FilePath)
	if err != nil {
		return err
	}

	return applyImportDataset(db, dataset, opts.IsReplace)
}

func loadImportDataset(filePath string) (AiTransferDataset, error) {
	cleanPath := filepath.Clean(filePath)
	ext := strings.ToLower(filepath.Ext(cleanPath))

	isZip := ext == ".zip"
	if isZip {
		return loadDatasetFromZip(cleanPath)
	}

	isSqlite := ext == ".db" || ext == ".sqlite"
	if isSqlite {
		return loadDatasetFromSqlite(cleanPath)
	}

	return loadDatasetFromJson(cleanPath)
}

func loadDatasetFromJson(filePath string) (AiTransferDataset, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return AiTransferDataset{}, apperror.WrapSimple(err, "cmdai.import.read_file")
	}

	var ds AiTransferDataset
	if err := json.Unmarshal(data, &ds); err != nil {
		return AiTransferDataset{}, apperror.WrapSimple(err, "cmdai.import.unmarshal_json")
	}

	return ds, nil
}

func loadDatasetFromZip(filePath string) (AiTransferDataset, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return AiTransferDataset{}, apperror.WrapSimple(err, "cmdai.import.open_zip")
	}
	defer r.Close()

	for _, f := range r.File {
		isDataset := f.Name == "dataset.json" || f.Name == "records.json"
		if isDataset {
			return readZipEntry(f)
		}
	}

	return AiTransferDataset{}, apperror.NewValidation("cmdai.import", "E3015", "zip archive missing dataset.json")
}

func readZipEntry(f *zip.File) (AiTransferDataset, error) {
	rc, err := f.Open()
	if err != nil {
		return AiTransferDataset{}, apperror.WrapSimple(err, "cmdai.import.open_zip_entry")
	}
	defer rc.Close()

	var ds AiTransferDataset
	if err := json.NewDecoder(rc).Decode(&ds); err != nil {
		return AiTransferDataset{}, apperror.WrapSimple(err, "cmdai.import.decode_zip_entry")
	}

	return ds, nil
}

func loadDatasetFromSqlite(filePath string) (AiTransferDataset, error) {
	srcDb, err := store.OpenAiAnalysisSplitDBAt(filePath)
	if err != nil {
		return AiTransferDataset{}, apperror.WrapSimple(err, "cmdai.import.open_sqlite")
	}
	defer srcDb.Close()

	return collectTransferDataset(srcDb, "")
}

func applyImportDataset(db *store.AiAnalysisSplitDB, dataset AiTransferDataset, isReplace bool) error {
	taskCount := 0
	fileCount := 0
	lineCount := 0

	for _, task := range dataset.Tasks {
		importedTaskId, err := importSingleTask(db, task, isReplace)
		if err != nil {
			return err
		}
		taskCount++

		fCount, lCount, err := importTaskFilesAndLines(db, task.AiTaskId, importedTaskId, dataset)
		if err != nil {
			return err
		}
		fileCount += fCount
		lineCount += lCount
	}

	fmt.Printf("Successfully imported %d tasks, %d files, and %d lines into AI analysis Split-DB.\n", taskCount, fileCount, lineCount)

	return nil
}

func resolveTaskUuidOrReplaceExisting(db *store.AiAnalysisSplitDB, existing *store.AiTask, task *store.AiTask, isReplace bool) {
	if existing == nil {
		return
	}

	if isReplace {
		_ = deleteTasksAndLines(db, []int64{existing.AiTaskId})

		return
	}

	task.TaskUuid = generateNonCollidingUuid(task.TaskUuid)
}

func importSingleTask(db *store.AiAnalysisSplitDB, task store.AiTask, isReplace bool) (int64, error) {
	existing, _ := db.GetTask(task.TaskUuid)
	resolveTaskUuidOrReplaceExisting(db, existing, &task, isReplace)

	newTask := task
	newTask.AiTaskId = 0

	return db.CreateTask(&newTask)
}

func generateNonCollidingUuid(baseUuid string) string {
	ts := time.Now().UnixNano() % 100000

	return fmt.Sprintf("%s-imp-%d", baseUuid, ts)
}

func importTaskFilesAndLines(db *store.AiAnalysisSplitDB, oldTaskId, newTaskId int64, ds AiTransferDataset) (int, int, error) {
	importedFiles := 0
	importedLines := 0

	for _, f := range ds.Files {
		isMatchingTask := f.AiTaskId == oldTaskId
		if !isMatchingTask {
			continue
		}

		oldFileId := f.AiTaskFileId
		newFile := f
		newFile.AiTaskFileId = 0
		newFile.AiTaskId = newTaskId

		newFileId, err := db.RecordFile(&newFile)
		if err != nil {
			return 0, 0, err
		}
		importedFiles++

		lCount, lErr := importFileLines(db, oldFileId, newFileId, ds)
		if lErr != nil {
			return 0, 0, lErr
		}
		importedLines += lCount
	}

	return importedFiles, importedLines, nil
}

func importFileLines(db *store.AiAnalysisSplitDB, oldFileId, newFileId int64, ds AiTransferDataset) (int, error) {
	importedLines := 0
	for _, l := range ds.Lines {
		isMatchingFile := l.AiTaskFileId == oldFileId
		if !isMatchingFile {
			continue
		}

		newLine := l
		newLine.AiTaskLineId = 0
		newLine.AiTaskFileId = newFileId

		if _, err := db.RecordLine(&newLine); err != nil {
			return 0, err
		}
		importedLines++
	}

	return importedLines, nil
}

func init() {
	initTransferFlags()

	AiAnalysisCmd.AddCommand(aiAnalysisExportCmd)
	AiAnalysisCmd.AddCommand(aiAnalysisImportCmd)
}

func initTransferFlags() {
	aiAnalysisExportCmd.Flags().StringVarP(&exportOptions.Format, "format", "f", "json", "Export format: json, zip, sqlite")
	aiAnalysisExportCmd.Flags().StringVarP(&exportOptions.Output, "out", "o", "", "Destination output file path")
	aiAnalysisExportCmd.Flags().StringVarP(&exportOptions.TaskId, "task", "t", "", "Filter export by specific task UUID")

	aiAnalysisImportCmd.Flags().StringVarP(&importOptions.FilePath, "file", "f", "", "Source file to import from (.json, .zip, .sqlite)")
	aiAnalysisImportCmd.Flags().BoolVar(&importOptions.IsMerge, "merge", true, "Merge imported tasks without overwriting existing data")
	aiAnalysisImportCmd.Flags().BoolVar(&importOptions.IsReplace, "replace", false, "Replace existing tasks when matching UUIDs collide")
}
