package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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
	if len(positional) == 0 {
		return nil, apperror.NewValidationError("usage: gitmap cluster exec <all|control|workers|<alias>> \"<command>\" [--sudo] [--parallel <n>]")
	}
	if isGitmapCommand(positional[0]) {
		opts.target = "all"
		opts.command = resolveGitmapCommandString(positional)
		return opts, nil
	}
	if len(positional) < 2 {
		return nil, apperror.NewValidationError("usage: gitmap cluster exec <all|control|workers|<alias>> \"<command>\" [--sudo] [--parallel <n>]")
	}
	opts.target = positional[0]
	opts.command = resolveClusterCommand(positional[1:])
	return opts, nil
}

func resolveClusterCommand(args []string) string {
	if len(args) > 0 && isGitmapCommand(args[0]) {
		return resolveGitmapCommandString(args)
	}
	return strings.Join(args, " ")
}

func showClusterExecHelp() error {
	fmt.Println("Usage: gitmap cluster exec <all|control|workers|<alias>> \"<command>\" [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -s, --sudo        Execute command with elevated sudo privileges")
	fmt.Println("  -p, --parallel    Number of concurrent worker threads (default 4)")
	fmt.Println("  -h, --help        Show this help message")
	return nil
}
