// Package cmdai — ai_train.go: LLM training dataset generation and Split-DB pruning for AI analysis.
package cmdai

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// LlmTrainingExample represents an instruction-tuning record with architectural reasoning.
type LlmTrainingExample struct {
	Id           string            `json:"id"`
	TaskId       string            `json:"task_id"`
	System       string            `json:"system"`
	Prompt       string            `json:"prompt"`
	Completion   string            `json:"completion"`
	Goal         string            `json:"goal"`
	Category     string            `json:"category"`
	Reasoning    string            `json:"reasoning"`
	FilesSummary []LlmFileDecision `json:"files_summary"`
	LineDiffs    []LlmLineDiff     `json:"line_diffs"`
	StepIndex    int               `json:"step_index"`
	TotalSteps   int               `json:"total_steps"`
	NextStep     string            `json:"next_step"`
}

// LlmFileDecision documents a file inspected, modified, or removed during an AI task.
type LlmFileDecision struct {
	RelPath   string `json:"rel_path"`
	Action    string `json:"action"`
	LineCount int    `json:"line_count"`
	Reasoning string `json:"reasoning"`
}

// LlmLineDiff documents granular line modifications and coding guideline rule citations.
type LlmLineDiff struct {
	RelPath       string `json:"rel_path"`
	LineNumber    int    `json:"line_number"`
	DiffKind      string `json:"diff_kind"`
	RuleViolation string `json:"rule_violation"`
	Reasoning     string `json:"reasoning"`
	Original      string `json:"original"`
	Proposed      string `json:"proposed"`
}

// LlmTrainCliOptions holds parsed CLI flags for the llm-train command.
type LlmTrainCliOptions struct {
	TaskId string
	Format string
	Output string
	Page   int
	Limit  int
}

// AiClearCliOptions holds parsed CLI flags for the clear command.
type AiClearCliOptions struct {
	Before     string
	Force      bool
	KeepRecent int
	TaskId     string
}

var (
	llmTrainOptions LlmTrainCliOptions
	aiClearOptions  AiClearCliOptions

	aiAnalysisLlmTrainCmd = &cobra.Command{
		Use:     "llm-train",
		Aliases: []string{"train", "llmtrain"},
		Short:   "Generate LLM instruction-tuning datasets from AI analysis decisions",
		RunE:    runLlmTrainCommand,
	}

	aiAnalysisClearCmd = &cobra.Command{
		Use:     "clear",
		Aliases: []string{"prune", "purge", "reset"},
		Short:   "Prune aged AI analysis task records and vacuum Split-DB storage",
		RunE:    runAiClearCommand,
	}
)

func runLlmTrainCommand(cmd *cobra.Command, args []string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.llm_train.open_db")
	}
	defer db.Close()

	return ExecuteLlmTrainExport(db, llmTrainOptions)
}

func runAiClearCommand(cmd *cobra.Command, args []string) error {
	db, err := store.OpenAiAnalysisSplitDB("")
	if err != nil {
		return apperror.WrapSimple(err, "cmdai.clear.open_db")
	}
	defer db.Close()

	return ExecuteAiAnalysisClear(db, aiClearOptions)
}

// ExecuteLlmTrainExport queries the database and renders the training dataset.
func ExecuteLlmTrainExport(db *store.AiAnalysisSplitDB, opts LlmTrainCliOptions) error {
	tasks, err := queryTrainTasks(db, opts)
	if err != nil {
		return err
	}

	examples, err := buildTrainingExamples(db, tasks)
	if err != nil {
		return err
	}

	payload, err := formatTrainingExamples(examples, opts.Format)
	if err != nil {
		return err
	}

	return outputTrainingDataset(payload, opts.Output)
}

func queryTrainTasks(db *store.AiAnalysisSplitDB, opts LlmTrainCliOptions) ([]store.AiTask, error) {
	hasSpecificTask := strings.TrimSpace(opts.TaskId) != ""
	if hasSpecificTask {
		return querySingleTrainTask(db, strings.TrimSpace(opts.TaskId))
	}

	return queryPaginatedTrainTasks(db, opts.Page, opts.Limit)
}

