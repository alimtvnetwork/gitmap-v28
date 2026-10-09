package cmdagent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type agentLsOptions struct {
	Limit  int
	Tree   bool
	IsJson bool
}

var agentLsOpts = agentLsOptions{Limit: 12, Tree: true}

var agentLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List recent parent agent tasks with subtask subtrees and IDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return toError(RunAgentLs(agentLsOpts))
	},
}

func init() {
	agentLsCmd.Flags().IntVar(&agentLsOpts.Limit, "limit", 12, "Max parent tasks to list (1-50)")
	agentLsCmd.Flags().BoolVar(&agentLsOpts.Tree, "tree", true, "Show subtask subtrees under each parent")
	agentLsCmd.Flags().BoolVar(&agentLsOpts.IsJson, "json", false, "Output results in JSON format")
	AgentCmd.AddCommand(agentLsCmd)
}

type agentLsSubtask struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Status string `json:"status"`
	Agent  string `json:"agent"`
}

type agentLsRow struct {
	ID             string           `json:"id"`
	Slug           string           `json:"slug"`
	Status         string           `json:"status"`
	ActiveAgents   int              `json:"active_agents"`
	SubtasksDone   int              `json:"subtasks_done"`
	SubtasksTotal  int              `json:"subtasks_total"`
	HasSubtaskData bool             `json:"has_subtask_data"`
	UpdatedAt      string           `json:"updated_at"`
	Subtasks       []agentLsSubtask `json:"subtasks,omitempty"`
	Tier2Broken    bool             `json:"tier2_broken,omitempty"`
}

type agentLsReport struct {
	ElapsedMs int64        `json:"elapsed_ms"`
	Rows      []agentLsRow `json:"rows"`
}

// RunAgentLs lists the most recent parent tasks with IDs, statuses, agent
// counts, subtask rollups, and optional subtask subtrees. Read-only; all
// SQLite, no network. Benchmark line prints first.
func RunAgentLs(opts agentLsOptions) *appfault.AppError {
	t0 := time.Now()

	limit := opts.Limit
	if limit < 1 {
		limit = 12
	}
	if limit > 50 {
		limit = 50
	}

	db, openErr := openObservabilityDB()
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer db.Close()

	rows, queryErr := db.Query(`SELECT ParentTaskId, TaskSlug, Status, RunDirectory, UpdatedAt FROM ParentTaskRegistry ORDER BY UpdatedAt DESC LIMIT ?`, limit)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return appfault.WrapExecution(queryErr, "RunAgentLs query")
	}
	defer rows.Close()

	report := agentLsReport{}
	for rows.Next() {
		var r agentLsRow
		var runDir string
		scanErr := rows.Scan(&r.ID, &r.Slug, &r.Status, &runDir, &r.UpdatedAt)
		hasScanErr := scanErr != nil
		if hasScanErr {
			return appfault.WrapExecution(scanErr, "RunAgentLs scan")
		}

		agents, agentsErr := countActiveAgents(db, r.ID)
		hasAgentsErr := agentsErr != nil
		if hasAgentsErr {
			return agentsErr
		}
		r.ActiveAgents = agents

		done, total, hasData := subtaskRollup(runDir)
		r.SubtasksDone = done
		r.SubtasksTotal = total
		r.HasSubtaskData = hasData

		if opts.Tree {
			subs, ok := loadLsSubtasks(runDir)
			if !ok {
				r.Tier2Broken = true
			} else {
				r.Subtasks = subs
			}
		}

		report.Rows = append(report.Rows, r)
	}

	report.ElapsedMs = time.Since(t0).Milliseconds()

	if opts.IsJson {
		return renderLsJSON(report)
	}
	renderLsTable(report)

	return nil
}

// loadLsSubtasks reads the subtask list from a Tier-2 DB. ok=false means the
// DB is missing or unreadable (caller shows BROKEN, never crashes).
func loadLsSubtasks(runDir string) (subs []agentLsSubtask, ok bool) {
	hasDB := tier2DBExists(runDir)
	if !hasDB {
		return nil, false
	}

	db, openErr := openTier2ReadOnly(runDir)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, false
	}
	defer db.Close()

	rows, queryErr := db.Query(`SELECT SubtaskId, TaskCode, Status, AssignedAgentRole FROM Subtask ORDER BY CreatedAt ASC, SubtaskId ASC`)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return nil, false
	}
	defer rows.Close()

	for rows.Next() {
		var s agentLsSubtask
		scanErr := rows.Scan(&s.ID, &s.Code, &s.Status, &s.Agent)
		hasScanErr := scanErr != nil
		if hasScanErr {
			return nil, false
		}
		subs = append(subs, s)
	}

	return subs, true
}

func renderLsJSON(report agentLsReport) *appfault.AppError {
	b, err := json.MarshalIndent(report, "", "  ")
	hasErr := err != nil
	if hasErr {
		return appfault.WrapSimple(err, "renderLsJSON")
	}
	fmt.Println(string(b))

	return nil
}

func renderLsTable(report agentLsReport) {
	fmt.Printf("[%dms] agent ls\n", report.ElapsedMs)
	hasRows := len(report.Rows) > 0
	if !hasRows {
		fmt.Println("No agent tasks.")
		return
	}

	fmt.Printf("%-36s  %-24s  %-9s  %-6s  %-8s  %s\n", "ID", "SLUG", "STATUS", "AGENTS", "SUBTASKS", "UPDATED")
	for _, r := range report.Rows {
		subtasks := "-"
		if r.HasSubtaskData {
			subtasks = fmt.Sprintf("%d/%d", r.SubtasksDone, r.SubtasksTotal)
		}
		fmt.Printf("%-36s  %-24s  %-9s  %-6d  %-8s  %s\n", r.ID, r.Slug, r.Status, r.ActiveAgents, subtasks, r.UpdatedAt)
		for _, s := range r.Subtasks {
			agent := s.Agent
			if strings.TrimSpace(agent) == "" {
				agent = "-"
			}
			fmt.Printf("  \u21b3 %s | %s | %s | %s\n", s.ID, s.Code, s.Status, agent)
		}
		if r.Tier2Broken {
			fmt.Printf("  \u21b3 [BROKEN] Tier-2 DB unavailable\n")
		}
	}
}
