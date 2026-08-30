package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runSSH handles the "ssh" subcommand and routes to sub-handlers.
func runSSH(args []string) error {
	checkHelp("ssh", args)
	if len(args) == 0 {
		runSSHGenerate(args)
		fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
		return nil
	}
	_ = dispatchSSH(context.Background(), args, nil)
	return nil
}

func dispatchSSH(ctx context.Context, args []string, parent *cobra.Command) error {
	if len(args) == 0 {
		runSSHGenerate(args)
		fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
		return nil
	}

	sub := args[0]
	switch sub {
	case "login", "login-install":
		return runSSHLogin(parent, args[1:], ctx)
	case "join", "sj":
		return runJoin(args[1:])
	case "alias":
		return runSSHAlias(parent, args[1:], ctx)
	case "exec", "se":
		return runSSHExec(args[1:])
	case "profiles", "profile", "p":
		return runProfile(args[1:])
	default:
		// Fallback for $username@ip and implicit aliases
		if sub == constants.SubCmdSSHCat || sub == constants.SubCmdSSHView || sub == constants.SubCmdSSHViewS {
			runSSHCat(args[1:])
			return nil
		}
		if sub == constants.SubCmdSSHCopy || sub == constants.SubCmdSSHCopyS {
			runSSHCopy(args[1:])
			return nil
		}
		if sub == constants.SubCmdSSHCreate {
			runSSHGenerate(args[1:])
			fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
			return nil
		}
		if sub == constants.SubCmdSSHList || sub == constants.SubCmdSSHListS {
			runSSHList(args[1:]...)
			return nil
		}
		if sub == constants.SubCmdSSHDelete || sub == constants.SubCmdSSHDeleteS {
			runSSHDelete(args[1:])
			return nil
		}
		if sub == constants.SubCmdSSHConfig {
			runSSHConfig()
			return nil
		}
		if sub == constants.SubCmdSSHStatus || sub == constants.SubCmdSSHStatusS {
			runSSHStatus(args[1:])
			return nil
		}

		// If it reaches here, treat as raw arguments for implicit alias or username@ip
		return runSSHLogin(parent, args, ctx)
	}
}
