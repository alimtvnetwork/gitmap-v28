// Package cmdai — ai_analysis_cmd.go: CLI subcommands and dispatcher for AI task analysis sessions.
package cmdai

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

var (
	// AiAnalysisCmd is the root Cobra command for gitmap ai-analysis.
	AiAnalysisCmd = &cobra.Command{
		Use:     "analysis",
		Aliases: []string{"ai-analysis", "trace", "session"},
		Short:   "AI task analysis session recording and granular reasoning tracer",
		RunE:    runDefaultAnalysisCmd,
	}

	analysisStartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a new AI analysis task session",
		RunE:  runAnalysisStartCmd,
	}

	analysisRecordCmd = &cobra.Command{
		Use:   "record",
		Short: "Record file observation, modification, or action with reasoning",
		RunE:  runAnalysisRecordCmd,
	}

	analysisLineCmd = &cobra.Command{
		Use:   "line",
		Short: "Record line diff mutation and rule violation reasoning",
		RunE:  runAnalysisLineCmd,
	}

	analysisStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Display status and summary for AI analysis sessions",
		RunE:  runAnalysisStatusCmd,
	}

	analysisListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List recorded AI analysis sessions",
		RunE:    runAnalysisListCmd,
	}

	analysisInspectCmd = &cobra.Command{
		Use:     "inspect [task_id]",
		Aliases: []string{"show", "view", "get"},
		Short:   "Inspect detailed file operations and reasoning for a task",
		RunE:    runAnalysisInspectCmd,
	}

	analysisFinishCmd = &cobra.Command{
		Use:   "finish",
		Short: "Mark an AI analysis task session as completed or failed",
		RunE:  runAnalysisFinishCmd,
	}

	taskFlag     string
	goalFlag     string
	descFlag     string
	categoryFlag string
	reasonFlag   string
	fileFlag     string
	actionFlag   string
	lineFlag     int
	diffFlag     string
	ruleFlag     string
	limitFlag    int
	jsonFlag     bool
	statusFlag   string
	finishSumm   string
	inspectTask  string
)

func runDefaultAnalysisCmd(cmd *cobra.Command, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		return cmd.Help()
	}

	return runAnalysisListCmd(cmd, args)
}

func runAnalysisStartCmd(cmd *cobra.Command, args []string) error {
	taskUuid := strings.TrimSpace(taskFlag)
	isEmptyTask := taskUuid == ""
	if isEmptyTask {
		return apperror.NewValidationError("flag --task is required")
	}

	goal := resolveTaskGoal()
	isEmptyGoal := goal == ""
	if isEmptyGoal {
		return apperror.NewValidationError("flag --goal is required")
	}

	return executeStartTask(taskUuid, goal)
}

func resolveTaskGoal() string {
	goal := strings.TrimSpace(goalFlag)
	hasGoal := goal != ""
	if hasGoal {
		return goal
	}

	return strings.TrimSpace(descFlag)
}

func executeStartTask(taskUuid, goal string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "ai_analysis.start.open_db")
	}
	defer db.Close()

	cat := resolveCategory(categoryFlag)
	task := &store.AiTask{
		TaskUuid:  taskUuid,
		Goal:      goal,
		Category:  cat,
		Reasoning: strings.TrimSpace(reasonFlag),
		Status:    "pending",
		IsActive:  true,
	}
	if _, err := db.CreateTask(task); err != nil {
		return err
	}

	fmt.Printf("Started AI analysis task %s: %s\n", taskUuid, goal)

	return nil
}

func resolveCategory(cat string) string {
	clean := strings.TrimSpace(cat)
	hasCategory := clean != ""
	if hasCategory {
		return clean
	}

	return "refactor"
}

func runAnalysisRecordCmd(cmd *cobra.Command, args []string) error {
	taskUuid := strings.TrimSpace(taskFlag)
	isEmptyTask := taskUuid == ""
	if isEmptyTask {
		return apperror.NewValidationError("flag --task is required")
	}

	filePath := strings.TrimSpace(fileFlag)
	isEmptyFile := filePath == ""
	if isEmptyFile {
		return apperror.NewValidationError("flag --file is required")
	}

	action := resolveAction(actionFlag)

	return executeRecordFile(taskUuid, filePath, action, strings.TrimSpace(reasonFlag))
}

func resolveAction(action string) string {
	clean := strings.TrimSpace(action)
	hasAction := clean != ""
	if hasAction {
		return clean
	}

	return "modified"
}

