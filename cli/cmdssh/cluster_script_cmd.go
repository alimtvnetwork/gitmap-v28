package cmdssh

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type clusterScriptOptions struct {
	target     string
	scriptPath string
	isSudo     bool
	parallel   int
	isShowHelp bool
}

var (
	executeNodeScriptFn     = ExecuteNodeScript
	dispatchClusterScriptFn = DispatchClusterScript
)

// ClusterScriptCmd represents the gitmap cluster run-script command.
var ClusterScriptCmd = &cobra.Command{
	Use:     "run-script <all|control|workers|<alias>> <script-path>",
	Aliases: []string{"script"},
	Short:   "Execute a local script file across cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClusterScriptCLI(args)
	},
}

func checkScriptBoolFlag(arg string, opts *clusterScriptOptions) bool {
	if isClusterHelpFlag(arg) {
		opts.isShowHelp = true
		return true
	}
	if isSudoFlag(arg) {
		opts.isSudo = true
		return true
	}
	return false
}

func parseScriptFlag(arg string, nextArg string, opts *clusterScriptOptions) (int, bool, error) {
	if checkScriptBoolFlag(arg, opts) {
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

func parseScriptArg(args []string, idx int, opts *clusterScriptOptions, pos *[]string) (int, error) {
	nextArg := ""
	if idx+1 < len(args) {
		nextArg = args[idx+1]
	}
	consumed, isFlag, err := parseScriptFlag(args[idx], nextArg, opts)
	if err != nil {
		return idx, err
	}
	if isFlag {
		return idx + consumed, nil
	}
	*pos = append(*pos, args[idx])
	return idx + 1, nil
}

func parseClusterScriptArgs(args []string) (*clusterScriptOptions, error) {
	opts := &clusterScriptOptions{parallel: 4}
	var positional []string
	idx := 0
	for idx < len(args) {
		nextIdx, err := parseScriptArg(args, idx, opts, &positional)
		if err != nil {
			return nil, err
		}
		idx = nextIdx
	}
	return resolveScriptPositional(opts, positional)
}

func resolveScriptPositional(opts *clusterScriptOptions, positional []string) (*clusterScriptOptions, error) {
	if opts.isShowHelp {
		return opts, nil
	}
	if len(positional) < 2 {
		return nil, apperror.NewValidationError("usage: gitmap cluster run-script <all|control|workers|<alias>> <script-path> [--sudo] [--parallel <n>]")
	}
	opts.target = positional[0]
	opts.scriptPath = positional[1]
	return opts, nil
}

func showClusterScriptHelp() error {
	fmt.Println("Usage: gitmap cluster run-script <all|control|workers|<alias>> <script-path> [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -s, --sudo        Execute script with elevated sudo privileges")
	fmt.Println("  -p, --parallel    Number of concurrent worker threads (default 4)")
	fmt.Println("  -h, --help        Show this help message")
	return nil
}

func readLocalScript(scriptPath string) ([]byte, error) {
	cleanPath := expandHome(scriptPath)
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, apperror.Wrap(err, "readLocalScript", map[string]any{"path": scriptPath})
	}
	if len(content) == 0 {
		return nil, apperror.NewValidationError(fmt.Sprintf("script file '%s' is empty", scriptPath))
	}
	return content, nil
}

func buildUploadCommand(scriptContent []byte, remotePath string) string {
	b64 := base64.StdEncoding.EncodeToString(scriptContent)
	return fmt.Sprintf("mkdir -p /tmp/on-the-fly-cmd && echo '%s' | base64 -d > %s && chmod +x %s", b64, remotePath, remotePath)
}

func cleanupRemoteScript(ctx context.Context, host store.SSHHost, remotePath string, password string, isSudo bool) {
	alias := resolveHostAlias(host)
	cleanupCmd := FormatClusterScriptCleanup(remotePath, password, isSudo)
	_ = runNodeDirectCmd(ctx, host, alias, cleanupCmd, password, time.Now())
}

func executeRemoteScript(ctx context.Context, host store.SSHHost, alias, remotePath, password string, isSudo bool, start time.Time) ClusterRunResult {
	execCmd := FormatClusterScriptCommand(remotePath, password, isSudo)
	return runNodeDirectCmd(ctx, host, alias, execCmd, password, start)
}

func uploadRemoteScript(ctx context.Context, host store.SSHHost, alias, uploadCmd, password string, start time.Time) ClusterRunResult {
	return runNodeDirectCmd(ctx, host, alias, uploadCmd, password, start)
}

func uploadAndPrepareScript(ctx context.Context, host store.SSHHost, alias, remotePath, password string, scriptContent []byte, start time.Time) (ClusterRunResult, bool) {
	uploadCmd := buildUploadCommand(scriptContent, remotePath)
	upRes := uploadRemoteScript(ctx, host, alias, uploadCmd, password, start)
	if upRes.ExitCode != 0 || upRes.Err != nil {
		return upRes, false
	}
	return upRes, true
}

// ExecuteNodeScript uploads, executes with optional sudo, and cleans up a script on a host.
func ExecuteNodeScript(ctx context.Context, host store.SSHHost, scriptContent []byte, isSudo bool) ClusterRunResult {
	startTime := time.Now()
	alias := resolveHostAlias(host)
	password, err := resolveNodePassword(host)
	if err != nil {
		return buildDecryptErrorResult(host, err, startTime)
	}
	remotePath := fmt.Sprintf("/tmp/on-the-fly-cmd/script-%d.sh", time.Now().UnixNano())
	upRes, isPrepared := uploadAndPrepareScript(ctx, host, alias, remotePath, password, scriptContent, startTime)
	if !isPrepared {
		return upRes
	}
	defer cleanupRemoteScript(ctx, host, remotePath, password, isSudo)
	return executeRemoteScript(ctx, host, alias, remotePath, password, isSudo, startTime)
}

func runScriptWorker(ctx context.Context, host store.SSHHost, scriptContent []byte, isSudo bool, idx int, results []ClusterRunResult, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()
	results[idx] = executeNodeScriptFn(ctx, host, scriptContent, isSudo)
}

// DispatchClusterScript executes a script across multiple cluster nodes with bounded concurrency.
func DispatchClusterScript(ctx context.Context, hosts []store.SSHHost, scriptContent []byte, isSudo bool, concurrency int) []ClusterRunResult {
	limit := resolveConcurrencyLimit(concurrency)
	results := make([]ClusterRunResult, len(hosts))
	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	for idx, host := range hosts {
		wg.Add(1)
		go runScriptWorker(ctx, host, scriptContent, isSudo, idx, results, sem, &wg)
	}
	wg.Wait()
	PrintClusterSummaryTable(results)
	return results
}

func executeClusterScript(ctx context.Context, opts *clusterScriptOptions) error {
	content, err := readLocalScript(opts.scriptPath)
	if err != nil {
		return err
	}
	hosts, err := resolveClusterHosts(ctx, opts.target)
	if err != nil {
		return err
	}
	results := dispatchClusterScriptFn(ctx, hosts, content, opts.isSudo, opts.parallel)
	if hasExecutionFailures(results) {
		return apperror.NewExecutionError("one or more cluster nodes failed script execution")
	}
	return nil
}

// RunClusterScriptCLI reads and executes a local script on targeted cluster nodes.
func RunClusterScriptCLI(args []string) error {
	opts, err := parseClusterScriptArgs(args)
	if err != nil {
		return err
	}
	if opts.isShowHelp {
		return showClusterScriptHelp()
	}
	return executeClusterScript(context.Background(), opts)
}
