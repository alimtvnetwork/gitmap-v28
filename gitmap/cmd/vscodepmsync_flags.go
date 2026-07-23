// Package cmd — vscodepmsync_flags.go: flag-parsing surface for
// `gitmap vscode-pm-sync` (alias `vpm`).
//
// Lives in its own file so the runner in vscodepmsync.go can stay
// under the 200-line strict-style budget while we keep growing the
// CLI. All four documented flags are parsed here:
//
//   - --dry-run         (bool)         preview, no write
//   - --mode <m>        (string)       union|replace|intersection
//   - --projects-json   (string)       absolute override path
//   - --tag <name>      (repeatable)   replace per-pair detected tags
//
// `--tag` is a custom flag.Value supporting BOTH repetition AND
// comma-separated values, so `--tag a --tag b,c` produces {a,b,c}.
// The order is preserved within a single Set call but de-duplicated
// across calls — first occurrence wins. The `gitmap` brand tag is
// NOT auto-prepended in this mode (per the constants_cli.go
// docstring): callers in full control means callers must pass it
// explicitly if they want it.
package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

// vscodePMSyncOpts is the parsed CLI surface for one runner invocation.
type vscodePMSyncOpts struct {
	DryRun         bool
	Mode           vscodepm.MergeMode
	ProjectsJSON   string
	TagOverride    []string
	HasTagOverride bool
}

// tagListValue implements flag.Value to capture --tag repetitions and
// comma-separated values into a deduplicated string slice. The
// "set-at-all" signal is preserved via wasSet so the runner can tell
// `--tag` (no values) apart from "flag never passed".
type tagListValue struct {
	values []string
	seen   map[string]bool
	wasSet bool
}

// newTagListValue returns a fresh, empty tagListValue ready for
// flag.Var registration.
func newTagListValue() *tagListValue {
	return &tagListValue{seen: map[string]bool{}}
}

// String renders the captured tags as a comma list — matches what
// the user would have typed.
func (t *tagListValue) String() string {
	if t == nil {
		return ""
	}

	return strings.Join(t.values, ",")
}

// Set splits raw on commas, trims whitespace, and appends each
// non-empty value not already captured. Empty entries (e.g. trailing
// comma) are dropped silently — they are never useful.
func (t *tagListValue) Set(raw string) error {
	t.wasSet = true

	for _, part := range strings.Split(raw, ",") {
		v := strings.TrimSpace(part)
		if v == "" || t.seen[v] {
			continue
		}

		t.seen[v] = true
		t.values = append(t.values, v)
	}

	return nil
}

// parseVSCodePMSyncFlags parses every documented flag. Unknown
// --mode values fail loud (per the zero-swallow rule) rather than
// silently defaulting to union. Returns ExitOnError-friendly errors
// the caller can write to stderr before exit(2).
func parseVSCodePMSyncFlags(args []string) (vscodePMSyncOpts, error) {
	fs := flag.NewFlagSet(constants.CmdVSCodePMSync, flag.ExitOnError)

	dryRun := fs.Bool(
		constants.FlagVSCodePMSyncDryRun, false,
		constants.FlagDescVSCodePMSyncDryRun,
	)
	modeRaw := fs.String(
		constants.FlagVSCodePMSyncMode, constants.VSCodePMSyncModeUnion,
		constants.FlagDescVSCodePMSyncMode,
	)
	projectsJSON := fs.String(
		constants.FlagVSCodePMSyncProjectsJSON, "",
		constants.FlagDescVSCodePMSyncProjectsJSON,
	)

	tags := newTagListValue()
	fs.Var(tags, constants.FlagVSCodePMSyncTag, constants.FlagDescVSCodePMSyncTag)

	if err := fs.Parse(args); err != nil {
		return vscodePMSyncOpts{}, fmt.Errorf("vscode-pm-sync: %w", err)
	}

	mode, err := vscodepm.ParseMergeMode(*modeRaw)
	if err != nil {
		return vscodePMSyncOpts{}, err
	}

	return vscodePMSyncOpts{
		DryRun:         *dryRun,
		Mode:           mode,
		ProjectsJSON:   *projectsJSON,
		TagOverride:    tags.values,
		HasTagOverride: tags.wasSet,
	}, nil
}