func executeRecordFile(taskUuid, filePath, action, reason string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "ai_analysis.record.open_db")
	}
	defer db.Close()

	task, err := db.GetTask(taskUuid)
	if err != nil {
		return apperror.NewNotFoundError("task not found: " + taskUuid)
	}

	file := &store.AiTaskFile{
		AiTaskId:   task.AiTaskId,
		RelPath:    filepath.ToSlash(filepath.Clean(filePath)),
		Action:     action,
		Reasoning:  reason,
		IsModified: action == "modified",
		IsRemoved:  action == "removed",
	}
	if _, err := db.RecordFile(file); err != nil {
		return err
	}

	fmt.Printf("Recorded file %s for task %s (action: %s)\n", file.RelPath, taskUuid, action)

	return nil
}

func runAnalysisLineCmd(cmd *cobra.Command, args []string) error {
	taskUuid := strings.TrimSpace(taskFlag)
	isEmptyTask := taskUuid == ""
	if isEmptyTask {
		return apperror.NewValidationError("flag --task is required")
	}

	filePath := strings.TrimSpace(fileFlag)
	isEmptyFile := filePath == ""
	if isEmptyFile {
		return apperror.NewValidationError("flag --file is required")
	}

	isInvalidLine := lineFlag <= 0
	if isInvalidLine {
		return apperror.NewValidationError("flag --line must be a positive line number")
	}

	diffKind := resolveDiffKind(diffFlag)

	return executeRecordLine(taskUuid, filePath, lineFlag, diffKind, strings.TrimSpace(reasonFlag))
}

func resolveDiffKind(diff string) string {
	clean := strings.TrimSpace(diff)
	hasDiff := clean != ""
	if hasDiff {
		return clean
	}

	return "replace"
}

func executeRecordLine(taskUuid, filePath string, lineNum int, diffKind, reason string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "ai_analysis.line.open_db")
	}
	defer db.Close()

	task, err := db.GetTask(taskUuid)
	if err != nil {
		return apperror.NewNotFoundError("task not found: " + taskUuid)
	}

	file, err := resolveOrCreateTaskFile(db, task.AiTaskId, filePath)
	if err != nil {
		return err
	}

	line := &store.AiTaskLine{
		AiTaskFileId:  file.AiTaskFileId,
		LineNumber:    lineNum,
		DiffKind:      diffKind,
		RuleViolation: strings.TrimSpace(ruleFlag),
		Reasoning:     reason,
		IsApplied:     true,
	}
	if _, err := db.RecordLine(line); err != nil {
		return err
	}

	fmt.Printf("Recorded line %d diff (%s) for %s in task %s\n", lineNum, diffKind, file.RelPath, taskUuid)

	return nil
}

func resolveOrCreateTaskFile(db *store.AiAnalysisSplitDB, taskId int64, filePath string) (*store.AiTaskFile, error) {
	cleanPath := filepath.ToSlash(filepath.Clean(filePath))
	file, err := db.GetTaskFileByRelPath(taskId, cleanPath)
	if err != nil {
		return nil, err
	}
	hasFile := file != nil
	if hasFile {
		return file, nil
	}

	newFile := &store.AiTaskFile{
		AiTaskId:   taskId,
		RelPath:    cleanPath,
		Action:     "modified",
		IsModified: true,
	}
	id, err := db.RecordFile(newFile)
	if err != nil {
		return nil, err
	}
	newFile.AiTaskFileId = id

	return newFile, nil
}

func runAnalysisStatusCmd(cmd *cobra.Command, args []string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "ai_analysis.status.open_db")
	}
	defer db.Close()

	taskUuid := strings.TrimSpace(taskFlag)
	hasTask := taskUuid != ""
	if hasTask {
		return renderTaskStatus(db, taskUuid, jsonFlag)
	}

	return renderOverallStatus(db, jsonFlag)
}

func renderTaskStatus(db *store.AiAnalysisSplitDB, taskUuid string, isJson bool) error {
	task, err := db.GetTask(taskUuid)
	if err != nil {
		return apperror.NewNotFoundError("task not found: " + taskUuid)
	}

	files, err := db.GetTaskFiles(task.AiTaskId)
	if err != nil {
		return err
	}

	if isJson {
		return printJsonOutput(map[string]any{"task": task, "files": files})
	}

	printTaskStatusCard(task, files)

	return nil
}

func renderOverallStatus(db *store.AiAnalysisSplitDB, isJson bool) error {
	summary, err := db.GetSummary()
	if err != nil {
		return err
	}

	if isJson {
		return printJsonOutput(summary)
	}

	printSummaryCards(summary)

	return nil
}

func runAnalysisListCmd(cmd *cobra.Command, args []string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "ai_analysis.list.open_db")
	}
	defer db.Close()

	limit := limitFlag
	if limit <= 0 {
		limit = 20
	}

	tasks, err := db.ListTasks(false, limit)
	if err != nil {
		return err
	}

	if jsonFlag {
		return printJsonOutput(tasks)
	}

	renderTasksTable(tasks)

	return nil
}

