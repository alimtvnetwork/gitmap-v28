package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// observabilitySchemaDDL creates the §3.5 observability tables in the Tier 1 central
// registry DB. Idempotent: safe regardless of worker ordering.
const observabilitySchemaDDL = `
CREATE TABLE IF NOT EXISTS FileClaim (
  ClaimId INTEGER PRIMARY KEY AUTOINCREMENT,
  ParentTaskSlug TEXT NOT NULL,
  SubtaskId INTEGER NOT NULL,
  AgentRole TEXT NULL,
  FilePath TEXT NOT NULL,
  ClaimKind TEXT NOT NULL DEFAULT 'write',
  Status TEXT NOT NULL DEFAULT 'active',
  CreatedAt TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS CollisionEvent (
  EventId INTEGER PRIMARY KEY AUTOINCREMENT,
  ParentTaskSlug TEXT NOT NULL,
  FilePath TEXT NOT NULL,
  OwnerA TEXT NOT NULL,
  OwnerB TEXT NOT NULL,
  DetectedAt TEXT NOT NULL,
  Resolved INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_fileclaim_path ON FileClaim(FilePath);
CREATE INDEX IF NOT EXISTS idx_fileclaim_parent ON FileClaim(ParentTaskSlug);
CREATE INDEX IF NOT EXISTS idx_collisionevent_parent ON CollisionEvent(ParentTaskSlug);
`

// openObservabilityDB opens the Tier 1 registry DB and ensures the observability schema. Shared by stats/heatmap/ps.
func openObservabilityDB() (*sql.DB, *appfault.AppError) {
	dbPath := store.ResolveMasterAgentDbPath("")
	db, openErr := OpenSqliteDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return nil, openErr
	}

	_, execErr := db.Exec(observabilitySchemaDDL)
	hasExecErr := execErr != nil
	if hasExecErr {
		_ = db.Close()
		return nil, appfault.WrapExecution(execErr, "failed to init observability schema")
	}

	return db, nil
}

// tier2DBExists reports whether a parent task's Tier 2 DB file exists, so read-only commands never create stray DB files.
func tier2DBExists(runDir string) bool {
	hasDir := strings.TrimSpace(runDir) != ""
	if !hasDir {
		return false
	}

	info, statErr := os.Stat(filepath.Join(runDir, Tier2TaskDBName))
	hasStatErr := statErr != nil

	return !hasStatErr && !info.IsDir()
}

// openTier2ReadOnly opens an existing Tier 2 DB for reads. The caller closes it.
func openTier2ReadOnly(runDir string) (*sql.DB, *appfault.AppError) {
	dbPath := filepath.ToSlash(filepath.Join(runDir, Tier2TaskDBName))

	return OpenSqliteDB(dbPath)
}

type agentStatsOptions struct {
	Days   int
	IsJson bool
}

var agentStatsOpts = agentStatsOptions{}

var agentStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show recently completed tasks, per-agent completions, and collision totals",
	RunE: func(cmd *cobra.Command, args []string) error {
		return toError(RunAgentStats(agentStatsOpts))
	},
}

func init() {
	agentStatsCmd.Flags().IntVar(&agentStatsOpts.Days, "days", 7, "Look back N days")
	agentStatsCmd.Flags().BoolVar(&agentStatsOpts.IsJson, "json", false, "Output results in JSON format")
	AgentCmd.AddCommand(agentStatsCmd)
}

type statsCompletedTask struct {
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	SubtasksDone    int    `json:"subtasks_done"`
	SubtasksTotal   int    `json:"subtasks_total"`
	CompletedAt     string `json:"completed_at"`
}

type statsAgentCount struct {
	AgentRole string `json:"agent_role"`
	Completed int    `json:"completed"`
}

type agentStatsReport struct {
	Days                   int                 `json:"days"`
	CompletedTasks         []statsCompletedTask `json:"completed_tasks"`
	TotalSubtasksCompleted int                 `json:"total_subtasks_completed"`
	PerAgentCompleted      []statsAgentCount   `json:"per_agent_completed"`
	TotalCollisions        int                 `json:"total_collisions"`
}

// RunAgentStats reports recently completed parent tasks, subtask completions,
// per-agent completion counts, and total collision events. Read-only.
func RunAgentStats(opts agentStatsOptions) *appfault.AppError {
	days := opts.Days
	hasDays := days > 0
	if !hasDays {
		days = 7
	}

	db, openErr := openObservabilityDB()
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer db.Close()

	cutoff := time.Now().UTC().AddDate(0, 0, -days).Format(time.RFC3339)

	tasks, tasksErr := queryCompletedTasks(db, cutoff)
	hasTasksErr := tasksErr != nil
	if hasTasksErr {
		return tasksErr
	}

	report := agentStatsReport{Days: days}
	agentTotals := map[string]int{}
	for _, r := range tasks {
		done, total, perAgent := countTier2Completions(r.RunDirForLookup)
		report.CompletedTasks = append(report.CompletedTasks, statsCompletedTask{
			Slug:          r.Slug,
			Name:          r.Name,
			SubtasksDone:  done,
			SubtasksTotal: total,
			CompletedAt:   r.CompletedAt,
		})
		report.TotalSubtasksCompleted += done
		for role, n := range perAgent {
			agentTotals[role] += n
		}
	}
	for role, n := range agentTotals {
		report.PerAgentCompleted = append(report.PerAgentCompleted, statsAgentCount{AgentRole: role, Completed: n})
	}
	sort.Slice(report.PerAgentCompleted, func(a, b int) bool {
		return report.PerAgentCompleted[a].Completed > report.PerAgentCompleted[b].Completed
	})

	collisions, collErr := countCollisions(db, cutoff)
	hasCollErr := collErr != nil
	if hasCollErr {
		return collErr
	}
	report.TotalCollisions = collisions

	if opts.IsJson {
		return renderAgentStatsJSON(report)
	}
	renderAgentStatsTable(report)

	return nil
}

