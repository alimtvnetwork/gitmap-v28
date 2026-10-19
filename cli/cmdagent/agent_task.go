package cmdagent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

var (
	taskInitOpts   TaskInitOptions
	taskLsOpts     TaskLsOptions
	taskStatusOpts TaskStatusOptions

	taskInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize parent agent task and register in master DB",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && len(taskInitOpts.Name) == 0 {
				taskInitOpts.Name = args[0]
			}
			return toError(RunTaskInit(taskInitOpts))
		},
	}

	taskLsCmd = &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List active and completed parent tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunTaskLs(taskLsOpts))
		},
	}

	taskStatusCmd = &cobra.Command{
		Use:     "status [task-id]",
		Aliases: []string{"st"},
		Short:   "Show parent task status and subtask rollups",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && len(taskStatusOpts.TaskId) == 0 {
				taskStatusOpts.TaskId = args[0]
			}
			return toError(RunTaskStatus(taskStatusOpts))
		},
	}
)

type TaskInitOptions struct {
	Name   string
	Budget int
	Dir    string
}

type TaskLsOptions struct {
	Limit  int
	IsAll  bool
	Status string
	IsJson bool
}

type TaskStatusOptions struct {
	TaskId string
	IsJson bool
}

func initTaskCommands() {
	initTaskInitFlags()
	initTaskLsFlags()
	initTaskStatusFlags()
	initTaskSlugCommands()
	initTaskExportFlags()

	TaskCmd.AddCommand(taskInitCmd)
	TaskCmd.AddCommand(taskLsCmd)
	TaskCmd.AddCommand(taskStatusCmd)
	TaskCmd.AddCommand(taskExportCmd)
	TaskCmd.AddCommand(taskImportCmd)
}

func initTaskInitFlags() {
	taskInitCmd.Flags().StringVarP(&taskInitOpts.Name, "name", "n", "", "Task name or slug")
	taskInitCmd.Flags().IntVarP(&taskInitOpts.Budget, "budget", "b", 300, "Total steps budget")
	taskInitCmd.Flags().StringVarP(&taskInitOpts.Dir, "dir", "d", "", "Custom run directory")
}

func initTaskLsFlags() {
	taskLsCmd.Flags().IntVarP(&taskLsOpts.Limit, "limit", "l", 20, "Maximum tasks to display")
	taskLsCmd.Flags().BoolVarP(&taskLsOpts.IsAll, "all", "a", false, "Show all tasks without limit")
	taskLsCmd.Flags().StringVarP(&taskLsOpts.Status, "status", "s", "", "Filter tasks by status (ACTIVE, COMPLETED, FAILED)")
	taskLsCmd.Flags().BoolVar(&taskLsOpts.IsJson, "json", false, "Output results in JSON format")
}