func runAnalysisInspectCmd(cmd *cobra.Command, args []string) error {
	targetTask := resolveInspectTaskId(args)
	isEmptyTask := targetTask == ""
	if isEmptyTask {
		return apperror.NewValidationError("task ID required to inspect")
	}

	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "ai_analysis.inspect.open_db")
	}
	defer db.Close()

	return renderTaskStatus(db, targetTask, jsonFlag)
}

func resolveInspectTaskId(args []string) string {
	hasArgs := len(args) > 0
	if hasArgs {
		return strings.TrimSpace(args[0])
	}
	hasInspectFlag := strings.TrimSpace(inspectTask) != ""
	if hasInspectFlag {
		return strings.TrimSpace(inspectTask)
	}

	return strings.TrimSpace(taskFlag)
}

func runAnalysisFinishCmd(cmd *cobra.Command, args []string) error {
	taskUuid := strings.TrimSpace(taskFlag)
	isEmptyTask := taskUuid == ""
	if isEmptyTask {
		return apperror.NewValidationError("flag --task is required")
	}

	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "ai_analysis.finish.open_db")
	}
	defer db.Close()

	st := resolveStatus(statusFlag)
	hasFailed := st == "failed"
	if err := db.UpdateTaskStatus(taskUuid, st, hasFailed); err != nil {
		return err
	}

	fmt.Printf("Task %s finished with status: %s\n", taskUuid, st)

	return nil
}

func resolveStatus(st string) string {
	clean := strings.TrimSpace(st)
	hasStatus := clean != ""
	if hasStatus {
		return clean
	}

	return "completed"
}

// DispatchAiAnalysis routes CLI arguments to native AI analysis commands.
func DispatchAiAnalysis(args []string) error {
	cleanArgs := stripAiAnalysisPrefix(args)
	isEmpty := len(cleanArgs) == 0
	if isEmpty {
		return runDefaultAnalysisCmd(AiAnalysisCmd, cleanArgs)
	}

	AiAnalysisCmd.SetArgs(cleanArgs)

	return AiAnalysisCmd.Execute()
}

func stripAiAnalysisPrefix(args []string) []string {
	isEmpty := len(args) == 0
	if isEmpty {
		return args
	}

	first := strings.ToLower(args[0])
	isPrefix := first == "ai-analysis" || first == "analysis"
	if isPrefix {
		return args[1:]
	}

	return args
}

func renderTasksTable(tasks []store.AiTask) {
	isEmpty := len(tasks) == 0
	if isEmpty {
		fmt.Println("No AI analysis tasks found.")

		return
	}

	cols := []termout.Column{
		{Title: "TASK ID", MinWidth: 20, Align: termout.AlignLeft},
		{Title: "STATUS", MinWidth: 10, Align: termout.AlignLeft},
		{Title: "CATEGORY", MinWidth: 10, Align: termout.AlignLeft},
		{Title: "GOAL", MinWidth: 30, Align: termout.AlignLeft},
		{Title: "FILES", MinWidth: 6, Align: termout.AlignRight},
		{Title: "LINES", MinWidth: 6, Align: termout.AlignRight},
		{Title: "CREATED", MinWidth: 19, Align: termout.AlignLeft},
	}

	rows := make([]termout.Row, 0, len(tasks))
	for _, t := range tasks {
		rows = append(rows, termout.Row{
			Cells: []string{
				t.TaskUuid,
				t.Status,
				t.Category,
				t.Goal,
				strconv.Itoa(t.TotalFiles),
				strconv.Itoa(t.TotalLines),
				formatDisplayTime(t.CreatedAt),
			},
		})
	}

	termout.PrintTable(termout.TableConfig{Columns: cols, Rows: rows, HasBorders: true})
}

func printSummaryCards(summary *store.AiAnalysisSummaryData) {
	fmt.Println("GitMap AI Analysis Engine Summary")
	fmt.Println("----------------------------------")
	fmt.Printf("Total Tasks     : %d\n", summary.TotalTasks)
	fmt.Printf("Active Tasks    : %d\n", summary.ActiveTasks)
	fmt.Printf("Completed Tasks : %d\n", summary.CompletedTasks)
	fmt.Printf("Failed Tasks    : %d\n", summary.FailedTasks)
	fmt.Printf("Total Files     : %d\n", summary.TotalFiles)
	fmt.Printf("Total Lines     : %d\n", summary.TotalLines)
}

func printTaskStatusCard(task *store.AiTask, files []store.AiTaskFile) {
	fmt.Printf("Task ID     : %s\n", task.TaskUuid)
	fmt.Printf("Goal        : %s\n", task.Goal)
	fmt.Printf("Status      : %s\n", task.Status)
	fmt.Printf("Category    : %s\n", task.Category)
	hasReasoning := task.Reasoning != ""
	if hasReasoning {
		fmt.Printf("Reasoning   : %s\n", task.Reasoning)
	}
	fmt.Printf("Files/Lines : %d / %d\n", task.TotalFiles, task.TotalLines)
	fmt.Printf("Created     : %s\n", formatDisplayTime(task.CreatedAt))
	fmt.Println()

	isEmpty := len(files) == 0
	if isEmpty {
		fmt.Println("No files recorded for this task.")

		return
	}

	renderFilesTable(files)
}

