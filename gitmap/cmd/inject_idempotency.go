// Package cmd — inject_idempotency.go: shared helpers used by `gitmap
// inject` and `gitmap open` to skip Desktop / VS Code re-registration
// when the per-tool stamp on the Repo row is already set.
//
// The `--force` (`-f`) flag bypasses both checks AND zeroes the stamps
// so the post-action UPDATE re-stamps to "now".
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// parseInjectForceFlag extracts --force/-f from the trailing flag set
// shared by `inject` and `open`. Unknown flags are surfaced via the
// usual flag.ExitOnError contract.
func parseInjectForceFlag(name string, args []string) bool {
	fs := flag.NewFlagSet(name, flag.ExitOnError)

	var force bool

	fs.BoolVar(&force, constants.FlagInjectForce, false, constants.FlagDescInjectForce)
	fs.BoolVar(&force, constants.FlagInjectForceShort, false, constants.FlagDescInjectForce)

	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		os.Exit(2)
	}

	return force
}

// loadInjectStamps best-effort reads the Repo row's per-tool stamps.
// Errors degrade to "never injected" (zero value) so a transient DB
// hiccup never blocks the user-visible Desktop / VS Code calls.
func loadInjectStamps(absPath string) store.InjectTimestamps {
	db, err := openDB()
	if err != nil {
		return store.InjectTimestamps{}
	}
	defer db.Close()

	ts, err := db.GetInjectTimestamps(absPath)
	if err != nil {
		return store.InjectTimestamps{}
	}

	return ts
}

// markInjected best-effort stamps the column for `kind`. Failures are
// warned but never fatal; the side-effect on Desktop/VS Code already
// happened, the stamp is just bookkeeping.
func markInjected(absPath string, kind constants.InjectKind) {
	db, err := openDB()
	if err != nil {
		return
	}
	defer db.Close()

	if err := db.MarkInjected(absPath, kind); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ inject: could not stamp timestamp: %v\n", err)
	}
}

// shouldRunDesktop reports whether GitHub Desktop registration should
// run, and prints the skip notice when it shouldn't. `force` always
// wins.
func shouldRunDesktop(absPath string, ts store.InjectTimestamps, force bool) bool {
	if force || len(ts.Desktop) == 0 {
		return true
	}

	fmt.Printf(constants.MsgInjectSkipDesktopFmt,
		constants.ColorCyan, constants.ColorReset,
		ts.Desktop,
		constants.ColorYellow, constants.ColorReset,
	)

	return false
}

// shouldRunVSCode mirrors shouldRunDesktop for the VS Code slot.
func shouldRunVSCode(absPath string, ts store.InjectTimestamps, force bool) bool {
	if force || len(ts.VSCode) == 0 {
		return true
	}

	fmt.Printf(constants.MsgInjectSkipVSCodeFmt,
		constants.ColorCyan, constants.ColorReset,
		ts.VSCode,
		constants.ColorYellow, constants.ColorReset,
	)

	return false
}
