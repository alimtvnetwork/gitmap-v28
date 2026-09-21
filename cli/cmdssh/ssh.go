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
	if checkSSHHelp(args) {
		return nil
	}
	if len(args) == 0 {
		runSSHGenerate(args)
		fmt.Fprint(os.Stdout, constants.MsgSSHAvailableCommands)
		return nil
	}
	return result.AsError(dispatchSSH(context.Background(), args, nil))
}

func dispatchKeyOpsSSH(sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "key", "keys", "ssh-key":
		if len(args) > 0 {
			switch args[0] {
			case "add", "install":
				return result.MatchWrapper(RunSSHAuthKeyAddCLI(args[1:]))
			case "rm", "remove", "delete":
				return result.MatchWrapper(runSSHDelete(args[1:]))
			case "manage", "list", "ls":
				return result.MatchWrapper(runSSHKeysManage(args[1:]))
			}
		}
	}

	return result.UnmatchedWrapper()
}

func dispatchCoreSSH(ctx context.Context, sub string, args []string, parent *cobra.Command) result.ErrorWrapper {
	if resNode := dispatchNodeSSH(ctx, sub, args, parent); resNode.IsMatched() {
		return resNode
	}
	if resKey := dispatchKeyOpsSSH(sub, args); resKey.IsMatched() {
		return resKey
	}

	return dispatchBasicCoreSSH(ctx, sub, args, parent)
}

func dispatchBasicCoreSSH(ctx context.Context, sub string, args []string, parent *cobra.Command) result.ErrorWrapper {
	switch sub {
	case "login", "login-install":
		return result.MatchWrapper(runSSHLogin(parent, args, ctx))
	case "alias":
		return result.MatchWrapper(runSSHAlias(parent, args, ctx))
	case "exec", "se":
		return result.MatchWrapper(runSSHExec(args))
	case "pass", "password":
		return result.MatchWrapper(RunSSHPassCLI(args))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchHistorySSH(sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "reset":
		return result.MatchWrapper(RunSSHResetCLI(args))
	case "undo":
		return result.MatchWrapper(RunSSHUndoCLI(args))
	case "redo":
		return result.MatchWrapper(RunSSHRedoCLI(args))
	case "restore":
		return result.MatchWrapper(RunSSHRestoreCLI(args))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchNodeSSH(ctx context.Context, sub string, args []string, parent *cobra.Command) result.ErrorWrapper {
	if resHist := dispatchHistorySSH(sub, args); resHist.IsMatched() {
		return resHist
	}
	switch sub {
	case "join", "sj", "add":
		return result.MatchWrapper(RunSSHJoinCLI(args))
	case "nodes", "node", "ls":
		return result.MatchWrapper(RunSSHNodesCLI(ctx, args))
	case "rm", "remove":
		return result.MatchWrapper(runSJRm(parent, args, ctx))
	case "ip":
		if len(args) > 0 && (args[0] == "rm" || args[0] == "remove" || args[0] == "delete") {
			return result.MatchWrapper(runSJRm(parent, args[1:], ctx))
		}
		return result.UnmatchedWrapper()
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchAuthOpsSSH(sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "auth-key-add", "ssh-key-add", "key-add", "add-key":
		return result.MatchWrapper(RunSSHAuthKeyAddCLI(args))
	case "auth-key", "copy-id", "fix-auth":
		if len(args) > 0 && (args[0] == "add" || args[0] == "install") {
			return result.MatchWrapper(RunSSHAuthKeyAddCLI(args[1:]))
		}
		return result.MatchWrapper(RunSSHAuthKeyDeployCLI(args))
	case "known-hosts", "knownhosts", "kh", "trust-list":
		return result.MatchWrapper(RunSSHKnownHostsCLI(args))
	case "trust":
		return result.MatchWrapper(RunSSHTrustCLI(args))
	case "untrust":
		return result.MatchWrapper(RunSSHUntrustCLI(args))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchPackageSSH(sub string, args []string) result.ErrorWrapper {
	if resAuth := dispatchAuthOpsSSH(sub, args); resAuth.IsMatched() {
		return resAuth
	}
	switch sub {
	case "install", "i":
		return result.MatchWrapper(runSSHInstallCLI(args))
	case "update", "u":
		return result.MatchWrapper(runSSHUpdateCLI(args))
	case "scan":
		return result.MatchWrapper(runSSHScanCLI(args))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchToolsSSH(ctx context.Context, sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "check", "health", "ping":
		return result.MatchWrapper(RunSJStatus(nil, args, ctx))
	case "agy":
		return result.MatchWrapper(runSSHAgyCLI(args))
	case "code":
		return result.MatchWrapper(runSSHCodeCLI(args))
	case "compare", "matrix":
		return result.MatchWrapper(runSSHCompareCLI(args))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchFilesSSH(sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "profiles", "profile", "p":
		return result.MatchWrapper(runSSHProfile(args))
	case "copy", "cp":
		return result.MatchWrapper(runSSHCopyCLI(args))
	case "mv", "move":
		return result.MatchWrapper(runSSHMvCLI(args))
	case "schedule", "schedules":
		return result.MatchWrapper(runSSHExec(append([]string{"schedule"}, args...)))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchMacroSSH(sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "macro", "m":
		return result.MatchWrapper(runSSHMacroCLI(args))
	case "macro-sync":
		return result.MatchWrapper(runSSHMacroSyncCLI(args))
	case "macro-export":
		return result.MatchWrapper(runSSHMacroExportCLI(args))
	case "macro-import":
		return result.MatchWrapper(runSSHMacroImportCLI(args))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchExportImportSSH(sub string, args []string) result.ErrorWrapper {
	switch sub {
	case "export-all", "exportall":
		return result.MatchWrapper(RunSSHExportAllCLI(args))
	case "import-all", "importall":
		return result.MatchWrapper(RunSSHImportAllCLI(args))
	default:
		return result.UnmatchedWrapper()
	}
}

func dispatchPrimarySSH(ctx context.Context, sub string, args []string, parent *cobra.Command) result.ErrorWrapper {
	if resCore := dispatchCoreSSH(ctx, sub, args, parent); resCore.IsMatched() {
		return resCore
	}
	if resPkg := dispatchPackageSSH(sub, args); resPkg.IsMatched() {
		return resPkg
	}
	if resTools := dispatchToolsSSH(ctx, sub, args); resTools.IsMatched() {
		return resTools
	}
	if resMacro := dispatchMacroSSH(sub, args); resMacro.IsMatched() {
		return resMacro
	}
	if resSync := dispatchExportImportSSH(sub, args); resSync.IsMatched() {
		return resSync
	}
	return dispatchFilesSSH(sub, args)
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

func dispatchFallbackOrLogin(ctx context.Context, sub string, args []string, parent *cobra.Command) result.ErrorWrapper {
	if isFallback := dispatchFallbackSSH(sub, args); isFallback {
		return result.SuccessWrapper()
	}
	return result.MatchWrapper(runSSHLogin(parent, append([]string{sub}, args...), ctx))
}

func dispatchSSH(ctx context.Context, args []string, parent *cobra.Command) result.ErrorWrapper {
	if len(args) == 0 {
		return result.MatchWrapper(handleEmptySSHArgs())
	}
	resPrimary := dispatchPrimarySSH(ctx, args[0], args[1:], parent)
	if resPrimary.IsMatched() {
		return resPrimary
	}
	return dispatchFallbackOrLogin(ctx, args[0], args[1:], parent)
}
