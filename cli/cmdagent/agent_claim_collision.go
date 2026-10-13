package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

func ensureCollisionSchema(conn *sql.DB) *appfault.AppError {
	for _, ddl := range []string{sqlCreateFileClaimTable, sqlCreateCollisionEventTable} {
		res := store.ExecWrapper(conn, ddl)
		if res.IsFailure {
			return appfault.WrapSimple(res.Error, "ensureCollisionSchema")
		}
	}

	return nil
}

// openCentralCollisionDb opens the Tier 1 master agent DB and ensures the
// collision-tracking schema exists.
func openCentralCollisionDb() (*sql.DB, *appfault.AppError) {
	masterDb := store.ResolveMasterAgentDbPath("")
	conn, openErr := store.InitMasterAgentDB(masterDb)
	if openErr != nil {
		return nil, openErr
	}
	if schemaErr := ensureCollisionSchema(conn); schemaErr != nil {
		_ = conn.Close()

		return nil, schemaErr
	}

	return conn, nil
}

// parseClaimFileList normalizes a comma-separated file list: trims spaces,
// converts backslashes, drops empties and duplicates (order preserved).
func parseClaimFileList(input string) []string {
	seen := map[string]bool{}
	var out []string
	for _, part := range strings.Split(input, ",") {
		p := strings.TrimSpace(strings.ReplaceAll(part, "\\", "/"))
		p = strings.TrimPrefix(p, "./")
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}

	return out
}

// parseOwnedFilesJson decodes an OwnedFilesJson payload into normalized paths.
func parseOwnedFilesJson(raw string) []string {
	var files []string
	if err := json.Unmarshal([]byte(raw), &files); err != nil {
		return nil
	}
	var out []string
	for _, f := range files {
		p := strings.TrimSpace(strings.ReplaceAll(f, "\\", "/"))
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}

// encodeOwnedFilesJson renders a file list as the OwnedFilesJson payload.
func encodeOwnedFilesJson(files []string) string {
	b, err := json.Marshal(files)
	if err != nil {
		return "[]"
	}

	return string(b)
}

// resolveParentSlugForClaims returns the canonical parent task slug for a
// task-id-or-slug input, falling back to the sanitized input.
func resolveParentSlugForClaims(taskId string) string {
	masterDb := store.ResolveMasterAgentDbPath("")
	task, findErr := store.FindParentTaskBySlug(masterDb, taskId)
	hasTask := findErr == nil && task != nil && len(strings.TrimSpace(task.TaskSlug)) > 0
	if hasTask {
		return task.TaskSlug
	}

	return store.SanitizeSlug(taskId)
}

// isSubtaskActiveForClaims reports whether a subtask status still counts for
// collision detection (finished subtasks release their boxes).
func isSubtaskActiveForClaims(status string) bool {
	switch status {
	case "DONE", "COMPLETED", "FAILED":
		return false
	}

	return true
}

// collisionOwner identifies one claimant of a file.
type collisionOwner struct {
	SubtaskId string
	TaskCode  string
	AgentRole string
	Status    string
}

func toCollisionOwner(s types.Subtask) collisionOwner {
	return collisionOwner{
		SubtaskId: s.SubtaskId,
		TaskCode:  s.TaskCode,
		AgentRole: s.AssignedAgentRole,
		Status:    s.Status,
	}
}

func (o collisionOwner) label() string {
	code := o.TaskCode
	if code == "" {
		code = o.SubtaskId
	}
	agent := o.AgentRole
	if agent == "" {
		agent = "-"
	}

	return fmt.Sprintf("%s (%s) [%s]", code, o.SubtaskId, agent)
}

// fileCollision is one file claimed by two or more active subtasks.
type fileCollision struct {
	File   string
	Owners []collisionOwner
}

// detectClaimOverlaps finds files in newFiles that are also claimed by other
// active subtasks (excluding the claiming subtask itself).
func detectClaimOverlaps(subtasks []types.Subtask, targetId string, newFiles []string) []fileCollision {
	others := map[string][]collisionOwner{}
	for _, s := range subtasks {
		if s.SubtaskId == targetId || !isSubtaskActiveForClaims(s.Status) {
			continue
		}
		for _, f := range parseOwnedFilesJson(s.OwnedFilesJson) {
			others[f] = append(others[f], toCollisionOwner(s))
		}
	}
	var out []fileCollision
	for _, f := range newFiles {
		if list, hit := others[f]; hit && len(list) > 0 {
			out = append(out, fileCollision{File: f, Owners: list})
		}
	}

	return out
}

// recordCollisionEvent writes a CollisionEvent row, skipping duplicates that
// are still unresolved.
func recordCollisionEvent(conn *sql.DB, parentSlug, file, ownerA, ownerB string) *appfault.AppError {
	var count int
	row := conn.QueryRow(`SELECT COUNT(*) FROM CollisionEvent WHERE ParentTaskSlug = ? AND FilePath = ? AND ((OwnerA = ? AND OwnerB = ?) OR (OwnerA = ? AND OwnerB = ?)) AND Resolved = 0`,
		parentSlug, file, ownerA, ownerB, ownerB, ownerA)
	if scanErr := row.Scan(&count); scanErr == nil && count > 0 {
		return nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	res := store.ExecWrapper(conn, `INSERT INTO CollisionEvent (ParentTaskSlug, FilePath, OwnerA, OwnerB, DetectedAt, Resolved) VALUES (?, ?, ?, ?, ?, 0)`,
		parentSlug, file, ownerA, ownerB, now)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "recordCollisionEvent")
	}

	return nil
}

