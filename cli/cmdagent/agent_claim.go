package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

// Collision-tracking DDL for the Tier 1 central registry DB. The tables are
// keyed by parent task slug so one collisions view spans every workstream.
const sqlCreateFileClaimTable = `CREATE TABLE IF NOT EXISTS FileClaim (
  ClaimId INTEGER PRIMARY KEY AUTOINCREMENT,
  ParentTaskSlug TEXT NOT NULL,
  SubtaskId INTEGER NOT NULL,
  AgentRole TEXT NULL,
  FilePath TEXT NOT NULL,
  ClaimKind TEXT NOT NULL DEFAULT 'write',
  Status TEXT NOT NULL DEFAULT 'active',
  CreatedAt TEXT NOT NULL
);`

const sqlCreateCollisionEventTable = `CREATE TABLE IF NOT EXISTS CollisionEvent (
  EventId INTEGER PRIMARY KEY AUTOINCREMENT,
  ParentTaskSlug TEXT NOT NULL,
  FilePath TEXT NOT NULL,
  OwnerA TEXT NOT NULL,
  OwnerB TEXT NOT NULL,
  DetectedAt TEXT NOT NULL,
  Resolved INTEGER NOT NULL DEFAULT 0
);`

// ensureCollisionSchema creates the FileClaim/CollisionEvent tables when missing.
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

// computeOverlaps returns every file claimed by two or more active subtasks,
// excluding excludeId (used at claim time so a subtask never collides with itself).
func computeOverlaps(subtasks []types.Subtask, excludeId string) []fileCollision {
	owners := map[string][]collisionOwner{}
	for _, s := range subtasks {
		if s.SubtaskId == excludeId || !isSubtaskActiveForClaims(s.Status) {
			continue
		}
		for _, f := range parseOwnedFilesJson(s.OwnedFilesJson) {
			owners[f] = append(owners[f], toCollisionOwner(s))
		}
	}
	var out []fileCollision
	for f, list := range owners {
		if len(list) >= 2 {
			out = append(out, fileCollision{File: f, Owners: list})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].File < out[j].File })

	return out
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

// TaskClaimedFiles returns the union of all OwnedFilesJson paths across a task's
// subtasks. Used by the commit flow for scoped staging. Returns empty (no error)
// when the task DB is missing or holds no claims.
func TaskClaimedFiles(taskSlug string) ([]string, *appfault.AppError) {
	slug := strings.TrimSpace(taskSlug)
	if slug == "" {
		return nil, nil
	}
	dbPath := store.ResolveTaskDbPath("", slug)
	if _, statErr := os.Stat(dbPath); statErr != nil {
		return nil, nil
	}
	conn, openErr := store.InitTaskDB(dbPath)
	if openErr != nil {
		return nil, openErr
	}
	defer conn.Close()
	rows, queryErr := conn.Query(`SELECT OwnedFilesJson FROM Subtask`)
	if queryErr != nil {
		return nil, appfault.WrapSimple(queryErr, "TaskClaimedFiles")
	}
	defer rows.Close()
	seen := map[string]bool{}
	var out []string
	for rows.Next() {
		var raw string
		if scanErr := rows.Scan(&raw); scanErr != nil {
			continue
		}
		for _, f := range parseOwnedFilesJson(raw) {
			if !seen[f] {
				seen[f] = true
				out = append(out, f)
			}
		}
	}

	return out, nil
}

type SubtaskClaimFilesOptions struct {
	TaskId  string
	Subtask string
	Files   string
	Agent   string
	IsJson  bool
}

var subtaskClaimFilesOpts SubtaskClaimFilesOptions

var subtaskClaimFilesCmd = &cobra.Command{
	Use:   "claim-files",
	Short: "Declare a subtask's file box (look-ahead) for collision tracking",
	RunE: func(cmd *cobra.Command, args []string) error {
		return toError(RunSubtaskClaimFiles(subtaskClaimFilesOpts))
	},
}

func init() {
	subtaskClaimFilesCmd.Flags().StringVarP(&subtaskClaimFilesOpts.TaskId, "task-id", "t", "", "Parent task ID or slug")
	subtaskClaimFilesCmd.Flags().StringVar(&subtaskClaimFilesOpts.Subtask, "subtask", "", "Subtask ID or task code")
	subtaskClaimFilesCmd.Flags().StringVar(&subtaskClaimFilesOpts.Files, "files", "", "Comma-separated relative file paths")
	subtaskClaimFilesCmd.Flags().StringVarP(&subtaskClaimFilesOpts.Agent, "agent", "a", "", "Agent role owning the box")
	subtaskClaimFilesCmd.Flags().BoolVar(&subtaskClaimFilesOpts.IsJson, "json", false, "Output JSON")
	SubtaskCmd.AddCommand(subtaskClaimFilesCmd)
}

