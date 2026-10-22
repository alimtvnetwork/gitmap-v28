package cmdagent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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

type SubtaskAddOptions struct {
	ParentId, JsonInput, Title, TaskCode, Slug, OwnedFiles string
}

type SubtaskClaimOptions struct {
	AgentRole, TaskId string
}

type SubtaskStartOptions struct {
	SubtaskId, TaskCode, Slug, AgentRole, TaskId string
}

type SubtaskCompleteOptions struct {
	SubtaskId, TaskCode, Slug, AgentRole, Evidence, TaskId string
}

type SubtaskFailOptions struct {
	SubtaskId, TaskCode, Slug, AgentRole, Reason, TaskId string
}

type SubtaskLsOptions struct {
	TaskId string
	IsJson bool
}

func registerTargetFlags(cmd *cobra.Command, taskID, role, code, slug *string) {
	cmd.Flags().StringVarP(role, "agent", "a", "", "Agent role")
	cmd.Flags().StringVarP(taskID, "task-id", "t", "", "Parent task ID or slug")
	cmd.Flags().StringVarP(code, "code", "c", "", "Subtask code (e.g. Subtask-01)")
	cmd.Flags().StringVarP(slug, "slug", "s", "", "Subtask slug")
}

func initSubtaskCommands() {
	f := subtaskAddCmd.Flags()
	f.StringVarP(&subtaskAddOpts.ParentId, "parent", "p", "", "Parent task ID or slug")
	f.StringVarP(&subtaskAddOpts.JsonInput, "json", "j", "", "JSON array or single subtask object")
	f.StringVar(&subtaskAddOpts.Title, "title", "", "Subtask title")
	f.StringVar(&subtaskAddOpts.TaskCode, "code", "", "Task code (e.g. Subtask-01)")
	f.StringVar(&subtaskAddOpts.Slug, "slug", "", "Subtask slug")
	f.StringVar(&subtaskAddOpts.OwnedFiles, "owned", "[]", "JSON array of owned files")

	subtaskClaimCmd.Flags().StringVarP(&subtaskClaimOpts.AgentRole, "agent", "a", "", "Agent role claiming subtask")
	subtaskClaimCmd.Flags().StringVarP(&subtaskClaimOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")

	registerTargetFlags(subtaskStartCmd, &subtaskStartOpts.TaskId, &subtaskStartOpts.AgentRole, &subtaskStartOpts.TaskCode, &subtaskStartOpts.Slug)
	registerTargetFlags(subtaskCompleteCmd, &subtaskCompleteOpts.TaskId, &subtaskCompleteOpts.AgentRole, &subtaskCompleteOpts.TaskCode, &subtaskCompleteOpts.Slug)
	subtaskCompleteCmd.Flags().StringVarP(&subtaskCompleteOpts.Evidence, "evidence", "e", "", "Completion evidence or summary")

	registerTargetFlags(subtaskFailCmd, &subtaskFailOpts.TaskId, &subtaskFailOpts.AgentRole, &subtaskFailOpts.TaskCode, &subtaskFailOpts.Slug)
	subtaskFailCmd.Flags().StringVarP(&subtaskFailOpts.Reason, "reason", "r", "", "Failure reason or RCA")

	subtaskLsCmd.Flags().StringVarP(&subtaskLsOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
	subtaskLsCmd.Flags().BoolVar(&subtaskLsOpts.IsJson, "json", false, "Output results in JSON format")

	SubtaskCmd.AddCommand(subtaskAddCmd, subtaskClaimCmd, subtaskStartCmd, subtaskCompleteCmd, subtaskFailCmd, subtaskLsCmd)
}

// RunSubtaskAdd parses and enqueues subtasks into the parent task DB.
func RunSubtaskAdd(opts SubtaskAddOptions) *appfault.AppError {
	subtasks, parseErr := parseSubtasksInput(opts)
	if parseErr != nil {
		return parseErr
	}
	dbPath, canonicalParentId, findErr := resolveSubtaskParent(opts.ParentId)
	if findErr != nil {
		return findErr
	}
	subtasks = normalizeSubtaskParents(subtasks, opts.ParentId, canonicalParentId)
	subtasks = applySubtaskSlugs(subtasks, opts)
	if addErr := store.AddSubtasks(dbPath, canonicalParentId, subtasks); addErr != nil {
		return addErr
	}
	fmt.Printf("Successfully added %d subtask(s) to %s\n", len(subtasks), dbPath)

	return nil
}

// RunSubtaskClaim claims the next available pending subtask for an agent.
func RunSubtaskClaim(opts SubtaskClaimOptions) *appfault.AppError {
	if len(strings.TrimSpace(opts.AgentRole)) == 0 {
		return appfault.NewValidationError("agent role is required (--agent)")
	}
	dbPath, canonicalId, findErr := resolveSubtaskParent(opts.TaskId)
	if findErr != nil {
		return findErr
	}
	sub, claimErr := claimSubtaskBothForms(dbPath, canonicalId, opts.TaskId, opts.AgentRole)
	if claimErr != nil {
		return claimErr
	}
	renderClaimResult(sub, opts.AgentRole, dbPath)

	return nil
}

func renderClaimResult(sub *types.Subtask, role, dbPath string) {
	out := buildClaimResultMap(sub, role, dbPath)
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

func buildClaimResultMap(sub *types.Subtask, role, dbPath string) map[string]any {
	if sub == nil {
		return map[string]any{"status": "empty", "message": "No pending subtasks available to claim"}
	}
	cwd, _ := os.Getwd()
	relDbPath := toRelativePath(cwd, dbPath)
	classifiedFiles := fsutil.ClassifyFileTargets(sub.OwnedFilesJson, "", cwd)

	return map[string]any{
		"status":            "claimed",
		"subtaskId":         sub.SubtaskId,
		"parentTaskId":      sub.ParentTaskId,
		"taskCode":          sub.TaskCode,
		"title":             sub.Title,
		"assignedAgentRole": role,
		"databasePath":      relDbPath,
		"files":             classifiedFiles,
		"ownedFilesJson":    sub.OwnedFilesJson,
	}
}

func resolveSubtaskTarget(taskID, subtaskID, code, slug string) (string, string, *appfault.AppError) {
	hasRef := len(strings.TrimSpace(subtaskID)) > 0 || len(strings.TrimSpace(code)) > 0 || len(strings.TrimSpace(slug)) > 0
	if !hasRef {
		return "", "", appfault.NewValidationError("subtask ID, --code, or --slug is required")
	}
	dbPath, findErr := store.FindTaskDbPath("", taskID)
	if findErr != nil {
		return "", "", findErr
	}
	resolvedId, resErr := ResolveSubtaskIdentifier(dbPath, subtaskID, code, slug)
	if resErr != nil {
		return "", "", appfault.NewValidationError(resErr.Error())
	}

	return dbPath, resolvedId, nil
}

// RunSubtaskStart marks a subtask as in-progress.
func RunSubtaskStart(opts SubtaskStartOptions) *appfault.AppError {
	dbPath, subID, err := resolveSubtaskTarget(opts.TaskId, opts.SubtaskId, opts.TaskCode, opts.Slug)
	if err != nil {
		return err
	}
	if startErr := store.StartSubtask(dbPath, subID, opts.AgentRole); startErr != nil {
		return startErr
	}
	fmt.Printf("Subtask %s started by %s\n", subID, opts.AgentRole)

	return nil
}

// RunSubtaskComplete marks a subtask as done with verification evidence.
func RunSubtaskComplete(opts SubtaskCompleteOptions) *appfault.AppError {
	dbPath, subID, err := resolveSubtaskTarget(opts.TaskId, opts.SubtaskId, opts.TaskCode, opts.Slug)
	if err != nil {
		return err
	}
	if compErr := store.CompleteSubtask(dbPath, subID, opts.AgentRole, opts.Evidence); compErr != nil {
		return compErr
	}
	releaseFileClaimsBySubtask(subID)
	fmt.Printf("Subtask %s completed by %s\n", subID, opts.AgentRole)

	return nil
}

// RunSubtaskFail records a subtask failure with diagnostic reason.
func RunSubtaskFail(opts SubtaskFailOptions) *appfault.AppError {
	dbPath, subID, err := resolveSubtaskTarget(opts.TaskId, opts.SubtaskId, opts.TaskCode, opts.Slug)
	if err != nil {
		return err
	}
	if failErr := store.FailSubtask(dbPath, subID, opts.AgentRole, opts.Reason); failErr != nil {
		return failErr
	}
	releaseFileClaimsBySubtask(subID)
	fmt.Printf("Subtask %s marked as FAILED by %s\n", subID, opts.AgentRole)

	return nil
}

// RunSubtaskLs lists all subtasks belonging to a parent task.
func RunSubtaskLs(opts SubtaskLsOptions) *appfault.AppError {
	dbPath, canonicalId, findErr := resolveSubtaskParent(opts.TaskId)
	if findErr != nil {
		return findErr
	}
	subtasks, listErr := listSubtasksBothForms(dbPath, canonicalId, opts.TaskId)
	if listErr != nil {
		return listErr
	}
	if opts.IsJson {
		cwd, _ := os.Getwd()
		for i := range subtasks {
			subtasks[i].Files = fsutil.ClassifyFileTargets(subtasks[i].OwnedFilesJson, "", cwd)
		}

		return renderSubtasksJson(subtasks)
	}
	renderSubtasksTable(subtasks)

	return nil
}
