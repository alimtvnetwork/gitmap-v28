package cmdssh

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var SJHistCmd = &cobra.Command{
	Use:     "history [count] [-offset] [filter]",
	Aliases: []string{"hist"},
	Short:   "Show SSH task and join history (default: 100, supports <count> and negative offset -<offset>)",
	Long: `Display recorded SSH tasks, node changes, and join history from the SSH split database.
Arguments:
  [count]    Number of records to retrieve (default: 100)
  [-offset]  Skip records from recent offset (e.g. -10 skips first 10)
  [filter]   Optional filter by target or host IP`,
	Example: `  gitmap ssh history
  gitmap ssh history 25
  gitmap ssh history -10
  gitmap ssh history 20 -10
  gitmap ssh history 192.168.1.10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSSHHistoryCLI(args)
	},
}

// RunSSHHistoryCLI handles the 'gitmap ssh history' CLI command.
func RunSSHHistoryCLI(args []string) error {
	limit, offset, filter := parseSSHHistoryArgs(args)
	return printSSHHistory(context.Background(), os.Stdout, limit, offset, filter)
}

func runSJHistory(cmd *cobra.Command, args []string, ctx context.Context) error {
	_ = cmd
	_ = ctx
	return RunSSHHistoryCLI(args)
}

func parseNegativeOffset(trimmed string) (int, bool) {
	if !strings.HasPrefix(trimmed, "-") || len(trimmed) <= 1 {
		return 0, false
	}
	val, err := strconv.Atoi(trimmed[1:])
	return val, err == nil
}

func parseSSHHistoryArgs(args []string) (int, int, string) {
	limit := 100
	offset := 0
	filter := ""

	for _, arg := range args {
		trimmed := strings.TrimSpace(arg)
		if val, isOffset := parseNegativeOffset(trimmed); isOffset {
			offset = val
			continue
		}
		if val, err := strconv.Atoi(trimmed); err == nil && val > 0 {
			limit = val
			continue
		}
		if trimmed != "" {
			filter = trimmed
		}
	}

	return limit, offset, filter
}

func printSSHHistory(ctx context.Context, out io.Writer, limit, offset int, filter string) error {
	tasks, errTasks := ListSSHHistoryTasks(ctx, limit, offset)
	if errTasks == nil && len(tasks) > 0 {
		return renderTaskHistoryTable(out, tasks, limit, offset, filter)
	}

	return printLegacySJHistory(ctx, out, limit, offset, filter)
}

func renderTaskHistoryTable(out io.Writer, tasks []SSHHistoryRecord, limit, offset int, filter string) error {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%sTASK_ID\tACTION\tTARGET\tSTATUS\tRESTORED\tCREATED_AT%s\n", constants.ColorCyan, constants.ColorReset)

	displayed := 0
	for _, t := range tasks {
		if filter != "" && !strings.Contains(t.Target, filter) && !strings.Contains(t.Action, filter) {
			continue
		}
		restored := "-"
		if t.RestoredAt != "" {
			restored = t.RestoredAt
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			t.TaskId,
			t.Action,
			t.Target,
			t.Status,
			restored,
			t.CreatedAt,
		)
		displayed++
	}

	if err := w.Flush(); err != nil {
		return apperror.WrapSimple(err, "flush history table")
	}

	fmt.Fprintf(out, "\n  Showing %d task(s) (limit: %d, offset: %d). Undo anytime: gitmap ssh undo\n", displayed, limit, offset)
	return nil
}

func printLegacySJHistory(ctx context.Context, out io.Writer, limit, offset int, filter string) error {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return apperror.New("printSJHistory", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to open db", "err": err.Error()})
	}
	defer dbConn.Close()

	_ = dbConn.Migrate()

	history, errList := store.ListSSHHistory(ctx, limit, offset, dbConn.SQL())
	if errList != nil {
		return apperror.New("printSJHistory", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to list history", "err": errList.Error()})
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%sID\tHOST_IP\tUSER\tJOINED_AT%s\n", constants.ColorCyan, constants.ColorReset)

	displayed := 0
	for _, h := range history {
		if filter != "" && h.HostIP != filter && h.User != filter {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			h.ID,
			h.HostIP,
			h.User,
			h.JoinedAt.Format("2006-01-02 15:04:05"),
		)
		displayed++
	}

	if err := w.Flush(); err != nil {
		return apperror.WrapSimple(err, "flush legacy history")
	}

	fmt.Fprintf(out, "\n  Showing %d entry(ies) (limit: %d, offset: %d).\n", displayed, limit, offset)
	return nil
}
