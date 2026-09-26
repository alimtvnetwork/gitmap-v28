package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var mostConvFileFlag string

var agyMostConvCmd = &cobra.Command{
	Use:   "most-conv",
	Short: "Groups and sorts repositories by highest count of conversations",
}

var agyMostConvLsCmd = &cobra.Command{
	Use:   "ls [N]",
	Short: "List repositories by conversation count",
	RunE: func(cmd *cobra.Command, args []string) error {
		n := 10
		n = parseMostConvArgCount(args, n)
		return runMostConv(n)
	},
}

func init() {
	agyMostConvLsCmd.Flags().StringVarP(&mostConvFileFlag, "file", "f", "", "Output results to JSON file")
	agyMostConvCmd.AddCommand(agyMostConvLsCmd)
	AgyCmd.AddCommand(agyMostConvCmd)
}

type ProjectConvCount struct {
	ProjectID string `json:"projectId"`
	Count     int    `json:"count"`
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
		return err
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer conn.Close()

	query := "SELECT project_id, count(*) as c FROM conversation_summaries WHERE project_id != '' GROUP BY project_id ORDER BY c DESC LIMIT ?"
	rows, err := conn.Query(query, n)
	if err != nil {
		return err
	}
	defer rows.Close()

	var results []ProjectConvCount
	for rows.Next() {
		var pid string
		var count int
		if scanErr := rows.Scan(&pid, &count); scanErr != nil {
			continue
		}
		results = append(results, ProjectConvCount{ProjectID: pid, Count: count})
	}

	if mostConvFileFlag != "" {
		data, _ := json.MarshalIndent(results, "", "  ")
		return os.WriteFile(mostConvFileFlag, data, 0644)
	}

	fmt.Println(constants.ColorCyan + "PROJECT ID" + strings.Repeat(" ", 26) + "CONVERSATIONS" + constants.ColorReset)
	for _, r := range results {
		fmt.Printf("%-36s %d\n", r.ProjectID, r.Count)
	}
	return nil
}
