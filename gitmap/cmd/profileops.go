package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runProfileCreate creates a new named profile.
func runProfileCreate(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrProfileCreateUsage)
		os.Exit(1)
	}

	name := args[0]
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	if profileExists(cfg.Profiles, name) {
		fmt.Fprintf(os.Stderr, constants.ErrProfileExists, name)
		os.Exit(1)
	}

	cfg.Profiles = append(cfg.Profiles, name)
	saveProfileOrExit(cfg)
	initProfileDB(name)

	fmt.Printf(constants.MsgProfileCreated, name)
}

// runProfileList displays all profiles with active marker.
func runProfileList() {
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	if len(cfg.Profiles) == 0 {
		fmt.Print(constants.MsgProfileEmpty)

		return
	}

	fmt.Println(constants.MsgProfileColumns)
	for _, p := range cfg.Profiles {
		tag := ""
		if p == cfg.Active {
			tag = constants.MsgProfileActiveTag
		}
		fmt.Printf(constants.MsgProfileRowFmt, p, tag)
	}
}

// runProfileSwitch changes the active profile.
func runProfileSwitch(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrProfileSwitchUsage)
		os.Exit(1)
	}

	name := args[0]
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	if !profileExists(cfg.Profiles, name) {
		fmt.Fprintf(os.Stderr, constants.ErrProfileNotFound, name)
		os.Exit(1)
	}

	cfg.Active = name
	saveProfileOrExit(cfg)

	fmt.Printf(constants.MsgProfileSwitched, name)
}

// runProfileDelete removes a profile (not the active or default).
func runProfileDelete(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrProfileDeleteUsage)
		os.Exit(1)
	}

	name := args[0]
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	validateProfileDelete(name, cfg)
	cfg.Profiles = removeProfile(cfg.Profiles, name)
	saveProfileOrExit(cfg)
	removeProfileDB(name)

	fmt.Printf(constants.MsgProfileDeleted, name)
}

// runProfileShow displays the currently active profile.
func runProfileShow() {
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)
	fmt.Printf(constants.MsgProfileActive, cfg.Active)
}