// supersedeFileClaims marks a subtask's previous active claims as superseded.
// The latest claim always wins; history stays queryable for the heatmap.
func supersedeFileClaims(conn *sql.DB, parentSlug, subtaskId string) *appfault.AppError {
	res := store.ExecWrapper(conn, `UPDATE FileClaim SET Status = 'superseded' WHERE ParentTaskSlug = ? AND SubtaskId = ? AND Status = 'active'`,
		parentSlug, subtaskId)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "supersedeFileClaims")
	}

	return nil
}

// releaseFileClaimsBySubtask marks a finished subtask's active claims as
// released, keeping the central FileClaim cache accurate for millisecond
// collision reads. Best-effort: failures are visible but never fail the caller.
func releaseFileClaimsBySubtask(subtaskId string) {
	central, openErr := openCentralCollisionDb()
	if openErr != nil {
		fmt.Fprintf(os.Stderr, "warning: could not release file claims: %v\n", openErr)

		return
	}
	defer central.Close()
	if _, execErr := central.Exec(`UPDATE FileClaim SET Status='released' WHERE SubtaskId=? AND Status='active'`, subtaskId); execErr != nil {
		fmt.Fprintf(os.Stderr, "warning: could not release file claims: %v\n", execErr)
	}
}

// insertFileClaims appends one active write claim row per file.
func insertFileClaims(conn *sql.DB, parentSlug string, sub *types.Subtask, agentRole string, files []string) *appfault.AppError {
	now := time.Now().UTC().Format(time.RFC3339)
	role := agentRole
	if role == "" {
		role = sub.AssignedAgentRole
	}
	for _, f := range files {
		res := store.ExecWrapper(conn, `INSERT INTO FileClaim (ParentTaskSlug, SubtaskId, AgentRole, FilePath, ClaimKind, Status, CreatedAt) VALUES (?, ?, ?, ?, 'write', 'active', ?)`,
			parentSlug, sub.SubtaskId, role, f, now)
		if res.IsFailure {
			return appfault.WrapSimple(res.Error, "insertFileClaims")
		}
	}

	return nil
}

// updateSubtaskOwnedFiles replaces a subtask's OwnedFilesJson (latest claim wins).
func updateSubtaskOwnedFiles(dbPath, subtaskId string, files []string) *appfault.AppError {
	conn, openErr := store.InitTaskDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer conn.Close()
	now := time.Now().UTC().Format(time.RFC3339)
	res := store.ExecWrapper(conn, `UPDATE Subtask SET OwnedFilesJson = ?, UpdatedAt = ? WHERE SubtaskId = ?`,
		encodeOwnedFilesJson(files), now, subtaskId)
	if res.IsFailure {
		return appfault.WrapSimple(res.Error, "updateSubtaskOwnedFiles")
	}

	return nil
}

// findSubtaskByIdOrCode locates a subtask by its SubtaskId or TaskCode.
func findSubtaskByIdOrCode(subtasks []types.Subtask, idOrCode string) *types.Subtask {
	for i := range subtasks {
		if subtasks[i].SubtaskId == idOrCode || subtasks[i].TaskCode == idOrCode {
			return &subtasks[i]
		}
	}

	return nil
}
