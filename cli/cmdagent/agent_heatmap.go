package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type agentHeatmapOptions struct {
	Parent string
	IsJson bool
}

var agentHeatmapOpts = agentHeatmapOptions{}

var agentHeatmapCmd = &cobra.Command{
	Use:   "heatmap",
	Short: "Rank files by write-claim count across agent tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		return toError(RunAgentHeatmap(agentHeatmapOpts))
	},
}

func init() {
	agentHeatmapCmd.Flags().StringVar(&agentHeatmapOpts.Parent, "parent", "", "Limit to one parent task slug")
	agentHeatmapCmd.Flags().BoolVar(&agentHeatmapOpts.IsJson, "json", false, "Output results in JSON format")
	AgentCmd.AddCommand(agentHeatmapCmd)
}

type heatmapRow struct {
	FilePath string `json:"file"`
	Claims   int    `json:"claims"`
	WriteClaims int `json:"write_claims"`
	ActionTouches int `json:"action_touches"`
}

type agentHeatmapReport struct {
	Parent string       `json:"parent"`
	Rows   []heatmapRow `json:"rows"`
}

// RunAgentHeatmap ranks files by write-claim count from the central FileClaim
// table plus per-task AgentActionLog touches. Read-only.
func RunAgentHeatmap(opts agentHeatmapOptions) *appfault.AppError {
	db, openErr := openObservabilityDB()
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return openErr
	}
	defer db.Close()

	parent := strings.TrimSpace(opts.Parent)
	counts := map[string]*heatmapRow{}

	claimErr := collectClaimCounts(db, parent, counts)
	hasClaimErr := claimErr != nil
	if hasClaimErr {
		return claimErr
	}

	actionErr := collectActionCounts(db, parent, counts)
	hasActionErr := actionErr != nil
	if hasActionErr {
		return actionErr
	}

	report := agentHeatmapReport{Parent: parent}
	for _, r := range counts {
		r.Claims = r.WriteClaims + r.ActionTouches
		report.Rows = append(report.Rows, *r)
	}
	sort.Slice(report.Rows, func(a, b int) bool {
		hasMore := report.Rows[a].Claims > report.Rows[b].Claims
		if report.Rows[a].Claims == report.Rows[b].Claims {
			return report.Rows[a].FilePath < report.Rows[b].FilePath
		}

		return hasMore
	})

	if opts.IsJson {
		return renderHeatmapJSON(report)
	}
	renderHeatmapTable(report)

	return nil
}

func collectClaimCounts(db *sql.DB, parent string, counts map[string]*heatmapRow) *appfault.AppError {
	query := `SELECT FilePath, ClaimKind, COUNT(*) FROM FileClaim`
	args := []any{}
	hasParent := parent != ""
	if hasParent {
		query += ` WHERE ParentTaskSlug = ?`
		args = append(args, parent)
	}
	query += ` GROUP BY FilePath, ClaimKind`

	rows, queryErr := db.Query(query, args...)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return appfault.WrapExecution(queryErr, "collectClaimCounts")
	}
	defer rows.Close()

	for rows.Next() {
		var path, kind string
		var n int
		hasScanErr := rows.Scan(&path, &kind, &n) != nil
		if hasScanErr {
			continue
		}
		row := heatmapRowFor(counts, path)
		isWrite := strings.EqualFold(kind, "write")
		if isWrite {
			row.WriteClaims += n
		}
	}

	return nil
}

// collectActionCounts merges per-task AgentActionLog.TargetFile touches into the
// heatmap. Tier 2 DBs that do not exist are skipped; the central FileClaim table
// remains the primary source.
func collectActionCounts(db *sql.DB, parent string, counts map[string]*heatmapRow) *appfault.AppError {
	runDirs, dirsErr := heatmapRunDirs(db, parent)
	hasDirsErr := dirsErr != nil
	if hasDirsErr {
		return dirsErr
	}

	for _, dir := range runDirs {
		hasDB := tier2DBExists(dir)
		if !hasDB {
			continue
		}
		tier2, openErr := openTier2ReadOnly(dir)
		hasOpenErr := openErr != nil
		if hasOpenErr {
			continue
		}

		rows, queryErr := tier2.Query(`SELECT TargetFile, COUNT(*) FROM AgentActionLog WHERE TargetFile IS NOT NULL AND TargetFile != '' GROUP BY TargetFile`)
		hasQueryErr := queryErr != nil
		if hasQueryErr {
			_ = tier2.Close()
			continue
		}
		for rows.Next() {
			var path string
			var n int
			hasScanErr := rows.Scan(&path, &n) != nil
			if hasScanErr {
				continue
			}
			heatmapRowFor(counts, path).ActionTouches += n
		}
		_ = rows.Close()
		_ = tier2.Close()
	}

	return nil
}

func heatmapRunDirs(db *sql.DB, parent string) ([]string, *appfault.AppError) {
	query := `SELECT RunDirectory FROM ParentTaskRegistry`
	args := []any{}
	hasParent := parent != ""
	if hasParent {
		query += ` WHERE TaskSlug = ?`
		args = append(args, parent)
	}

	rows, queryErr := db.Query(query, args...)
	hasQueryErr := queryErr != nil
	if hasQueryErr {
		return nil, appfault.WrapExecution(queryErr, "heatmapRunDirs")
	}
	defer rows.Close()

	var dirs []string
	for rows.Next() {
		var dir string
		hasScanErr := rows.Scan(&dir) != nil
		if hasScanErr {
			continue
		}
		dirs = append(dirs, dir)
	}

	return dirs, nil
}

func heatmapRowFor(counts map[string]*heatmapRow, path string) *heatmapRow {
	row, hasRow := counts[path]
	if !hasRow {
		row = &heatmapRow{FilePath: path}
		counts[path] = row
	}

	return row
}

func renderHeatmapJSON(report agentHeatmapReport) *appfault.AppError {
	b, err := json.MarshalIndent(report, "", "  ")
	hasErr := err != nil
	if hasErr {
		return appfault.WrapSimple(err, "renderHeatmapJSON")
	}
	fmt.Println(string(b))

	return nil
}

func renderHeatmapTable(report agentHeatmapReport) {
	scope := "all tasks"
	hasParent := report.Parent != ""
	if hasParent {
		scope = "parent " + report.Parent
	}
	fmt.Printf("File heatmap — %s\n\n", scope)

	hasRows := len(report.Rows) > 0
	if !hasRows {
		fmt.Println("No file claims recorded.")
		return
	}

	fmt.Printf("%-4s  %-6s  %s\n", "RANK", "CLAIMS", "FILE")
	for i, r := range report.Rows {
		fmt.Printf("%-4d  %-6d  %s\n", i+1, r.Claims, r.FilePath)
	}
}