func querySingleTrainTask(db *store.AiAnalysisSplitDB, taskUuid string) ([]store.AiTask, error) {
	task, err := db.GetTask(taskUuid)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.train.get_task")
	}

	return []store.AiTask{*task}, nil
}

func queryPaginatedTrainTasks(db *store.AiAnalysisSplitDB, page, limit int) ([]store.AiTask, error) {
	resolvedLimit := resolvePageLimit(limit)
	resolvedOffset := resolvePageOffset(page, resolvedLimit)

	query := `SELECT 
		AiTaskId, TaskUuid, Goal, Status, Category, Reasoning,
		TotalFiles, TotalLines, IsActive, HasFailed, Description,
		Notes, Comments, StartedAt, CompletedAt, CreatedAt, UpdatedAt
	FROM AiTask 
	ORDER BY AiTaskId ASC 
	LIMIT ? OFFSET ?`

	rows, err := db.Conn().Query(query, resolvedLimit, resolvedOffset)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.train.query_tasks")
	}
	defer rows.Close()

	return scanTrainTasks(rows)
}

func resolvePageLimit(limit int) int {
	isPositive := limit > 0
	if isPositive {
		return limit
	}

	return 10
}

func resolvePageOffset(page, limit int) int {
	isMultiPage := page > 1
	if isMultiPage {
		return (page - 1) * limit
	}

	return 0
}

