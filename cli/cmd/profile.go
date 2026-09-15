package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchromeprofile"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runProfile handles the "profile" subcommand routing.
func runProfile(args []string) error {
	if len(args) < 1 {
		checkHelp("profile", args)

		return runProfileList()
	}

	if isHelpFlag(args[0]) {
		checkHelp("profile", args)

		return nil
	}

	subCmd := args[0]
	tailArgs := args[1:]

	return routeProfileSub(subCmd, tailArgs)
}

// routeProfileSub routes to the appropriate profile subcommand.
func routeProfileSub(subCmd string, tailArgs []string) error {
	resGit := routeGitProfileSub(subCmd, tailArgs)
	if resGit.IsMatched() {
		return resGit.AppError()
	}

	resDB := routeDBProfileSub(subCmd, tailArgs)
	if resDB.IsMatched() {
		return resDB.AppError()
	}

	resChrome := routeChromeProfileSub(subCmd, tailArgs)
	if resChrome.IsMatched() {
		return resChrome.AppError()
	}

	resInstall := routeInstallProfileSub(subCmd, tailArgs)
	if resInstall.IsMatched() {
		return resInstall.AppError()
	}

	fmt.Fprint(os.Stderr, constants.ErrProfileUsage)

	return apperror.NewSimple("fatal error", "E9000")
}

func routeInstallProfileSub(subCmd string, tailArgs []string) result.ErrorWrapper {
	if isInstallSubCmd(subCmd) {
		return result.MatchWrapper(cmdinstall.RunInstall(append([]string{"profile"}, tailArgs...)))
	}

	if isProfileTreeSubCmd(subCmd) {
		return result.MatchWrapper(cmdinstall.RunInstall(append([]string{"profile", "tree"}, tailArgs...)))
	}

	if cmdinstall.IsInstallProfile(subCmd) {
		return result.MatchWrapper(cmdinstall.RunInstall(append([]string{"profile", subCmd}, tailArgs...)))
	}

	return result.UnmatchedWrapper()
}

func isInstallSubCmd(subCmd string) bool {
	return subCmd == constants.CmdInstall || subCmd == constants.CmdInstallAlias
}

func isProfileTreeSubCmd(subCmd string) bool {
	return subCmd == "tree" || subCmd == "--tree" || subCmd == "-t"
}

func routeGitProfileSub(subCmd string, tailArgs []string) result.ErrorWrapper {
	if subCmd == "git" || subCmd == "accounts" || subCmd == "set-default" {
		return result.MatchWrapper(runProfiles(append([]string{subCmd}, tailArgs...)))
	}

	return result.UnmatchedWrapper()
}

func routeDBProfileSub(subCmd string, tailArgs []string) result.ErrorWrapper {
	switch subCmd {
	case constants.CmdProfileCreate:
		return result.MatchWrapper(runProfileCreate(tailArgs))
	case constants.CmdProfileList, "ls", "status":
		return result.MatchWrapper(runProfileList())
	case constants.CmdProfileSwitch:
		return result.MatchWrapper(runProfileSwitch(tailArgs))
	default:
		return routeDBProfileExtra(subCmd, tailArgs)
	}
}

func routeDBProfileExtra(subCmd string, tailArgs []string) result.ErrorWrapper {
	if subCmd == constants.CmdProfileDelete {
		return result.MatchWrapper(runProfileDelete(tailArgs))
	}

	if subCmd == constants.CmdProfileShow {
		return result.MatchWrapper(runProfileShow())
	}

	return result.UnmatchedWrapper()
}

func routeChromeProfileSub(subCmd string, tailArgs []string) result.ErrorWrapper {
	resImport := routeChromeImportSub(subCmd, tailArgs)
	if resImport.IsMatched() {
		return resImport
	}

	return routeChromeExportSub(subCmd, tailArgs)
}

func routeChromeImportSub(subCmd string, tailArgs []string) result.ErrorWrapper {
	switch subCmd {
	case constants.CmdProfileImport, "cpi", "profile-import":
		return result.MatchWrapper(cmdchromeprofile.RunProfileImport(tailArgs))
	case constants.CmdProfileImportAll, "cpi-all", "all-profile-import", "import-all-profiles":
		return result.MatchWrapper(cmdchromeprofile.RunImportAll(tailArgs))
	case constants.CmdProfileInspect, constants.CmdProfilePreview, constants.CmdProfileCheck,
		constants.CmdProfileImportCheck, "check-import":
		return result.MatchWrapper(cmdchromeprofile.RunProfileImportCheck(tailArgs))
	}

	return result.UnmatchedWrapper()
}

func routeChromeExportSub(subCmd string, tailArgs []string) result.ErrorWrapper {
	switch subCmd {
	case constants.CmdProfileExport, "cpe", "profile-export":
		return result.MatchWrapper(cmdchromeprofile.RunProfileExport(tailArgs))
	case constants.CmdProfileExportAll, "cpe-all", "all-profile-export", "export-all-profiles":
		return result.MatchWrapper(cmdchromeprofile.RunExportAll(tailArgs))
	}

	return result.UnmatchedWrapper()
}
