package cmd

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagent"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdci"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfleet"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdflatten"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhealth"
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

	handled, err := tryDispatchHyphen275(command, args)
	if handled {
		return true, err
	}

	handled, err = tryDispatchTwoWord275(command, args)
	if handled {
		return true, err
	}

	// Single-word aliases.
	return dispatch275Alias(command, args)
}

func tryDispatchHyphen275(command string, args []string) (bool, error) {
	if !strings.Contains(command, "-") {
		return false, nil
	}
	hyphenCmd := strings.ReplaceAll(command, "-", " ")
	return dispatch275TwoWord(hyphenCmd, args)
}

func tryDispatchTwoWord275(command string, args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	sub := strings.ToLower(args[0])
	twoWord := command + " " + sub
	return dispatch275TwoWord(twoWord, args[1:])
}

// dispatch275TwoWord routes the canonical two-word space-separated commands.
func dispatch275TwoWord(twoWord string, args []string) (bool, error) {
	switch twoWord {
	// Group A: GitHub / GitLab Operations.
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
	case "ci watch":
		if appErr := cmdmerge.RunCiWatch(args); appErr != nil {
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
	case "spec progress":
		if appErr := cmdspec.RunSpecProgress(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group B: CI Diagnostics & Fix Automation.
	case "ci quick":
		if appErr := cmdci.RunCiQuick(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "ci check":
		if appErr := cmdci.RunCiCheck(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "ci explain":
		if appErr := cmdci.RunCiExplain(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "agent task":
		if appErr := cmdagent.RunAgentTask(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "agent report":
		if appErr := cmdagent.RunAgentReport(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "agent logs":
		if appErr := cmdagent.RunAgentLogs(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "agent clear":
		if appErr := cmdagent.RunAgentClear(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "agent stats":
		if appErr := cmdagent.RunAgentStats(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group C: Prompts & Deployment Health.
	case "prompts check":
		if appErr := cmdprompts.RunPromptsCheck(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "prompts drift":
		if appErr := cmdprompts.RunPromptsDrift(args); appErr != nil {
			return true, appErr
		}
		return true, nil
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
	case "release win-scan":
		if appErr := cmdrelease.RunReleaseWinScan(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "release train":
		if appErr := cmdrelease.RunReleaseTrain(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group D: Themes & Navigation.
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
	// Group E: Engine Internals & Diagnostics.
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
	case "db doctor":
		if appErr := cmddb.RunDbDoctor(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "hooks install":
		if appErr := cmdhooks.RunHooksInstall(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "health check":
		if appErr := cmdhealth.RunHealth(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	}

	// Program 275 commands with flags.
	isAgmQueueAudit := strings.HasPrefix(twoWord, "agm queue-audit")
	if isAgmQueueAudit {
		return true, cmdagm.RunAgmQueueAudit(args[1:])
	}

	// Three-word form: `gitmap split large files`.
	isSplitLargeFiles := twoWord == "split large" && len(args) > 0 && strings.EqualFold(args[0], "files")
	if isSplitLargeFiles {
		return true, cmdsplit.RunSplitLargeFilesCmd(args[1:])
	}

	return false, nil
}