func renderFilesTable(files []store.AiTaskFile) {
	cols := []termout.Column{
		{Title: "ID", MinWidth: 5, Align: termout.AlignRight},
		{Title: "REL PATH", MinWidth: 35, Align: termout.AlignLeft},
		{Title: "ACTION", MinWidth: 10, Align: termout.AlignLeft},
		{Title: "REASONING", MinWidth: 30, Align: termout.AlignLeft},
	}
	rows := make([]termout.Row, 0, len(files))
	for _, f := range files {
		rows = append(rows, termout.Row{
			Cells: []string{
				strconv.FormatInt(f.AiTaskFileId, 10),
				f.RelPath,
				f.Action,
				f.Reasoning,
			},
		})
	}
	termout.PrintTable(termout.TableConfig{Columns: cols, Rows: rows, HasBorders: true})
}

func printJsonOutput(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "ai_analysis.json_marshal")
	}
	fmt.Println(string(data))

	return nil
}

func formatDisplayTime(raw string) string {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}

	return t.Format("2006-01-02 15:04:05")
}

var initAnalysisOnce sync.Once

func initAnalysisCommands() {
	initAnalysisOnce.Do(func() {
		AiAnalysisCmd.AddCommand(analysisStartCmd)
		AiAnalysisCmd.AddCommand(analysisRecordCmd)
		AiAnalysisCmd.AddCommand(analysisLineCmd)
		AiAnalysisCmd.AddCommand(analysisStatusCmd)
		AiAnalysisCmd.AddCommand(analysisListCmd)
		AiAnalysisCmd.AddCommand(analysisInspectCmd)
		AiAnalysisCmd.AddCommand(analysisFinishCmd)

		initAnalysisFlags()
	})
}

func init() {
	AiCmd.AddCommand(AiAnalysisCmd)
	initAnalysisCommands()
}

func initAnalysisFlags() {
	analysisStartCmd.Flags().StringVarP(&taskFlag, "task", "t", "", "Task session identifier")
	analysisStartCmd.Flags().StringVarP(&goalFlag, "goal", "g", "", "Task goal or objective")
	analysisStartCmd.Flags().StringVarP(&descFlag, "desc", "d", "", "Task description (alias for goal)")
	analysisStartCmd.Flags().StringVarP(&categoryFlag, "category", "c", "refactor", "Task category")
	analysisStartCmd.Flags().StringVarP(&reasonFlag, "reason", "r", "", "Strategic reasoning explanation")

	analysisRecordCmd.Flags().StringVarP(&taskFlag, "task", "t", "", "Task session identifier")
	analysisRecordCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Target file path")
	analysisRecordCmd.Flags().StringVarP(&actionFlag, "action", "a", "modified", "File action (modified, created, removed, analyzed)")
	analysisRecordCmd.Flags().StringVarP(&reasonFlag, "reason", "r", "", "Decision reasoning explanation")

	analysisLineCmd.Flags().StringVarP(&taskFlag, "task", "t", "", "Task session identifier")
	analysisLineCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Target file path")
	analysisLineCmd.Flags().IntVarP(&lineFlag, "line", "l", 0, "Line number")
	analysisLineCmd.Flags().StringVarP(&diffFlag, "diff", "d", "replace", "Diff kind (insert, delete, replace)")
	analysisLineCmd.Flags().StringVarP(&ruleFlag, "rule", "u", "", "Rule violation identifier")
	analysisLineCmd.Flags().StringVarP(&reasonFlag, "reason", "r", "", "Reasoning explanation for line mutation")

	analysisStatusCmd.Flags().StringVarP(&taskFlag, "task", "t", "", "Task session identifier")
	analysisStatusCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output results in JSON format")

	analysisListCmd.Flags().IntVarP(&limitFlag, "limit", "n", 20, "Maximum tasks to display")
	analysisListCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output results in JSON format")

	analysisInspectCmd.Flags().StringVarP(&inspectTask, "task", "t", "", "Task session identifier")
	analysisInspectCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output results in JSON format")

	analysisFinishCmd.Flags().StringVarP(&taskFlag, "task", "t", "", "Task session identifier")
	analysisFinishCmd.Flags().StringVar(&statusFlag, "status", "completed", "Final status")
	analysisFinishCmd.Flags().StringVarP(&finishSumm, "summary", "s", "", "Task outcome reasoning summary")
}
