package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runGroup handles the "group" subcommand and routes to sub-handlers.
func runGroup(args []string) error {
	checkHelp("group", args)
	if len(args) == 0 {
		showActiveGroup()

		return nil
	}
	dispatchGroup(args[0], args[1:])
	return nil
}

// showActiveGroup prints the currently active group.
func showActiveGroup() {
	db, err := openDB()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrListDBFailed)
		return
	}
	defer db.Close()

	value := db.GetSetting(constants.SettingActiveGroup)
	if len(value) == 0 {
		fmt.Fprintln(os.Stderr, constants.MsgGroupNoActive)

		return
	}
	fmt.Printf(constants.MsgGroupActiveShow, value)
	printHints(activeGroupHints())
}

// dispatchGroup routes group subcommands to their handlers.
func dispatchGroup(sub string, args []string) {
	if sub == constants.CmdGroupCreate {
		runGroupCreate(args)

		return
	}
	if sub == constants.CmdGroupAdd {
		runGroupAdd(args)

		return
	}
	if sub == constants.CmdGroupRemove {
		runGroupRemove(args)

		return
	}
	if sub == constants.CmdGroupList {
		runGroupList()

		return
	}
	if sub == constants.CmdGroupShow {
		runGroupShow(args)

		return
	}
	if sub == constants.CmdGroupDelete {
		runGroupDelete(args)

		return
	}
	if dispatchGroupScoped(sub, args) {
		return
	}

	activateGroup(sub)
}

// dispatchGroupScoped handles pull/status/exec on the active group.
func dispatchGroupScoped(sub string, args []string) bool {
	if sub == constants.CmdMGPull {
		runActiveGroupPull()

		return true
	}
	if sub == constants.CmdMGStatus {
		runActiveGroupStatus()

		return true
	}
	if sub == constants.CmdMGExec {
		runActiveGroupExec(args)

		return true
	}
	if sub == constants.CmdMGClear {
		clearActiveGroup()

		return true
	}

	return false
}

// activateGroup sets a group as the active group.
func activateGroup(name string) {
	db, err := openDB()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrListDBFailed)
		return
	}
	defer db.Close()

	_, gErr := db.ShowGroup(name)
	if gErr != nil {
		apperror.NewSimple(constants.ErrBareFmt, "E9000")
		return
	}

	if err := db.SetSetting(constants.SettingActiveGroup, name); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not save active group setting: %v\n", err)
	}
	fmt.Printf(constants.MsgGroupActivated, name)
	printHints(activeGroupHints())
}

// clearActiveGroup removes the active group selection.
func clearActiveGroup() {
	db, err := openDB()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrListDBFailed)
		return
	}
	defer db.Close()

	if err := db.DeleteSetting(constants.SettingActiveGroup); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clear active group setting: %v\n", err)
	}
	fmt.Println(constants.MsgGroupCleared)
}
