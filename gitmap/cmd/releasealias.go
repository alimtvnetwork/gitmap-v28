package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runReleaseAlias implements `gitmap release-alias <alias> <version>`.
//
// Behavior:
//  1. Resolve the alias -> absolute repo path via SQLite.
//  2. (--pull) `git pull --ff-only` inside that path.
//  3. Auto-stash dirty changes (unless --no-stash).
//  4. chdir into the repo and invoke the existing runRelease pipeline.
//  5. Pop the auto-stash on the way out.
func runReleaseAlias(args []string, forcePull bool) {
	checkHelp(constants.CmdReleaseAlias, args)

	alias, version, pull, noStash, dryRun := parseRAArgs(args, forcePull)
	target := resolveReleaseAliasPath(alias)
	performReleaseAlias(target, alias, version, pull, noStash, dryRun)
}

// parseRAArgs extracts positional <alias> <version> + flags.
func parseRAArgs(args []string, forcePull bool) (string, string, bool, bool, bool) {
	fs := flag.NewFlagSet(constants.CmdReleaseAlias, flag.ExitOnError)
	pull := fs.Bool(constants.FlagRAPull, forcePull, "git pull --ff-only before releasing")
	noStash := fs.Bool(constants.FlagRANoStash, false, "abort on dirty tree instead of auto-stashing")
	dryRun := fs.Bool(constants.FlagRADryRun, false, "preview without releasing")

	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		os.Exit(2)
	}

	rest := fs.Args()
	if len(rest) != 2 {
		fmt.Fprintln(os.Stderr, constants.ErrRAUsage)
		os.Exit(2)
	}

	return rest[0], rest[1], *pull || forcePull, *noStash, *dryRun
}

// resolveReleaseAliasPath looks up the alias and returns the absolute path.
func resolveReleaseAliasPath(alias string) string {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	resolved, err := db.ResolveAlias(alias)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRAUnknownAliasFmt, alias, alias)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	return resolved.AbsolutePath
}

// performReleaseAlias runs the chdir + pull + stash + release sequence.
func performReleaseAlias(target, alias, version string, pull, noStash, dryRun bool) {
	originalDir, _ := os.Getwd()
	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRAChdirFailedFmt, target, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	fmt.Printf(constants.MsgRAReleasingFmt, alias, target, version)

	if pull {
		runReleaseAliasPull(target)
	}

	if !noStash {
		stashLabel := autoStashIfDirty(target, alias, version)
		if stashLabel != "" {
			defer popAutoStash(target, stashLabel)
		}
	}

	invokeAliasRelease(version, dryRun)
}

// invokeAliasRelease assembles the args expected by runRelease and dispatches.
func invokeAliasRelease(version string, dryRun bool) {
	releaseArgs := []string{version}
	if dryRun {
		releaseArgs = append(releaseArgs, "--dry-run")
	}

	runRelease(releaseArgs)
}
