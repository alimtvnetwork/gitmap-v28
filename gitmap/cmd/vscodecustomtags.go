// Package cmd — vscodecustomtags.go: global CLI surface for tuning
// VS Code Project Manager tag detection.
//
// Three repeatable flags (each accepting comma-lists) let users
// override the defaults that ship in
// constants.AutoTagMarkers / AutoTagOrder:
//
//	--vscode-tag <name>             always add to every entry
//	--vscode-tag-skip <name>        drop this auto-detected tag
//	--vscode-tag-marker <file>=<tag> register marker→tag rule
//
// Like `--vscode-sync-disabled`, the flags are stripped from argv
// before any subcommand sees its flagset and persisted into
// GITMAP_VSCODE_TAG_{ADD,SKIP,MARKER} for the current process. The
// vscodepm.DetectTagsCustom helper consumes those env vars, so every
// caller that swapped DetectTags → DetectTagsCustom inherits the
// rules without per-call wiring.
package cmd

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// stripVSCodeTagFlags peels every recognized tag-customization flag
// off args, accumulating values into the matching env vars (joined
// by EnvVSCodeTagSeparator so embedded commas survive). Returns the
// cleaned argv. Unrecognized flags pass through untouched.
func stripVSCodeTagFlags(args []string) []string {
	out := make([]string, 0, len(args))
	add := envExisting(constants.EnvVSCodeTagAdd)
	skip := envExisting(constants.EnvVSCodeTagSkip)
	marker := envExisting(constants.EnvVSCodeTagMarker)

	for i := 0; i < len(args); i++ {
		consumed, valuesAdd, valuesSkip, valuesMarker := matchTagFlag(args, i)
		if consumed == 0 {
			out = append(out, args[i])

			continue
		}
		add = append(add, valuesAdd...)
		skip = append(skip, valuesSkip...)
		marker = append(marker, valuesMarker...)
		i += consumed - 1
	}

	persistTagEnv(constants.EnvVSCodeTagAdd, add)
	persistTagEnv(constants.EnvVSCodeTagSkip, skip)
	persistTagEnv(constants.EnvVSCodeTagMarker, marker)

	return out
}

// matchTagFlag inspects args[i] (and possibly args[i+1]) for any of
// the three tag flags. Returns how many tokens were consumed (0, 1
// or 2) and the values that should be routed into each env list.
// Supports both `--flag value` and `--flag=value` shapes.
func matchTagFlag(args []string, i int) (consumed int, add, skip, marker []string) {
	tok := args[i]
	for _, name := range tagFlagNames() {
		hit, value, took := matchSingleTagFlag(args, i, name)
		if !hit {
			continue
		}
		values := splitCommaList(value)
		switch name {
		case constants.FlagVSCodeTag:
			return took, values, nil, nil
		case constants.FlagVSCodeTagSkip:
			return took, nil, values, nil
		case constants.FlagVSCodeTagMarker:
			return took, nil, nil, values
		}
	}
	_ = tok

	return 0, nil, nil, nil
}

// matchSingleTagFlag tests args[i] against one flag name in both
// short (-name) and long (--name) shapes, with or without `=value`.
func matchSingleTagFlag(args []string, i int, name string) (hit bool, value string, consumed int) {
	tok := args[i]
	long := "--" + name
	short := "-" + name
	for _, prefix := range []string{long, short} {
		if tok == prefix {
			if i+1 >= len(args) {
				return true, "", 1
			}

			return true, args[i+1], 2
		}
		if strings.HasPrefix(tok, prefix+"=") {
			return true, strings.TrimPrefix(tok, prefix+"="), 1
		}
	}

	return false, "", 0
}

// tagFlagNames returns the three recognized flag names in stable order.
func tagFlagNames() []string {
	return []string{
		constants.FlagVSCodeTag,
		constants.FlagVSCodeTagSkip,
		constants.FlagVSCodeTagMarker,
	}
}

// splitCommaList trims and drops empties for one raw flag value.
func splitCommaList(raw string) []string {
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

// envExisting reads an env var that already encodes a tag list (set
// by a prior gitmap call or the user's shell rc) so flag values
// merge with — rather than overwrite — pre-existing config.
func envExisting(name string) []string {
	raw := os.Getenv(name)
	if raw == "" {
		return nil
	}

	return strings.Split(raw, constants.EnvVSCodeTagSeparator)
}

// persistTagEnv writes values back to the named env var using the
// ASCII-unit-separator joiner. Unset when empty so DetectTagsCustom
// short-circuits cleanly.
func persistTagEnv(name string, values []string) {
	if len(values) == 0 {
		os.Unsetenv(name)

		return
	}
	os.Setenv(name, strings.Join(values, constants.EnvVSCodeTagSeparator))
}
