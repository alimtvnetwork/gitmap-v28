package cmdssh

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var SJLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List joined SSH machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJLs(cmd, args, cmd.Context())
	},
}

var SJNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "List joined SSH machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJLs(cmd, args, cmd.Context())
	},
}

// runSJLs handles the 'gitmap sj ls' command.
//
//nolint:revive
func runSJLs(cmd *cobra.Command, args []string, ctx context.Context) error {
	return printSJList(ctx, os.Stdout, 0)
}

func fetchSJHosts(ctx context.Context) ([]store.SSHHost, error) {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return nil, apperror.WrapSimple(err, "fetchSJHosts_OpenDB")
	}
	defer dbConn.Close()
	if err := dbConn.Migrate(); err != nil {
		return nil, apperror.WrapSimple(err, "fetchSJHosts_MigrateDB")
	}
	syncSSHHostsFromConnections(ctx, dbConn.SQL())
	return store.ListHosts(ctx, dbConn.SQL())
}

func limitSJHosts(hosts []store.SSHHost, max int) []store.SSHHost {
	if max > 0 && len(hosts) > max {
		return hosts[:max]
	}
	return hosts
}

// printSJList fetches all hosts and formats them using RenderSSHHostsTable.
func printSJList(ctx context.Context, out io.Writer, max int) error {
	hosts, err := fetchSJHosts(ctx)
	if err != nil {
		return apperror.New("printSJList", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to list hosts", "err": err.Error()})
	}
	limitedHosts := limitSJHosts(hosts, max)
	return RenderSSHHostsTable(out, limitedHosts)
}
