package cmdssh

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runSSH handles the "ssh" subcommand and routes to sub-handlers.
func runSSH(args []string) error {
	checkSSHHelp(args)
	if len(args) == 0 {
		runSSHGenerate(args)
		fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
		return nil
	}
	return dispatchSSH(context.Background(), args, nil)
}

func dispatchPrimarySSH(ctx context.Context, sub string, args []string, parent *cobra.Command) (error, bool) {
	switch sub {
	case "login", "login-install":
		return runSSHLogin(parent, args, ctx), true
	case "join", "sj":
		return RunSSHJoinCLI(args), true
	case "alias":
		return runSSHAlias(parent, args, ctx), true
	case "exec", "se":
		return runSSHExec(args), true
	case "profiles", "profile", "p":
		return runSSHProfile(args), true
	}
	return nil, false
}

func runSSHProfile(args []string) error {
	if ProfileRunner != nil {
		return ProfileRunner(args)
	}
	return nil
}

func isSSHViewSub(sub string) bool {
	return sub == constants.SubCmdSSHCat || sub == constants.SubCmdSSHView || sub == constants.SubCmdSSHViewS
}

func isSSHCopySub(sub string) bool {
	return sub == constants.SubCmdSSHCopy || sub == constants.SubCmdSSHCopyS
}

func isSSHListSub(sub string) bool {
	return sub == constants.SubCmdSSHList || sub == constants.SubCmdSSHListS
}

func isSSHDeleteSub(sub string) bool {
	return sub == constants.SubCmdSSHDelete || sub == constants.SubCmdSSHDeleteS
}

func isSSHStatusSub(sub string) bool {
	return sub == constants.SubCmdSSHStatus || sub == constants.SubCmdSSHStatusS
}

func dispatchConfigOrStatus(sub string, args []string) bool {
	if sub == constants.SubCmdSSHConfig {
		runSSHConfig()
		return true
	}
	if isSSHStatusSub(sub) {
		runSSHStatus(args)
		return true
	}
	return false
}

func dispatchOtherFallbackSSH(sub string, args []string) bool {
	if isSSHListSub(sub) {
		runSSHList(args...)
		return true
	}
	if isSSHDeleteSub(sub) {
		runSSHDelete(args)
		return true
	}
	return dispatchConfigOrStatus(sub, args)
}

func dispatchCreateSSH(sub string, args []string) bool {
	if sub == constants.SubCmdSSHCreate {
		runSSHGenerate(args)
		fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
		return true
	}
	return false
}

func dispatchFallbackSSH(sub string, args []string) bool {
	if isSSHViewSub(sub) {
		runSSHCat(args)
		return true
	}
	if isSSHCopySub(sub) {
		runSSHCopy(args)
		return true
	}
	if isCreated := dispatchCreateSSH(sub, args); isCreated {
		return true
	}
	return dispatchOtherFallbackSSH(sub, args)
}

func handleEmptySSHArgs() error {
	runSSHGenerate(nil)
	fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
	return nil
}

func dispatchSSH(ctx context.Context, args []string, parent *cobra.Command) error {
	if len(args) == 0 {
		return handleEmptySSHArgs()
	}
	sub := args[0]
	if err, isPrimary := dispatchPrimarySSH(ctx, sub, args[1:], parent); isPrimary {
		return err
	}
	if isFallback := dispatchFallbackSSH(sub, args[1:]); isFallback {
		return nil
	}
	return runSSHLogin(parent, args, ctx)
}
