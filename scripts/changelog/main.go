// Command changelog regenerates CHANGELOG.md and src/data/changelog.ts
// from git commits since the most-recent annotated tag (or all commits if
// no tag exists yet).
//
// Modes (set via -mode flag):
//
//	write — overwrite the two changelog sources with regenerated content.
//	check — write the regenerated content to a temp file and diff against
//	        the on-disk versions; exit non-zero on any drift. Used by CI.
//
// Conventional Commits prefixes (feat:, fix:, docs:, chore:, refactor:,
// perf:, test:, build:, ci:, style:, revert:) are grouped into named
// sections. Commits without a recognised prefix are reported on stderr
// and skipped so a single sloppy "Changes" subject cannot pollute the
// release notes.
package main

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/runner"
)

func main() {
	args, err := runner.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "changelog: %v\n", err)
		os.Exit(2)
	}

	exitCode, err := runner.Execute(args, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changelog: %v\n", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}