func scanTrainTasks(rows *sql.Rows) ([]store.AiTask, error) {
	var tasks []store.AiTask
	for rows.Next() {
		var t store.AiTask
		var isActiveInt, hasFailedInt int
		err := rows.Scan(
			&t.AiTaskId, &t.TaskUuid, &t.Goal, &t.Status, &t.Category, &t.Reasoning,
			&t.TotalFiles, &t.TotalLines, &isActiveInt, &hasFailedInt, &t.Description,
			&t.Notes, &t.Comments, &t.StartedAt, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "cmdai.train.scan_task")
		}
		t.IsActive = isActiveInt == 1
		t.HasFailed = hasFailedInt == 1
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func buildTrainingExamples(db *store.AiAnalysisSplitDB, tasks []store.AiTask) ([]LlmTrainingExample, error) {
	examples := make([]LlmTrainingExample, 0, len(tasks))
	total := len(tasks)

	for i, task := range tasks {
		ex, err := buildSingleTrainingExample(db, task, i+1, total, tasks)
		if err != nil {
			return nil, err
		}
		examples = append(examples, ex)
	}

	return examples, nil
}

func buildSingleTrainingExample(db *store.AiAnalysisSplitDB, task store.AiTask, idx, total int, allTasks []store.AiTask) (LlmTrainingExample, error) {
	files, err := db.GetTaskFiles(task.AiTaskId)
	if err != nil {
		return LlmTrainingExample{}, apperror.WrapSimple(err, "cmdai.train.get_files")
	}

	fileDecisions := make([]LlmFileDecision, 0, len(files))
	var lineDiffs []LlmLineDiff

	for _, f := range files {
		fileDecisions = append(fileDecisions, LlmFileDecision{
			RelPath:   f.RelPath,
			Action:    f.Action,
			LineCount: f.LineCount,
			Reasoning: f.Reasoning,
		})

		diffs, diffErr := collectFileLineDiffs(db, f)
		if diffErr != nil {
			return LlmTrainingExample{}, diffErr
		}
		lineDiffs = append(lineDiffs, diffs...)
	}

	nextStep := resolveNextStepMessage(idx, total, allTasks)
	prompt := fmt.Sprintf("Hey, learn this: What was the decision-making for task %s? Why was it made, which files and lines were modified, and what was the architectural reasoning?", task.TaskUuid)
	completion := composeCompletionText(task, fileDecisions, lineDiffs, nextStep)

	return LlmTrainingExample{
		Id:           fmt.Sprintf("gitmap-task-%s", task.TaskUuid),
		TaskId:       task.TaskUuid,
		System:       "You are an autonomous senior software engineer explaining architectural decisions, code modifications, and guideline enforcement.",
		Prompt:       prompt,
		Completion:   completion,
		Goal:         task.Goal,
		Category:     task.Category,
		Reasoning:    task.Reasoning,
		FilesSummary: fileDecisions,
		LineDiffs:    lineDiffs,
		StepIndex:    idx,
		TotalSteps:   total,
		NextStep:     nextStep,
	}, nil
}

func collectFileLineDiffs(db *store.AiAnalysisSplitDB, f store.AiTaskFile) ([]LlmLineDiff, error) {
	lines, err := db.GetTaskLines(f.AiTaskFileId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.train.get_lines")
	}

	diffs := make([]LlmLineDiff, 0, len(lines))
	for _, l := range lines {
		diffs = append(diffs, LlmLineDiff{
			RelPath:       f.RelPath,
			LineNumber:    l.LineNumber,
			DiffKind:      l.DiffKind,
			RuleViolation: l.RuleViolation,
			Reasoning:     l.Reasoning,
			Original:      l.OriginalContent,
			Proposed:      l.ProposedContent,
		})
	}

	return diffs, nil
}

func resolveNextStepMessage(idx, total int, tasks []store.AiTask) string {
	hasNext := idx < total
	if hasNext {
		nextTask := tasks[idx]
		return fmt.Sprintf("Next step: Proceed to task %s (Goal: %q).", nextTask.TaskUuid, nextTask.Goal)
	}

	return fmt.Sprintf("Next step: All %d tasks in sequence processed. Training sequence complete.", total)
}

func composeCompletionText(task store.AiTask, files []LlmFileDecision, diffs []LlmLineDiff, nextStep string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Decision-Making Summary for Task: %s\n", task.TaskUuid))
	sb.WriteString(fmt.Sprintf("Goal: %s\n", task.Goal))
	sb.WriteString(fmt.Sprintf("Category: %s | Status: %s\n", task.Category, task.Status))

	hasReasoning := strings.TrimSpace(task.Reasoning) != ""
	if hasReasoning {
		sb.WriteString(fmt.Sprintf("Architectural Reasoning: %s\n", task.Reasoning))
	}

	sb.WriteString("\nFiles Inspected & Modified:\n")
	for _, f := range files {
		sb.WriteString(fmt.Sprintf("- [%s] %s (lines: %d): %s\n", f.Action, f.RelPath, f.LineCount, f.Reasoning))
	}

	hasDiffs := len(diffs) > 0
	if hasDiffs {
		sb.WriteString("\nGranular Line Decisions:\n")
		for _, d := range diffs {
			sb.WriteString(fmt.Sprintf("- %s:%d [%s] Rule: %s — %s\n", d.RelPath, d.LineNumber, d.DiffKind, d.RuleViolation, d.Reasoning))
		}
	}

	sb.WriteString(fmt.Sprintf("\n%s\n", nextStep))

	return sb.String()
}

func formatTrainingExamples(examples []LlmTrainingExample, format string) ([]byte, error) {
	cleanFmt := strings.ToLower(strings.TrimSpace(format))
	isText := cleanFmt == "text" || cleanFmt == "markdown" || cleanFmt == "md"
	if isText {
		return formatExamplesAsText(examples), nil
	}

	return formatExamplesAsJsonl(examples)
}

func formatExamplesAsJsonl(examples []LlmTrainingExample) ([]byte, error) {
	var buf bytes.Buffer
	for _, ex := range examples {
		lineBytes, err := json.Marshal(ex)
		if err != nil {
			return nil, apperror.WrapSimple(err, "cmdai.train.marshal_jsonl")
		}
		buf.Write(lineBytes)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

func formatExamplesAsText(examples []LlmTrainingExample) []byte {
	var sb strings.Builder
	sb.WriteString("# AI Analysis LLM Sequential Training Log\n\n")

	for _, ex := range examples {
		sb.WriteString(fmt.Sprintf("## Task Decision [%d/%d]: %s\n\n", ex.StepIndex, ex.TotalSteps, ex.TaskId))
		sb.WriteString(fmt.Sprintf("**User Prompt:**\n%s\n\n", ex.Prompt))
		sb.WriteString(fmt.Sprintf("**Assistant Rationale:**\n%s\n\n", ex.Completion))
		sb.WriteString("---\n\n")
	}

	return []byte(sb.String())
}

func outputTrainingDataset(payload []byte, outPath string) error {
	hasOutPath := strings.TrimSpace(outPath) != ""
	if !hasOutPath {
		fmt.Print(string(payload))

		return nil
	}

	dir := filepath.Dir(outPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, "cmdai.train.mkdir")
	}

	if err := os.WriteFile(outPath, payload, 0644); err != nil {
		return apperror.WrapSimple(err, "cmdai.train.write_file")
	}

	fmt.Printf("Wrote %d bytes of LLM training dataset to %s\n", len(payload), outPath)

	return nil
}

// ExecuteAiAnalysisClear prunes aged tasks and reclaims disk space with VACUUM.
func ExecuteAiAnalysisClear(db *store.AiAnalysisSplitDB, opts AiClearCliOptions) error {
	taskIds, err := resolveTaskIdsToClear(db, opts)
	if err != nil {
		return err
	}

	isEmpty := len(taskIds) == 0
	if isEmpty {
		fmt.Println("No matching AI analysis tasks found to clear.")

		return nil
	}

	if err := deleteTasksAndLines(db, taskIds); err != nil {
		return err
	}

	if err := vacuumDatabase(db); err != nil {
		return err
	}

	fmt.Printf("Successfully cleared %d task record(s) and reclaimed storage via SQLite VACUUM.\n", len(taskIds))

	return nil
}

func resolveTaskIdsToClear(db *store.AiAnalysisSplitDB, opts AiClearCliOptions) ([]int64, error) {
	hasTaskId := strings.TrimSpace(opts.TaskId) != ""
	if hasTaskId {
		return resolveSingleTaskForClear(db, strings.TrimSpace(opts.TaskId))
	}

	hasKeepRecent := opts.KeepRecent > 0
	if hasKeepRecent {
		return resolveKeepRecentTasksToClear(db, opts.KeepRecent)
	}

	hasBefore := strings.TrimSpace(opts.Before) != ""
	if hasBefore {
		return resolveBeforeTasksToClear(db, strings.TrimSpace(opts.Before))
	}

	if !opts.Force {
		return nil, apperror.NewValidation("cmdai.clear", "E3012", "must specify --before, --keep-recent, --task, or --force to clear all")
	}

	return resolveAllTaskIds(db)
}

func resolveSingleTaskForClear(db *store.AiAnalysisSplitDB, taskUuid string) ([]int64, error) {
	task, err := db.GetTask(taskUuid)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.clear.get_task")
	}

	return []int64{task.AiTaskId}, nil
}

func resolveKeepRecentTasksToClear(db *store.AiAnalysisSplitDB, keepCount int) ([]int64, error) {
	query := `SELECT AiTaskId FROM AiTask ORDER BY AiTaskId DESC LIMIT -1 OFFSET ?`
	rows, err := db.Conn().Query(query, keepCount)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.clear.keep_recent_query")
	}
	defer rows.Close()

	return scanIdList(rows)
}

