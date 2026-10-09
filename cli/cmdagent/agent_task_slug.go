package cmdagent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

var (
	taskEnqueueOpts   TaskEnqueueOptions
	taskProgressOpts  TaskProgressOptions
	taskPendingOpts   TaskPendingOptions
	taskRecentOpts    TaskRecentOptions
	taskCompletedOpts TaskCompletedOptions

	taskEnqueueCmd = &cobra.Command{
		Use:   "enqueue",
		Short: "Get-or-create a parent task by slug (idempotent)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunTaskEnqueue(taskEnqueueOpts))
		},
	}

	taskProgressCmd = &cobra.Command{
		Use:   "progress",
		Short: "Show task progress and related previous tasks by slug",
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunTaskProgress(taskProgressOpts))
		},
	}

	taskPendingCmd = &cobra.Command{
		Use:   "pending",
		Short: "List or count pending subtasks across agent tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunTaskPending(taskPendingOpts))
		},
	}

	taskRecentCmd = &cobra.Command{
		Use:   "recent",
		Short: "List recently created parent tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunTaskRecent(taskRecentOpts))
		},
	}

	taskCompletedCmd = &cobra.Command{
		Use:   "completed",
		Short: "List completed root-level parent tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunTaskCompleted(taskCompletedOpts))
		},
	}
)

type TaskEnqueueOptions struct {
	Slug   string
	Name   string
	Budget int
	Dir    string
	IsJson bool
}

type TaskProgressOptions struct {
	Slug   string
	IsJson bool
}

type TaskPendingOptions struct {
	IsCount bool
	TaskId  string
	IsJson  bool
}

type TaskRecentOptions struct {
	Limit  int
	IsJson bool
}

type TaskCompletedOptions struct {
	Limit  int
	IsJson bool
}

func initTaskSlugCommands() {
	initTaskEnqueueFlags()
	initTaskProgressFlags()
	initTaskPendingFlags()
	initTaskRecentFlags()
	initTaskCompletedFlags()

	TaskCmd.AddCommand(taskEnqueueCmd)
	TaskCmd.AddCommand(taskProgressCmd)
	TaskCmd.AddCommand(taskPendingCmd)
	TaskCmd.AddCommand(taskRecentCmd)
	TaskCmd.AddCommand(taskCompletedCmd)
}

func initTaskEnqueueFlags() {
	taskEnqueueCmd.Flags().StringVar(&taskEnqueueOpts.Slug, "slug", "", "Task slug / ID (Title Case accepted; sanitized on store)")
	taskEnqueueCmd.Flags().StringVarP(&taskEnqueueOpts.Name, "name", "n", "", "Task display name (defaults to slug)")
	taskEnqueueCmd.Flags().IntVarP(&taskEnqueueOpts.Budget, "budget", "b", 300, "Total steps budget")
	taskEnqueueCmd.Flags().StringVarP(&taskEnqueueOpts.Dir, "dir", "d", "", "Custom run directory")
	taskEnqueueCmd.Flags().BoolVar(&taskEnqueueOpts.IsJson, "json", false, "Output results in JSON format")
}

func initTaskProgressFlags() {
	taskProgressCmd.Flags().StringVar(&taskProgressOpts.Slug, "slug", "", "Task slug / ID")
	taskProgressCmd.Flags().BoolVar(&taskProgressOpts.IsJson, "json", false, "Output results in JSON format")
}

func initTaskPendingFlags() {
	taskPendingCmd.Flags().BoolVar(&taskPendingOpts.IsCount, "count", false, "Print only the pending count")
	taskPendingCmd.Flags().StringVarP(&taskPendingOpts.TaskId, "task-id", "t", "", "Scope to one parent task ID or slug")
	taskPendingCmd.Flags().BoolVar(&taskPendingOpts.IsJson, "json", false, "Output results in JSON format")
}

func initTaskRecentFlags() {
	taskRecentCmd.Flags().IntVarP(&taskRecentOpts.Limit, "limit", "l", 10, "Maximum tasks to display")
	taskRecentCmd.Flags().BoolVar(&taskRecentOpts.IsJson, "json", false, "Output results in JSON format")
}

func initTaskCompletedFlags() {
	taskCompletedCmd.Flags().IntVarP(&taskCompletedOpts.Limit, "limit", "l", 20, "Maximum tasks to display")
	taskCompletedCmd.Flags().BoolVar(&taskCompletedOpts.IsJson, "json", false, "Output results in JSON format")
}

