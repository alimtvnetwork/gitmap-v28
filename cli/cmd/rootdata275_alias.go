package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfleet"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhealth"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhelp"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpe"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtheme"
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

func dispatch275Db(twoWord string, args []string) (bool, error) {
	switch twoWord {
	case "db doctor":
		if appErr := cmddb.RunDbDoctor(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db listing":
		if appErr := cmddb.RunDbListing(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db add":
		if appErr := cmddb.RunDbAdd(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db edit":
		if appErr := cmddb.RunDbEdit(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db find":
		if appErr := cmddb.RunDbFind(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db view":
		if appErr := cmddb.RunDbView(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db rest":
		if appErr := cmddb.RunDbRest(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db record-del", "db record del", "db dbrd":
		if appErr := cmddb.RunDbRecordDel(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	case "db rm", "db delete":
		if appErr := cmddb.RunDbRm(args); appErr != nil {
			return true, appErr
		}
		return true, nil
	}
	return false, nil
}


