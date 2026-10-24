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
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

func runPull(args []string) error {
	if len(args) > 0 && strings.EqualFold(args[0], constants.CmdAll) {
		return dispatchPullAllWithArgs(args[1:])
	}

	return cmdpull.RunPull(args)
}

func dispatchPullAllWithArgs(subArgs []string) error {
	if len(subArgs) > 0 && strutil.EqualFoldAny(subArgs[0], "ssh", constants.FlagSSH) {
		return cmdpull.RunPullAll(append([]string{constants.FlagSSH}, subArgs[1:]...))
	}

	return cmdpull.RunPullAll(subArgs)
}

func runPush(args []string) error {
	if len(args) > 0 && (args[0] == constants.CmdFix || args[0] == constants.CmdPushFix) {
		return cmdpushfix.RunPushFix(args[1:])
	}

	return cmdpull.RunPush(args)
}
