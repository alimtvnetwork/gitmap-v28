// Package completion generates shell tab-completion scripts for gitmap.
//
//go:generate go run ./internal/gencommands
package completion

import (
	"fmt"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// manualExtras is a stop-gap for command names whose Cmd* constants cannot
// be discovered by the marker-driven generator (e.g. computed strings,
// build-tag-gated files).
//
// As of v2.99.0 this slice is empty: inclusion is controlled locally inside
// each constants_*.go file via the marker comments documented in
// internal/gencommands/main.go:
//
//   - Add `// gitmap:cmd top-level` immediately above a `const (...)` block
//     to opt every Cmd* string constant in that block into completion.
//   - Add `// gitmap:cmd skip` as the trailing line comment of an individual
//     ValueSpec to exclude it (e.g. for subcommand IDs like "create" used by
//     `gitmap group`).
//
// Run `go generate ./completion/...` after editing constants to refresh
// allcommands_generated.go. Domain owners never need to edit the generator.
var manualExtras = []string{"db", "start-fresh", "find-duplicates", "pe", "ee", "pd"}

// CustomGenerator allows package cmd to provide dynamic Cobra completion scripts.
var CustomGenerator func(shell string) (string, error)

// Generate returns the completion script for the given shell.
func Generate(shell string) (string, error) {
	if script, isGenerated := tryCustomGenerator(shell); isGenerated {
		return script, nil
	}

	switch shell {
	case constants.ShellPowerShell:
		return generatePowerShell(), nil
	case constants.ShellBash:
		return generateBash(), nil
	case constants.ShellZsh:
		return generateZsh(), nil
	case constants.ShellFish:
		return generateFish(), nil
	default:
		return "", fmt.Errorf(constants.ErrCompUnknownShell, shell)
	}
}

func tryCustomGenerator(shell string) (string, bool) {
	if CustomGenerator == nil {
		return "", false
	}
	script, err := CustomGenerator(shell)
	if err != nil || len(script) == 0 {
		return "", false
	}
	return script, true
}

// AllCommands returns every command name and alias offered by tab-completion.
//
// The list is the union of generatedCommands (auto-extracted from the
// constants package by internal/gencommands) and manualExtras. Run
// `go generate ./completion/...` after adding new Cmd* constants to refresh
// the generated portion.
func AllCommands() []string {
	seen := make(map[string]bool, len(generatedCommands)+len(manualExtras))

	for _, v := range generatedCommands {
		seen[v] = true
	}

	for _, v := range manualExtras {
		seen[v] = true
	}

	out := make([]string, 0, len(seen))
	for v := range seen {
		out = append(out, v)
	}

	sort.Strings(out)

	return out
}
