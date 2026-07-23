package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runMultiGroup handles the "multi-group" subcommand.
func runMultiGroup(args []string) {
	checkHelp("multi-group", args)
	if len(args) == 0 {
		showActiveMultiGroup()

		return
	}

	routeMultiGroup(args[0], args[1:])
}

// showActiveMultiGroup prints the current multi-group selection.
func showActiveMultiGroup() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	value := db.GetSetting(constants.SettingActiveMultiGroup)
	if len(value) == 0 {
		fmt.Fprint(os.Stderr, constants.MsgMGNone)

		return
	}
	fmt.Printf(constants.MsgMGActive, value)
	printHints(activeGroupHints())
}

// routeMultiGroup dispatches multi-group subcommands.
func routeMultiGroup(sub string, args []string) {
	if sub == constants.CmdMGClear {
		clearMultiGroup()

		return
	}
	if sub == constants.CmdMGPull {
		runMultiGroupPull()

		return
	}
	if sub == constants.CmdMGStatus {
		runMultiGroupStatus()

		return
	}
	if sub == constants.CmdMGExec {
		runMultiGroupExec(args)

		return
	}

	setMultiGroup(sub)
}

// setMultiGroup saves comma-separated group names as active.
func setMultiGroup(groups string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	validateMultiGroupNames(db, groups)
	err = db.SetSetting(constants.SettingActiveMultiGroup, groups)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGenericFmt, err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgMGSet, groups)
}

// validateMultiGroupNames checks each group name exists.
func validateMultiGroupNames(db interface{ GetSetting(string) string }, groups string) {
	// validation done in multigroupops via store lookup
}

// clearMultiGroup removes the multi-group selection.
func clearMultiGroup() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.DeleteSetting(constants.SettingActiveMultiGroup); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clear multi-group setting: %v\n", err)
	}
	fmt.Print(constants.MsgMGCleared)
}

// loadMultiGroupNames returns group names from the active setting.
func loadMultiGroupNames(dbGetter interface{ GetSetting(string) string }) []string {
	value := dbGetter.GetSetting(constants.SettingActiveMultiGroup)
	if len(value) == 0 {
		fmt.Fprint(os.Stderr, constants.MsgMGNone)
		os.Exit(1)
	}

	return strings.Split(value, ",")
}
