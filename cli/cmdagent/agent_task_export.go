package cmdagent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

type TaskExportOptions struct {
	TaskId string
	Format string
	Out    string
}

type TaskImportOptions struct {
	File   string
	Format string
}

type TaskExportBundle struct {
	ParentTask types.ParentTask `json:"parentTask" yaml:"parentTask"`
	Subtasks   []types.Subtask  `json:"subtasks" yaml:"subtasks"`
}

var (
	taskExportOpts TaskExportOptions
	taskImportOpts TaskImportOptions

	taskExportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export parent tasks and subtasks to JSON or YAML",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && len(taskExportOpts.TaskId) == 0 {
				taskExportOpts.TaskId = args[0]
			}

			return toError(RunTaskExport(taskExportOpts))
		},
	}

	taskImportCmd = &cobra.Command{
		Use:   "import [file]",
		Short: "Import tasks and subtasks from JSON or YAML file into SQLite task DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && len(taskImportOpts.File) == 0 {
				taskImportOpts.File = args[0]
			}

			return toError(RunTaskImport(taskImportOpts))
		},
	}
)

func initTaskExportFlags() {
	taskExportCmd.Flags().StringVarP(&taskExportOpts.TaskId, "task-id", "t", "", "Parent task ID or slug (optional, default active/all)")
	taskExportCmd.Flags().StringVarP(&taskExportOpts.Format, "format", "f", "json", "Output format (json, yaml)")
	taskExportCmd.Flags().StringVarP(&taskExportOpts.Out, "out", "o", "", "Output file path (optional, stdout if omitted)")

	taskImportCmd.Flags().StringVarP(&taskImportOpts.File, "file", "f", "", "Input file path (required)")
	taskImportCmd.Flags().StringVar(&taskImportOpts.Format, "format", "", "Format (json, yaml, auto-detect if omitted)")
}

// RunTaskExport serializes task and subtask records into JSON or YAML.
func RunTaskExport(opts TaskExportOptions) *appfault.AppError {
	format := resolveExportFormat(opts.Format, opts.Out)
	bundles, fetchErr := CollectExportBundles(opts.TaskId)
	if fetchErr != nil {
		return fetchErr
	}

	if len(bundles) == 0 {
		return appfault.NewValidationError("no tasks found to export")
	}

	var payload any = bundles
	if len(bundles) == 1 {
		payload = bundles[0]
	}

	data, serErr := serializePayload(payload, format)
	if serErr != nil {
		return serErr
	}

	return writeExportData(data, opts.Out, len(bundles), format)
}

// CollectExportBundles retrieves task and subtask bundles for export.
func CollectExportBundles(taskId string) ([]TaskExportBundle, *appfault.AppError) {
	trimmed := strings.TrimSpace(taskId)
	if len(trimmed) > 0 {
		return collectSingleTaskBundle(trimmed)
	}

	return collectAllTaskBundles()
}

func collectSingleTaskBundle(target string) ([]TaskExportBundle, *appfault.AppError) {
	masterDb := store.ResolveMasterAgentDbPath("")
	task, _ := store.FindParentTaskBySlug(masterDb, target)

	if task == nil {
		dbPath, findErr := store.FindTaskDbPath("", target)
		if findErr == nil {
			summary, sumErr := store.GetTaskStatus(dbPath, target)
			if sumErr == nil && len(summary.ParentTask.ParentTaskId) > 0 {
				task = &summary.ParentTask
			}
		}
	}

	if task == nil {
		all, _ := store.ListParentTasks(masterDb, 1000, "")
		for _, t := range all {
			if t.ParentTaskId == target || t.TaskSlug == target {
				task = &t
				break
			}
		}
	}

	if task == nil {
		return nil, appfault.NewValidationError(fmt.Sprintf("task %q not found", target))
	}

	subtasks, _ := store.ListSubtasks(task.RootDbPath, task.ParentTaskId)
	if len(subtasks) == 0 {
		subtasks, _ = store.ListSubtasks(task.RootDbPath, task.TaskSlug)
	}

	return []TaskExportBundle{
		{
			ParentTask: *task,
			Subtasks:   subtasks,
		},
	}, nil
}