func resolveBeforeTasksToClear(db *store.AiAnalysisSplitDB, beforeStr string) ([]int64, error) {
	cutoff, err := parseBeforeCutoff(beforeStr)
	if err != nil {
		return nil, err
	}

	cutoffStr := cutoff.Format(time.RFC3339)
	query := `SELECT AiTaskId FROM AiTask WHERE CreatedAt < ?`
	rows, err := db.Conn().Query(query, cutoffStr)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.clear.before_query")
	}
	defer rows.Close()

	return scanIdList(rows)
}

func parseBeforeCutoff(raw string) (time.Time, error) {
	clean := strings.TrimSpace(raw)
	dur, durErr := parseRelativeDuration(clean)
	isDuration := durErr == nil && dur > 0
	if isDuration {
		return time.Now().UTC().Add(-dur), nil
	}

	dateFormats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range dateFormats {
		parsed, parseErr := time.Parse(layout, clean)
		isParsed := parseErr == nil
		if isParsed {
			return parsed.UTC(), nil
		}
	}

	return time.Time{}, apperror.NewValidation("cmdai.clear", "E3013", fmt.Sprintf("invalid before date format: %q", raw))
}

func parseCustomUnitDuration(raw, suffix string, multiplier time.Duration) (time.Duration, bool) {
	if !strings.HasSuffix(raw, suffix) {
		return 0, false
	}

	val, err := strconv.Atoi(strings.TrimSuffix(raw, suffix))
	if err == nil && val > 0 {
		return time.Duration(val) * multiplier, true
	}

	return 0, false
}

