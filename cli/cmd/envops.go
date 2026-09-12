package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

type envSetFlags struct {
	system, verbose, dryRun bool
	shell                   string
}

type envCommonFlags struct {
	system, dryRun bool
	shell          string
}

func parseEnvSetFlags(args []string) (string, string, envSetFlags) {
	fs := flag.NewFlagSet("env-set", flag.ExitOnError)
	var f envSetFlags
	fs.BoolVar(&f.system, constants.FlagEnvSystem, false, constants.FlagDescEnvSystem)
	fs.StringVar(&f.shell, constants.FlagEnvShell, "", constants.FlagDescEnvShell)
	fs.BoolVar(&f.verbose, constants.FlagEnvVerbose, false, constants.FlagDescEnvVerbose)
	fs.BoolVar(&f.dryRun, constants.FlagEnvDryRun, false, constants.FlagDescEnvDryRun)
	fs.Parse(args)

	return fs.Arg(0), fs.Arg(1), f
}

func persistEnvVariable(name, value string) *apperror.AppError {
	registry, appErr := loadEnvRegistry()
	if appErr != nil {
		return appErr
	}

	registry = upsertEnvVariable(registry, name, value)

	return saveEnvRegistry(registry)
}

func applyEnvSet(name, value string, f envSetFlags) error {
	if f.dryRun {
		fmt.Printf(constants.MsgEnvDrySet, name, value)

		return nil
	}

	if err := setEnvPersistent(name, value, f.system, f.shell); err != nil {
		return err
	}

	if appErr := persistEnvVariable(name, value); appErr != nil {
		return appErr
	}

	fmt.Printf(constants.MsgEnvSet, name, value)

	return nil
}

// runEnvSet sets an environment variable persistently.
func runEnvSet(args []string) error {
	name, value, flags := parseEnvSetFlags(args)
	if appErr := validateEnvName(name); appErr != nil {
		return appErr
	}

	if appErr := validateEnvValue(value); appErr != nil {
		return appErr
	}

	return applyEnvSet(name, value, flags)
}

func fetchAndDisplayEnv(name string) *apperror.AppError {
	registry, appErr := loadEnvRegistry()
	if appErr != nil {
		return appErr
	}

	entry, findErr := findEnvVariable(registry, name)
	if findErr != nil {
		return findErr
	}

	fmt.Printf(constants.MsgEnvGetFmt, entry.Name, entry.Value)

	return nil
}

// runEnvGet retrieves a managed environment variable value.
func runEnvGet(args []string) error {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrEnvNameRequired)

		return apperror.NewSimple("fatal error", "E9000")
	}

	if appErr := fetchAndDisplayEnv(args[0]); appErr != nil {
		return appErr
	}

	return nil
}

func parseEnvCommonFlags(cmdName string, args []string) (string, envCommonFlags) {
	fs := flag.NewFlagSet(cmdName, flag.ExitOnError)
	var f envCommonFlags
	fs.BoolVar(&f.system, constants.FlagEnvSystem, false, constants.FlagDescEnvSystem)
	fs.StringVar(&f.shell, constants.FlagEnvShell, "", constants.FlagDescEnvShell)
	fs.BoolVar(&f.dryRun, constants.FlagEnvDryRun, false, constants.FlagDescEnvDryRun)
	fs.Parse(args)

	return fs.Arg(0), f
}

func deleteEnvFromRegistry(name string) *apperror.AppError {
	registry, appErr := loadEnvRegistry()
	if appErr != nil {
		return appErr
	}

	registry = removeEnvVariable(registry, name)

	return saveEnvRegistry(registry)
}

func applyEnvDelete(name string, f envCommonFlags) error {
	if f.dryRun {
		fmt.Printf(constants.MsgEnvDryDelete, name)

		return nil
	}

	if err := deleteEnvPersistent(name, f.system, f.shell); err != nil {
		return err
	}

	if appErr := deleteEnvFromRegistry(name); appErr != nil {
		return appErr
	}

	fmt.Printf(constants.MsgEnvDeleted, name)

	return nil
}