func collectAllTaskBundles() ([]TaskExportBundle, *appfault.AppError) {
	masterDb := store.ResolveMasterAgentDbPath("")
	tasks, listErr := store.ListParentTasks(masterDb, 1000, "")
	if listErr != nil {
		return nil, listErr
	}

	var bundles []TaskExportBundle
	for _, t := range tasks {
		subtasks, _ := store.ListSubtasks(t.RootDbPath, t.ParentTaskId)
		if len(subtasks) == 0 {
			subtasks, _ = store.ListSubtasks(t.RootDbPath, t.TaskSlug)
		}

		bundles = append(bundles, TaskExportBundle{
			ParentTask: t,
			Subtasks:   subtasks,
		})
	}

	if len(bundles) == 0 {
		dbPath, findErr := store.FindTaskDbPath("", "")
		if findErr == nil {
			summary, sumErr := store.GetTaskStatus(dbPath, "")
			if sumErr == nil && len(summary.ParentTask.ParentTaskId) > 0 {
				bundles = append(bundles, TaskExportBundle{
					ParentTask: summary.ParentTask,
					Subtasks:   summary.Subtasks,
				})
			}
		}
	}

	return bundles, nil
}

func resolveExportFormat(format, outPath string) string {
	f := strings.ToLower(strings.TrimSpace(format))
	if f == "yaml" || f == "yml" {
		return "yaml"
	}

	if f == "json" {
		return "json"
	}

	ext := strings.ToLower(filepath.Ext(outPath))
	if ext == ".yaml" || ext == ".yml" {
		return "yaml"
	}

	return "json"
}

func serializePayload(payload any, format string) ([]byte, *appfault.AppError) {
	if format == "yaml" {
		data, err := yaml.Marshal(payload)
		if err != nil {
			return nil, appfault.WrapSimple(err, "yaml.Marshal")
		}

		return data, nil
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, appfault.WrapSimple(err, "json.Marshal")
	}

	return data, nil
}

func writeExportData(data []byte, outPath string, count int, format string) *appfault.AppError {
	trimmed := strings.TrimSpace(outPath)
	if len(trimmed) == 0 {
		fmt.Println(string(data))

		return nil
	}

	dir := filepath.Dir(trimmed)
	_ = os.MkdirAll(dir, 0755)

	writeErr := os.WriteFile(trimmed, data, 0644)
	if writeErr != nil {
		return appfault.WrapSimple(writeErr, "writeExportData")
	}

	fmt.Printf("Successfully exported %d task(s) to %s (%s)\n", count, trimmed, format)

	return nil
}

// RunTaskImport deserializes and registers tasks and subtasks into SQLite task DB.
func RunTaskImport(opts TaskImportOptions) *appfault.AppError {
	trimmedFile := strings.TrimSpace(opts.File)
	if len(trimmedFile) == 0 {
		return appfault.NewValidationError("import file path is required (--file)")
	}

	data, readErr := os.ReadFile(trimmedFile)
	if readErr != nil {
		return appfault.WrapSimple(readErr, "RunTaskImport.read")
	}

	bundles, parseErr := parseImportPayload(data, opts.Format, trimmedFile)
	if parseErr != nil {
		return parseErr
	}

	return registerImportedBundles(bundles, trimmedFile)
}

func parseImportPayload(data []byte, format, filePath string) ([]TaskExportBundle, *appfault.AppError) {
	f := strings.ToLower(strings.TrimSpace(format))
	if f == "" {
		ext := strings.ToLower(filepath.Ext(filePath))
		if ext == ".yaml" || ext == ".yml" {
			f = "yaml"
		} else if ext == ".json" {
			f = "json"
		}
	}

	switch f {
	case "yaml", "yml":
		return parseYamlBundles(data)
	case "json":
		return parseJsonBundles(data)
	default:
		bundles, err := parseJsonBundles(data)
		if err == nil && len(bundles) > 0 {
			return bundles, nil
		}

		return parseYamlBundles(data)
	}
}

