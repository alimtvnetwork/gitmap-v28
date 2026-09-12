package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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
	if err, isHandled := routeGitProfileSub(subCmd, tailArgs); isHandled {
		return err
	}

	if err, isHandled := routeDBProfileSub(subCmd, tailArgs); isHandled {
		return err
	}

	if err, isHandled := routeChromeProfileSub(subCmd, tailArgs); isHandled {
		return err
	}

	fmt.Fprint(os.Stderr, constants.ErrProfileUsage)

	return apperror.NewSimple("fatal error", "E9000")
}

func routeGitProfileSub(subCmd string, tailArgs []string) (error, bool) {
	if subCmd == "git" || subCmd == "accounts" || subCmd == "set-default" {
		return runProfiles(append([]string{subCmd}, tailArgs...)), true
	}

	return nil, false
}

func routeDBProfileSub(subCmd string, tailArgs []string) (error, bool) {
	switch subCmd {
	case constants.CmdProfileCreate:
		return runProfileCreate(tailArgs), true
	case constants.CmdProfileList, "ls", "status":
		return runProfileList(), true
	case constants.CmdProfileSwitch:
		return runProfileSwitch(tailArgs), true
	default:
		return routeDBProfileExtra(subCmd, tailArgs)
	}
}

func routeDBProfileExtra(subCmd string, tailArgs []string) (error, bool) {
	if subCmd == constants.CmdProfileDelete {
		return runProfileDelete(tailArgs), true
	}

	if subCmd == constants.CmdProfileShow {
		return runProfileShow(), true
	}

	return nil, false
}

func routeChromeProfileSub(subCmd string, tailArgs []string) (error, bool) {
	if err, isHandled := routeChromeImportSub(subCmd, tailArgs); isHandled {
		return err, true
	}

	if err, isHandled := routeChromeExportSub(subCmd, tailArgs); isHandled {
		return err, true
	}

	return nil, false
}

func routeChromeImportSub(subCmd string, tailArgs []string) (error, bool) {
	switch subCmd {
	case constants.CmdProfileImport, "cpi", "profile-import":
		return runChromeProfileImport(tailArgs), true
	case constants.CmdProfileImportAll, "cpi-all", "all-profile-import", "import-all-profiles":
		return runChromeImportAll(tailArgs), true
	case constants.CmdProfileInspect, constants.CmdProfilePreview, constants.CmdProfileCheck,
		constants.CmdProfileImportCheck, "check-import":
		return runChromeProfileImportCheck(tailArgs), true
	}

	return nil, false
}

func routeChromeExportSub(subCmd string, tailArgs []string) (error, bool) {
	switch subCmd {
	case constants.CmdProfileExport, "cpe", "profile-export":
		return runChromeProfileExport(tailArgs), true
	case constants.CmdProfileExportAll, "cpe-all", "all-profile-export", "export-all-profiles":
		return runChromeExportAll(tailArgs), true
	}

	return nil, false
}
