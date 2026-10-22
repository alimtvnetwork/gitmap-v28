package cmdagent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

// RunSubtaskAdd parses and enqueues subtasks into the parent task DB.
func RunSubtaskAdd(opts SubtaskAddOptions) *appfault.AppError {
	subtasks, dbPath, canonicalParentId, err := resolveAndPrepareSubtasks(opts)
	if err != nil {
		return err
	}
	addErr := store.AddSubtasks(dbPath, canonicalParentId, subtasks)
	if addErr != nil {
		return addErr
	}

	return renderSubtaskAddOutput(subtasks, dbPath, canonicalParentId, opts.IsJson)
}

func resolveAndPrepareSubtasks(opts SubtaskAddOptions) ([]types.Subtask, string, string, *appfault.AppError) {
	subtasks, parseErr := parseSubtasksInput(opts)
	if parseErr != nil {
		return nil, "", "", parseErr
	}
	dbPath, canonicalParentId, findErr := resolveSubtaskParent(opts.ParentId)
	if findErr != nil {
		return nil, "", "", findErr
	}
	subtasks = prepareSubtasksForAdd(subtasks, opts.ParentId, canonicalParentId, opts)

	return subtasks, dbPath, canonicalParentId, nil
}

func prepareSubtasksForAdd(subtasks []types.Subtask, rawParent, canonicalParent string, opts SubtaskAddOptions) []types.Subtask {
	subtasks = normalizeSubtaskParents(subtasks, rawParent, canonicalParent)

	return applySubtaskSlugs(subtasks, opts)
}

func renderSubtaskAddOutput(subtasks []types.Subtask, dbPath, canonicalParentId string, isJson bool) *appfault.AppError {
	if isJson {
		return renderSubtaskAddEnvelope(subtasks, dbPath, canonicalParentId)
	}
	fmt.Printf("Successfully added %d subtask(s) to %s\n", len(subtasks), dbPath)

	return nil
}

// RunSubtaskClaim claims the next available pending subtask for an agent.
func RunSubtaskClaim(opts SubtaskClaimOptions) *appfault.AppError {
	err := validateClaimAgentRole(opts.AgentRole)
	if err != nil {
		return err
	}
	dbPath, canonicalId, findErr := resolveSubtaskParent(opts.TaskId)
	if findErr != nil {
		return findErr
	}

	return claimAndRenderSubtask(dbPath, canonicalId, opts)
}

func validateClaimAgentRole(role string) *appfault.AppError {
	if len(strings.TrimSpace(role)) == 0 {
		return appfault.NewValidationError("agent role is required (--agent)")
	}

	return nil
}

func claimAndRenderSubtask(dbPath, canonicalId string, opts SubtaskClaimOptions) *appfault.AppError {
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

	return buildClaimedDetails(sub, role, dbPath, cwd)
}

func buildClaimedDetails(sub *types.Subtask, role, dbPath, cwd string) map[string]any {
	return map[string]any{
		"status":            "claimed",
		"subtaskId":         sub.SubtaskId,
		"parentTaskId":      sub.ParentTaskId,
		"taskCode":          sub.TaskCode,
		"title":             sub.Title,
		"assignedAgentRole": role,
		"databasePath":      toRelativePath(cwd, dbPath),
		"files":             fsutil.ClassifyFileTargets(sub.OwnedFilesJson, "", cwd),
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

	return resolveSubtaskIdInDb(dbPath, subtaskID, code, slug)
}

func resolveSubtaskIdInDb(dbPath, subtaskID, code, slug string) (string, string, *appfault.AppError) {
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

	return finishSubtaskStart(dbPath, subID, opts)
}

func finishSubtaskStart(dbPath, subID string, opts SubtaskStartOptions) *appfault.AppError {
	startErr := store.StartSubtask(dbPath, subID, opts.AgentRole)
	if startErr != nil {
		return startErr
	}

	return reportSubtaskStart(dbPath, subID, opts)
}

func reportSubtaskStart(dbPath, subID string, opts SubtaskStartOptions) *appfault.AppError {
	if opts.IsJson {
		return renderSubtaskLifecycleJSON("in_progress", subID, opts.AgentRole, dbPath, "")
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

	return finishSubtaskComplete(dbPath, subID, opts)
}

func finishSubtaskComplete(dbPath, subID string, opts SubtaskCompleteOptions) *appfault.AppError {
	compErr := store.CompleteSubtask(dbPath, subID, opts.AgentRole, opts.Evidence)
	if compErr != nil {
		return compErr
	}
	releaseFileClaimsBySubtask(subID)

	return reportSubtaskComplete(dbPath, subID, opts)
}

func reportSubtaskComplete(dbPath, subID string, opts SubtaskCompleteOptions) *appfault.AppError {
	if opts.IsJson {
		return renderSubtaskLifecycleJSON("completed", subID, opts.AgentRole, dbPath, opts.Evidence)
	}
	fmt.Printf("Subtask %s completed by %s\n", subID, opts.AgentRole)

	return nil
}

// RunSubtaskFail records a subtask failure with diagnostic reason.
func RunSubtaskFail(opts SubtaskFailOptions) *appfault.AppError {
	dbPath, subID, err := resolveSubtaskTarget(opts.TaskId, opts.SubtaskId, opts.TaskCode, opts.Slug)
	if err != nil {
		return err
	}

	return finishSubtaskFail(dbPath, subID, opts)
}

func finishSubtaskFail(dbPath, subID string, opts SubtaskFailOptions) *appfault.AppError {
	failErr := store.FailSubtask(dbPath, subID, opts.AgentRole, opts.Reason)
	if failErr != nil {
		return failErr
	}
	releaseFileClaimsBySubtask(subID)

	return reportSubtaskFail(dbPath, subID, opts)
}

func reportSubtaskFail(dbPath, subID string, opts SubtaskFailOptions) *appfault.AppError {
	if opts.IsJson {
		return renderSubtaskLifecycleJSON("failed", subID, opts.AgentRole, dbPath, opts.Reason)
	}
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

	return renderSubtaskLsOutput(subtasks, dbPath, canonicalId, opts.IsJson)
}

func renderSubtaskLsOutput(subtasks []types.Subtask, dbPath, canonicalId string, isJson bool) *appfault.AppError {
	if isJson {
		return renderSubtasksJsonEnvelope(subtasks, dbPath, canonicalId)
	}
	renderSubtasksTable(subtasks)

	return nil
}