// RunSubtaskClaimFiles records a subtask's file box: updates OwnedFilesJson in
// the per-task DB (latest claim wins), mirrors rows into the central FileClaim
// table, and reports any overlaps with other active subtasks.
func RunSubtaskClaimFiles(opts SubtaskClaimFilesOptions) *appfault.AppError {
	taskId := opts.TaskId
	subtaskRef := opts.Subtask
	filesInput := opts.Files
	agentRole := opts.Agent
	asJson := opts.IsJson
	if strings.TrimSpace(taskId) == "" {
		return appfault.NewValidationError("--task-id is required")
	}
	if strings.TrimSpace(subtaskRef) == "" {
		return appfault.NewValidationError("--subtask is required (subtask ID or task code)")
	}
	files := parseClaimFileList(filesInput)
	if len(files) == 0 {
		return appfault.NewValidationError("--files is required (comma-separated relative paths)")
	}
	dbPath, canonicalId, findErr := resolveSubtaskParent(taskId)
	if findErr != nil {
		return findErr
	}
	subtasks, listErr := listSubtasksBothForms(dbPath, canonicalId, taskId)
	if listErr != nil {
		return listErr
	}
	target := findSubtaskByIdOrCode(subtasks, strings.TrimSpace(subtaskRef))
	if target == nil {
		return appfault.NewValidationError(fmt.Sprintf("subtask %q not found under task %q", subtaskRef, taskId))
	}
	if updErr := updateSubtaskOwnedFiles(dbPath, target.SubtaskId, files); updErr != nil {
		return updErr
	}
	parentSlug := resolveParentSlugForClaims(taskId)
	central, openErr := openCentralCollisionDb()
	if openErr != nil {
		return openErr
	}
	defer central.Close()
	if supErr := supersedeFileClaims(central, parentSlug, target.SubtaskId); supErr != nil {
		return supErr
	}
	role := agentRole
	if role == "" {
		role = target.AssignedAgentRole
	}
	if insErr := insertFileClaims(central, parentSlug, target, role, files); insErr != nil {
		return insErr
	}
	collisions := detectClaimOverlaps(subtasks, target.SubtaskId, files)
	for _, c := range collisions {
		for _, o := range c.Owners {
			if recErr := recordCollisionEvent(central, parentSlug, c.File, targetLabel(target), o.label()); recErr != nil {
				return recErr
			}
		}
	}
	renderClaimFilesResult(parentSlug, target, files, collisions, asJson)

	return nil
}

func targetLabel(s *types.Subtask) string {
	return toCollisionOwner(*s).label()
}

func renderClaimFilesResult(parentSlug string, target *types.Subtask, files []string, collisions []fileCollision, asJson bool) {
	if asJson {
		out := map[string]any{
			"status":         "claimed",
			"parentTaskSlug": parentSlug,
			"subtaskId":      target.SubtaskId,
			"taskCode":       target.TaskCode,
			"files":          files,
			"collisions":     collisionsToMaps(collisions),
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))

		return
	}
	fmt.Printf("File box claimed for %s (%s): %d file(s)\n", target.TaskCode, target.SubtaskId, len(files))
	for _, f := range files {
		fmt.Printf("  + %s\n", f)
	}
	renderCollisionList(collisions)
}

func collisionsToMaps(collisions []fileCollision) []map[string]any {
	var out []map[string]any
	for _, c := range collisions {
		var owners []string
		for _, o := range c.Owners {
			owners = append(owners, o.label())
		}
		out = append(out, map[string]any{"file": c.File, "owners": owners})
	}

	return out
}

func renderCollisionList(collisions []fileCollision) {
	if len(collisions) == 0 {
		fmt.Println("No collisions: box is clean.")
		return
	}
	fmt.Printf("COLLISIONS DETECTED (%d file(s) overlap with active subtasks):\n", len(collisions))
	for _, c := range collisions {
		fmt.Printf("  ! %s\n", c.File)
		for _, o := range c.Owners {
			fmt.Printf("      - %s\n", o.label())
		}
	}
}
