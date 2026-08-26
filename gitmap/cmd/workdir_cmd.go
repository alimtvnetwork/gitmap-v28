// Package cmd — workdir_cmd.go defines the top-level workdir CLI entrypoint.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runWorkDir handles `gitmap workdir` / `gitmap wd` CLI commands.
func runWorkDir(args []string) {
	db, errDB := store.OpenDefault()
	if errDB == nil {
		_, _ = db.SQL().Exec(store.SQLCreateWorkDirsTable)
		db.Close()
	}

	opts := parseWorkDirFlags(args)
	switch opts.Action {
	case "ls", "list":
		_ = runWorkDirLs()
	case "add":
		_ = runWorkDirAdd(opts.Target, opts.Label)
	case "rm", "remove", "delete":
		_ = runWorkDirRm(opts.Target)
	case "set-default", "default":
		_ = runWorkDirSetDefault(opts.Target)
	default:
		fmt.Printf("Usage: gitmap workdir [ls|add <path>|rm <path|id>|set-default <path|id>]\n")
		os.Exit(1)
	}
}
