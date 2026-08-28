package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runProfileCreate creates a new named profile.
func runProfileCreate(args []string) error {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrProfileCreateUsage)
		return apperror.NewSimple("fatal error", "E9000")
	}

	name := args[0]
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	if profileExists(cfg.Profiles, name) {
		return apperror.NewSimple(constants.ErrProfileExists, "E9000")
	}

	cfg.Profiles = append(cfg.Profiles, name)
	saveProfileOrExit(cfg)
	initProfileDB(name)

	fmt.Printf(constants.MsgProfileCreated, name)
	return nil
}

// runProfileList displays all profiles with active marker.
func runProfileList() error {
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	if len(cfg.Profiles) == 0 {
		fmt.Print(constants.MsgProfileEmpty)

		return nil
	}

	fmt.Println(constants.MsgProfileColumns)
	for _, p := range cfg.Profiles {
		tag := ""
		if p == cfg.Active {
			tag = constants.MsgProfileActiveTag
		}
		fmt.Printf(constants.MsgProfileRowFmt, p, tag)
	}
	return nil
}

// runProfileSwitch changes the active profile.
func runProfileSwitch(args []string) error {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrProfileSwitchUsage)
		return apperror.NewSimple("fatal error", "E9000")
	}

	name := args[0]
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	if !profileExists(cfg.Profiles, name) {
		return apperror.NewSimple(constants.ErrProfileNotFound, "E9000")
	}

	cfg.Active = name
	saveProfileOrExit(cfg)

	fmt.Printf(constants.MsgProfileSwitched, name)
	return nil
}

// runProfileDelete removes a profile (not the active or default).
func runProfileDelete(args []string) error {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrProfileDeleteUsage)
		return apperror.NewSimple("fatal error", "E9000")
	}

	name := args[0]
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)

	validateProfileDelete(name, cfg)
	cfg.Profiles = removeProfile(cfg.Profiles, name)
	saveProfileOrExit(cfg)
	removeProfileDB(name)

	fmt.Printf(constants.MsgProfileDeleted, name)
	return nil
}

// runProfileShow displays the currently active profile.
func runProfileShow() error {
	cfg := store.LoadProfileConfig(constants.DefaultOutputFolder)
	fmt.Printf(constants.MsgProfileActive, cfg.Active)
	return nil
}
