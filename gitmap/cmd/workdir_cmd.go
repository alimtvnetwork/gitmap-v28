// Package cmd — workdir_cmd.go defines the top-level workdir CLI entrypoint.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runWorkDir handles `gitmap workdir` / `gitmap wd` CLI commands.
func runWorkDir(args []string) error {
	checkHelp("workdir", args)

	ensureWorkDirsTableExists()

	opts := parseWorkDirFlags(args)

	return dispatchWorkDirAction(opts)
}

func ensureWorkDirsTableExists() {
	db, errDB := store.OpenDefault()
	if errDB == nil {
		_, _ = db.SQL().Exec(store.SQLCreateWorkDirsTable)
		db.Close()
	}
}

func dispatchWorkDirAction(opts workDirOptions) error {
	switch opts.Action {
	case "ls", "list":
		return runWorkDirLs()
	case "add":
		return runWorkDirAdd(opts.Target, opts.Label)
	case "rm", "remove", "delete":
		return runWorkDirRm(opts.Target)
	case "set-default", "set":
		return runWorkDirSetDefault(opts.Target)
	case "default":
		return runWorkDirDefault(opts.Target)
	case "path", "get", "current":
		return runWorkDirPath()
	case "help", "-h", "--help":
		printWorkDirUsage()

		return nil
	default:
		if isDirectoryPath(opts.Action) {
			return runWorkDirSetDefault(opts.Action)
		}

		printWorkDirUsage()

		return nil
	}
}

func isDirectoryPath(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

func printWorkDirUsage() {
	fmt.Println("Usage: gitmap workdir [ls|add <path>|rm <path|id>|set <path|id>|default [path]|path]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  ls                         List all registered work directories")
	fmt.Println("  add <path> [--label <l>]   Register a work directory")
	fmt.Println("  rm <path|id>               Remove a registered work directory")
	fmt.Println("  set <path|id>              Set active default work directory")
	fmt.Println("  default [path]             Display or set default work directory")
	fmt.Println("  path                       Print absolute path of active default work directory")
}
