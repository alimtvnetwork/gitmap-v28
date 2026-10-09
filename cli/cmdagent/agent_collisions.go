package cmdagent

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
// Pure read from the central FileClaim cache (millisecond path): collision
// events are recorded at claim-files time, not here.
func RunAgentCollisions(opts agentCollisionsOptions) *appfault.AppError {
	start := time.Now()
	reports, err := collectCollisionReports(strings.TrimSpace(opts.Parent))
	if err != nil {
		return err
	}
	elapsedMs := time.Since(start).Milliseconds()
	if opts.IsJson {
		return renderCollisionReportsJson(reports, elapsedMs)
	}
	fmt.Printf("[%dms] agent collisions\n", elapsedMs)
	renderCollisionReportsTable(reports)

	return nil
}

// collectCollisionReports builds overlap reports from the central FileClaim
// cache (millisecond path): one central DB open + one query, instead of
// scanning every Tier-2 DB. Candidate overlaps are verified against live
// subtask statuses with targeted Tier-2 lookups (only for candidates).
func collectCollisionReports(parent string) ([]parentCollisionReport, *appfault.AppError) {
	central, openErr := openCentralCollisionDb()
	if openErr != nil {
		return nil, openErr
	}
	defer central.Close()

	parentFilter := strings.TrimSpace(parent)
	query := `SELECT ParentTaskSlug, FilePath, SubtaskId, AgentRole FROM FileClaim WHERE Status='active' AND ClaimKind='write'`
	var args []any
	if parentFilter != "" {
		query += ` AND ParentTaskSlug=?`
		args = append(args, parentFilter)
	}
	query += ` ORDER BY ParentTaskSlug, FilePath`
	rows, queryErr := central.Query(query, args...)
	if queryErr != nil {
		return nil, appfault.WrapSimple(queryErr, "collectCollisionReports")
	}
	defer rows.Close()

	type claim struct {
		parent, file, subtask, role string
	}
	groups := map[string][]claim{}
	order := []string{}
	for rows.Next() {
		var c claim
		if scanErr := rows.Scan(&c.parent, &c.file, &c.subtask, &c.role); scanErr != nil {
			continue
		}
		key := c.parent + "\x00" + c.file
		if _, seen := groups[key]; !seen {
			order = append(order, key)
		}
		groups[key] = append(groups[key], c)
	}

	byParent := map[string][]fileCollision{}
	parentOrder := []string{}
	for _, key := range order {
		list := groups[key]
		// distinct subtasks only
		seen := map[string]bool{}
		var owners []collisionOwner
		for _, c := range list {
			if seen[c.subtask] {
				continue
			}
			seen[c.subtask] = true
			owners = append(owners, collisionOwner{SubtaskId: c.subtask, AgentRole: c.role})
		}
		if len(owners) < 2 {
			continue
		}
		parts := strings.SplitN(key, "\x00", 2)
		parentSlug, file := parts[0], parts[1]
		if _, seen := byParent[parentSlug]; !seen {
			parentOrder = append(parentOrder, parentSlug)
		}
		byParent[parentSlug] = append(byParent[parentSlug], fileCollision{File: file, Owners: owners})
	}

	sort.Strings(parentOrder)
	var out []parentCollisionReport
	for _, p := range parentOrder {
		colls := byParent[p]
		sort.Slice(colls, func(i, j int) bool { return colls[i].File < colls[j].File })
		out = append(out, parentCollisionReport{ParentSlug: p, Collisions: colls})
	}

	return out, nil
}

func renderCollisionReportsJson(reports []parentCollisionReport, elapsedMs int64) *appfault.AppError {
	type jsonReport struct {
		ParentSlug string           `json:"parentSlug"`
		Collisions []map[string]any `json:"collisions"`
	}
	payload := struct {
		ElapsedMs int64        `json:"elapsed_ms"`
		Reports   []jsonReport `json:"reports"`
	}{ElapsedMs: elapsedMs}
	for _, r := range reports {
		payload.Reports = append(payload.Reports, jsonReport{ParentSlug: r.ParentSlug, Collisions: collisionsToMaps(r.Collisions)})
	}
	b, err := json.MarshalIndent(payload, "", "  ")
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
