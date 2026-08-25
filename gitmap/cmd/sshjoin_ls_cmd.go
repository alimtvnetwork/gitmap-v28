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

// runSJLs handles the 'gitmap sj ls' command.
func runSJLs(cmd *cobra.Command, args []string, ctx context.Context) error {
	return printSJList(ctx, os.Stdout, 0) // Passing 0 for no limit, though max is not strictly defined in requirements
}

// printSJList fetches all hosts and formats them using text/tabwriter.
func printSJList(ctx context.Context, out io.Writer, max int) error {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return apperror.New("printSJList", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to open db", "err": err.Error()})
	}
	defer dbConn.Close()

	hosts, err := store.ListHosts(ctx, dbConn.SQL())
	if err != nil {
		return apperror.New("printSJList", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to list hosts", "err": err.Error()})
	}

	if max > 0 && len(hosts) > max {
		hosts = hosts[:max]
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tALIAS\tIP\tUSERNAME\tCREATED_AT")
	
	for _, host := range hosts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			host.ID,
			host.Alias,
			host.IP,
			host.Username,
			host.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	if err := w.Flush(); err != nil {
		return apperror.New("printSJList", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to flush tabwriter", "err": err.Error()})
	}

	return nil
}
