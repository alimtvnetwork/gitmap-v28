// Package runner glues the changelog sub-packages (gitlog, group, render,
// writer, drift) into the two top-level modes: write and check.
package runner

import (
	"flag"
	"fmt"
)

// Mode selects between regeneration and CI drift detection.
type Mode string

const (
	// ModeWrite overwrites the on-disk changelog files.
	ModeWrite Mode = "write"
	// ModeCheck regenerates into memory and exits non-zero on drift.
	ModeCheck Mode = "check"
)

// Args is the parsed CLI surface for the changelog command. Bundled
// into a struct so we can grow the flag set without breaking every
// caller's signature each release.
//
// Field semantics:
//
//   - Mode: write | check (see ModeWrite / ModeCheck).
//   - Version: explicit version label for the new entry. Falls back to
//     `<lastTag>+next` or `vNEXT` when empty.
//   - RepoRoot: repository root containing CHANGELOG.md and
//     src/data/changelog.ts.
//   - Since: optional revision (tag, branch, or hash) used as the
//     lower boundary instead of the highest-semver tag. Enables
//     partial regenerations like `make changelog SINCE=v3.90.0`.
//   - ReleaseTag: optional tag used as the **upper** boundary instead
//     of HEAD. When set, both bounds are explicit and the entry covers
//     exactly `<since>..<release-tag>`. Combined with --since this is
//     how you backfill a missed historical release.
type Args struct {
	Mode       Mode
	Version    string
	RepoRoot   string
	Since      string
	ReleaseTag string
}

// ParseArgs parses command-line flags into Args.
func ParseArgs(args []string) (Args, error) {
	fs := flag.NewFlagSet("changelog", flag.ContinueOnError)
	mode := fs.String("mode", "write", "write | check")
	version := fs.String("version", "", "Version label for the new entry")
	repoRoot := fs.String("repo", ".", "Repository root")
	since := fs.String("since", "", "Lower-bound revision (tag/branch/hash); overrides auto-detected last tag")
	releaseTag := fs.String("release-tag", "", "Upper-bound revision (defaults to HEAD); useful for backfilling a historical release")

	err := fs.Parse(args)
	if err != nil {
		return Args{}, err
	}

	parsed := Mode(*mode)
	if parsed != ModeWrite && parsed != ModeCheck {
		return Args{}, fmt.Errorf("invalid -mode %q (want write or check)", *mode)
	}

	return Args{Mode: parsed, Version: *version, RepoRoot: *repoRoot, Since: *since, ReleaseTag: *releaseTag}, nil
}
