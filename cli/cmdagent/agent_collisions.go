package cmdagent

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

type agentCollisionsOptions struct {
	Parent string
	IsJson bool
}

var agentCollisionsOpts agentCollisionsOptions

var agentCollisionsCmd = &cobra.Command{
	Use:     "collisions",
	Aliases: []string{"overlap", "overlaps"},
	Short:   "Report overlapping file-box claims across active subtasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && len(agentCollisionsOpts.Parent) == 0 {
			agentCollisionsOpts.Parent = args[0]
		}

		return toError(RunAgentCollisions(agentCollisionsOpts))
	},
}

func init() {
	agentCollisionsCmd.Flags().StringVarP(&agentCollisionsOpts.Parent, "parent", "p", "", "Parent task slug (all parents when empty)")
	agentCollisionsCmd.Flags().BoolVar(&agentCollisionsOpts.IsJson, "json", false, "Output JSON")
	AgentCmd.AddCommand(agentCollisionsCmd)
}

// parentCollisionReport groups a parent task's overlaps for rendering.
type parentCollisionReport struct {
	ParentSlug string
	DbPath     string
	Collisions []fileCollision
}

// RunAgentCollisions reports live file-box overlaps. With --parent it covers one
// parent task; without it covers every parent task DB under temp-agents.
// Detected overlaps are recorded as CollisionEvent rows (deduplicated while
// unresolved). Overlaps are reported, never blocked.
func RunAgentCollisions(opts agentCollisionsOptions) *appfault.AppError {
	reports, err := collectCollisionReports(strings.TrimSpace(opts.Parent))
	if err != nil {
		return err
	}
	central, openErr := openCentralCollisionDb()
	if openErr != nil {
		return openErr
	}
	defer central.Close()
	for _, r := range reports {
		for _, c := range r.Collisions {
			for i := 0; i+1 < len(c.Owners); i++ {
				for j := i + 1; j < len(c.Owners); j++ {
					if recErr := recordCollisionEvent(central, r.ParentSlug, c.File, c.Owners[i].label(), c.Owners[j].label()); recErr != nil {
						return recErr
					}
				}
			}
		}
	}
	if opts.IsJson {
		return renderCollisionReportsJson(reports)
	}
	renderCollisionReportsTable(reports)

	return nil
}

// collectCollisionReports builds overlap reports for one parent or all parents.
func collectCollisionReports(parent string) ([]parentCollisionReport, *appfault.AppError) {
	if parent != "" {
		return collectOneParentReport(parent)
	}

	return collectAllParentsReports()
}

func collectOneParentReport(parent string) ([]parentCollisionReport, *appfault.AppError) {
	dbPath, canonicalId, findErr := resolveSubtaskParent(parent)
	if findErr != nil {
		return nil, findErr
	}
	subtasks, listErr := listSubtasksBothForms(dbPath, canonicalId, parent)
	if listErr != nil {
		return nil, listErr
	}

	return []parentCollisionReport{{
		ParentSlug: resolveParentSlugForClaims(parent),
		DbPath:     dbPath,
		Collisions: computeOverlaps(subtasks, ""),
	}}, nil
}

// collectAllParentsReports scans every per-task DB under temp-agents.
func collectAllParentsReports() ([]parentCollisionReport, *appfault.AppError) {
	tempDir := store.ResolveAgentTempDir("")
	entries, readErr := os.ReadDir(tempDir)
	if readErr != nil {
		return nil, appfault.WrapSimple(readErr, "collectAllParentsReports")
	}
	var out []parentCollisionReport
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dbPath := tempDir + "/" + e.Name() + "/agent-task.db"
		if _, statErr := os.Stat(dbPath); statErr != nil {
			continue
		}
		subtasks, listErr := listAllSubtasksInDb(dbPath)
		if listErr != nil {
			continue
		}
		out = append(out, parentCollisionReport{
			ParentSlug: e.Name(),
			DbPath:     dbPath,
			Collisions: computeOverlaps(subtasks, ""),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ParentSlug < out[j].ParentSlug })

	return out, nil
}

// listAllSubtasksInDb lists every subtask row in a per-task DB without a parent
// filter (a Tier 2 DB only holds one task's subtasks).
func listAllSubtasksInDb(dbPath string) ([]types.Subtask, *appfault.AppError) {
	conn, openErr := store.InitTaskDB(dbPath)
	if openErr != nil {
		return nil, openErr
	}
	defer conn.Close()
	rows, queryErr := conn.Query(`SELECT SubtaskId, ParentTaskId, TaskCode, TaskSlug, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask`)
	if queryErr != nil {
		return nil, appfault.WrapSimple(queryErr, "listAllSubtasksInDb")
	}
	defer rows.Close()
	var out []types.Subtask
	for rows.Next() {
		var s types.Subtask
		if scanErr := rows.Scan(&s.SubtaskId, &s.ParentTaskId, &s.TaskCode, &s.TaskSlug, &s.Title, &s.AssignedAgentRole, &s.OwnedFilesJson, &s.Status, &s.Evidence, &s.CreatedAt, &s.UpdatedAt); scanErr != nil {
			continue
		}
		out = append(out, s)
	}

	return out, nil
}

func renderCollisionReportsJson(reports []parentCollisionReport) *appfault.AppError {
	type jsonReport struct {
		ParentSlug string           `json:"parentSlug"`
		Collisions []map[string]any `json:"collisions"`
	}
	var out []jsonReport
	for _, r := range reports {
		out = append(out, jsonReport{ParentSlug: r.ParentSlug, Collisions: collisionsToMaps(r.Collisions)})
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return appfault.WrapSimple(err, "renderCollisionReportsJson")
	}
	fmt.Println(string(b))

	return nil
}

func renderCollisionReportsTable(reports []parentCollisionReport) {
	total := 0
	for _, r := range reports {
		total += len(r.Collisions)
	}
	if total == 0 {
		fmt.Println("No collisions: all active file boxes are disjoint.")
		return
	}
	fmt.Printf("COLLISIONS: %d overlapping file(s) across %d parent task(s)\n", total, len(reports))
	for _, r := range reports {
		if len(r.Collisions) == 0 {
			continue
		}
		fmt.Printf("\n[%s]\n", r.ParentSlug)
		renderCollisionList(r.Collisions)
	}
}
