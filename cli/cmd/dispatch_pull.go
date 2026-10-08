// Package cmd — dispatch_pull.go: thin pull/push subcommand routing.
//
// Kept from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D).
// The "all" and "fix" sub-branch decisions are thin dispatch and stay in
// package cmd; implementations live in cmdpull / cmdpushfix.
package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpushfix"
)

func runPull(args []string) error {
	if len(args) > 0 && strings.EqualFold(args[0], "all") {
		return dispatchPullAllWithArgs(args[1:])
	}
	return cmdpull.RunPull(args)
}

func dispatchPullAllWithArgs(subArgs []string) error {
	if len(subArgs) > 0 && (strings.EqualFold(subArgs[0], "ssh") || strings.EqualFold(subArgs[0], "--ssh")) {
		return cmdpull.RunPullAll(append([]string{"--ssh"}, subArgs[1:]...))
	}
	return cmdpull.RunPullAll(subArgs)
}

func runPush(args []string) error {
	if len(args) > 0 && (args[0] == "fix" || args[0] == "push-fix") {
		return cmdpushfix.RunPushFix(args[1:])
	}
	return cmdpull.RunPush(args)
}
