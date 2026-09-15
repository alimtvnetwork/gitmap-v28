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
	if resGit.Data {
		return resGit.AppError()
	}

	resDB := routeDBProfileSub(subCmd, tailArgs)
	if resDB.Data {
		return resDB.AppError()
	}

	resChrome := routeChromeProfileSub(subCmd, tailArgs)
	if resChrome.Data {
		return resChrome.AppError()
	}

	resInstall := routeInstallProfileSub(subCmd, tailArgs)
	if resInstall.Data {
		return resInstall.AppError()
	}

	fmt.Fprint(os.Stderr, constants.ErrProfileUsage)

	return apperror.NewSimple("fatal error", "E9000")
}

func routeInstallProfileSub(subCmd string, tailArgs []string) result.Result[bool] {
	if isInstallSubCmd(subCmd) {
		return result.RouteMatched(cmdinstall.RunInstall(append([]string{"profile"}, tailArgs...)))
	}

	if isProfileTreeSubCmd(subCmd) {
		return result.RouteMatched(cmdinstall.RunInstall(append([]string{"profile", "tree"}, tailArgs...)))
	}

	if cmdinstall.IsInstallProfile(subCmd) {
		return result.RouteMatched(cmdinstall.RunInstall(append([]string{"profile", subCmd}, tailArgs...)))
	}

	return result.RouteUnmatched()
}

func isInstallSubCmd(subCmd string) bool {
	return subCmd == constants.CmdInstall || subCmd == constants.CmdInstallAlias
}

func isProfileTreeSubCmd(subCmd string) bool {
	return subCmd == "tree" || subCmd == "--tree" || subCmd == "-t"
}

func routeGitProfileSub(subCmd string, tailArgs []string) result.Result[bool] {
	if subCmd == "git" || subCmd == "accounts" || subCmd == "set-default" {
		return result.RouteMatched(runProfiles(append([]string{subCmd}, tailArgs...)))
	}

	return result.RouteUnmatched()
}

func routeDBProfileSub(subCmd string, tailArgs []string) result.Result[bool] {
	switch subCmd {
	case constants.CmdProfileCreate:
		return result.RouteMatched(runProfileCreate(tailArgs))
	case constants.CmdProfileList, "ls", "status":
		return result.RouteMatched(runProfileList())
	case constants.CmdProfileSwitch:
		return result.RouteMatched(runProfileSwitch(tailArgs))
	default:
		return routeDBProfileExtra(subCmd, tailArgs)
	}
}

func routeDBProfileExtra(subCmd string, tailArgs []string) result.Result[bool] {
	if subCmd == constants.CmdProfileDelete {
		return result.RouteMatched(runProfileDelete(tailArgs))
	}

	if subCmd == constants.CmdProfileShow {
		return result.RouteMatched(runProfileShow())
	}

	return result.RouteUnmatched()
}

func routeChromeProfileSub(subCmd string, tailArgs []string) result.Result[bool] {
	resImport := routeChromeImportSub(subCmd, tailArgs)
	if resImport.Data {
		return resImport
	}

	return routeChromeExportSub(subCmd, tailArgs)
}

func routeChromeImportSub(subCmd string, tailArgs []string) result.Result[bool] {
	switch subCmd {
	case constants.CmdProfileImport, "cpi", "profile-import":
		return result.RouteMatched(cmdchromeprofile.RunProfileImport(tailArgs))
	case constants.CmdProfileImportAll, "cpi-all", "all-profile-import", "import-all-profiles":
		return result.RouteMatched(cmdchromeprofile.RunImportAll(tailArgs))
	case constants.CmdProfileInspect, constants.CmdProfilePreview, constants.CmdProfileCheck,
		constants.CmdProfileImportCheck, "check-import":
		return result.RouteMatched(cmdchromeprofile.RunProfileImportCheck(tailArgs))
	}

	return result.RouteUnmatched()
}

func routeChromeExportSub(subCmd string, tailArgs []string) result.Result[bool] {
	switch subCmd {
	case constants.CmdProfileExport, "cpe", "profile-export":
		return result.RouteMatched(cmdchromeprofile.RunProfileExport(tailArgs))
	case constants.CmdProfileExportAll, "cpe-all", "all-profile-export", "export-all-profiles":
		return result.RouteMatched(cmdchromeprofile.RunExportAll(tailArgs))
	}

	return result.RouteUnmatched()
}
