package cmdagent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

var slugSeparatorRuns = regexp.MustCompile(`[ \-]+`)

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
		Use:   "start <subtask-id>",
		Short: "Mark a claimed subtask as in-progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				subtaskStartOpts.SubtaskId = args[0]
			}
			return toError(RunSubtaskStart(subtaskStartOpts))
		},
	}

	subtaskCompleteCmd = &cobra.Command{
		Use:   "complete <subtask-id>",
		Short: "Mark a subtask as done with verification evidence",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				subtaskCompleteOpts.SubtaskId = args[0]
			}
			return toError(RunSubtaskComplete(subtaskCompleteOpts))
		},
	}

	subtaskFailCmd = &cobra.Command{
		Use:   "fail <subtask-id>",
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
	ParentId   string
	JsonInput  string
	Title      string
	TaskCode   string
	Slug       string
	OwnedFiles string
}

type SubtaskClaimOptions struct {
	AgentRole string
	TaskId    string
}

type SubtaskStartOptions struct {
	SubtaskId string
	AgentRole string
	TaskId    string
}

type SubtaskCompleteOptions struct {
	SubtaskId string
	AgentRole string
	Evidence  string
	TaskId    string
}

type SubtaskFailOptions struct {
	SubtaskId string
	AgentRole string
	Reason    string
	TaskId    string
}

type SubtaskLsOptions struct {
	TaskId string
	IsJson bool
}

func initSubtaskCommands() {
	initSubtaskAddFlags()
	initSubtaskClaimFlags()
	initSubtaskStartFlags()
	initSubtaskCompleteFlags()
	initSubtaskFailFlags()
	initSubtaskLsFlags()

	SubtaskCmd.AddCommand(subtaskAddCmd)
	SubtaskCmd.AddCommand(subtaskClaimCmd)
	SubtaskCmd.AddCommand(subtaskStartCmd)
	SubtaskCmd.AddCommand(subtaskCompleteCmd)
	SubtaskCmd.AddCommand(subtaskFailCmd)
	SubtaskCmd.AddCommand(subtaskLsCmd)
}

func initSubtaskAddFlags() {
	subtaskAddCmd.Flags().StringVarP(&subtaskAddOpts.ParentId, "parent", "p", "", "Parent task ID or slug")
	subtaskAddCmd.Flags().StringVarP(&subtaskAddOpts.JsonInput, "json", "j", "", "JSON array or single subtask object")
	subtaskAddCmd.Flags().StringVar(&subtaskAddOpts.Title, "title", "", "Subtask title")
	subtaskAddCmd.Flags().StringVar(&subtaskAddOpts.TaskCode, "code", "", "Task code (e.g. Subtask-01)")
	subtaskAddCmd.Flags().StringVar(&subtaskAddOpts.Slug, "slug", "", "Subtask slug (Title Case accepted; sanitized on store)")
	subtaskAddCmd.Flags().StringVar(&subtaskAddOpts.OwnedFiles, "owned", "[]", "JSON array of owned files")
}

