package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/spf13/cobra"
)

var SJHistCmd = &cobra.Command{
	Use:   "history",
	Short: "Show SSH join history",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJHistory(cmd, args, cmd.Context())
	},
}

// runSJHistory handles the 'gitmap sj history' command.
func runSJHistory(cmd *cobra.Command, args []string, ctx context.Context) error {
	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}
	return printSJHistory(ctx, os.Stdout, filter)
}

// printSJHistory fetches history and formats it.
func printSJHistory(ctx context.Context, out io.Writer, filter string) error {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return apperror.New("printSJHistory", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to open db", "err": err.Error()})
	}
	defer dbConn.Close()

	// ListSSHHistory gets up to 100 history items for display
	history, err := store.ListSSHHistory(ctx, 100, 0, dbConn.SQL())
	if err != nil {
		return apperror.New("printSJHistory", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to list history", "err": err.Error()})
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tHOST_IP\tUSER\tJOINED_AT")
	
	for _, h := range history {
		if filter != "" && h.HostIP != filter && h.User != filter {
			continue // simple filter
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			h.ID,
			h.HostIP,
			h.User,
			h.JoinedAt.Format("2006-01-02 15:04:05"),
		)
	}

	if err := w.Flush(); err != nil {
		return apperror.New("printSJHistory", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to flush tabwriter", "err": err.Error()})
	}

	return nil
}
