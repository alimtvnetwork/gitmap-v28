package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runCDSetDefault sets a default path for a repo name.
func runCDSetDefault(args []string) error {
	if len(args) < 2 {
		fmt.Fprint(os.Stderr, constants.ErrCDSetDefaultUsage)
		return apperror.NewSimple("fatal error", "E9000")
	}

	name := args[0]
	path := args[1]

	defaults := store.LoadCDDefaults(constants.DefaultOutputFolder)
	defaults[name] = path

	saveCDDefaultsOrExit(defaults)
	fmt.Printf(constants.MsgCDDefaultSet, name, path)
	return nil
}

// runCDClearDefault removes the default path for a repo name.
func runCDClearDefault(args []string) error {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrCDClearDefaultUsage)
		return apperror.NewSimple("fatal error", "E9000")
	}

	name := args[0]
	defaults := store.LoadCDDefaults(constants.DefaultOutputFolder)

	if _, ok := defaults[name]; !ok {
		return apperror.NewSimple(constants.ErrCDDefaultNotFound, "E9000")
	}

	delete(defaults, name)
	saveCDDefaultsOrExit(defaults)
	fmt.Printf(constants.MsgCDDefaultCleared, name)
	return nil
}

// loadCDDefault returns the default path for a repo, or empty string.
func loadCDDefault(name string) string {
	defaults := store.LoadCDDefaults(constants.DefaultOutputFolder)

	return defaults[name]
}

// saveCDDefaultsOrExit saves the cd-defaults.json, exiting on error.
func saveCDDefaultsOrExit(defaults map[string]string) *apperror.AppError {
	err := store.SaveCDDefaults(constants.DefaultOutputFolder, defaults)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrGenericFmt)
	}
	return nil
}