// runEnvDelete removes a managed environment variable.
func runEnvDelete(args []string) error {
	name, flags := parseEnvCommonFlags("env-delete", args)
	if appErr := validateEnvName(name); appErr != nil {
		return appErr
	}

	return applyEnvDelete(name, flags)
}

func printEnvVariables(vars []model.EnvVariable) {
	fmt.Print(constants.MsgEnvListHeader)
	for _, v := range vars {
		fmt.Printf(constants.MsgEnvListRow, v.Name, v.Value)
	}
}

// runEnvList prints all managed environment variables.
func runEnvList() error {
	registry, appErr := loadEnvRegistry()
	if appErr != nil {
		return appErr
	}

	if len(registry.Variables) == 0 {
		fmt.Print(constants.MsgEnvListEmpty)

		return nil
	}

	printEnvVariables(registry.Variables)

	return nil
}

func persistEnvPath(dir string) *apperror.AppError {
	registry, appErr := loadEnvRegistry()
	if appErr != nil {
		return appErr
	}

	registry.Paths = append(registry.Paths, model.EnvPathEntry{Path: dir})

	return saveEnvRegistry(registry)
}

func applyEnvPathAdd(dir string, f envCommonFlags) error {
	if f.dryRun {
		fmt.Printf(constants.MsgEnvDryPath, dir)

		return nil
	}

	if err := addPathPersistent(dir, f.system, f.shell); err != nil {
		return err
	}

	if appErr := persistEnvPath(dir); appErr != nil {
		return appErr
	}

	fmt.Printf(constants.MsgEnvPathAdded, dir)

	return nil
}

// runEnvPathAdd adds a directory to the system PATH.
func runEnvPathAdd(args []string) error {
	dir, flags := parseEnvCommonFlags("env-path-add", args)
	if appErr := validateEnvPathDir(dir); appErr != nil {
		return appErr
	}

	registry, appErr := loadEnvRegistry()
	if appErr != nil {
		return appErr
	}

	if appErr := checkEnvPathNotDuplicate(registry, dir); appErr != nil {
		return appErr
	}

	return applyEnvPathAdd(dir, flags)
}

func validateEnvPathRemove(dir string) *apperror.AppError {
	if dir == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvPathRequired)

		return apperror.NewSimple("fatal error", "E9000")
	}

	return nil
}

func deleteEnvPathFromRegistry(dir string) *apperror.AppError {
	registry, appErr := loadEnvRegistry()
	if appErr != nil {
		return appErr
	}

	registry = removeEnvPath(registry, dir)

	return saveEnvRegistry(registry)
}

func applyEnvPathRemove(dir string, f envCommonFlags) error {
	if f.dryRun {
		fmt.Printf(constants.MsgEnvDryDelete, dir)

		return nil
	}

	if err := removePathPersistent(dir, f.system, f.shell); err != nil {
		return err
	}

	if appErr := deleteEnvPathFromRegistry(dir); appErr != nil {
		return appErr
	}

	fmt.Printf(constants.MsgEnvPathRemoved, dir)

	return nil
}

// runEnvPathRemove removes a directory from the system PATH.
func runEnvPathRemove(args []string) error {
	dir, flags := parseEnvCommonFlags("env-path-remove", args)
	if appErr := validateEnvPathRemove(dir); appErr != nil {
		return appErr
	}

	return applyEnvPathRemove(dir, flags)
}

func printEnvPaths(paths []model.EnvPathEntry) {
	fmt.Print(constants.MsgEnvPathHeader)
	for _, p := range paths {
		fmt.Printf(constants.MsgEnvPathRow, p.Path)
	}
}

// runEnvPathList prints all managed PATH entries.
func runEnvPathList() error {
	registry, appErr := loadEnvRegistry()
	if appErr != nil {
		return appErr
	}

	if len(registry.Paths) == 0 {
		fmt.Print(constants.MsgEnvPathEmpty)

		return nil
	}

	printEnvPaths(registry.Paths)

	return nil
}