// RunTaskEnqueue implements get-or-create by slug: an existing slug takes the
// check path (reports progress, no duplicate); a new slug is enqueued.
func RunTaskEnqueue(opts TaskEnqueueOptions) *appfault.AppError {
	rawSlug := strings.TrimSpace(opts.Slug)
	hasName := len(strings.TrimSpace(opts.Name)) > 0
	if len(rawSlug) == 0 && hasName {
		rawSlug = opts.Name
	}
	if len(rawSlug) == 0 {
		return appfault.NewValidationError("task slug is required (--slug)")
	}
	slug := store.SanitizeSlug(rawSlug)
	masterDb := store.ResolveMasterAgentDbPath("")
	existing, findErr := store.FindParentTaskBySlug(masterDb, slug)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	hasExisting := existing != nil
	if hasExisting {
		return renderEnqueueCheck(existing, opts.IsJson)
	}
	name := strings.TrimSpace(opts.Name)
	if len(name) == 0 {
		name = rawSlug
	}
	taskId := generateTaskId(slug)
	runDir, dbPath := resolveTaskInitPaths(opts.Dir, slug)
	task := buildParentTask(taskId, slug, name, runDir, dbPath, opts.Budget)
	upsertErr := store.UpsertParentTask(task.RootDbPath, task)
	hasUpsertErr := upsertErr != nil
	if hasUpsertErr {
		return upsertErr
	}
	regErr := store.RegisterParentTask(masterDb, task)
	hasRegErr := regErr != nil
	if hasRegErr {
		return regErr
	}
	renderEnqueueCreated(task, opts.IsJson)

	return nil
}

func renderEnqueueCheck(task *types.ParentTask, isJson bool) *appfault.AppError {
	summary, sumErr := loadTaskSummary(task)
	hasSumErr := sumErr != nil
	if hasSumErr {
		return sumErr
	}
	if isJson {
		out := map[string]any{
			"mode":       "exists",
			"parentTask": task,
			"pending":    summary.PendingCount,
			"inProgress": summary.InProgressCount,
			"done":       summary.CompletedCount,
			"failed":     summary.FailedCount,
			"total":      len(summary.Subtasks),
		}
		b, merr := json.MarshalIndent(out, "", "  ")
		if merr != nil {
			return appfault.WrapSimple(merr, "renderEnqueueCheck")
		}
		fmt.Println(string(b))

		return nil
	}
	fmt.Printf("Task already exists (check path, not enqueued):\n")
	fmt.Printf("  Slug:     %s\n", task.TaskSlug)
	fmt.Printf("  Task ID:  %s\n", task.ParentTaskId)
	fmt.Printf("  Status:   %s\n", task.Status)
	fmt.Printf("  Subtasks: %d total (%d pending, %d in progress, %d done, %d failed)\n",
		len(summary.Subtasks), summary.PendingCount, summary.InProgressCount,
		summary.CompletedCount, summary.FailedCount)

	return nil
}

func renderEnqueueCreated(task types.ParentTask, isJson bool) {
	if isJson {
		out := map[string]any{"mode": "enqueued", "parentTask": task}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))

		return
	}
	fmt.Printf("Enqueued new task:\n")
	fmt.Printf("  Slug:     %s\n", task.TaskSlug)
	fmt.Printf("  Task ID:  %s\n", task.ParentTaskId)
	fmt.Printf("  Name:     %s\n", task.TaskName)
	fmt.Printf("  Run dir:  %s\n", task.RunDirectory)
}

// loadTaskSummary builds the subtask rollup for a task, matching subtasks
// stored under either the canonical task ID or the slug form.
func loadTaskSummary(task *types.ParentTask) (*types.TaskStatusSummary, *appfault.AppError) {
	byId, idErr := store.ListSubtasks(task.RootDbPath, task.ParentTaskId)
	hasIdErr := idErr != nil
	if hasIdErr {
		return nil, idErr
	}
	bySlug, slugErr := store.ListSubtasks(task.RootDbPath, task.TaskSlug)
	hasSlugErr := slugErr != nil
	if hasSlugErr {
		return nil, slugErr
	}
	merged := mergeSubtasksDedupe(byId, bySlug)

	return summarizeSubtasks(*task, merged), nil
}

func mergeSubtasksDedupe(lists ...[]types.Subtask) []types.Subtask {
	seen := make(map[string]bool)
	var merged []types.Subtask
	for _, list := range lists {
		for _, sub := range list {
			hasSeen := seen[sub.SubtaskId]
			if hasSeen {
				continue
			}
			seen[sub.SubtaskId] = true
			merged = append(merged, sub)
		}
	}

	return merged
}

func summarizeSubtasks(task types.ParentTask, subs []types.Subtask) *types.TaskStatusSummary {
	s := &types.TaskStatusSummary{ParentTask: task, Subtasks: subs}
	for _, sub := range subs {
		switch sub.Status {
		case "PENDING":
			s.PendingCount++
		case "IN_PROGRESS", "WORKING":
			s.InProgressCount++
		case "DONE", "COMPLETED":
			s.CompletedCount++
		case "FAILED":
			s.FailedCount++
		}
	}
	s.IsCompleted = len(subs) > 0 && s.CompletedCount == len(subs)

	return s
}

