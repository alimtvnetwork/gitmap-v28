// Package cmd — workdir_cmd.go defines the top-level workdir CLI entrypoint.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runWorkDir handles `gitmap workdir` / `gitmap wd` CLI commands.
func runWorkDir(args []string) error {
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
	case "set-default", "default", "set":
		_ = runWorkDirSetDefault(opts.Target)
	default:
		fmt.Printf("Usage: gitmap workdir [ls|add <path>|rm <path|id>|set <path|id>|default]\n")
		return apperror.New("fatal error", "E9000", nil)
	}
	return nil
}