func parseJsonBundles(data []byte) ([]TaskExportBundle, *appfault.AppError) {
	var list []TaskExportBundle
	if err := json.Unmarshal(data, &list); err == nil && len(list) > 0 {
		return list, nil
	}

	var single TaskExportBundle
	if err := json.Unmarshal(data, &single); err == nil && isBundlePopulated(single) {
		return []TaskExportBundle{single}, nil
	}

	return nil, appfault.NewValidationError("unable to parse JSON task import: expected task bundle or list of bundles")
}

func parseYamlBundles(data []byte) ([]TaskExportBundle, *appfault.AppError) {
	var list []TaskExportBundle
	if err := yaml.Unmarshal(data, &list); err == nil && len(list) > 0 {
		return list, nil
	}

	var single TaskExportBundle
	if err := yaml.Unmarshal(data, &single); err == nil && isBundlePopulated(single) {
		return []TaskExportBundle{single}, nil
	}

	return nil, appfault.NewValidationError("unable to parse YAML task import: expected task bundle or list of bundles")
}

func isBundlePopulated(b TaskExportBundle) bool {
	return len(b.ParentTask.ParentTaskId) > 0 || len(b.ParentTask.TaskName) > 0 || len(b.Subtasks) > 0
}

func registerImportedBundles(bundles []TaskExportBundle, sourceFile string) *appfault.AppError {
	importedTasks := 0
	importedSubtasks := 0

	for _, bundle := range bundles {
		task := bundle.ParentTask
		hasTask := len(task.ParentTaskId) > 0 || len(task.TaskName) > 0

		if hasTask {
			task = normalizeImportedTask(task)
			if upsertErr := store.UpsertParentTask(task.RootDbPath, task); upsertErr != nil {
				return upsertErr
			}

			masterDb := store.ResolveMasterAgentDbPath("")
			if regErr := store.RegisterParentTask(masterDb, task); regErr != nil {
				return regErr
			}

			importedTasks++
		}

		if len(bundle.Subtasks) > 0 {
			targetDb := task.RootDbPath
			targetParentId := task.ParentTaskId

			if len(targetDb) == 0 {
				firstSub := bundle.Subtasks[0]
				targetParentId = firstSub.ParentTaskId
				slug := store.SanitizeSlug(targetParentId)
				targetDb = store.ResolveTaskDbPath("", slug)
			}

			if addErr := store.AddSubtasks(targetDb, targetParentId, bundle.Subtasks); addErr != nil {
				return addErr
			}

			importedSubtasks += len(bundle.Subtasks)
		}
	}

	fmt.Printf("Successfully imported %d task(s) and %d subtask(s) from %s\n", importedTasks, importedSubtasks, sourceFile)

	return nil
}

func normalizeImportedTask(task types.ParentTask) types.ParentTask {
	now := time.Now().UTC().Format(time.RFC3339)

	if len(task.TaskSlug) == 0 {
		task.TaskSlug = store.SanitizeSlug(task.TaskName)
	}

	if len(task.ParentTaskId) == 0 {
		task.ParentTaskId = generateTaskId(task.TaskSlug)
	}

	if len(task.RunDirectory) == 0 {
		task.RunDirectory = store.ResolveParentTaskDir("", task.TaskSlug)
	}

	if len(task.RootDbPath) == 0 {
		task.RootDbPath = store.ResolveTaskDbPath("", task.TaskSlug)
	}

	if len(task.Status) == 0 {
		task.Status = "ACTIVE"
	}

	if task.TotalStepsBudget <= 0 {
		task.TotalStepsBudget = 300
	}

	if len(task.CreatedAt) == 0 {
		task.CreatedAt = now
	}

	if len(task.UpdatedAt) == 0 {
		task.UpdatedAt = now
	}

	return task
}
