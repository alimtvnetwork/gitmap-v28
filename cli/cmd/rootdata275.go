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
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhealth"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhelp"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpe"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompts"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspec"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtheme"
)

// dispatch275 routes the 30 new commands from program 275 (spec 275).
// It handles both two-word forms (e.g. `gitmap mr list`) and single-word
// aliases (e.g. `gitmap mrl`). Returns (false, nil) if not handled.
//
// NOTE: Every call site uses the explicit nil-check pattern because
// *appfault.AppError is a typed nil — a bare `return RunXxx()` would
// convert nil to a non-nil error interface.
func dispatch275(command string) (bool, error) {
	args := os.Args[2:]

	// Hyphenated forms: `gitmap repo-export` -> `dispatch275TwoWord("repo export", ...)`
	if strings.Contains(command, "-") {
		hyphenCmd := strings.ReplaceAll(command, "-", " ")
		if handled, err := dispatch275TwoWord(hyphenCmd, args); handled {
			return true, err
		}
	}

	// Two-word forms: command + first arg.
	if len(args) > 0 {
		sub := strings.ToLower(args[0])
		twoWord := command + " " + sub

		if handled, err := dispatch275TwoWord(twoWord, args[1:]); handled {
			return true, err
		}
	}

	// Single-word aliases.
	if handled, err := dispatch275Alias(command, args); handled {
		return true, err
	}

	return false, nil
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
	// Group D: install + release (also routed in-package; belt-and-braces).
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
		// "agm queue audit" is three words; handled via "agm queue" prefix.
		if appErr := cmdagm.RunAgmQueueAudit(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	}

	// Three-word form: `gitmap agm queue audit`.
	if twoWord == "agm queue" && len(args) > 0 && strings.ToLower(args[0]) == "audit" {
		if appErr := cmdagm.RunAgmQueueAudit(args[1:]); appErr != nil {
			return true, appErr
		}

		return true, nil
	}

	return false, nil
}

// dispatch275Alias handles single-word aliases for the 275 commands.
func dispatch275Alias(command string, args []string) (bool, error) {
	switch command {
	// Group A aliases.
	case "mrl":
		if appErr := cmdmerge.RunMrList(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "mrv":
		if appErr := cmdmerge.RunMrView(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "mrm":
		if appErr := cmdmerge.RunMrMerge(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "mrco":
		if appErr := cmdmerge.RunMrCheckout(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	// Group B aliases.
	case "repox":
		if appErr := cmdrepo.RunRepoExport(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "repoi":
		if appErr := cmdrepo.RunRepoImport(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "scand":
		if appErr := cmdscan.RunScanDiff(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "scanverify":
		if appErr := cmdscan.RunScanVerify(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	// Group C aliases.
	case "ht":
		if appErr := cmdhelp.RunHelpTopic(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "howto":
		if appErr := cmdhelp.RunHowto(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "themes":
		if appErr := cmdtheme.RunTheme(append([]string{"list"}, args...)); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "themep":
		if appErr := cmdtheme.RunThemePreview(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	// Group E aliases.
	case "pet":
		if appErr := cmdpe.RunPeTriage(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "perr":
		if appErr := cmdpe.RunPeRerun(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "nj":
		if appErr := cmdfleet.RunNodesJump(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "nv":
		if appErr := cmdfleet.RunNodesVersion(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "errl":
		if appErr := cmderrors.RunErrorsLast(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "dbdoc":
		if appErr := cmddb.RunDbDoctor(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "hkin":
		if appErr := cmdhooks.RunHooksInstall(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	case "hlth", "health":
		if appErr := cmdhealth.RunHealth(args); appErr != nil {
			return true, appErr
		}

		return true, nil
	}

	return false, nil
}