type completedTaskRow struct {
	Slug            string
	Name            string
	CompletedAt     string
	RunDirForLookup string
}

func queryCompletedTasks(db *sql.DB, cutoff string) ([]completedTaskRow, *appfault.AppError) {
	rows, queryErr := db.Query(`SELECT TaskSlug, TaskName, RunDirectory, UpdatedAt FROM ParentTaskRegistry WHERE Status = 'COMPLETED' AND UpdatedAt >= ? ORDER BY UpdatedAt DESC`, cutoff)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return nil, appfault.WrapExecution(queryErr, "queryCompletedTasks")
	}
	defer rows.Close()

	var out []completedTaskRow
	for rows.Next() {
		var r completedTaskRow
		scanErr := rows.Scan(&r.Slug, &r.Name, &r.RunDirForLookup, &r.CompletedAt)
		hasScanErr := scanErr != nil
		if hasScanErr {
			return nil, appfault.WrapExecution(scanErr, "queryCompletedTasks scan")
		}
		out = append(out, r)
	}

	return out, nil
}

// countTier2Completions counts completed/total subtasks and per-role completions from a parent task's Tier 2 DB.
func countTier2Completions(runDir string) (done int, total int, perAgent map[string]int) {
	perAgent = map[string]int{}
	hasDB := tier2DBExists(runDir)
	if !hasDB {
		return 0, 0, perAgent
	}

	db, openErr := openTier2ReadOnly(runDir)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return 0, 0, perAgent
	}
	defer db.Close()

	totalErr := db.QueryRow(`SELECT COUNT(*) FROM Subtask`).Scan(&total)
	doneErr := db.QueryRow(`SELECT COUNT(*) FROM Subtask WHERE Status IN ('DONE','COMPLETED')`).Scan(&done)
	hasCountErr := totalErr != nil || doneErr != nil
	if hasCountErr {
		return 0, 0, perAgent
	}

	rows, queryErr := db.Query(`SELECT AssignedAgentRole, COUNT(*) FROM Subtask WHERE Status IN ('DONE','COMPLETED') AND AssignedAgentRole IS NOT NULL AND AssignedAgentRole != '' GROUP BY AssignedAgentRole`)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return done, total, perAgent
	}
	defer rows.Close()
	for rows.Next() {
		var role string
		var n int
		hasScanErr := rows.Scan(&role, &n) != nil
		if hasScanErr {
			continue
		}
		perAgent[role] = n
	}

	return done, total, perAgent
}

func countCollisions(db *sql.DB, cutoff string) (int, *appfault.AppError) {
	var n int
	scanErr := db.QueryRow(`SELECT COUNT(*) FROM CollisionEvent WHERE DetectedAt >= ?`, cutoff).Scan(&n)
	hasScanErr := scanErr != nil
	if hasScanErr {
		return 0, appfault.WrapExecution(scanErr, "countCollisions")
	}

	return n, nil
}

func renderAgentStatsJSON(report agentStatsReport) *appfault.AppError {
	b, err := json.MarshalIndent(report, "", "  ")
	hasErr := err != nil
	if hasErr {
		return appfault.WrapSimple(err, "renderAgentStatsJSON")
	}
	fmt.Println(string(b))

	return nil
}

func renderAgentStatsTable(report agentStatsReport) {
	fmt.Printf("Agent stats — last %d day(s)\n\n", report.Days)

	hasTasks := len(report.CompletedTasks) > 0
	if !hasTasks {
		fmt.Println("No completed tasks in range.")
	} else {
		fmt.Printf("%-36s  %-9s  %s\n", "SLUG", "SUBTASKS", "COMPLETED AT")
		for _, t := range report.CompletedTasks {
			fmt.Printf("%-36s  %4d/%-4d  %s\n", t.Slug, t.SubtasksDone, t.SubtasksTotal, t.CompletedAt)
		}
	}

	fmt.Printf("\nSubtasks completed: %d\n", report.TotalSubtasksCompleted)
	fmt.Printf("Collisions detected: %d\n", report.TotalCollisions)

	hasAgents := len(report.PerAgentCompleted) > 0
	if hasAgents {
		fmt.Printf("\n%-24s  %s\n", "AGENT ROLE", "COMPLETED")
		for _, a := range report.PerAgentCompleted {
			fmt.Printf("%-24s  %d\n", a.AgentRole, a.Completed)
		}
	}
}