// RunTaskProgress shows progress for a slug plus related previous tasks.
func RunTaskProgress(opts TaskProgressOptions) *appfault.AppError {
	slug := strings.TrimSpace(opts.Slug)
	if len(slug) == 0 {
		return appfault.NewValidationError("task slug is required (--slug)")
	}
	masterDb := store.ResolveMasterAgentDbPath("")
	task, findErr := store.FindParentTaskBySlug(masterDb, slug)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	if task == nil {
		return appfault.NewValidationError(fmt.Sprintf("no task found for slug %q", slug))
	}
	summary, sumErr := loadTaskSummary(task)
	hasSumErr := sumErr != nil
	if hasSumErr {
		return sumErr
	}
	related, relErr := store.ListParentTasks(masterDb, 6, "")
	hasRelErr := relErr != nil
	if hasRelErr {
		return relErr
	}
	related = filterOutTask(related, task.ParentTaskId, 5)
	if opts.IsJson {
		out := map[string]any{
			"parentTask":   task,
			"summary":      summary,
			"relatedTasks": related,
		}
		b, merr := json.MarshalIndent(out, "", "  ")
		if merr != nil {
			return appfault.WrapSimple(merr, "RunTaskProgress")
		}
		fmt.Println(string(b))

		return nil
	}
	renderStatusText(summary)
	fmt.Printf("\nRelated previous tasks:\n")
	renderTasksTable(related)

	return nil
}

func filterOutTask(tasks []types.ParentTask, excludeId string, limit int) []types.ParentTask {
	var kept []types.ParentTask
	for _, t := range tasks {
		isExcluded := t.ParentTaskId == excludeId
		if isExcluded {
			continue
		}
		kept = append(kept, t)
		hasEnough := len(kept) >= limit
		if hasEnough {
			break
		}
	}

	return kept
}

// RunTaskPending lists or counts pending subtasks across agent tasks.
func RunTaskPending(opts TaskPendingOptions) *appfault.AppError {
	masterDb := store.ResolveMasterAgentDbPath("")
	tasks, listErr := store.ListParentTasks(masterDb, 1000, "")
	hasListErr := listErr != nil
	if hasListErr {
		return listErr
	}
	tasks = scopeTasksById(tasks, opts.TaskId)
	rows := collectPendingRows(tasks)
	if opts.IsCount {
		fmt.Printf("%d\n", len(rows))
		return nil
	}
	if opts.IsJson {
		b, merr := json.MarshalIndent(rows, "", "  ")
		if merr != nil {
			return appfault.WrapSimple(merr, "RunTaskPending")
		}
		fmt.Println(string(b))

		return nil
	}
	renderPendingTable(rows)

	return nil
}

type pendingRow struct {
	TaskSlug  string `json:"taskSlug"`
	TaskId    string `json:"taskId"`
	SubtaskId string `json:"subtaskId"`
	Code      string `json:"code"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
}

func scopeTasksById(tasks []types.ParentTask, taskId string) []types.ParentTask {
	trimmed := strings.TrimSpace(taskId)
	if len(trimmed) == 0 {
		return tasks
	}
	clean := store.SanitizeSlug(trimmed)
	var scoped []types.ParentTask
	for _, t := range tasks {
		matches := t.ParentTaskId == trimmed || t.TaskSlug == clean
		if matches {
			scoped = append(scoped, t)
		}
	}

	return scoped
}

func collectPendingRows(tasks []types.ParentTask) []pendingRow {
	var rows []pendingRow
	for _, t := range tasks {
		subs, listErr := store.ListPendingSubtasks(t.RootDbPath, t.ParentTaskId, t.TaskSlug)
		hasListErr := listErr != nil
		if hasListErr {
			continue
		}
		for _, sub := range subs {
			rows = append(rows, pendingRow{
				TaskSlug:  t.TaskSlug,
				TaskId:    t.ParentTaskId,
				SubtaskId: sub.SubtaskId,
				Code:      sub.TaskCode,
				Slug:      sub.TaskSlug,
				Title:     sub.Title,
			})
		}
	}

	return rows
}

func renderPendingTable(rows []pendingRow) {
	hasRows := len(rows) > 0
	if !hasRows {
		fmt.Println("No pending subtasks.")
		return
	}
	fmt.Printf("%-24s  %-14s  %-28s  %s\n", "TASK SLUG", "CODE", "SUBTASK SLUG", "TITLE")
	for _, r := range rows {
		fmt.Printf("%-24s  %-14s  %-28s  %s\n", r.TaskSlug, r.Code, r.Slug, r.Title)
	}
	fmt.Printf("\n%d pending subtask(s)\n", len(rows))
}

// RunTaskRecent lists recently created parent tasks.
func RunTaskRecent(opts TaskRecentOptions) *appfault.AppError {
	masterDb := store.ResolveMasterAgentDbPath("")
	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}
	tasks, listErr := store.ListParentTasks(masterDb, limit, "")
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

// RunTaskCompleted lists completed root-level parent tasks.
func RunTaskCompleted(opts TaskCompletedOptions) *appfault.AppError {
	masterDb := store.ResolveMasterAgentDbPath("")
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	tasks, listErr := store.ListParentTasks(masterDb, limit, "COMPLETED")
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
