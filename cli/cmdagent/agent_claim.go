package cmdagent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
