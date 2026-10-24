package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

var slugSeparatorRuns = regexp.MustCompile(`[-_\s]+`)

func validateSubtaskIdentifierParams(id, code, slug string) (string, string, string, error) {
	cleanId, cleanCode, cleanSlug := strings.TrimSpace(id), strings.TrimSpace(code), strings.TrimSpace(slug)
	hasAny := cleanId != "" || cleanCode != "" || cleanSlug != ""
	if !hasAny {
		return "", "", "", appfault.NewValidationError("subtask identifier, code, or slug is required")
	}

	return cleanId, cleanCode, cleanSlug, nil
}

// ResolveSubtaskIdentifier resolves an identifier, code, or slug to a canonical SubtaskId.
func ResolveSubtaskIdentifier(dbPath, identifier, code, slug string) (string, error) {
	cleanId, cleanCode, cleanSlug, valErr := validateSubtaskIdentifierParams(identifier, code, slug)
	if valErr != nil {
		return "", valErr
	}
	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return "", appfault.WrapSimple(openRes.AppError(), fmt.Sprintf("open subtask database at %s", dbPath))
	}
	conn := openRes.Value
	defer conn.Close()

	return querySubtaskCanonicalId(conn, dbPath, cleanId, cleanCode, cleanSlug)
}

func querySubtaskCanonicalId(conn *sql.DB, dbPath, id, code, slug string) (string, error) {
	targetId, targetCode, targetSlug := resolveLookupParams(id, code, slug)
	query := "SELECT SubtaskId FROM Subtask WHERE SubtaskId = ? OR TaskCode = ? COLLATE NOCASE OR TaskSlug = ? COLLATE NOCASE LIMIT 1"
	var canonicalId string
	err := conn.QueryRow(query, targetId, targetCode, targetSlug).Scan(&canonicalId)
	if err != nil {
		disp := chooseDisplayRef(id, code, slug)
		return "", appfault.NewNotFoundError(fmt.Sprintf("subtask '%s' not found in %s", disp, dbPath))
	}

	return canonicalId, nil
}

func resolveLookupParams(id, code, slug string) (string, string, string) {
	fallback := chooseDisplayRef(id, code, slug)

	return pickParam(id, fallback), pickParam(code, fallback), pickParam(slug, fallback)
}

func pickParam(val, fallback string) string {
	if val != "" {
		return val
	}

	return fallback
}

func chooseDisplayRef(id, code, slug string) string {
	if id != "" {
		return id
	}
	if code != "" {
		return code
	}

	return slug
}

// resolveSubtaskParent resolves the parent reference to its Tier 2 DB path and canonical task ID.
func resolveSubtaskParent(parentId string) (string, string, *appfault.AppError) {
	trimmed := strings.TrimSpace(parentId)
	path, taskID, hasMatch := resolveParentSlugMatch(trimmed)
	if hasMatch {
		return path, taskID, nil
	}
	dbPath, findErr := store.FindTaskDbPath("", parentId)
	if findErr != nil {
		return "", "", findErr
	}

	return dbPath, parentId, nil
}

func resolveParentSlugMatch(trimmed string) (string, string, bool) {
	if len(trimmed) == 0 {
		return "", "", false
	}
	masterDb := store.ResolveMasterAgentDbPath("")
	task, slugErr := store.FindParentTaskBySlug(masterDb, trimmed)
	if slugErr == nil && task != nil {
		return task.RootDbPath, task.ParentTaskId, true
	}

	return "", "", false
}

// claimSubtaskBothForms tries canonical task ID first, then slug form.
func claimSubtaskBothForms(dbPath, canonicalId, rawTaskId, agentRole string) (*types.Subtask, *appfault.AppError) {
	sub, claimErr := store.ClaimSubtask(dbPath, canonicalId, agentRole)
	if claimErr != nil || sub != nil {
		return sub, claimErr
	}
	slugForm := store.SanitizeSlug(rawTaskId)
	if slugForm == canonicalId || len(strings.TrimSpace(rawTaskId)) == 0 {
		return nil, nil
	}

	return store.ClaimSubtask(dbPath, slugForm, agentRole)
}

// listSubtasksBothForms lists subtasks stored under either canonical task ID or slug form.
func listSubtasksBothForms(dbPath, canonicalId, rawTaskId string) ([]types.Subtask, *appfault.AppError) {
	byId, idErr := store.ListSubtasks(dbPath, canonicalId)
	if idErr != nil {
		return nil, idErr
	}
	subtasks := resolveScopedOrAllSubtasks(dbPath, canonicalId, rawTaskId, byId)
	if len(subtasks) == 0 {
		all, _ := store.ListSubtasks(dbPath, "")

		return all, nil
	}

	return subtasks, nil
}

