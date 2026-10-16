package cmdnodes

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

// isCommandToken reports whether a token is a command alias for running-only.
func isCommandToken(tok string) bool {
	return tok == "nodes-running-only" ||
		tok == "nodesrunningonly" ||
		tok == "nodes-running" ||
		tok == "running-only" ||
		tok == "nro" ||
		tok == "nodes-online" ||
		tok == "nodesonline" ||
		tok == "nodes-open" ||
		tok == "nodes-active" ||
		tok == "running" ||
		tok == "online" ||
		tok == "open" ||
		tok == "ro"
}

// isRunningOnlyFlag reports whether a token is an online/running filter flag.
func isRunningOnlyFlag(tok string) bool {
	return tok == "--running-only" || tok == "-ro" || tok == "--open-only" || tok == "--online-only"
}

// stripLeadingCommandTokens strips leading command tokens like "nodes-running-only", "nro", "running-only".
func stripLeadingCommandTokens(args []string) []string {
	for len(args) > 0 {
		tok := strings.ToLower(args[0])
		if isCommandToken(tok) {
			args = args[1:]
			continue
		}
		break
	}
	return args
}

// dispatchRunningOnlyDelegation handles delegation commands by forwarding with appropriate flags.
func dispatchRunningOnlyDelegation(tok string, args []string) (bool, error) {
	switch tok {
	case "exec", "run":
		return true, cmdssh.RunSSHExec(append([]string{"--running-only"}, args...))
	case "tasks", "task", "queue":
		return true, RunNodesTasksCLI(append([]string{"--open-only"}, args...))
	case "resources", "res", "top", "metrics":
		return true, RunNodesResourcesCLI(append([]string{"--open-only"}, args...))
	case "ping", "probe":
		return true, RunUnifiedNodesPingCLI(append([]string{"--running-only"}, args...))
	case "pull", "pa", "pull-all":
		return true, RunNodesPullAll(append([]string{"--running-only"}, args...))
	case "commits":
		return true, RunNodesCommits(append([]string{"--open-only"}, args...))
	case "cp":
		return true, RunNodesCommitPush(append([]string{"--open-only"}, args...))
	case "cpf":
		return true, RunNodesCommitPushFeature(append([]string{"--open-only"}, args...))
	case "cpb":
		return true, RunNodesCommitPushBug(append([]string{"--open-only"}, args...))
	case "cpr":
		return true, RunNodesCommitPushRelease(append([]string{"--open-only"}, args...))
	case "commit-fix", "commitfix":
		return true, RunNodesCommitFix(append([]string{"--open-only"}, args...))
	case "pending-commits", "pc", "pendingcommits":
		return true, RunNodesPendingCommits(append([]string{"--open-only"}, args...))
	case "scan":
		return true, RunNodesScan(append([]string{"--open-only"}, args...))
	case "rescan":
		return true, RunNodesRescan(append([]string{"--open-only"}, args...))
	default:
		return false, nil
	}
}

// handleFlagPrefixedDelegation detects delegation commands when flags precede the subcommand.
func handleFlagPrefixedDelegation(args []string) (bool, error) {
	hasEnoughArgs := len(args) > 1
	if !hasEnoughArgs {
		return false, nil
	}
	isFlag := isRunningOnlyFlag(strings.ToLower(args[0]))
	if !isFlag {
		return false, nil
	}
	subTok := strings.ToLower(args[1])
	return dispatchRunningOnlyDelegation(subTok, args[2:])
}

// isNodesRunningOnlyRequest reports whether the first argument is a running-only subcommand.
func isNodesRunningOnlyRequest(args []string) bool {
	hasArgs := len(args) > 0
	if !hasArgs {
		return false
	}
	tok := strings.ToLower(args[0])
	return isCommandToken(tok)
}

// hasNodesRunningOnlyFlag reports whether any argument is a running-only or open-only flag.
func hasNodesRunningOnlyFlag(args []string) bool {
	for _, arg := range args {
		if isRunningOnlyFlag(strings.ToLower(arg)) {
			return true
		}
	}
	return false
}

// FilterOnlineNodes returns a new slice containing only nodes where IsOnline is true.
func FilterOnlineNodes(nodes []UnifiedFleetNode) []UnifiedFleetNode {
	var online []UnifiedFleetNode
	for _, n := range nodes {
		if n.IsOnline {
			online = append(online, n)
		}
	}
	return online
}

// filterOnlineUnifiedNodes returns a new slice containing only nodes where IsOnline is true.
func filterOnlineUnifiedNodes(nodes []UnifiedFleetNode) []UnifiedFleetNode {
	return FilterOnlineNodes(nodes)
}

// HasOnlineNodes returns true if at least one online node is present.
func HasOnlineNodes(nodes []UnifiedFleetNode) bool {
	for _, n := range nodes {
		if n.IsOnline {
			return true
		}
	}
	return false
}

// RunNodesRunningOnlyCLI handles `gitmap nodes-running-only` and `gitmap nro`.
func RunNodesRunningOnlyCLI(args []string) error {
	cleanArgs := stripLeadingCommandTokens(args)
	if len(cleanArgs) == 0 {
		return runUnifiedNodesList([]string{"--open-only"})
	}

	isHelp := isNodesHelpRequest(cleanArgs)
	if isHelp {
		return printUnifiedNodesHelp()
	}

	tok := strings.ToLower(cleanArgs[0])
	isDelegated, err := dispatchRunningOnlyDelegation(tok, cleanArgs[1:])
	if isDelegated {
		return err
	}

	isPrefixedDelegation, prefErr := handleFlagPrefixedDelegation(cleanArgs)
	if isPrefixedDelegation {
		return prefErr
	}

	return runUnifiedNodesList(append([]string{"--open-only"}, cleanArgs...))
}