func parseRelativeDuration(raw string) (time.Duration, error) {
	lower := strings.ToLower(raw)
	if dur, isUnitMatch := parseCustomUnitDuration(lower, "d", 24*time.Hour); isUnitMatch {
		return dur, nil
	}

	if dur, isUnitMatch := parseCustomUnitDuration(lower, "w", 7*24*time.Hour); isUnitMatch {
		return dur, nil
	}

	return time.ParseDuration(lower)
}

func resolveAllTaskIds(db *store.AiAnalysisSplitDB) ([]int64, error) {
	query := `SELECT AiTaskId FROM AiTask ORDER BY AiTaskId ASC`
	rows, err := db.Conn().Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "cmdai.clear.all_query")
	}
	defer rows.Close()

	return scanIdList(rows)
}

func scanIdList(rows *sql.Rows) ([]int64, error) {
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, apperror.WrapSimple(err, "cmdai.clear.scan_id")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func deleteTasksAndLines(db *store.AiAnalysisSplitDB, taskIds []int64) error {
	for _, id := range taskIds {
		delLinesQuery := `DELETE FROM AiTaskLine WHERE AiTaskFileId IN (SELECT AiTaskFileId FROM AiTaskFile WHERE AiTaskId = ?)`
		lineRes := store.ExecWrapper(db.Conn(), delLinesQuery, id)
		if lineRes.IsFailure {
			return apperror.WrapSimple(lineRes.Error, "cmdai.clear.delete_lines")
		}

		delFilesQuery := `DELETE FROM AiTaskFile WHERE AiTaskId = ?`
		fileRes := store.ExecWrapper(db.Conn(), delFilesQuery, id)
		if fileRes.IsFailure {
			return apperror.WrapSimple(fileRes.Error, "cmdai.clear.delete_files")
		}

		delTaskQuery := `DELETE FROM AiTask WHERE AiTaskId = ?`
		taskRes := store.ExecWrapper(db.Conn(), delTaskQuery, id)
		if taskRes.IsFailure {
			return apperror.WrapSimple(taskRes.Error, "cmdai.clear.delete_task")
		}
	}

	return nil
}

func vacuumDatabase(db *store.AiAnalysisSplitDB) error {
	res := store.ExecWrapper(db.Conn(), "VACUUM;")
	if res.IsFailure {
		return apperror.WrapSimple(res.Error, "cmdai.clear.vacuum")
	}

	return nil
}

func init() {
	initLlmTrainFlags()
	initAiClearFlags()

	AiAnalysisCmd.AddCommand(aiAnalysisLlmTrainCmd)
	AiAnalysisCmd.AddCommand(aiAnalysisClearCmd)
}

func initLlmTrainFlags() {
	aiAnalysisLlmTrainCmd.Flags().StringVarP(&llmTrainOptions.TaskId, "task", "t", "", "Filter training dataset by specific task UUID")
	aiAnalysisLlmTrainCmd.Flags().StringVarP(&llmTrainOptions.Format, "format", "f", "jsonl", "Output format (jsonl, text, markdown)")
	aiAnalysisLlmTrainCmd.Flags().StringVarP(&llmTrainOptions.Output, "out", "o", "", "Destination output file path")
	aiAnalysisLlmTrainCmd.Flags().IntVar(&llmTrainOptions.Page, "page", 1, "Page number for training example navigation")
	aiAnalysisLlmTrainCmd.Flags().IntVar(&llmTrainOptions.Limit, "limit", 10, "Number of tasks per training batch")
}

func initAiClearFlags() {
	aiAnalysisClearCmd.Flags().StringVar(&aiClearOptions.Before, "before", "", "Prune task sessions older than timestamp or duration (e.g. 14d, 2026-01-01)")
	aiAnalysisClearCmd.Flags().BoolVar(&aiClearOptions.Force, "force", false, "Force destructive purge confirmation")
	aiAnalysisClearCmd.Flags().IntVar(&aiClearOptions.KeepRecent, "keep-recent", 0, "Retain the N most recent tasks and prune older ones")
	aiAnalysisClearCmd.Flags().StringVarP(&aiClearOptions.TaskId, "task", "t", "", "Prune a specific task session by UUID")
}
