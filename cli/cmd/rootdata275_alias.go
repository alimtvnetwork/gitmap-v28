package cmd

import (
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
	case "ciwa":
		if appErr := cmdmerge.RunCiWatch(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "rscan":
		if appErr := cmdscan.RunScanDiff(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "sdiff":
		if appErr := cmdscan.RunScanDiff(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "rexp":
		if appErr := cmdrepo.RunRepoExport(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "rimp":
		if appErr := cmdrepo.RunRepoImport(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "sprog":
		if appErr := cmdspec.RunSpecProgress(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group C aliases.
	case "pdrift":
		if appErr := cmdprompts.RunPromptsDrift(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "insv":
		if appErr := cmdinstall.RunInstallVerify(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "relv":
		if appErr := cmdrelease.RunReleaseVerify(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "relw":
		if appErr := cmdrelease.RunReleaseWinScan(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "relt":
		if appErr := cmdrelease.RunReleaseTrain(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	// Group D aliases.
	case "themel":
		if appErr := cmdtheme.RunTheme(append([]string{"list"}, args...)); appErr != nil {
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
	case "repo-flatten", "flatten":
		if err := cmdflatten.RunFlattenCmd(args); err != nil {
			return true, err
		}
		return true, nil
	case "split-large-files":
		if err := cmdsplit.RunSplitLargeFilesCmd(args); err != nil {
			return true, err
		}
		return true, nil
	case "lfs-convert":
		if err := cmdlfs.RunLFSConvertCmd(args); err != nil {
			return true, err
		}
		return true, nil
	case "undo-flatten", "repo-undo-flatten":
		target := ""
		if len(args) > 0 {
			target = args[0]
		}
		if err := cmdundo.RunUndoFlatten(target); err != nil {
			return true, err
		}
		return true, nil
	}

	return false, nil
}
