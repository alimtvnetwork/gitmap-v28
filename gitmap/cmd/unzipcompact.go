// Package cmd — `gitmap unzip-compact` (alias `uzc`).
//
// Resolves an input source (local archive, HTTP(S) URL, or auto-detect a
// single archive in the current folder) and runs the compact-extract
// algorithm into either the user-supplied destination folder or the
// current working directory.
//
// Listing mode (--list / -l) skips extraction and prints the archive's
// entry table to stderr.
package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/archive"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runUnzipCompact is the dispatch entrypoint for `unzip-compact` / `uzc`.
func runUnzipCompact(args []string) error {
	checkHelp(constants.CmdUnzipCompact, args)

	listMode, positional := parseUnzipCompactFlags(args)
	src, dest, err := resolveUnzipInputs(positional)
	if err != nil {
		appErr := apperror.NewWithDetails(
			"cmd.unzipcompact.inputs",
			"E1126",
			fmt.Sprintf(constants.ErrArchiveNoSource, constants.CmdUnzipCompact),
			"cmd.unzipcompact",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}

	ctx := context.Background()
	resolved, err := archive.ResolveSource(ctx, src)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.unzipcompact.resolve",
			"E1127",
			"failed to resolve archive source",
			"cmd.unzipcompact",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"src": src},
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	defer archive.CleanupResolved(resolved)

	if listMode {
		runListMode(ctx, resolved.LocalPath)
		return nil
	}

	executeCompactExtract(ctx, resolved, src, dest)
	return nil
}

// parseUnzipCompactFlags pulls --list off the args and returns the rest
// in original order.
func parseUnzipCompactFlags(args []string) (listMode bool, positional []string) {
	fs := flag.NewFlagSet(constants.CmdUnzipCompact, flag.ContinueOnError)
	fs.BoolVar(&listMode, constants.FlagArchiveList, false, constants.FlagDescArchiveList)
	fs.BoolVar(&listMode, constants.FlagArchiveListShrt, listMode, constants.FlagDescArchiveList)
	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.unzipcompact.parseFlags",
			"E1128",
			"failed to parse unzip-compact flags",
			"cmd.unzipcompact",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)
	}

	return listMode, fs.Args()
}

// resolveUnzipInputs implements the spec's positional rules:
//
//	uzc                       → auto-detect single archive in cwd
//	uzc <src>                 → src + cwd
//	uzc <src> <dest>          → src + dest
func resolveUnzipInputs(positional []string) (src, dest string, err error) {
	cwd, _ := os.Getwd()
	switch len(positional) {
	case 0:
		picked, perr := archive.AutoDetectSingleArchive(cwd)
		if perr != nil {
			return "", "", perr
		}
		fmt.Fprintf(os.Stderr, constants.MsgArchiveAutoPicked+"\n", picked)

		return picked, cwd, nil
	case 1:
		return positional[0], cwd, nil
	default:
		return positional[0], positional[1], nil
	}
}

// runListMode prints the entry table for the archive at path and exits.
func runListMode(ctx context.Context, path string) error {
	entries, format, err := archive.ListEntries(ctx, path)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.unzipcompact.listEntries",
			"E1129",
			"failed to list archive entries",
			"cmd.unzipcompact",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": path},
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	fmt.Fprintf(os.Stderr, constants.MsgArchiveListHeader+"\n", path, format, len(entries))
	for _, e := range entries {
		fmt.Fprintf(os.Stderr, constants.MsgArchiveListEntry+"\n", e.Path, e.Size)
	}
	return nil
}

// executeCompactExtract is the success-path: open DB, persist a
// pre-flight history row, run the extraction, finalize the row.
func executeCompactExtract(ctx context.Context, resolved archive.ResolvedSource, originalSrc, dest string) {
	db, dbErr := openDB()
	var historyID int64
	if dbErr == nil {
		defer db.Close()
	}
	migErr := error(nil)
	if dbErr == nil {
		migErr = db.Migrate()
	}
	if dbErr == nil && migErr == nil {
		historyID = startArchiveRow(db, constants.ArchiveCmdUnzipCompact, []string{originalSrc}, "")
	}

	fmt.Fprintf(os.Stderr, constants.MsgArchiveBanner+"\n", constants.CmdUnzipCompact, constants.Version)
	fmt.Fprintf(os.Stderr, constants.MsgArchiveExtractStart+"\n", resolved.LocalPath, dest)

	res, err := archive.CompactExtract(ctx, resolved.LocalPath, dest)
	finishArchiveRow(db, historyID, res.OutputDir, string(res.Format), res.UsedTempDir, err)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.unzipcompact.compactExtract",
			"E1130",
			"compact extract failed",
			"cmd.unzipcompact",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": resolved.LocalPath, "dest": dest},
		)
		cliexit.HandleError(appErr, 1)
		return
	}

	fmt.Fprintf(os.Stderr, constants.MsgArchiveExtractDone+"\n", res.OutputDir, res.EntriesWritten, res.Format)
	if res.FlattenedLayers > 0 {
		fmt.Fprintf(os.Stderr, constants.MsgArchiveCompactFlatten+"\n", res.FlattenedLayers)
	}
}

// startArchiveRow inserts an in-flight ArchiveHistory row and surfaces
// (but does not fail on) write errors. The archive command keeps running
// even when history persistence fails.
func startArchiveRow(db *store.DB, cmd string, inputs []string, mode string) int64 {
	id, err := db.StartArchiveHistory(cmd, inputs, mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnArchiveHistoryWrite+"\n", err)

		return 0
	}

	return id
}

// finishArchiveRow updates an in-flight row with the final outcome. Safe
// to call with id == 0 (no-op) when the start failed.
func finishArchiveRow(db *store.DB, id int64, outputPath, format string, usedTemp bool, runErr error) {
	if db == nil || id == 0 {
		return
	}
	status := constants.ArchiveStatusSuccess
	errMsg := ""
	if runErr != nil {
		status = constants.ArchiveStatusFailed
		errMsg = runErr.Error()
	}
	if err := db.FinishArchiveHistory(id, outputPath, format, status, errMsg, usedTemp); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnArchiveHistoryWrite+"\n", err)

		return
	}
	fmt.Fprintf(os.Stderr, constants.MsgArchiveHistoryRecorded+"\n", id, status)
}
