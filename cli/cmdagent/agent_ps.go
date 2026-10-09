package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type agentPsOptions struct {
	IsJson bool
}

var agentPsOpts = agentPsOptions{}

var agentPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running agent tasks with agent counts and subtask progress",
	RunE: func(cmd *cobra.Command, args []string) error {
		return toError(RunAgentPs(agentPsOpts))
	},
}

func init() {
	agentPsCmd.Flags().BoolVar(&agentPsOpts.IsJson, "json", false, "Output results in JSON format")
	AgentCmd.AddCommand(agentPsCmd)
}

type agentPsRow struct {
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	ActiveAgents   int    `json:"active_agents"`
	SubtasksDone   int    `json:"subtasks_done"`
	SubtasksTotal  int    `json:"subtasks_total"`
	HasSubtaskData bool   `json:"has_subtask_data"`
	UpdatedAt      string `json:"updated_at"`
}

type agentPsReport struct {
	ElapsedMs int64        `json:"elapsed_ms"`
	Rows      []agentPsRow `json:"rows"`
}

// terminalAgentStatuses are AgentRegistry statuses that mean the agent is no
// longer running. Anything else counts as active.
var terminalAgentStatuses = []string{"IDLE", "COMPLETED", "FAILED", "CRASHED"}

// RunAgentPs lists running parent tasks with active agent counts and subtask
// progress rollups. Read-only.
func RunAgentPs(opts agentPsOptions) *appfault.AppError {
	start := time.Now()
	db, openErr := openObservabilityDB()
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer db.Close()

	rows, queryErr := db.Query(`SELECT ParentTaskId, TaskSlug, TaskName, Status, RunDirectory, UpdatedAt FROM ParentTaskRegistry WHERE Status = 'ACTIVE' ORDER BY UpdatedAt DESC`)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return appfault.WrapExecution(queryErr, "RunAgentPs query")
	}
	defer rows.Close()

	report := agentPsReport{}
	for rows.Next() {
		var id string
		var r agentPsRow
		var runDir string
		scanErr := rows.Scan(&id, &r.Slug, &r.Name, &r.Status, &runDir, &r.UpdatedAt)
		hasScanErr := scanErr != nil
		if hasScanErr {
			return appfault.WrapExecution(scanErr, "RunAgentPs scan")
		}

		agents, agentsErr := countActiveAgents(db, id)
		hasAgentsErr := agentsErr != nil
		if hasAgentsErr {
			return agentsErr
		}
		r.ActiveAgents = agents

		done, total, hasData := subtaskRollup(runDir)
		r.SubtasksDone = done
		r.SubtasksTotal = total
		r.HasSubtaskData = hasData

		report.Rows = append(report.Rows, r)
	}

	report.ElapsedMs = time.Since(start).Milliseconds()
	if opts.IsJson {
		return renderPsJSON(report)
	}
	fmt.Printf("[%dms] agent ps\n", report.ElapsedMs)
	renderPsTable(report)

	return nil
}

func countActiveAgents(db *sql.DB, parentTaskId string) (int, *appfault.AppError) {
	placeholders := strings.Repeat("?,", len(terminalAgentStatuses))
	placeholders = strings.TrimSuffix(placeholders, ",")
	args := make([]any, 0, len(terminalAgentStatuses)+1)
	args = append(args, parentTaskId)
	for _, s := range terminalAgentStatuses {
		args = append(args, s)
	}

	var n int
	scanErr := db.QueryRow(`SELECT COUNT(*) FROM AgentRegistry WHERE ParentTaskId = ? AND Status NOT IN (`+placeholders+`)`, args...).Scan(&n)
	hasScanErr := scanErr != nil
	if hasScanErr {
		return 0, appfault.WrapExecution(scanErr, "countActiveAgents")
	}

	return n, nil
}

// subtaskRollup returns done/total subtask counts from a parent task's Tier 2
// DB. Missing DBs report hasData=false and are skipped gracefully.
func subtaskRollup(runDir string) (done int, total int, hasData bool) {
	hasDB := tier2DBExists(runDir)
	if !hasDB {
		return 0, 0, false
	}

	db, openErr := openTier2ReadOnly(runDir)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return 0, 0, false
	}
	defer db.Close()

	totalErr := db.QueryRow(`SELECT COUNT(*) FROM Subtask`).Scan(&total)
	doneErr := db.QueryRow(`SELECT COUNT(*) FROM Subtask WHERE Status IN ('DONE','COMPLETED')`).Scan(&done)
	hasErr := totalErr != nil || doneErr != nil
	if hasErr {
		return 0, 0, false
	}

	return done, total, true
}

func renderPsJSON(report agentPsReport) *appfault.AppError {
	b, err := json.MarshalIndent(report, "", "  ")
	hasErr := err != nil
	if hasErr {
		return appfault.WrapSimple(err, "renderPsJSON")
	}
	fmt.Println(string(b))

	return nil
}

func renderPsTable(report agentPsReport) {
	hasRows := len(report.Rows) > 0
	if !hasRows {
		fmt.Println("No running agent tasks.")
		return
	}

	fmt.Printf("%-32s  %-8s  %-6s  %-9s  %s\n", "SLUG", "STATUS", "AGENTS", "SUBTASKS", "UPDATED")
	for _, r := range report.Rows {
		subtasks := "-"
		if r.HasSubtaskData {
			subtasks = fmt.Sprintf("%d/%d", r.SubtasksDone, r.SubtasksTotal)
		}
		name := r.Slug
		fmt.Printf("%-32s  %-8s  %-6d  %-9s  %s\n", name, r.Status, r.ActiveAgents, subtasks, r.UpdatedAt)
	}
}
