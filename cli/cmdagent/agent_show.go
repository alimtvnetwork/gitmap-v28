package cmdagent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type agentShowOptions struct {
	IsJson bool
}

var agentShowOpts = agentShowOptions{}

var agentShowCmd = &cobra.Command{
	Use:   "show <id-or-slug>",
	Short: "Show full detail for one parent agent task with subtask subtree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return toError(RunAgentShow(args[0], agentShowOpts))
	},
}

func init() {
	agentShowCmd.Flags().BoolVar(&agentShowOpts.IsJson, "json", false, "Output results in JSON format")
	AgentCmd.AddCommand(agentShowCmd)
}

type agentShowSubtask struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Title      string `json:"title"`
	Agent      string `json:"agent"`
	Status     string `json:"status"`
	FilesCount int    `json:"files_claimed"`
	Evidence   string `json:"evidence_excerpt"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type agentShowReport struct {
	ElapsedMs      int64              `json:"elapsed_ms"`
	ID             string             `json:"id"`
	Slug           string             `json:"slug"`
	Name           string             `json:"name"`
	Status         string             `json:"status"`
	RunDir         string             `json:"run_dir"`
	BudgetSteps    int                `json:"budget_steps"`
	CompletedSteps int                `json:"completed_steps"`
	SpawnedAgents  int                `json:"spawned_agents"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
	Subtasks       []agentShowSubtask `json:"subtasks,omitempty"`
	Tier2Broken    bool               `json:"tier2_broken,omitempty"`
}

// RunAgentShow shows full detail for one parent task (by ID or slug) with its
// subtask subtree. Read-only; all SQLite, no network. Benchmark line first.
func RunAgentShow(idOrSlug string, opts agentShowOptions) *appfault.AppError {
	t0 := time.Now()

	trimmed := strings.TrimSpace(idOrSlug)
	hasInput := trimmed != ""
	if !hasInput {
		return appfault.NewValidationError("task ID or slug is required")
	}

	db, openErr := openObservabilityDB()
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer db.Close()

	var r agentShowReport
	var runDir string
	row := db.QueryRow(`SELECT ParentTaskId, TaskSlug, TaskName, RunDirectory, Status, TotalStepsBudget, CompletedSteps, SpawnedAgentCount, CreatedAt, UpdatedAt FROM ParentTaskRegistry WHERE ParentTaskId = ? OR TaskSlug = ? LIMIT 1`, trimmed, trimmed)
	scanErr := row.Scan(&r.ID, &r.Slug, &r.Name, &runDir, &r.Status, &r.BudgetSteps, &r.CompletedSteps, &r.SpawnedAgents, &r.CreatedAt, &r.UpdatedAt)
	hasScanErr := scanErr != nil
	if hasScanErr {
		return appfault.NewValidationError(fmt.Sprintf("no task found for %q", trimmed))
	}
	r.RunDir = runDir

	subs, ok := loadShowSubtasks(runDir)
	if !ok {
		r.Tier2Broken = true
	} else {
		r.Subtasks = subs
	}

	r.ElapsedMs = time.Since(t0).Milliseconds()

	if opts.IsJson {
		return renderShowJSON(r)
	}
	renderShowTable(r)

	return nil
}

// loadShowSubtasks reads full subtask detail from a Tier-2 DB. ok=false means
// missing/unreadable (caller shows BROKEN, never crashes).
func loadShowSubtasks(runDir string) (subs []agentShowSubtask, ok bool) {
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

	rows, queryErr := db.Query(`SELECT SubtaskId, TaskCode, Title, AssignedAgentRole, OwnedFilesJson, Status, Evidence, CreatedAt, UpdatedAt FROM Subtask ORDER BY CreatedAt ASC, SubtaskId ASC`)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return nil, false
	}
	defer rows.Close()

	for rows.Next() {
		var s agentShowSubtask
		var ownedJson, evidence string
		scanErr := rows.Scan(&s.ID, &s.Code, &s.Title, &s.Agent, &ownedJson, &s.Status, &evidence, &s.CreatedAt, &s.UpdatedAt)
		hasScanErr := scanErr != nil
		if hasScanErr {
			return nil, false
		}
		s.FilesCount = countOwnedFiles(ownedJson)
		s.Evidence = excerpt(evidence, 120)
		subs = append(subs, s)
	}

	return subs, true
}

// countOwnedFiles counts entries in an OwnedFilesJson array. Malformed input
// counts as zero (never crashes).
func countOwnedFiles(raw string) int {
	trimmed := strings.TrimSpace(raw)
	hasContent := trimmed != ""
	if !hasContent {
		return 0
	}
	var files []string
	unmarshalErr := json.Unmarshal([]byte(trimmed), &files)
	hasUnmarshalErr := unmarshalErr != nil
	if hasUnmarshalErr {
		return 0
	}

	return len(files)
}

// excerpt shortens text to max runes, appending … when truncated.
func excerpt(text string, max int) string {
	trimmed := strings.TrimSpace(text)
	runes := []rune(trimmed)
	fits := len(runes) <= max
	if fits {
		return trimmed
	}

	return string(runes[:max]) + "…"
}

func renderShowJSON(r agentShowReport) *appfault.AppError {
	b, err := json.MarshalIndent(r, "", "  ")
	hasErr := err != nil
	if hasErr {
		return appfault.WrapSimple(err, "renderShowJSON")
	}
	fmt.Println(string(b))

	return nil
}

func renderShowTable(r agentShowReport) {
	fmt.Printf("[%dms] agent show %s\n", r.ElapsedMs, r.Slug)
	fmt.Printf("ID:        %s\n", r.ID)
	fmt.Printf("Slug:      %s\n", r.Slug)
	fmt.Printf("Name:      %s\n", r.Name)
	fmt.Printf("Status:    %s\n", r.Status)
	fmt.Printf("Run dir:   %s\n", r.RunDir)
	fmt.Printf("Budget:    %d steps (%d completed)\n", r.BudgetSteps, r.CompletedSteps)
	fmt.Printf("Agents:    %d spawned\n", r.SpawnedAgents)
	fmt.Printf("Created:   %s\n", r.CreatedAt)
	fmt.Printf("Updated:   %s\n", r.UpdatedAt)
	fmt.Printf("Subtasks (%d):\n", len(r.Subtasks))
	for _, s := range r.Subtasks {
		agent := s.Agent
		if strings.TrimSpace(agent) == "" {
			agent = "-"
		}
		fmt.Printf("  \u21b3 %s [%s] %s (%s)\n", s.ID, s.Code, s.Status, agent)
		fmt.Printf("    Title: %s\n", s.Title)
		fmt.Printf("    Files: %d claimed\n", s.FilesCount)
		hasEvidence := strings.TrimSpace(s.Evidence) != ""
		if hasEvidence {
			fmt.Printf("    Evidence: %s\n", s.Evidence)
		}
	}
	if r.Tier2Broken {
		fmt.Printf("  \u21b3 [BROKEN] Tier-2 DB unavailable\n")
	}
}
