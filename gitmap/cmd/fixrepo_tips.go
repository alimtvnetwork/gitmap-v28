package cmd

// User-facing hints emitted by `gitmap fix-repo` (v5.45.0+):
//   - fixRepoFlagHint: appended to E_BAD_FLAG so users immediately see
//     the accepted flag surface (including the new bare-digit span).
//   - emitFixRepoTips: trailing block printed after a successful sweep
//     reminding the user how to undo, dry-run, narrow scope, etc.

import (
	"fmt"
	"os"
)

// fixRepoFlagHint returns the multi-line accepted-flag reference that
// the bad-flag error appends. Kept as a single string so it can be
// reused verbatim from `gitmap help fix-repo`.
func fixRepoFlagHint() string {
	return "  accepted flags:\n" +
		"    -2 | -3 | -5 | -N | N        span: rewrite the last N prior versions (default: -2)\n" +
		"    --all                        rewrite every prior version v1..v(K-1) -> vK\n" +
		"    --dry-run    (-DryRun)       preview only; no file is written\n" +
		"    --verbose    (-Verbose)      print every modified file with replacement count\n" +
		"    --strict     (-Strict)       post-rewrite `go test` on touched Go packages\n" +
		"    --restrict no-version | -r nv\n" +
		"                                 skip the v1->v2 bare-base sweep ({base}-vN only)\n" +
		"    --config <path>              override fix-repo.config.json location\n" +
		"  examples:\n" +
		"    gitmap fix-repo               # default: last 2 prior versions\n" +
		"    gitmap fix-repo 4             # widen window to last 4 prior versions\n" +
		"    gitmap fix-repo --all --dry-run\n" +
		"    gitmap fix-repo -2 -r nv      # skip bare-base sweep on v1->v2"
}

// emitFixRepoTips prints the post-run options block. Always rendered
// to stderr so it never contaminates dry-run stdout consumed by
// scripts. Includes the gitmap-undo reminder which is the single most
// common follow-up after an unintended rewrite.
func emitFixRepoTips(opts fixRepoOptions, changed int) {
	if opts.isDryRun {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "tips (dry-run):")
		fmt.Fprintln(os.Stderr, "  - re-run without --dry-run to apply the rewrite")
		fmt.Fprintln(os.Stderr, "  - narrow scope: --restrict no-version | -r nv  (skip bare-base on v1->v2)")
		fmt.Fprintln(os.Stderr, "  - widen window: gitmap fix-repo N   (e.g. 3, 4, 7)  or  --all")
		return
	}
	if changed == 0 {
		return
	}
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "next steps:")
	fmt.Fprintln(os.Stderr, "  - undo:        gitmap undo            (restores the snapshot just written)")
	fmt.Fprintln(os.Stderr, "  - list backups: gitmap undo --list")
	fmt.Fprintln(os.Stderr, "  - preview only: gitmap fix-repo --dry-run")
	fmt.Fprintln(os.Stderr, "  - narrow scope: gitmap fix-repo -r nv  (skip bare {base} sweep on v1->v2)")
	fmt.Fprintln(os.Stderr, "  - widen window: gitmap fix-repo N      (e.g. 3, 4, 7)  or  --all")
}