func resolveScopedOrAllSubtasks(dbPath, canonicalId, rawTaskId string, byId []types.Subtask) []types.Subtask {
	slugForm := store.SanitizeSlug(rawTaskId)
	if slugForm == canonicalId || len(strings.TrimSpace(rawTaskId)) == 0 {
		return byId
	}
	bySlug, slugErr := store.ListSubtasks(dbPath, slugForm)
	if slugErr != nil {
		return byId
	}

	return mergeSubtasksDedupe(byId, bySlug)
}

func normalizeSubtaskParents(subtasks []types.Subtask, rawParentId, canonicalParentId string) []types.Subtask {
	for i := range subtasks {
		current := subtasks[i].ParentTaskId
		if len(current) == 0 || current == rawParentId {
			subtasks[i].ParentTaskId = canonicalParentId
		}
	}

	return subtasks
}

func applySubtaskSlugs(subtasks []types.Subtask, opts SubtaskAddOptions) []types.Subtask {
	parentSlug := resolveParentSlugForSubtasks(opts.ParentId)
	for i := range subtasks {
		subtasks[i].TaskSlug = resolveSubtaskSlug(subtasks[i].TaskSlug, opts.Slug, parentSlug, subtasks[i].TaskCode)
	}

	return subtasks
}

func resolveSubtaskSlug(explicit, flagSlug, parentSlug, code string) string {
	if len(strings.TrimSpace(flagSlug)) > 0 {
		return sanitizeSlugTight(flagSlug)
	}
	if len(strings.TrimSpace(explicit)) > 0 {
		return sanitizeSlugTight(explicit)
	}
	base := sanitizeSlugTight(parentSlug)
	cleanCode := store.SanitizeSlug(code)
	if len(base) == 0 {
		return cleanCode
	}

	return base + "-" + cleanCode
}

func sanitizeSlugTight(input string) string {
	collapsed := slugSeparatorRuns.ReplaceAllString(strings.TrimSpace(input), "-")

	return store.SanitizeSlug(collapsed)
}

func resolveParentSlugForSubtasks(parentId string) string {
	if len(strings.TrimSpace(parentId)) == 0 {
		return ""
	}
	masterDb := store.ResolveMasterAgentDbPath("")
	task, findErr := store.FindParentTaskBySlug(masterDb, parentId)
	if findErr != nil || task == nil {
		return ""
	}

	return task.TaskSlug
}

func parseSubtasksInput(opts SubtaskAddOptions) ([]types.Subtask, *appfault.AppError) {
	if len(strings.TrimSpace(opts.JsonInput)) > 0 {
		return parseSubtasksJson(opts.JsonInput, opts.ParentId)
	}
	if len(strings.TrimSpace(opts.Title)) > 0 {
		return buildSingleSubtaskFromFlags(opts), nil
	}

	return nil, appfault.NewValidationError("subtask JSON (--json) or title (--title) is required")
}

func parseSubtasksJson(raw, parentId string) ([]types.Subtask, *appfault.AppError) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "[") {
		return parseSubtaskArray(trimmed, parentId)
	}

	return parseSingleSubtask(trimmed, parentId)
}

func parseSubtaskArray(trimmed, parentId string) ([]types.Subtask, *appfault.AppError) {
	var list []types.Subtask
	err := json.Unmarshal([]byte(trimmed), &list)
	if err != nil {
		return nil, appfault.WrapSimple(err, "parseSubtasksJson.array")
	}
	normalizeSubtaskList(list, parentId)

	return list, nil
}

func parseSingleSubtask(trimmed, parentId string) ([]types.Subtask, *appfault.AppError) {
	var single types.Subtask
	err := json.Unmarshal([]byte(trimmed), &single)
	if err != nil {
		return nil, appfault.WrapSimple(err, "parseSubtasksJson.single")
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
	if len(s.ParentTaskId) == 0 {
		s.ParentTaskId = parentId
	}
	if len(s.SubtaskId) == 0 {
		s.SubtaskId = generateSubtaskId(s.TaskCode)
	}
	if len(s.Status) == 0 {
		s.Status = "PENDING"
	}
}

func generateSubtaskId(code string) string {
	ts := time.Now().UTC().Format("150405")
	cleanCode := store.SanitizeSlug(code)

	return fmt.Sprintf("sub-%s-%s", ts, cleanCode)
}

func resolveSubtaskCode(code string) string {
	if len(code) == 0 {
		return "task"
	}

	return code
}

func buildSingleSubtaskFromFlags(opts SubtaskAddOptions) []types.Subtask {
	code := resolveSubtaskCode(opts.TaskCode)
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
