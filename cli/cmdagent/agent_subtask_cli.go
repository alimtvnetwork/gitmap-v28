package cmdagent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

var (
	subtaskAddOpts      SubtaskAddOptions
	subtaskClaimOpts    SubtaskClaimOptions
	subtaskStartOpts    SubtaskStartOptions
	subtaskCompleteOpts SubtaskCompleteOptions
	subtaskFailOpts     SubtaskFailOptions
	subtaskLsOpts       SubtaskLsOptions

	subtaskAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Enqueue subtasks to parent task database",
		RunE: func(cmd *cobra.Command, args []string) error {
			InterceptSubtaskAddJSON(args, &subtaskAddOpts)
			if len(args) > 0 && len(subtaskAddOpts.JsonInput) == 0 {
				subtaskAddOpts.JsonInput = args[0]
			}
			return toError(RunSubtaskAdd(subtaskAddOpts))
		},
	}

	subtaskClaimCmd = &cobra.Command{
		Use:   "claim",
		Short: "Atomically claim the next available pending subtask",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && len(subtaskClaimOpts.AgentRole) == 0 {
				subtaskClaimOpts.AgentRole = args[0]
			}
			return toError(RunSubtaskClaim(subtaskClaimOpts))
		},
	}

	subtaskStartCmd = &cobra.Command{
		Use:   "start [subtask-id]",
		Short: "Mark a claimed subtask as in-progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				subtaskStartOpts.SubtaskId = args[0]
			}
			return toError(RunSubtaskStart(subtaskStartOpts))
		},
	}

	subtaskCompleteCmd = &cobra.Command{
		Use:   "complete [subtask-id]",
		Short: "Mark a subtask as done with verification evidence",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				subtaskCompleteOpts.SubtaskId = args[0]
			}
			return toError(RunSubtaskComplete(subtaskCompleteOpts))
		},
	}

	subtaskFailCmd = &cobra.Command{
		Use:   "fail [subtask-id]",
		Short: "Mark a subtask as failed with root cause analysis",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				subtaskFailOpts.SubtaskId = args[0]
			}
			return toError(RunSubtaskFail(subtaskFailOpts))
		},
	}

	subtaskLsCmd = &cobra.Command{
		Use:     "ls [task-id]",
		Aliases: []string{"list"},
		Short:   "List subtasks belonging to a parent task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && len(subtaskLsOpts.TaskId) == 0 {
				subtaskLsOpts.TaskId = args[0]
			}
			return toError(RunSubtaskLs(subtaskLsOpts))
		},
	}
)

func registerTargetFlags(cmd *cobra.Command, taskID, role, code, slug *string) {
	cmd.Flags().StringVarP(role, "agent", "a", "", "Agent role")
	cmd.Flags().StringVarP(taskID, "task-id", "t", "", "Parent task ID or slug")
	cmd.Flags().StringVarP(code, "code", "c", "", "Subtask code (e.g. Subtask-01)")
	cmd.Flags().StringVarP(slug, "slug", "s", "", "Subtask slug")
}

func initSubtaskCommands() {
	initSubtaskAddFlags()
	initSubtaskClaimFlags()
	initSubtaskLifecycleFlags()
	initSubtaskLsFlags()
	SubtaskCmd.AddCommand(
		subtaskAddCmd,
		subtaskClaimCmd,
		subtaskStartCmd,
		subtaskCompleteCmd,
		subtaskFailCmd,
		subtaskLsCmd,
	)
}

func initSubtaskAddFlags() {
	f := subtaskAddCmd.Flags()
	f.StringVarP(&subtaskAddOpts.ParentId, "parent", "p", "", "Parent task ID or slug")
	f.StringVarP(&subtaskAddOpts.JsonInput, "data", "d", "", "JSON array or single subtask object")
	f.StringVar(&subtaskAddOpts.JsonInput, "input", "", "JSON array or single subtask object")
	f.StringVar(&subtaskAddOpts.Title, "title", "", "Subtask title")
	f.StringVar(&subtaskAddOpts.TaskCode, "code", "", "Task code (e.g. Subtask-01)")
	f.StringVar(&subtaskAddOpts.Slug, "slug", "", "Subtask slug")
	f.StringVar(&subtaskAddOpts.OwnedFiles, "owned", "[]", "JSON array of owned files")
	f.BoolVarP(&subtaskAddOpts.IsJson, "json", "j", false, "Output results in JSON format")
}

