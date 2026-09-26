package cmdagy

import (
	"database/sql"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

var mostConvFileFlag string

var agyMostConvCmd = &cobra.Command{
	Use:     "most-conv",
	Aliases: []string{"most-conversation", "most-conversations"},
	Short:   "Groups and sorts repositories by highest count of conversations",
}

var agyMostConvLsCmd = &cobra.Command{
	Use:   "ls [N]",
	Short: "List repositories by conversation count with status and queued counts",
	RunE: func(cmd *cobra.Command, args []string) error {
		n := 8
		n = parseMostConvArgCount(args, n)
		return runMostConv(n)
	},
}

func init() {
	agyMostConvLsCmd.Flags().StringVarP(&mostConvFileFlag, "file", "f", "", "Output results to JSON file (default most-repo-conversations.json)")
	if flag := agyMostConvLsCmd.Flags().Lookup("file"); flag != nil {
		flag.NoOptDefVal = "most-repo-conversations.json"
	}
	agyMostConvCmd.AddCommand(agyMostConvLsCmd)
	AgyCmd.AddCommand(agyMostConvCmd)
}

func parseMostConvArgCount(args []string, defaultN int) int {
	if len(args) == 0 {
		return defaultN
	}
	val, err := strconv.Atoi(args[0])
	if err == nil && val > 0 {
		return val
	}
	return defaultN
}

func runMostConv(n int) error {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return apperror.WrapSimple(err, "get summaries db path")
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return apperror.WrapSimple(openErr, "open summaries db")
	}
	defer conn.Close()

	query := `SELECT project_id, MAX(workspace_uris) as ws, COUNT(*) as c, 
		MAX(CASE WHEN not_fully_idle != 0 THEN 1 ELSE 0 END) as has_running,
		MAX(last_modified_time) as last_mod
		FROM conversation_summaries 
		WHERE project_id != '' 
		GROUP BY project_id 
		ORDER BY c DESC 
		LIMIT ?`
	rows, qErr := conn.Query(query, n)
	if qErr != nil {
		return apperror.WrapSimple(qErr, "query most conversations")
	}
	defer rows.Close()

	var inspectRows []AgyConvInspectRow
	for rows.Next() {
		var pid, wsRaw, lastMod sql.NullString
		var count, hasRunning sql.NullInt64

		if scanErr := rows.Scan(&pid, &wsRaw, &count, &hasRunning, &lastMod); scanErr != nil {
			continue
		}

		wsPath := extractCleanWorkspaceFromURIs(wsRaw.String)
		isRunning := hasRunning.Int64 != 0
		status := "IDLE"
		if isRunning {
			status = "RUNNING"
		}

		queuedCount := countWorkspaceQueued(wsPath)

		inspectRows = append(inspectRows, AgyConvInspectRow{
			ID:        pid.String,
			Name:      pid.String,
			Path:      wsPath,
			Messages:  int(count.Int64),
			Status:    status,
			Queued:    queuedCount,
			IsRunning: isRunning,
			LastMod:   lastMod.String,
		})
	}

	sortInspectRows(inspectRows)
	assignInspectRowSequences(inspectRows)
	CacheInspectRowSequences(inspectRows, "most-conv")

	if mostConvFileFlag != "" {
		return ExportInspectRowsToJSONFile(inspectRows, mostConvFileFlag)
	}

	RenderInspectRowsTable(inspectRows, false)
	return nil
}
