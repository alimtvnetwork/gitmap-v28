package cmdssh

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type clusterExecOptions struct {
	target     string
	command    string
	isSudo     bool
	parallel   int
	isShowHelp bool
}

var (
	openClusterDBFunc    = openDB
	dispatchClusterRunFn = DispatchClusterRun
)

// ClusterExecCmd represents the gitmap cluster exec command.
var ClusterExecCmd = &cobra.Command{
	Use:   "exec <all|control|workers|<alias>> <command>",
	Short: "Execute a command across cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClusterExecCLI(args)
	},
}

func isSudoFlag(arg string) bool {
	return arg == "--sudo" || arg == "-s"
}

func isClusterHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}

func parseParallelArgValue(nextArg string, arg string) (int, int, error) {
	if nextArg == "" {
		return 0, 0, apperror.NewValidationError("flag requires an argument: " + arg)
	}
	val, err := strconv.Atoi(nextArg)
	return val, 2, err
}

func parseParallelFlag(arg string, nextArg string) (int, int, error) {
	if strings.HasPrefix(arg, "--parallel=") {
		val, err := strconv.Atoi(strings.TrimPrefix(arg, "--parallel="))
		return val, 1, err
	}
	if strings.HasPrefix(arg, "-p=") {
		val, err := strconv.Atoi(strings.TrimPrefix(arg, "-p="))
		return val, 1, err
	}
	if arg == "--parallel" || arg == "-p" {
		return parseParallelArgValue(nextArg, arg)
	}
	return 0, 0, nil
}

func checkBoolFlag(arg string, isHelp *bool, isSudo *bool) bool {
	if isClusterHelpFlag(arg) {
		*isHelp = true
		return true
	}
	if isSudoFlag(arg) {
		*isSudo = true
		return true
	}
	return false
}

func parseParallelOrHelp(arg string, nextArg string, opts *clusterExecOptions) (int, bool, error) {
	if checkBoolFlag(arg, &opts.isShowHelp, &opts.isSudo) {
		return 1, true, nil
	}
	val, consumed, err := parseParallelFlag(arg, nextArg)
	if err != nil {
		return 0, false, apperror.NewValidationError(fmt.Sprintf("invalid parallel value: %v", err))
	}
	if consumed > 0 {
		opts.parallel = val
		return consumed, true, nil
	}
	return 0, false, nil
}

func parseExecArg(args []string, idx int, opts *clusterExecOptions, pos *[]string) (int, error) {
	nextArg := ""
	if idx+1 < len(args) {
		nextArg = args[idx+1]
	}
	consumed, isFlag, err := parseParallelOrHelp(args[idx], nextArg, opts)
	if err != nil {
		return idx, err
	}
	if isFlag {
		return idx + consumed, nil
	}
	*pos = append(*pos, args[idx])
	return idx + 1, nil
}

func parseClusterExecArgs(args []string) (*clusterExecOptions, error) {
	opts := &clusterExecOptions{parallel: 4}
	var positional []string
	idx := 0
	for idx < len(args) {
		nextIdx, err := parseExecArg(args, idx, opts, &positional)
		if err != nil {
			return nil, err
		}
		idx = nextIdx
	}
	return resolveExecPositional(opts, positional)
}

func resolveExecPositional(opts *clusterExecOptions, positional []string) (*clusterExecOptions, error) {
	if opts.isShowHelp {
		return opts, nil
	}
	if len(positional) < 2 {
		return nil, apperror.NewValidationError("usage: gitmap cluster exec <all|control|workers|<alias>> \"<command>\" [--sudo] [--parallel <n>]")
	}
	opts.target = positional[0]
	opts.command = strings.Join(positional[1:], " ")
	return opts, nil
}

func showClusterExecHelp() error {
	fmt.Println("Usage: gitmap cluster exec <all|control|workers|<alias>> \"<command>\" [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -s, --sudo        Execute command with elevated sudo privileges")
	fmt.Println("  -p, --parallel    Number of concurrent worker threads (default 4)")
	fmt.Println("  -h, --help        Show this help message")
	return nil
}

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
