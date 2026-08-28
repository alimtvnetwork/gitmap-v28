package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
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

func applyEnvSet(name, value string, f envSetFlags) error {
	if f.dryRun {
		fmt.Printf(constants.MsgEnvDrySet, name, value)
		return nil
	}
	if err := setEnvPersistent(name, value, f.system, f.shell); err != nil {
		return err
	}
	registry := loadEnvRegistry()
	registry = upsertEnvVariable(registry, name, value)
	saveEnvRegistry(registry)
	fmt.Printf(constants.MsgEnvSet, name, value)
	return nil
}

// runEnvSet sets an environment variable persistently.
func runEnvSet(args []string) error {
	name, value, flags := parseEnvSetFlags(args)
	validateEnvName(name)
	validateEnvValue(value)
	return applyEnvSet(name, value, flags)
}

// runEnvGet retrieves a managed environment variable value.
func runEnvGet(args []string) error {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrEnvNameRequired)
		return apperror.NewSimple("fatal error", "E9000")
	}

	name := args[0]
	registry := loadEnvRegistry()
	entry := findEnvVariable(registry, name)

	fmt.Printf(constants.MsgEnvGetFmt, entry.Name, entry.Value)
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

func applyEnvDelete(name string, f envCommonFlags) error {
	if f.dryRun {
		fmt.Printf(constants.MsgEnvDryDelete, name)
		return nil
	}
	if err := deleteEnvPersistent(name, f.system, f.shell); err != nil {
		return err
	}
	registry := loadEnvRegistry()
	registry = removeEnvVariable(registry, name)
	saveEnvRegistry(registry)
	fmt.Printf(constants.MsgEnvDeleted, name)
	return nil
}

// runEnvDelete removes a managed environment variable.
func runEnvDelete(args []string) error {
	name, flags := parseEnvCommonFlags("env-delete", args)
	validateEnvName(name)
	return applyEnvDelete(name, flags)
}

// runEnvList prints all managed environment variables.
func runEnvList() error {
	registry := loadEnvRegistry()
	if len(registry.Variables) == 0 {
		fmt.Print(constants.MsgEnvListEmpty)
		return nil
	}
	fmt.Print(constants.MsgEnvListHeader)
	for _, v := range registry.Variables {
		fmt.Printf(constants.MsgEnvListRow, v.Name, v.Value)
	}
	return nil
}

func applyEnvPathAdd(dir string, f envCommonFlags) error {
	if f.dryRun {
		fmt.Printf(constants.MsgEnvDryPath, dir)
		return nil
	}
	if err := addPathPersistent(dir, f.system, f.shell); err != nil {
		return err
	}
	registry := loadEnvRegistry()
	registry.Paths = append(registry.Paths, model.EnvPathEntry{Path: dir})
	saveEnvRegistry(registry)
	fmt.Printf(constants.MsgEnvPathAdded, dir)
	return nil
}

// runEnvPathAdd adds a directory to the system PATH.
func runEnvPathAdd(args []string) error {
	dir, flags := parseEnvCommonFlags("env-path-add", args)
	validateEnvPathDir(dir)
	registry := loadEnvRegistry()
	checkEnvPathNotDuplicate(registry, dir)
	return applyEnvPathAdd(dir, flags)
}

func validateEnvPathRemove(dir string) {
	if dir == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvPathRequired)
		apperror.NewSimple("fatal error", "E9000")
		return
	}
}

func applyEnvPathRemove(dir string, f envCommonFlags) error {
	if f.dryRun {
		fmt.Printf(constants.MsgEnvDryDelete, dir)
		return nil
	}
	if err := removePathPersistent(dir, f.system, f.shell); err != nil {
		return err
	}
	registry := loadEnvRegistry()
	registry = removeEnvPath(registry, dir)
	saveEnvRegistry(registry)
	fmt.Printf(constants.MsgEnvPathRemoved, dir)
	return nil
}

// runEnvPathRemove removes a directory from the system PATH.
func runEnvPathRemove(args []string) error {
	dir, flags := parseEnvCommonFlags("env-path-remove", args)
	validateEnvPathRemove(dir)
	return applyEnvPathRemove(dir, flags)
}

// runEnvPathList prints all managed PATH entries.
func runEnvPathList() error {
	registry := loadEnvRegistry()
	if len(registry.Paths) == 0 {
		fmt.Print(constants.MsgEnvPathEmpty)
		return nil
	}
	fmt.Print(constants.MsgEnvPathHeader)
	for _, p := range registry.Paths {
		fmt.Printf(constants.MsgEnvPathRow, p.Path)
	}
	return nil
}