func initSubtaskClaimFlags() {
	subtaskClaimCmd.Flags().StringVarP(&subtaskClaimOpts.AgentRole, "agent", "a", "", "Agent role claiming subtask")
	subtaskClaimCmd.Flags().StringVarP(&subtaskClaimOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
}

func initSubtaskStartFlags() {
	subtaskStartCmd.Flags().StringVarP(&subtaskStartOpts.AgentRole, "agent", "a", "", "Agent role starting subtask")
	subtaskStartCmd.Flags().StringVarP(&subtaskStartOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
}

func initSubtaskCompleteFlags() {
	subtaskCompleteCmd.Flags().StringVarP(&subtaskCompleteOpts.AgentRole, "agent", "a", "", "Agent role completing subtask")
	subtaskCompleteCmd.Flags().StringVarP(&subtaskCompleteOpts.Evidence, "evidence", "e", "", "Completion evidence or summary")
	subtaskCompleteCmd.Flags().StringVarP(&subtaskCompleteOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
}

func initSubtaskFailFlags() {
	subtaskFailCmd.Flags().StringVarP(&subtaskFailOpts.AgentRole, "agent", "a", "", "Agent role failing subtask")
	subtaskFailCmd.Flags().StringVarP(&subtaskFailOpts.Reason, "reason", "r", "", "Failure reason or RCA")
	subtaskFailCmd.Flags().StringVarP(&subtaskFailOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
}

func initSubtaskLsFlags() {
	subtaskLsCmd.Flags().StringVarP(&subtaskLsOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
	subtaskLsCmd.Flags().BoolVar(&subtaskLsOpts.IsJson, "json", false, "Output results in JSON format")
}

// RunSubtaskAdd parses and enqueues subtasks into the parent task DB.
func RunSubtaskAdd(opts SubtaskAddOptions) *appfault.AppError {
	subtasks, parseErr := parseSubtasksInput(opts)
	hasParseErr := parseErr != nil
	if hasParseErr {
		return parseErr
	}
	dbPath, canonicalParentId, findErr := resolveSubtaskParent(opts.ParentId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	subtasks = normalizeSubtaskParents(subtasks, opts.ParentId, canonicalParentId)
	subtasks = applySubtaskSlugs(subtasks, opts)
	addErr := store.AddSubtasks(dbPath, canonicalParentId, subtasks)
	hasAddErr := addErr != nil
	if hasAddErr {
		return addErr
	}
	fmt.Printf("Successfully added %d subtask(s) to %s\n", len(subtasks), dbPath)

	return nil
}

// resolveSubtaskParent resolves the parent reference to its Tier 2 DB path and
// canonical task ID. Slug input (any case) matches via the sanitized Tier 1
// lookup; anything else falls back to the existing ID-or-latest resolution.
func resolveSubtaskParent(parentId string) (string, string, *appfault.AppError) {
	trimmed := strings.TrimSpace(parentId)
	hasInput := len(trimmed) > 0
	if hasInput {
		masterDb := store.ResolveMasterAgentDbPath("")
		task, slugErr := store.FindParentTaskBySlug(masterDb, trimmed)
		hasFound := slugErr == nil && task != nil
		if hasFound {
			return task.RootDbPath, task.ParentTaskId, nil
		}
	}
	dbPath, findErr := store.FindTaskDbPath("", parentId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return "", "", findErr
	}

	return dbPath, parentId, nil
}

// normalizeSubtaskParents rewrites the stored ParentTaskId to the canonical
// task ID when it is empty or carries the raw --parent input form.
func normalizeSubtaskParents(subtasks []types.Subtask, rawParentId, canonicalParentId string) []types.Subtask {
	for i := range subtasks {
		current := subtasks[i].ParentTaskId
		needsRewrite := len(current) == 0 || current == rawParentId
		if needsRewrite {
			subtasks[i].ParentTaskId = canonicalParentId
		}
	}

	return subtasks
}

// applySubtaskSlugs resolves the TaskSlug for each subtask: explicit --slug or
// JSON taskSlug wins (sanitized); otherwise derived from the parent task slug
// and the subtask code.
func applySubtaskSlugs(subtasks []types.Subtask, opts SubtaskAddOptions) []types.Subtask {
	parentSlug := resolveParentSlugForSubtasks(opts.ParentId)
	for i := range subtasks {
		subtasks[i].TaskSlug = resolveSubtaskSlug(subtasks[i].TaskSlug, opts.Slug, parentSlug, subtasks[i].TaskCode)
	}

	return subtasks
}

func resolveSubtaskSlug(explicit, flagSlug, parentSlug, code string) string {
	hasFlag := len(strings.TrimSpace(flagSlug)) > 0
	if hasFlag {
		return sanitizeSlugTight(flagSlug)
	}
	hasExplicit := len(strings.TrimSpace(explicit)) > 0
	if hasExplicit {
		return sanitizeSlugTight(explicit)
	}
	base := sanitizeSlugTight(parentSlug)
	cleanCode := store.SanitizeSlug(code)
	hasBase := len(base) > 0
	if !hasBase {
		return cleanCode
	}

	return base + "-" + cleanCode
}

// sanitizeSlugTight collapses runs of spaces/hyphens before the shared
// sanitizer, so prompt-form slugs like "SEO Writing Task - Task 01" become
// "seo-writing-task-task-01" instead of "seo-writing-task---task-01".
// The shared store.SanitizeSlug is intentionally left untouched: existing
// stored slugs were produced by it and lookups must keep matching them.
func sanitizeSlugTight(input string) string {
	collapsed := slugSeparatorRuns.ReplaceAllString(strings.TrimSpace(input), "-")

	return store.SanitizeSlug(collapsed)
}

func resolveParentSlugForSubtasks(parentId string) string {
	hasId := len(strings.TrimSpace(parentId)) > 0
	if !hasId {
		return ""
	}
	masterDb := store.ResolveMasterAgentDbPath("")
	task, findErr := store.FindParentTaskBySlug(masterDb, parentId)
	hasMiss := findErr != nil || task == nil
	if hasMiss {
		return ""
	}

	return task.TaskSlug
}

func parseSubtasksInput(opts SubtaskAddOptions) ([]types.Subtask, *appfault.AppError) {
	hasJson := len(strings.TrimSpace(opts.JsonInput)) > 0
	if hasJson {
		return parseSubtasksJson(opts.JsonInput, opts.ParentId)
	}
	hasTitle := len(strings.TrimSpace(opts.Title)) > 0
	if hasTitle {
		return buildSingleSubtaskFromFlags(opts), nil
	}

	return nil, appfault.NewValidationError("subtask JSON (--json) or title (--title) is required")
}

func parseSubtasksJson(raw, parentId string) ([]types.Subtask, *appfault.AppError) {
	trimmed := strings.TrimSpace(raw)
	isArr := strings.HasPrefix(trimmed, "[")
	if isArr {
		return parseSubtaskArray(trimmed, parentId)
	}

	return parseSingleSubtask(trimmed, parentId)
}

func parseSubtaskArray(trimmed, parentId string) ([]types.Subtask, *appfault.AppError) {
	var list []types.Subtask
	err := json.Unmarshal([]byte(trimmed), &list)
	if err != nil {
		return nil, appfault.WrapSimple(err, "parseSubtaskArray")
	}
	normalizeSubtaskList(list, parentId)

	return list, nil
}

func parseSingleSubtask(trimmed, parentId string) ([]types.Subtask, *appfault.AppError) {
	var single types.Subtask
	err := json.Unmarshal([]byte(trimmed), &single)
	if err != nil {
		return nil, appfault.WrapSimple(err, "parseSingleSubtask")
	}
	normalizeSingleSubtask(&single, parentId)

	return []types.Subtask{single}, nil
}

func normalizeSubtaskList(list []types.Subtask, parentId string) {
	for i := range list {
		normalizeSingleSubtask(&list[i], parentId)
	}
}

func normalizeSingleSubtask(s *types.Subtask, parentId string) {
	hasParent := len(s.ParentTaskId) > 0
	if !hasParent {
		s.ParentTaskId = parentId
	}
	hasId := len(s.SubtaskId) > 0
	if !hasId {
		s.SubtaskId = generateSubtaskId(s.TaskCode)
	}
	hasStatus := len(s.Status) > 0
	if !hasStatus {
		s.Status = "PENDING"
	}
}

func generateSubtaskId(code string) string {
	ts := time.Now().UTC().Format("150405")
	cleanCode := store.SanitizeSlug(code)

	return fmt.Sprintf("sub-%s-%s", ts, cleanCode)
}

func buildSingleSubtaskFromFlags(opts SubtaskAddOptions) []types.Subtask {
	code := opts.TaskCode
	hasCode := len(code) > 0
	if !hasCode {
		code = "task"
	}
	sub := types.Subtask{
		SubtaskId:      generateSubtaskId(code),
		ParentTaskId:   opts.ParentId,
		TaskCode:       code,
		Title:          opts.Title,
		OwnedFilesJson: opts.OwnedFiles,
		Status:         "PENDING",
	}

	return []types.Subtask{sub}
}

// RunSubtaskClaim claims the next available pending subtask for an agent.
func RunSubtaskClaim(opts SubtaskClaimOptions) *appfault.AppError {
	hasRole := len(strings.TrimSpace(opts.AgentRole)) > 0
	if !hasRole {
		return appfault.NewValidationError("agent role is required (--agent)")
	}
	dbPath, canonicalId, findErr := resolveSubtaskParent(opts.TaskId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	sub, claimErr := claimSubtaskBothForms(dbPath, canonicalId, opts.TaskId, opts.AgentRole)
	hasClaimErr := claimErr != nil
	if hasClaimErr {
		return claimErr
	}
	renderClaimResult(sub, opts.AgentRole)

	return nil
}

// claimSubtaskBothForms tries the canonical task ID first, then the slug form,
// so subtasks stored under either parent reference remain claimable.
func claimSubtaskBothForms(dbPath, canonicalId, rawTaskId, agentRole string) (*types.Subtask, *appfault.AppError) {
	sub, claimErr := store.ClaimSubtask(dbPath, canonicalId, agentRole)
	hasClaimErr := claimErr != nil
	if hasClaimErr {
		return nil, claimErr
	}
	hasSub := sub != nil
	if hasSub {
		return sub, nil
	}
	slugForm := store.SanitizeSlug(rawTaskId)
	isSame := slugForm == canonicalId || len(strings.TrimSpace(rawTaskId)) == 0
	if isSame {
		return nil, nil
	}

	return store.ClaimSubtask(dbPath, slugForm, agentRole)
}

func renderClaimResult(sub *types.Subtask, role string) {
	hasSub := sub != nil
	if !hasSub {
		renderClaimEmpty()
		return
	}
	renderClaimPayload(sub, role)
}

func renderClaimEmpty() {
	out := map[string]any{"status": "empty", "message": "No pending subtasks available to claim"}
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

func renderClaimPayload(sub *types.Subtask, role string) {
	out := map[string]any{
		"status":            "claimed",
		"subtaskId":         sub.SubtaskId,
		"parentTaskId":      sub.ParentTaskId,
		"taskCode":          sub.TaskCode,
		"title":             sub.Title,
		"assignedAgentRole": role,
		"ownedFilesJson":    sub.OwnedFilesJson,
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

// RunSubtaskStart marks a subtask as in-progress.
func RunSubtaskStart(opts SubtaskStartOptions) *appfault.AppError {
	hasSub := len(strings.TrimSpace(opts.SubtaskId)) > 0
	if !hasSub {
		return appfault.NewValidationError("subtask ID is required")
	}
	dbPath, findErr := store.FindTaskDbPath("", opts.TaskId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	startErr := store.StartSubtask(dbPath, opts.SubtaskId, opts.AgentRole)
	hasStartErr := startErr != nil
	if hasStartErr {
		return startErr
	}
	fmt.Printf("Subtask %s started by %s\n", opts.SubtaskId, opts.AgentRole)

	return nil
}

// RunSubtaskComplete marks a subtask as done with verification evidence.
func RunSubtaskComplete(opts SubtaskCompleteOptions) *appfault.AppError {
	hasSub := len(strings.TrimSpace(opts.SubtaskId)) > 0
	if !hasSub {
		return appfault.NewValidationError("subtask ID is required")
	}
	dbPath, findErr := store.FindTaskDbPath("", opts.TaskId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	compErr := store.CompleteSubtask(dbPath, opts.SubtaskId, opts.AgentRole, opts.Evidence)
	hasCompErr := compErr != nil
	if hasCompErr {
		return compErr
	}
	releaseFileClaimsBySubtask(opts.SubtaskId)
	fmt.Printf("Subtask %s completed by %s\n", opts.SubtaskId, opts.AgentRole)

	return nil
}

// RunSubtaskFail records a subtask failure with diagnostic reason.
func RunSubtaskFail(opts SubtaskFailOptions) *appfault.AppError {
	hasSub := len(strings.TrimSpace(opts.SubtaskId)) > 0
	if !hasSub {
		return appfault.NewValidationError("subtask ID is required")
	}
	dbPath, findErr := store.FindTaskDbPath("", opts.TaskId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	failErr := store.FailSubtask(dbPath, opts.SubtaskId, opts.AgentRole, opts.Reason)
	hasFailErr := failErr != nil
	if hasFailErr {
		return failErr
	}
	releaseFileClaimsBySubtask(opts.SubtaskId)
	fmt.Printf("Subtask %s marked as FAILED by %s\n", opts.SubtaskId, opts.AgentRole)

	return nil
}

// RunSubtaskLs lists all subtasks belonging to a parent task.
func RunSubtaskLs(opts SubtaskLsOptions) *appfault.AppError {
	dbPath, canonicalId, findErr := resolveSubtaskParent(opts.TaskId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	subtasks, listErr := listSubtasksBothForms(dbPath, canonicalId, opts.TaskId)
	hasListErr := listErr != nil
	if hasListErr {
		return listErr
	}
	if opts.IsJson {
		return renderSubtasksJson(subtasks)
	}
	renderSubtasksTable(subtasks)

	return nil
}

// listSubtasksBothForms lists subtasks stored under either the canonical task
// ID or the slug form of the parent reference.
func listSubtasksBothForms(dbPath, canonicalId, rawTaskId string) ([]types.Subtask, *appfault.AppError) {
	byId, idErr := store.ListSubtasks(dbPath, canonicalId)
	hasIdErr := idErr != nil
	if hasIdErr {
		return nil, idErr
	}
	slugForm := store.SanitizeSlug(rawTaskId)
	isSame := slugForm == canonicalId || len(strings.TrimSpace(rawTaskId)) == 0
	if isSame {
		return byId, nil
	}
	bySlug, slugErr := store.ListSubtasks(dbPath, slugForm)
	hasSlugErr := slugErr != nil
	if hasSlugErr {
		return nil, slugErr
	}

	return mergeSubtasksDedupe(byId, bySlug), nil
}

func renderSubtasksJson(subtasks []types.Subtask) *appfault.AppError {
	b, err := json.MarshalIndent(subtasks, "", "  ")
	if err != nil {
		return appfault.WrapSimple(err, "renderSubtasksJson")
	}
	fmt.Println(string(b))

	return nil
}

func renderSubtasksTable(subtasks []types.Subtask) {
	hasSub := len(subtasks) > 0
	if !hasSub {
		fmt.Println("No subtasks found.")
		return
	}
	fmt.Printf("%-24s  %-14s  %-12s  %-15s  %s\n", "SUBTASK ID", "CODE", "STATUS", "AGENT", "TITLE")
	for _, s := range subtasks {
		agent := s.AssignedAgentRole
		hasAgent := len(agent) > 0
		if !hasAgent {
			agent = "-"
		}
		fmt.Printf("%-24s  %-14s  %-12s  %-15s  %s\n", s.SubtaskId, s.TaskCode, s.Status, agent, s.Title)
	}
}
