// Package cmd — `gitmap zip` (alias `z`).
//
// Resolves N heterogeneous sources (folders, archive URLs, git repos)
// into local paths, then runs CreateArchive into the user-supplied
// --out path. Compression mode is chosen via mutually exclusive flags
// (--best / --fast / --standard, with -s as a synonym for standard).
//
// Each invocation writes one ArchiveHistory row (in-flight at start,
// finalized at end) so a partial failure still leaves a forensic trace.
package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/archive"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runZip is the dispatch entrypoint for `zip` / `z`.
func runZip(args []string) {
	checkHelp(constants.CmdZip, args)

	opts, sources, err := parseZipFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}
	if opts.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "  ✗ "+constants.ErrArchiveCreateNeedsOut)
		os.Exit(1)
	}
	if len(sources) == 0 {
		fmt.Fprintf(os.Stderr, "  ✗ "+constants.ErrArchiveNoSource+"\n", constants.CmdZip)
		os.Exit(1)
	}

	ctx := context.Background()
	resolved, err := resolveAllSources(ctx, sources)
	defer cleanupAllSources(resolved)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		cleanupAllSources(resolved)
		exitWith(1)
	}

	opts.Sources = resolvedToPaths(resolved)
	executeZip(ctx, opts, sources)
}

// zipFlags is the parsed flag set the cmd layer feeds into CreateArchive.
type zipFlags struct {
	archive.CreateOptions
	best     bool
	fast     bool
	standard bool
	include  string
	exclude  string
}

// parseZipFlags returns CreateOptions plus the positional source list.
func parseZipFlags(args []string) (archive.CreateOptions, []string, error) {
	fs := flag.NewFlagSet(constants.CmdZip, flag.ContinueOnError)
	z := zipFlags{}

	fs.StringVar(&z.OutputPath, constants.FlagArchiveOut, "", constants.FlagDescArchiveOut)
	fs.StringVar(&z.OutputPath, constants.FlagArchiveOutShort, z.OutputPath, constants.FlagDescArchiveOut)
	fs.BoolVar(&z.best, constants.FlagArchiveBest, false, constants.FlagDescArchiveBest)
	fs.BoolVar(&z.fast, constants.FlagArchiveFast, false, constants.FlagDescArchiveFast)
	fs.BoolVar(&z.standard, constants.FlagArchiveStandard, false, constants.FlagDescArchiveStandard)
	fs.BoolVar(&z.standard, constants.FlagArchiveStdShort, z.standard, constants.FlagDescArchiveStandard)
	fs.StringVar(&z.include, constants.FlagArchiveInclude, "", constants.FlagDescArchiveInclude)
	fs.StringVar(&z.exclude, constants.FlagArchiveExclude, "", constants.FlagDescArchiveExclude)

	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		return archive.CreateOptions{}, nil, err
	}

	mode, err := resolveCompressionMode(z.best, z.fast, z.standard)
	if err != nil {
		return archive.CreateOptions{}, nil, err
	}
	z.Mode = mode
	z.Includes = splitCSV(z.include)
	z.Excludes = splitCSV(z.exclude)

	return z.CreateOptions, fs.Args(), nil
}

// resolveCompressionMode enforces "at most one of --best / --fast /
// --standard" and defaults to standard when no flag is set.
func resolveCompressionMode(best, fast, standard bool) (archive.CompressionMode, error) {
	count := 0
	for _, b := range []bool{best, fast, standard} {
		if b {
			count++
		}
	}
	if count > 1 {
		return "", fmt.Errorf(constants.ErrArchiveBadCompression, constants.CmdZip)
	}
	switch {
	case best:
		return archive.ModeBest, nil
	case fast:
		return archive.ModeFast, nil
	}

	return archive.ModeStandard, nil
}

// splitCSV trims and splits a comma-separated glob list, dropping empty
// segments so "a, ,b" yields ["a", "b"].
func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}

// resolveAllSources turns each user input into a ResolvedSource. On the
// first failure the previously-resolved entries are still returned so
// the caller's deferred cleanupAllSources can wipe their temp dirs.
func resolveAllSources(ctx context.Context, sources []string) ([]archive.ResolvedSource, error) {
	out := make([]archive.ResolvedSource, 0, len(sources))
	for _, s := range sources {
		fmt.Fprintf(os.Stderr, constants.MsgArchiveResolving+"\n", 1)
		r, err := archive.ResolveSource(ctx, s)
		if err != nil {
			return out, err
		}
		fmt.Fprintf(os.Stderr, constants.MsgArchiveResolved+"\n", s, r.LocalPath)
		out = append(out, r)
	}

	return out, nil
}

// cleanupAllSources removes every temp dir created by resolveAllSources.
func cleanupAllSources(rs []archive.ResolvedSource) {
	for _, r := range rs {
		archive.CleanupResolved(r)
	}
}

// resolvedToPaths flattens to the absolute LocalPath list CreateArchive
// expects.
func resolvedToPaths(rs []archive.ResolvedSource) []string {
	out := make([]string, len(rs))
	for i, r := range rs {
		out[i] = r.LocalPath
	}

	return out
}

// executeZip runs CreateArchive with history persistence. originalSrcs
// is preserved (not the resolved temp paths) so the history row records
// what the user actually typed.
func executeZip(ctx context.Context, opts archive.CreateOptions, originalSrcs []string) {
	db, dbErr := openDB()
	var historyID int64
	if dbErr == nil {
		defer db.Close()
		if migErr := db.Migrate(); migErr == nil {
			historyID = startArchiveRow(db, constants.ArchiveCmdZip, originalSrcs, string(opts.Mode))
		}
	}

	fmt.Fprintf(os.Stderr, constants.MsgArchiveBanner+"\n", constants.CmdZip, constants.Version)
	fmt.Fprintf(os.Stderr, constants.MsgArchiveCreateStart+"\n", opts.OutputPath, opts.Mode, len(opts.Sources))

	res, err := archive.CreateArchive(ctx, opts)
	finishArchiveRow(db, historyID, res.OutputPath, string(res.Format), false, err)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, constants.MsgArchiveCreateDone+"\n", res.OutputPath)
}
