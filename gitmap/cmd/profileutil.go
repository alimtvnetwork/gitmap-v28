package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// profileExists checks if a profile name is in the list.
func profileExists(profiles []string, name string) bool {
	for _, p := range profiles {
		if p == name {
			return true
		}
	}

	return false
}

// removeProfile removes a name from the profile list.
func removeProfile(profiles []string, name string) []string {
	var result []string

	for _, p := range profiles {
		if p != name {
			result = append(result, p)
		}
	}

	return result
}

// validateProfileDelete checks delete is allowed.
func validateProfileDelete(name string, cfg model.ProfileConfig) {
	if name == constants.DefaultProfileName {
		fmt.Fprint(os.Stderr, constants.ErrProfileDeleteDefault)
		os.Exit(1)
	}

	if name == cfg.Active {
		fmt.Fprint(os.Stderr, constants.ErrProfileDeleteActive)
		os.Exit(1)
	}

	if !profileExists(cfg.Profiles, name) {
		fmt.Fprintf(os.Stderr, constants.ErrProfileNotFound, name)
		os.Exit(1)
	}
}

// saveProfileOrExit saves the profile config, exiting on error.
func saveProfileOrExit(cfg model.ProfileConfig) {
	err := store.SaveProfileConfig(constants.DefaultOutputFolder, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProfileConfig, err)
		os.Exit(1)
	}
}

// initProfileDB creates and migrates the database for a new profile.
func initProfileDB(name string) {
	db, err := store.OpenDefaultProfile(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not initialize profile database for %s: %v\n", name, err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Profile DB migration failed: %v\n", err)
	}
}

// removeProfileDB deletes the database file for a profile.
func removeProfileDB(name string) {
	dbFile := store.ProfileDBFile(name)
	path := filepath.Join(constants.DefaultOutputFolder, constants.DBDir, dbFile)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not remove profile DB %s: %v\n", path, err)
	}
}
