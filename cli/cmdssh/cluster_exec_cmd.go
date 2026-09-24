package cmdssh

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type clusterExecOptions struct {
	target      string
	command     string
	isSudo      bool
	hasForceAll bool
	parallel    int
	isShowHelp  bool
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

func isForceAllFlag(arg string) bool {
	return arg == "--force-all" || arg == "-f" || arg == "force-all"
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

func checkBoolFlag(arg string, opts *clusterExecOptions) bool {
	if isClusterHelpFlag(arg) {
		opts.isShowHelp = true
		return true
	}
	if isSudoFlag(arg) {
		opts.isSudo = true
		return true
	}
	if isForceAllFlag(arg) {
		opts.hasForceAll = true
		return true
	}
	return false
}

func parseParallelOrHelp(arg string, nextArg string, opts *clusterExecOptions) (int, bool, error) {
	if checkBoolFlag(arg, opts) {
		return 1, true, nil
	}
	val, consumed, err := parseParallelFlag(arg, nextArg)
	if err != nil {
		return 0, false, apperror.WrapValidation(err, "invalid parallel value")
	}
	if consumed > 0 {
		opts.parallel = val
		return consumed, true, nil
	}
	return 0, false, nil
}
