// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/render"
)

// runChangelog handles the 'changelog' command.
func runChangelog(args []string) {
	checkHelp("changelog", args)
	cleaned, mode := ParsePrettyFlag(args)
	pretty := render.Decide(mode, render.StdoutIsTerminal(), true)
	version, latest, limit, openFile, source := parseChangelogFlags(cleaned)
	version, openFile = resolveChangelogAlias(version, openFile)
	if openFile {
		handleChangelogOpen(latest, version)
	}
	if !latest && len(version) == 0 && openFile {
		return
	}

	dispatchChangelogOutput(version, latest, limit, source, pretty)
}

// resolveChangelogAlias detects if the version arg is actually a file-open alias.
func resolveChangelogAlias(version string, openFile bool) (string, bool) {
	if strings.EqualFold(version, constants.ChangelogFile) || strings.EqualFold(version, constants.CmdChangelogMD) {
		return "", true
	}

	return version, openFile
}

// handleChangelogOpen opens the changelog file and exits on error.
func handleChangelogOpen(latest bool, version string) {
	err := openChangelogFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChangelogOpen, constants.ChangelogFile, err)
		os.Exit(1)
	}
	if !latest && len(version) == 0 {
		os.Exit(0)
	}
}

// dispatchChangelogOutput prints the appropriate changelog entries.
// `pretty` controls ANSI rendering for headers + bullet bodies; pass
// false to emit terminal-safe plain text (no escape codes anywhere).
func dispatchChangelogOutput(version string, latest bool, limit int, source string, pretty bool) {
	entries, err := release.ReadChangelog()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChangelogRead, constants.ChangelogFile, err)
		os.Exit(1)
	}
	entries = filterChangelogBySource(entries, source)
	if latest {
		printChangelogEntries(entries, 1, pretty)

		return
	}
	if len(version) > 0 {
		printSingleVersion(entries, version, pretty)

		return
	}
	printChangelogEntries(entries, limit, pretty)
}

// filterChangelogBySource keeps only entries whose version exists in the DB with the given source.
func filterChangelogBySource(entries []release.ChangelogEntry, source string) []release.ChangelogEntry {
	if source == "" {
		return entries
	}

	sources := loadChangelogSourceMap()
	var filtered []release.ChangelogEntry
	for _, e := range entries {
		tag := release.NormalizeVersion(e.Version)
		if sources[tag] == source {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

// loadChangelogSourceMap reads the Releases table to build a tag→source map.
func loadChangelogSourceMap() map[string]string {
	db, err := openDB()
	if err != nil {
		return map[string]string{}
	}
	defer db.Close()

	releases, err := db.ListReleases()
	if err != nil {
		return map[string]string{}
	}

	m := make(map[string]string, len(releases))
	for _, r := range releases {
		m[r.Tag] = r.Source
	}

	return m
}

// printSingleVersion finds and prints one version's changelog.
func printSingleVersion(entries []release.ChangelogEntry, version string, pretty bool) {
	entry, found := release.FindChangelogEntry(entries, version)
	if !found {
		fmt.Fprintf(os.Stderr, constants.ErrChangelogVersionNotFound, release.NormalizeVersion(version))
		os.Exit(1)
	}
	printChangelogEntry(entry, pretty)
}

// parseChangelogFlags parses flags for the changelog command. The
// --pretty / --no-pretty flag is stripped by the caller before reaching
// here, so this FlagSet stays focused on the command's own flags.
func parseChangelogFlags(args []string) (version string, latest bool, limit int, openFile bool, source string) {
	fs := flag.NewFlagSet(constants.CmdChangelog, flag.ExitOnError)
	latestFlag := fs.Bool("latest", false, constants.FlagDescLatest)
	limitFlag := fs.Int("limit", 5, constants.FlagDescLimit)
	openFlag := fs.Bool("open", false, constants.FlagDescOpenChangelog)
	sourceFlag := fs.String("source", "", constants.FlagDescSource)
	_ = fs.Parse(args)

	version = ""
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}
	if *limitFlag < 1 {
		*limitFlag = 1
	}

	return version, *latestFlag, *limitFlag, *openFlag, *sourceFlag
}

// printChangelogEntries prints the newest N changelog entries.
func printChangelogEntries(entries []release.ChangelogEntry, limit int, pretty bool) {
	if limit > len(entries) {
		limit = len(entries)
	}
	for i := 0; i < limit; i++ {
		printChangelogEntry(entries[i], pretty)
	}
}

// printChangelogEntry prints a single changelog entry. When pretty is
// true the rich renderer (colored header + inline-markdown bullets +
// word wrapping) is used; when false, the same layout is emitted with
// every ANSI escape sequence suppressed — output is safe to redirect
// into a file or pipe into tools that choke on color codes.
func printChangelogEntry(entry release.ChangelogEntry, pretty bool) {
	renderChangelogEntry(entry, pretty)
}

// openChangelogFile opens CHANGELOG.md with the default OS app.
func openChangelogFile() error {
	absPath, err := filepath.Abs(constants.ChangelogFile)
	if err != nil {
		return err
	}

	return runOpenCommand(absPath)
}

// runOpenCommand executes the platform-specific open command.
func runOpenCommand(path string) error {
	if runtime.GOOS == constants.OSWindows {
		cmd := exec.Command(constants.CmdWindowsShell, constants.CmdArgSlashC, constants.CmdArgStart, constants.CmdArgEmpty, path)

		return cmd.Run()
	}
	if runtime.GOOS == constants.OSDarwin {
		cmd := exec.Command(constants.CmdOpen, path)

		return cmd.Run()
	}
	cmd := exec.Command(constants.CmdXdgOpen, path)

	return cmd.Run()
}
