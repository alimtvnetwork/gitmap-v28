package cmd

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagent"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdci"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdflatten"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfleet"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhelp"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlfs"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpe"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompts"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspec"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsplit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtheme"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo"
)

// dispatch275 routes the 30 new commands from program 275 (spec 275).
// It handles both two-word forms (e.g. `gitmap mr list`) and single-word
// aliases (e.g. `gitmap mrl`). Returns (false, nil) if not handled.
func dispatch275(command string) (bool, error) {
	args := os.Args[2:]

	if handled, err := tryDispatchHyphenated(command, args); handled {
		return true, err
	}

	if handled, err := tryDispatchTwoWordArg(command, args); handled {
		return true, err
	}

	if handled, err := dispatch275Alias(command, args); handled {
		return true, err
	}

	return false, nil
}

func tryDispatchHyphenated(command string, args []string) (bool, error) {
	if !strings.Contains(command, "-") {
		return false, nil
	}
	hyphenCmd := strings.ReplaceAll(command, "-", " ")

	return dispatch275TwoWord(hyphenCmd, args)
}

func tryDispatchTwoWordArg(command string, args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	sub := strings.ToLower(args[0])
	twoWord := command + " " + sub

	return dispatch275TwoWord(twoWord, args[1:])
}

// dispatch275TwoWord handles `gitmap <word> <sub>` forms.
func dispatch275TwoWord(twoWord string, args []string) (bool, error) {
	switch twoWord {
	// Group A: mr family (also routed via DispatchMoveMerge; belt-and-braces).
	case "mr list":
		if appErr := cmdmerge.RunMrList(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "mr view":
		if appErr := cmdmerge.RunMrView(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "mr merge":
		if appErr := cmdmerge.RunMrMerge(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "mr checkout":
		if appErr := cmdmerge.RunMrCheckout(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group B: repo + scan.
	case "repo export":
		if appErr := cmdrepo.RunRepoExport(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "repo import":
		if appErr := cmdrepo.RunRepoImport(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "repo flatten":
		if err := cmdflatten.RunFlattenCmd(args); err != nil {
			return true, err
		}
		return true, nil
	case "repo undo-flatten", "repo undo flatten", "undo flatten":
		target := ""
		if len(args) > 0 {
			target = args[0]
		}
		if err := cmdundo.RunUndoFlatten(target); err != nil {
			return true, err
		}
		return true, nil
	case "lfs convert":
		if err := cmdlfs.RunLFSConvertCmd(args); err != nil {
			return true, err
		}
		return true, nil
	case "split large files":
		if err := cmdsplit.RunSplitLargeFilesCmd(args); err != nil {
			return true, err
		}
		return true, nil
	case "scan diff":
		if appErr := cmdscan.RunScanDiff(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "scan verify":
		if appErr := cmdscan.RunScanVerify(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group C: help + theme.
	case "help topic":
		if appErr := cmdhelp.RunHelpTopic(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "theme list":
		if appErr := cmdtheme.RunTheme(append([]string{"list"}, args...)); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "theme preview":
		if appErr := cmdtheme.RunThemePreview(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group D: install + release.
	case "install verify":
		if appErr := cmdinstall.RunInstallVerify(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "install urls":
		if appErr := cmdinstall.RunInstallURLs(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "release verify":
		if appErr := cmdrelease.RunReleaseVerify(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "release winscan":
		if appErr := cmdrelease.RunReleaseWinScan(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "release train":
		if appErr := cmdrelease.RunReleaseTrain(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group E: pipeline + fleet + hygiene.
	case "pe triage":
		if appErr := cmdpe.RunPeTriage(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "pe rerun":
		if appErr := cmdpe.RunPeRerun(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "nodes jump":
		if appErr := cmdfleet.RunNodesJump(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "nodes version":
		if appErr := cmdfleet.RunNodesVersion(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "errors last":
		if appErr := cmderrors.RunErrorsLast(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db doctor", "db listing", "db add", "db edit", "db find", "db view":
		return dispatch275Db(twoWord, args)
	case "hooks install":
		if appErr := cmdhooks.RunHooksInstall(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "prompts drift":
		if appErr := cmdprompts.RunPromptsDrift(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "agent report":
		if appErr := cmdagent.RunAgentReport(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "ci open":
		if appErr := cmdci.RunCiOpen(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "spec progress":
		if appErr := cmdspec.RunSpecProgress(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "agm queue audit":
		if appErr := cmdagm.RunAgmQueueAudit(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	}

	if isAgmQueueAudit(twoWord, args) {
		return dispatchAgmQueueAudit(args[1:])
	}

	if isSplitLargeFiles(twoWord, args) {
		return true, cmdsplit.RunSplitLargeFilesCmd(args[1:])
	}

	return false, nil
}

func isAgmQueueAudit(twoWord string, args []string) bool {
	return twoWord == "agm queue" && len(args) > 0 && strings.EqualFold(args[0], "audit")
}

func dispatchAgmQueueAudit(args []string) (bool, error) {
	if appErr := cmdagm.RunAgmQueueAudit(args); appErr != nil {
		return true, appErr
	}
	return true, nil
}

func isSplitLargeFiles(twoWord string, args []string) bool {
	return twoWord == "split large" && len(args) > 0 && strings.EqualFold(args[0], "files")
}
