// Package cmdagy — agy_history_cmd.go displays recorded Antigravity task history.
package cmdagy

import (
	"database/sql"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var agyHistoryCmd = &cobra.Command{
	Use:     "history",
	Aliases: []string{"hist"},
	Short:   "Show Antigravity task history and recent mutations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyHistory(args)
	},
}

func runAgyHistory(args []string) error {
	db, err := store.OpenSectionTasksDB("agy", "")
	if err != nil {
		return err
	}
	defer db.Close()
	return renderTaskHistoryRows(db)
}

func renderTaskHistoryRows(db *store.SectionTasksDB) error {
	rows, qErr := db.Conn().Query("SELECT TaskId, Action, Target, CreatedAt FROM TaskHistory ORDER BY TaskHistoryId DESC LIMIT 50")
	if qErr != nil {
		return qErr
	}
	defer rows.Close()
	count := printHistoryTable(rows)
	printHistoryFooter(count)
	return nil
}

func printHistoryTable(rows *sql.Rows) int {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "  %sTASK ID\tACTION\tTARGET\tCREATED%s\n", constants.ColorWhite, constants.ColorReset)
	count := 0
	for rows.Next() {
		var taskID, action, target, createdAt string
		if err := rows.Scan(&taskID, &action, &target, &createdAt); err == nil {
			fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", taskID, action, target, createdAt)
			count++
		}
	}
	w.Flush()
	return count
}

func printHistoryFooter(count int) {
	if count == 0 {
		fmt.Printf("\n  %sNo recorded Antigravity task history found.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	fmt.Printf("\n  %s%d Antigravity task history record(s) listed.%s\n\n", constants.ColorDim, count, constants.ColorReset)
}
