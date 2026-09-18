package cmdssh

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func hasExecutionFailures(results []ClusterRunResult) bool {
	for _, res := range results {
		if res.ExitCode != 0 || res.Err != nil {
			return true
		}
	}
	return false
}

func resolveClusterHosts(ctx context.Context, target string) ([]store.SSHHost, error) {
	dbConn, err := openClusterDBFunc()
	if err != nil {
		return nil, apperror.WrapSimple(err, "openClusterDB")
	}
	defer dbConn.Close()
	hosts, err := store.ListHostsByTarget(ctx, target, dbConn.Conn())
	if err != nil {
		return nil, apperror.Wrap(err, "ListHostsByTarget", map[string]any{"target": target})
	}
	if len(hosts) == 0 {
		return nil, apperror.NewNotFoundError(fmt.Sprintf("no cluster hosts found for target '%s'", target))
	}
	return hosts, nil
}

func executeClusterExec(ctx context.Context, opts *clusterExecOptions) error {
	hosts, err := resolveClusterHosts(ctx, opts.target)
	if err != nil {
		return err
	}
	results := dispatchClusterRunFn(ctx, hosts, opts.command, opts.isSudo, opts.parallel)
	if hasExecutionFailures(results) {
		return apperror.NewExecutionError("one or more cluster nodes failed execution")
	}
	return nil
}

// RunClusterExecCLI executes remote command on targeted cluster nodes.
func RunClusterExecCLI(args []string) error {
	opts, err := parseClusterExecArgs(args)
	if err != nil {
		return err
	}
	if opts.isShowHelp {
		return showClusterExecHelp()
	}
	return executeClusterExec(context.Background(), opts)
}
