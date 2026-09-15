package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runGroup handles the "group" subcommand and routes to sub-handlers.
func runGroup(args []string) error {
	checkHelp("group", args)
	if len(args) == 0 {
		return showActiveGroup()
	}

	return dispatchGroup(args[0], args[1:])
}

func displayActiveGroup(value string) {
	if len(value) == 0 {
		fmt.Fprintln(os.Stderr, constants.MsgGroupNoActive)

		return
	}

	fmt.Printf(constants.MsgGroupActiveShow, value)
	printHints(activeGroupHints())
}

// showActiveGroup prints the currently active group.
func showActiveGroup() *apperror.AppError {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()

	displayActiveGroup(db.GetSetting(constants.SettingActiveGroup))

	return nil
}

func dispatchGroupCRUD(sub string, args []string) result.Result[bool] {
	if sub == constants.CmdGroupCreate {
		return result.RouteMatched(runGroupCreate(args))
	}

	if sub == constants.CmdGroupAdd {
		return result.RouteMatched(runGroupAdd(args))
	}

	if sub == constants.CmdGroupRemove {
		return result.RouteMatched(runGroupRemove(args))
	}

	if sub == constants.CmdGroupList {
		return result.RouteMatched(runGroupList())
	}

	return result.RouteUnmatched()
}

// dispatchGroup routes group subcommands to their handlers.
func dispatchGroup(sub string, args []string) error {
	resCRUD := dispatchGroupCRUD(sub, args)
	if resCRUD.Data {
		return resCRUD.AppError()
	}

	if sub == constants.CmdGroupShow {
		return runGroupShow(args)
	}

	if sub == constants.CmdGroupDelete {
		return runGroupDelete(args)
	}

	resScoped := dispatchGroupScoped(sub, args)
	if resScoped.Data {
		return resScoped.AppError()
	}

	return activateGroup(sub)
}

// dispatchGroupScoped handles pull/status/exec on the active group.
func dispatchGroupScoped(sub string, args []string) result.Result[bool] {
	if sub == constants.CmdMGPull {
		return result.RouteMatched(runActiveGroupPull())
	}

	if sub == constants.CmdMGStatus {
		return result.RouteMatched(runActiveGroupStatus())
	}

	if sub == constants.CmdMGExec {
		return result.RouteMatched(runActiveGroupExec(args))
	}

	if sub == constants.CmdMGClear {
		return result.RouteMatched(clearActiveGroup())
	}

	return result.RouteUnmatched()
}

func persistActiveGroupSetting(db *store.DB, name string) {
	if err := db.SetSetting(constants.SettingActiveGroup, name); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not save active group setting: %v\n", err)
	}

	fmt.Printf(constants.MsgGroupActivated, name)
	printHints(activeGroupHints())
}

// activateGroup sets a group as the active group.
func activateGroup(name string) *apperror.AppError {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()

	if _, gErr := db.ShowGroup(name); gErr != nil {
		return apperror.WrapSimple(gErr, constants.ErrBareFmt)
	}

	persistActiveGroupSetting(db, name)

	return nil
}

// clearActiveGroup removes the active group selection.
func clearActiveGroup() *apperror.AppError {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()

	if err := db.DeleteSetting(constants.SettingActiveGroup); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clear active group setting: %v\n", err)
	}

	fmt.Println(constants.MsgGroupCleared)

	return nil
}
