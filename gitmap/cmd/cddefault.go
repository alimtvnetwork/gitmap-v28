package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runCDSetDefault sets a default path for a repo name.
func runCDSetDefault(args []string) {
	if len(args) < 2 {
		fmt.Fprint(os.Stderr, constants.ErrCDSetDefaultUsage)
		os.Exit(1)
	}

	name := args[0]
	path := args[1]

	defaults := store.LoadCDDefaults(constants.DefaultOutputFolder)
	defaults[name] = path

	saveCDDefaultsOrExit(defaults)
	fmt.Printf(constants.MsgCDDefaultSet, name, path)
}

// runCDClearDefault removes the default path for a repo name.
func runCDClearDefault(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrCDClearDefaultUsage)
		os.Exit(1)
	}

	name := args[0]
	defaults := store.LoadCDDefaults(constants.DefaultOutputFolder)

	if _, ok := defaults[name]; !ok {
		fmt.Fprintf(os.Stderr, constants.ErrCDDefaultNotFound, name)
		os.Exit(1)
	}

	delete(defaults, name)
	saveCDDefaultsOrExit(defaults)
	fmt.Printf(constants.MsgCDDefaultCleared, name)
}

// loadCDDefault returns the default path for a repo, or empty string.
func loadCDDefault(name string) string {
	defaults := store.LoadCDDefaults(constants.DefaultOutputFolder)

	return defaults[name]
}

// saveCDDefaultsOrExit saves the cd-defaults.json, exiting on error.
func saveCDDefaultsOrExit(defaults map[string]string) {
	err := store.SaveCDDefaults(constants.DefaultOutputFolder, defaults)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGenericFmt, err)
		os.Exit(1)
	}
}