func initSubtaskClaimFlags() {
	subtaskClaimCmd.Flags().StringVarP(&subtaskClaimOpts.AgentRole, "agent", "a", "", "Agent role claiming subtask")
	subtaskClaimCmd.Flags().StringVarP(&subtaskClaimOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
}

func initSubtaskLifecycleFlags() {
	registerTargetFlags(subtaskStartCmd, &subtaskStartOpts.TaskId, &subtaskStartOpts.AgentRole, &subtaskStartOpts.TaskCode, &subtaskStartOpts.Slug)
	subtaskStartCmd.Flags().BoolVarP(&subtaskStartOpts.IsJson, "json", "j", false, "Output result in JSON format")
	registerTargetFlags(subtaskCompleteCmd, &subtaskCompleteOpts.TaskId, &subtaskCompleteOpts.AgentRole, &subtaskCompleteOpts.TaskCode, &subtaskCompleteOpts.Slug)
	subtaskCompleteCmd.Flags().StringVarP(&subtaskCompleteOpts.Evidence, "evidence", "e", "", "Completion evidence or summary")
	subtaskCompleteCmd.Flags().BoolVarP(&subtaskCompleteOpts.IsJson, "json", "j", false, "Output result in JSON format")
	registerTargetFlags(subtaskFailCmd, &subtaskFailOpts.TaskId, &subtaskFailOpts.AgentRole, &subtaskFailOpts.TaskCode, &subtaskFailOpts.Slug)
	subtaskFailCmd.Flags().StringVarP(&subtaskFailOpts.Reason, "reason", "r", "", "Failure reason or RCA")
	subtaskFailCmd.Flags().BoolVarP(&subtaskFailOpts.IsJson, "json", "j", false, "Output result in JSON format")
}

func initSubtaskLsFlags() {
	subtaskLsCmd.Flags().StringVarP(&subtaskLsOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
	subtaskLsCmd.Flags().BoolVarP(&subtaskLsOpts.IsJson, "json", "j", false, "Output results in JSON format")
}

func renderSubtasksTable(subtasks []types.Subtask) {
	if len(subtasks) == 0 {
		fmt.Println("No subtasks found.")
		return
	}
	cfg := termout.TableConfig{
		HeaderColor: constants.ColorYellow,
		BorderColor: constants.ColorDim,
		Columns:     buildSubtaskColumns(),
		Rows:        buildSubtaskRows(subtasks),
	}
	termout.PrintTable(cfg)
}

func buildSubtaskColumns() []termout.Column {
	return []termout.Column{
		{Title: "SUBTASK ID", MinWidth: 20, Align: termout.AlignLeft},
		{Title: "CODE", MinWidth: 12, Align: termout.AlignLeft},
		{Title: "STATUS", MinWidth: 12, Align: termout.AlignLeft},
		{Title: "AGENT", MinWidth: 15, Align: termout.AlignLeft},
		{Title: "TITLE", MinWidth: 20, Align: termout.AlignLeft},
	}
}

func buildSubtaskRows(subtasks []types.Subtask) []termout.Row {
	rows := make([]termout.Row, 0, len(subtasks))
	for _, s := range subtasks {
		rows = append(rows, buildSingleSubtaskRow(s))
	}

	return rows
}

func buildSingleSubtaskRow(s types.Subtask) termout.Row {
	agent := s.AssignedAgentRole
	if len(agent) == 0 {
		agent = "-"
	}

	return termout.Row{
		Cells: []string{s.SubtaskId, s.TaskCode, s.Status, agent, s.Title},
		Color: resolveSubtaskStatusColor(s.Status),
	}
}

func resolveSubtaskStatusColor(status string) string {
	switch strings.ToUpper(status) {
	case "COMPLETED", "DONE":
		return constants.ColorGreen
	case "IN_PROGRESS", "RUNNING":
		return constants.ColorCyan
	case "FAILED":
		return constants.ColorRed
	default:
		return constants.ColorWhite
	}
}

func renderSubtasksJsonEnvelope(subtasks []types.Subtask, dbPath, parentId string) *appfault.AppError {
	cwd, _ := os.Getwd()
	classifySubtaskFiles(subtasks, cwd)
	envelope := buildSubtaskLsEnvelope(subtasks, dbPath, parentId, cwd)
	b, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return appfault.WrapSimple(err, "renderSubtasksJsonEnvelope")
	}
	fmt.Println(string(b))

	return nil
}

func buildSubtaskLsEnvelope(subtasks []types.Subtask, dbPath, parentId, cwd string) SubtaskLsEnvelope {
	return SubtaskLsEnvelope{
		Status:       "success",
		DatabasePath: toRelativePath(cwd, dbPath),
		ParentTaskId: parentId,
		Summary:      computeSubtaskSummary(subtasks),
		Subtasks:     subtasks,
	}
}

func classifySubtaskFiles(subtasks []types.Subtask, cwd string) {
	for i := range subtasks {
		subtasks[i].Files = fsutil.ClassifyFileTargets(subtasks[i].OwnedFilesJson, "", cwd)
	}
}

func computeSubtaskSummary(subtasks []types.Subtask) SubtaskStatusSummary {
	summary := SubtaskStatusSummary{TotalCount: len(subtasks)}
	for _, s := range subtasks {
		accumulateStatusCount(&summary, s.Status)
	}

	return summary
}

func accumulateStatusCount(summary *SubtaskStatusSummary, status string) {
	switch strings.ToUpper(status) {
	case "PENDING":
		summary.PendingCount++
	case "IN_PROGRESS", "RUNNING":
		summary.InProgressCount++
	case "COMPLETED", "DONE":
		summary.CompletedCount++
	case "FAILED":
		summary.FailedCount++
	}
}

func renderSubtaskAddEnvelope(subtasks []types.Subtask, dbPath, parentId string) *appfault.AppError {
	cwd, _ := os.Getwd()
	envelope := buildSubtaskAddEnvelope(subtasks, dbPath, parentId, cwd)
	b, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return appfault.WrapSimple(err, "renderSubtaskAddEnvelope")
	}
	fmt.Println(string(b))

	return nil
}

func buildSubtaskAddEnvelope(subtasks []types.Subtask, dbPath, parentId, cwd string) SubtaskAddEnvelope {
	return SubtaskAddEnvelope{
		Status:       "success",
		DatabasePath: toRelativePath(cwd, dbPath),
		ParentTaskId: parentId,
		AddedCount:   len(subtasks),
		Subtasks:     subtasks,
	}
}

func renderSubtaskLifecycleJSON(status, subID, role, dbPath, details string) *appfault.AppError {
	cwd, _ := os.Getwd()
	env := SubtaskLifecycleEnvelope{
		Status:       status,
		SubtaskId:    subID,
		AgentRole:    role,
		DatabasePath: toRelativePath(cwd, dbPath),
		Details:      details,
	}
	b, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return appfault.WrapSimple(err, "renderSubtaskLifecycleJSON")
	}
	fmt.Println(string(b))

	return nil
}