func initTaskStatusFlags() {
	taskStatusCmd.Flags().StringVarP(&taskStatusOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
	taskStatusCmd.Flags().BoolVar(&taskStatusOpts.IsJson, "json", false, "Output results in JSON format")
}

// RunTaskInit creates and initializes a parent task hierarchy across Tier 1 and Tier 2.
func RunTaskInit(opts TaskInitOptions) *appfault.AppError {
	hasName := len(strings.TrimSpace(opts.Name)) > 0
	if !hasName {
		return appfault.NewValidationError("task name is required (--name)")
	}
	slug := store.SanitizeSlug(opts.Name)
	taskId := generateTaskId(slug)
	runDir, dbPath := resolveTaskInitPaths(opts.Dir, slug)
	task := buildParentTask(taskId, slug, opts.Name, runDir, dbPath, opts.Budget)

	return persistTaskInit(task)
}

func resolveTaskInitPaths(dir, slug string) (string, string) {
	hasDir := len(strings.TrimSpace(dir)) > 0
	if hasDir {
		runDir := filepath.ToSlash(dir)
		dbPath := filepath.ToSlash(filepath.Join(runDir, "agent-task.db"))
		return runDir, dbPath
	}
	runDir := store.ResolveParentTaskDir("", slug)
	dbPath := store.ResolveTaskDbPath("", slug)

	return runDir, dbPath
}

func generateTaskId(slug string) string {
	ts := time.Now().UTC().Format("20060102150405")

	return fmt.Sprintf("task-%s-%s", ts, slug)
}

func buildParentTask(taskId, slug, name, runDir, dbPath string, budget int) types.ParentTask {
	now := time.Now().UTC().Format(time.RFC3339)
	if budget <= 0 {
		budget = 300
	}

	return types.ParentTask{
		ParentTaskId:      taskId,
		TaskSlug:          slug,
		TaskName:          name,
		RunDirectory:      runDir,
		RootDbPath:        dbPath,
		Status:            "ACTIVE",
		TotalStepsBudget:  budget,
		CompletedSteps:    0,
		SpawnedAgentCount: 0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func persistTaskInit(task types.ParentTask) *appfault.AppError {
	upsertErr := store.UpsertParentTask(task.RootDbPath, task)
	hasUpsertErr := upsertErr != nil
	if hasUpsertErr {
		return upsertErr
	}
	masterDb := store.ResolveMasterAgentDbPath("")
	regErr := store.RegisterParentTask(masterDb, task)
	hasRegErr := regErr != nil
	if hasRegErr {
		return regErr
	}
	renderTaskInitResult(task)

	return nil
}

// RenderTaskInitResult outputs JSON metadata for initialized parent tasks.
func RenderTaskInitResult(task types.ParentTask) {
	renderTaskInitResult(task)
}

func renderTaskInitResult(task types.ParentTask) {
	cwd, _ := os.Getwd()
	relDbPath := toRelativePath(cwd, task.RootDbPath)
	masterDb := store.ResolveMasterAgentDbPath("")
	relMasterDb := toRelativePath(cwd, masterDb)

	rootDbInfo := map[string]any{
		"masterDbPath": relMasterDb,
		"taskDbPath":   relDbPath,
		"schema":       "ParentTaskRegistry, AgentRegistry, GlobalLifecycleMetrics",
	}

	queryCommands := []string{
		"gitmap agent task ls",
		"gitmap agent subtask list --task-id " + task.ParentTaskId,
		fmt.Sprintf("sqlite3 %s \"SELECT * FROM tasks;\"", relDbPath),
		fmt.Sprintf("sqlite3 %s \"SELECT * FROM parent_tasks WHERE parent_task_id = '%s';\"", relMasterDb, task.ParentTaskId),
	}

	out := map[string]any{
		"status":        "success",
		"parentTaskId":  task.ParentTaskId,
		"taskSlug":      task.TaskSlug,
		"taskName":      task.TaskName,
		"rootDbPath":    task.RootDbPath,
		"databasePath":  relDbPath,
		"masterDbPath":  relMasterDb,
		"rootDbInfo":    rootDbInfo,
		"queryCommands": queryCommands,
		"runDirectory":  task.RunDirectory,
		"budget":        task.TotalStepsBudget,
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

// RunTaskLs queries and displays registered parent tasks from Tier 1 master DB.
func RunTaskLs(opts TaskLsOptions) *appfault.AppError {
	masterDb := store.ResolveMasterAgentDbPath("")
	limit := resolveTaskLsLimit(opts.Limit, opts.IsAll)
	tasks, listErr := store.ListParentTasks(masterDb, limit, opts.Status)
	hasListErr := listErr != nil
	if hasListErr {
		return listErr
	}
	if opts.IsJson {
		return renderTasksJson(tasks)
	}
	renderTasksTable(tasks)

	return nil
}

func resolveTaskLsLimit(limit int, isAll bool) int {
	if isAll {
		return 1000
	}
	if limit <= 0 {
		return 20
	}

	return limit
}

func renderTasksJson(tasks []types.ParentTask) *appfault.AppError {
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return appfault.WrapSimple(err, "renderTasksJson")
	}
	fmt.Println(string(b))

	return nil
}

func renderTasksTable(tasks []types.ParentTask) {
	hasTasks := len(tasks) > 0
	if !hasTasks {
		fmt.Println("No agent parent tasks found.")
		return
	}
	fmt.Printf("%-32s  %-10s  %-24s  %-12s  %s\n", "TASK ID", "STATUS", "SLUG", "STEPS", "DIRECTORY")
	for _, t := range tasks {
		steps := fmt.Sprintf("%d/%d", t.CompletedSteps, t.TotalStepsBudget)
		fmt.Printf("%-32s  %-10s  %-24s  %-12s  %s\n", t.ParentTaskId, t.Status, t.TaskSlug, steps, t.RunDirectory)
	}
}

