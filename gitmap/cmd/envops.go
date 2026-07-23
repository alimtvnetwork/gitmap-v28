package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runEnvSet sets an environment variable persistently.
func runEnvSet(args []string) {
	fs := flag.NewFlagSet("env-set", flag.ExitOnError)

	var system, verbose, dryRun bool
	var shell string

	fs.BoolVar(&system, constants.FlagEnvSystem, false, constants.FlagDescEnvSystem)
	fs.StringVar(&shell, constants.FlagEnvShell, "", constants.FlagDescEnvShell)
	fs.BoolVar(&verbose, constants.FlagEnvVerbose, false, constants.FlagDescEnvVerbose)
	fs.BoolVar(&dryRun, constants.FlagEnvDryRun, false, constants.FlagDescEnvDryRun)
	fs.Parse(args)

	name := fs.Arg(0)
	value := fs.Arg(1)

	validateEnvName(name)
	validateEnvValue(value)

	if dryRun {
		fmt.Printf(constants.MsgEnvDrySet, name, value)

		return
	}

	setEnvPersistent(name, value, system, shell)
	registry := loadEnvRegistry()
	registry = upsertEnvVariable(registry, name, value)
	saveEnvRegistry(registry)

	fmt.Printf(constants.MsgEnvSet, name, value)
}

// runEnvGet retrieves a managed environment variable value.
func runEnvGet(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrEnvNameRequired)
		os.Exit(1)
	}

	name := args[0]
	registry := loadEnvRegistry()
	entry := findEnvVariable(registry, name)

	fmt.Printf(constants.MsgEnvGetFmt, entry.Name, entry.Value)
}

// runEnvDelete removes a managed environment variable.
func runEnvDelete(args []string) {
	fs := flag.NewFlagSet("env-delete", flag.ExitOnError)

	var system, dryRun bool
	var shell string

	fs.BoolVar(&system, constants.FlagEnvSystem, false, constants.FlagDescEnvSystem)
	fs.StringVar(&shell, constants.FlagEnvShell, "", constants.FlagDescEnvShell)
	fs.BoolVar(&dryRun, constants.FlagEnvDryRun, false, constants.FlagDescEnvDryRun)
	fs.Parse(args)

	name := fs.Arg(0)
	validateEnvName(name)

	if dryRun {
		fmt.Printf(constants.MsgEnvDryDelete, name)

		return
	}

	deleteEnvPersistent(name, system, shell)
	registry := loadEnvRegistry()
	registry = removeEnvVariable(registry, name)
	saveEnvRegistry(registry)

	fmt.Printf(constants.MsgEnvDeleted, name)
}

// runEnvList prints all managed environment variables.
func runEnvList() {
	registry := loadEnvRegistry()

	if len(registry.Variables) == 0 {
		fmt.Print(constants.MsgEnvListEmpty)

		return
	}

	fmt.Print(constants.MsgEnvListHeader)

	for _, v := range registry.Variables {
		fmt.Printf(constants.MsgEnvListRow, v.Name, v.Value)
	}
}

// runEnvPathAdd adds a directory to the system PATH.
func runEnvPathAdd(args []string) {
	fs := flag.NewFlagSet("env-path-add", flag.ExitOnError)

	var system, dryRun bool
	var shell string

	fs.BoolVar(&system, constants.FlagEnvSystem, false, constants.FlagDescEnvSystem)
	fs.StringVar(&shell, constants.FlagEnvShell, "", constants.FlagDescEnvShell)
	fs.BoolVar(&dryRun, constants.FlagEnvDryRun, false, constants.FlagDescEnvDryRun)
	fs.Parse(args)

	dir := fs.Arg(0)
	validateEnvPathDir(dir)

	registry := loadEnvRegistry()
	checkEnvPathNotDuplicate(registry, dir)

	if dryRun {
		fmt.Printf(constants.MsgEnvDryPath, dir)

		return
	}

	addPathPersistent(dir, system, shell)
	registry.Paths = append(registry.Paths, model.EnvPathEntry{Path: dir})
	saveEnvRegistry(registry)

	fmt.Printf(constants.MsgEnvPathAdded, dir)
}

// runEnvPathRemove removes a directory from the system PATH.
func runEnvPathRemove(args []string) {
	fs := flag.NewFlagSet("env-path-remove", flag.ExitOnError)

	var system, dryRun bool
	var shell string

	fs.BoolVar(&system, constants.FlagEnvSystem, false, constants.FlagDescEnvSystem)
	fs.StringVar(&shell, constants.FlagEnvShell, "", constants.FlagDescEnvShell)
	fs.BoolVar(&dryRun, constants.FlagEnvDryRun, false, constants.FlagDescEnvDryRun)
	fs.Parse(args)

	dir := fs.Arg(0)
	if dir == "" {
		fmt.Fprint(os.Stderr, constants.ErrEnvPathRequired)
		os.Exit(1)
	}

	if dryRun {
		fmt.Printf(constants.MsgEnvDryDelete, dir)

		return
	}

	removePathPersistent(dir, system, shell)
	registry := loadEnvRegistry()
	registry = removeEnvPath(registry, dir)
	saveEnvRegistry(registry)

	fmt.Printf(constants.MsgEnvPathRemoved, dir)
}

// runEnvPathList prints all managed PATH entries.
func runEnvPathList() {
	registry := loadEnvRegistry()

	if len(registry.Paths) == 0 {
		fmt.Print(constants.MsgEnvPathEmpty)

		return
	}

	fmt.Print(constants.MsgEnvPathHeader)

	for _, p := range registry.Paths {
		fmt.Printf(constants.MsgEnvPathRow, p.Path)
	}
}
