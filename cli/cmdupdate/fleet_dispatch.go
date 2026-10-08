package cmdupdate

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// IsFleetUpdateCommand reports whether the command invocation routes to fleet update.
func IsFleetUpdateCommand(cmd string, args []string) bool {
	if isFleetUpdateCmdToken(cmd) {
		return true
	}
	if cmd != "update" {
		return false
	}
	return isFleetUpdateArgs(args)
}

func isFleetUpdateCmdToken(cmd string) bool {
	switch cmd {
	case "ua", "update-all", "updateall", "uaz", "update-all-zip", "updateallzip":
		return true
	default:
		return false
	}
}

func isFleetUpdateArgs(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(args[0])
	if first == "ls" || first == "list" {
		return true
	}
	if isAllFleetToken(first) || hasZipFlag(args) || hasIncludeOthersFlag(args) {
		return true
	}
	if hasExceptFlag(args) {
		return true
	}
	return hasFleetNodeTarget(args)
}

func hasZipFlag(args []string) bool {
	for _, a := range args {
		if isZipToken(a) {
			return true
		}
	}
	return false
}

func hasIncludeOthersFlag(args []string) bool {
	for _, a := range args {
		if isIncludeOthersToken(a) {
			return true
		}
	}
	return false
}

func isAllFleetToken(token string) bool {
	return token == "--all" || token == "-all" || token == "all" ||
		token == "all-nodes" || token == "allnodes" || token == "-a"
}

func hasExceptFlag(args []string) bool {
	for _, a := range args {
		if isExceptParam(a) {
			return true
		}
	}
	return false
}

func hasFleetNodeTarget(args []string) bool {
	for _, a := range args {
		if isAllFleetToken(a) || strings.HasPrefix(a, "--node=") || strings.HasPrefix(a, "--remote=") {
			return true
		}
	}
	return false
}

// RunFleetUpdateDispatch routes to the appropriate fleet update action.
func RunFleetUpdateDispatch(cmd string, args []string) error {
	if isFleetZipCmdToken(cmd) {
		return ExecuteFleetUpdate(append([]string{"--all", "--zip"}, args...))
	}
	if cmd == "ua" || cmd == "update-all" || cmd == "updateall" {
		return ExecuteFleetUpdate(append([]string{"--all"}, args...))
	}
	if len(args) > 0 && isLSKeyword(args[0]) {
		return ExecuteFleetUpdateLS(args[1:])
	}
	return ExecuteFleetUpdate(args)
}

func isFleetZipCmdToken(cmd string) bool {
	return cmd == "uaz" || cmd == "update-all-zip" || cmd == "updateallzip"
}

func isLSKeyword(token string) bool {
	low := strings.ToLower(token)
	return low == "ls" || low == "list"
}

// ExecuteFleetUpdate updates applications across cluster nodes in parallel.
func ExecuteFleetUpdate(args []string) error {
	opts := parseFleetUpdateOptions(args)
	targets, err := resolveFleetTargets(opts.IncludeOthers)
	if err != nil {
		return apperror.WrapSimple(err, "ExecuteFleetUpdate.resolveFleetTargets")
	}

	exclusionSet := parseFleetExclusionSet(opts.Except)
	filteredTargets, excludedCount := filterFleetTargets(targets, opts.Target, exclusionSet)
	if len(filteredTargets) == 0 {
		printFleetNoTargetsBanner(opts.Pkg, excludedCount)
		return nil
	}

	results := executeParallelFleetUpdate(filteredTargets, opts)
	renderFleetUpdateSummary(results, opts.Pkg, excludedCount)
	return nil
}

func parseFleetUpdateOptions(args []string) FleetUpdateOptions {
	opts := FleetUpdateOptions{
		Pkg: "gitmap",
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		processUpdateFlag(a, args, &i, &opts)
	}
	if opts.IsAll && opts.Pkg == "gitmap" {
		opts.Pkg = "all"
	}
	return opts
}

func processUpdateFlag(arg string, args []string, index *int, opts *FleetUpdateOptions) {
	if isAllFleetToken(arg) {
		opts.IsAll = true
		return
	}
	if isZipToken(arg) {
		opts.IsZip = true
		return
	}
	if isIncludeOthersToken(arg) {
		opts.IncludeOthers = true
		return
	}
	if isExceptParam(arg) {
		consumeExceptFlagUpdate(args, index, opts)
		return
	}
	if isTargetParam(arg) {
		consumeTargetFlagUpdate(args, index, opts)
		return
	}
	if arg == "--dry-run" {
		opts.IsDryRun = true
		return
	}
	if arg == "-f" || arg == "--force" {
		opts.IsForce = true
		return
	}
	processPositionalPkg(arg, opts)
}

func processPositionalPkg(arg string, opts *FleetUpdateOptions) {
	if strings.HasPrefix(arg, "-") || isFleetUpdateCmdToken(arg) || arg == "update" || arg == "zip" {
		return
	}
	if opts.Pkg == "gitmap" || opts.Pkg == "all" {
		opts.Pkg = arg
	}
}

func isZipToken(arg string) bool {
	low := strings.ToLower(arg)
	return low == "--zip" || low == "-zip" || low == "zip"
}

func isIncludeOthersToken(arg string) bool {
	low := strings.ToLower(arg)
	return low == "--include-others" || low == "--include-other" || low == "--includeothers" || low == "--includeother"
}

func isExceptParam(arg string) bool {
	return arg == "--except" || arg == "--excep" || arg == "--exclude" ||
		strings.HasPrefix(arg, "--except=") || strings.HasPrefix(arg, "--excep=") || strings.HasPrefix(arg, "--exclude=")
}

func isTargetParam(arg string) bool {
	return arg == "-t" || arg == "--target" || strings.HasPrefix(arg, "--target=") || strings.HasPrefix(arg, "--node=") || strings.HasPrefix(arg, "--remote=")
}

func consumeExceptFlagUpdate(args []string, index *int, opts *FleetUpdateOptions) {
	arg := args[*index]
	if strings.Contains(arg, "=") {
		opts.Except = strings.SplitN(arg, "=", 2)[1]
		return
	}
	var tokens []string
	for *index+1 < len(args) {
		next := args[*index+1]
		isFlag := strings.HasPrefix(next, "-")
		if isFlag {
			break
		}
		*index++
		tokens = append(tokens, next)
	}
	opts.Except = strings.Join(tokens, ",")
}

func consumeTargetFlagUpdate(args []string, index *int, opts *FleetUpdateOptions) {
	arg := args[*index]
	if strings.Contains(arg, "=") {
		opts.Target = strings.SplitN(arg, "=", 2)[1]
		return
	}
	if *index+1 < len(args) {
		*index++
		opts.Target = args[*index]
	}
}

func parseFleetExclusionSet(rawExcept string) map[string]bool {
	set := make(map[string]bool)
	if rawExcept == "" {
		return set
	}
	parts := strings.Split(rawExcept, ",")
	for _, p := range parts {
		cleaned := strings.ToLower(strings.TrimSpace(p))
		if cleaned != "" {
			set[cleaned] = true
		}
	}
	return set
}
