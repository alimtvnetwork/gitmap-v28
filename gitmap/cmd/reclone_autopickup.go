package cmd

// Auto-discovery of a scan artifact for `gitmap reclone` when the
// user invoked the command with no <file> positional argument and
// no --manifest value.
//
// Convention: `gitmap scan` always writes its artifacts under
// `<scan-root>/.gitmap/output/`. By default we treat the process
// CWD as that root; --scan-root <dir> overrides it so the same
// `reclone` invocation can target a different tree (CI scripts,
// scheduled jobs, ad-hoc inspection from a sibling directory).
//
// Probe order is fixed: JSON first because it is the richest /
// least-lossy representation, CSV as a fallback for environments
// that only kept the spreadsheet view.
//
// We deliberately do NOT walk upward or scan sibling directories:
// silent path-magic across a tree would be surprising and would make
// "which manifest fed this reclone?" un-answerable from the command
// line alone. The MsgCloneNowAutoPickup stderr line documents the
// chosen path so the run is fully reproducible.

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// resolveCloneNowSource picks the input file for `gitmap reclone`
// from three possible sources, in priority order:
//
//  1. --manifest <path>   (explicit, highest priority)
//  2. positional <file>   (legacy form, kept for back-compat)
//  3. auto-pickup         (<scan-root>/.gitmap/output/gitmap.{json,csv})
//
// Supplying BOTH --manifest AND a positional file is a usage error
// (exit 2): rather than silently preferring one, we refuse so the
// run is unambiguous and reproducible. --scan-root is consulted
// ONLY by the auto-pickup branch — it is a no-op when an explicit
// path is provided, so the CLI never has two competing roots.
func resolveCloneNowSource(fs *flag.FlagSet, manifest, scanRoot string) string {
	if manifest != "" && fs.NArg() >= 1 {
		fmt.Fprintf(os.Stderr, constants.MsgCloneNowManifestConflict,
			fs.Arg(0), manifest)
		os.Exit(2)
	}
	if manifest != "" {

		return manifest
	}
	if fs.NArg() >= 1 {

		return fs.Arg(0)
	}
	picked, ok := autoPickupRecloneManifest(scanRoot)
	if !ok {
		printAutoPickupMiss(scanRoot)
		os.Exit(2)
	}
	fmt.Fprintf(os.Stderr, constants.MsgCloneNowAutoPickup, picked)

	return picked
}

// printAutoPickupMiss emits the right miss message depending on
// whether the user explicitly steered auto-pickup at a custom root.
// Echoing the resolved root back on the --scan-root path makes the
// typo / wrong-directory case immediately obvious in scrollback.
func printAutoPickupMiss(scanRoot string) {
	if scanRoot == "" {
		fmt.Fprintln(os.Stderr, constants.MsgCloneNowMissingArg)

		return
	}
	fmt.Fprintf(os.Stderr, constants.MsgCloneNowMissingArgScanRoot+"\n", scanRoot)
}

// autoPickupRecloneManifest returns the first existing scan-artifact
// path under `<scanRoot>/.gitmap/output/` and ok=true, or "", false
// if neither candidate file is present. An empty scanRoot means
// "probe relative to the current process CWD", preserving the
// original behavior. Only regular files count -- a directory at the
// candidate path is treated as a miss rather than silently used.
func autoPickupRecloneManifest(scanRoot string) (string, bool) {
	base := filepath.Join(scanRoot, constants.GitMapDir, constants.OutputDirName)
	candidates := []string{
		filepath.Join(base, constants.DefaultJSONFile),
		filepath.Join(base, constants.DefaultCSVFile),
	}
	for _, path := range candidates {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.Mode().IsRegular() {

			return path, true
		}
	}

	return "", false
}
