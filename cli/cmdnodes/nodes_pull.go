package cmdnodes

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

// RunNodesPullAllFn is wired by cmd/di_hooks_wave2.go to delegate to cmdpull.RunPullAll with --ssh.
var RunNodesPullAllFn func(args []string) error

// IsNodesPullAllRequest checks if the given arguments represent a nodes pull-all invocation.
func IsNodesPullAllRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	sub := strings.ToLower(args[0])
	if sub == "pull" && len(args) > 1 && strings.EqualFold(args[1], "errors") {
		return false
	}
	return sub == "pa" || sub == "pull-all" || sub == "pullall" || sub == "pull" || sub == "pas" || sub == "pull-all-ssh"
}

// RunNodesPullAll dispatches fleet-wide batch pull across all nodes and local machine (alias for gitmap pa --ssh).
func RunNodesPullAll(args []string) error {
	if RunNodesPullAllFn != nil {
		return RunNodesPullAllFn(args)
	}
	return cmdssh.RunSSHPullAllFleet(args)
}
