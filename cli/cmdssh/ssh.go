package cmdssh

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runSSH handles the "ssh" subcommand and routes to sub-handlers.
func runSSH(args []string) error {
	checkSSHHelp(args)
	if len(args) == 0 {
		runSSHGenerate(args)
		fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
		return nil
	}
	return result.AsError(dispatchSSH(context.Background(), args, nil))
}

func dispatchPrimarySSH(ctx context.Context, sub string, args []string, parent *cobra.Command) result.ErrorWrapper {
	switch sub {
	case "login", "login-install":
		return result.MatchWrapper(runSSHLogin(parent, args, ctx))
	case "join", "sj":
		return result.MatchWrapper(RunSSHJoinCLI(args))
	case "alias":
		return result.MatchWrapper(runSSHAlias(parent, args, ctx))
	case "exec", "se":
		return result.MatchWrapper(runSSHExec(args))
	case "install", "i":
		return result.MatchWrapper(runSSHInstallCLI(args))
	case "update", "u":
		return result.MatchWrapper(runSSHUpdateCLI(args))
	case "scan":
		return result.MatchWrapper(runSSHScanCLI(args))
	case "check", "health", "ping":
		return result.MatchWrapper(RunSJStatus(nil, args, ctx))
	case "agy":
		return result.MatchWrapper(runSSHAgyCLI(args))
	case "code":
		return result.MatchWrapper(runSSHCodeCLI(args))
	case "compare", "matrix":
		return result.MatchWrapper(runSSHCompareCLI(args))
	case "profiles", "profile", "p":
		return result.MatchWrapper(runSSHProfile(args))
	default:
		return result.UnmatchedWrapper()
	}
}

func runSSHProfile(args []string) error {
	if ProfileRunner != nil {
		return ProfileRunner(args)
	}
	return nil
}

func handleEmptySSHArgs() error {
	runSSHGenerate(nil)
	fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
	return nil
}

func dispatchSSH(ctx context.Context, args []string, parent *cobra.Command) result.ErrorWrapper {
	if len(args) == 0 {
		return result.MatchWrapper(handleEmptySSHArgs())
	}
	sub := args[0]
	resPrimary := dispatchPrimarySSH(ctx, sub, args[1:], parent)
	if resPrimary.IsMatched() {
		return resPrimary
	}

	if isFallback := dispatchFallbackSSH(sub, args[1:]); isFallback {
		return result.SuccessWrapper()
	}

	return result.MatchWrapper(runSSHLogin(parent, args, ctx))
}
