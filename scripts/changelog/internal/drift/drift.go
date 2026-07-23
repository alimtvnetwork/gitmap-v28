// Package drift implements `make changelog-check`: it regenerates the
// changelog into a temp copy of each file and reports drift so CI can
// gate PRs that forgot to commit the regeneration.
package drift

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/render"
	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/writer"
)

const driftExitCode = 3

// Check returns driftExitCode (3) when either changelog file would
// change, 0 when both already include the rendered entry, and a
// non-nil error on any I/O problem.
func Check(mdPath, tsPath string, entry render.Entry) (int, error) {
	mdDrift, err := wouldChange(mdPath, render.Markdown(entry))
	if err != nil {
		return 0, err
	}

	tsDrift, err := wouldChange(tsPath, render.TypeScript(entry))
	if err != nil {
		return 0, err
	}

	if !mdDrift && !tsDrift {
		fmt.Fprintln(os.Stdout, "changelog: up to date.")

		return 0, nil
	}

	reportDrift(mdPath, mdDrift)
	reportDrift(tsPath, tsDrift)

	return driftExitCode, nil
}

func wouldChange(path, fragment string) (bool, error) {
	current, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	return !contains(string(current), fragment), nil
}

func contains(haystack, needle string) bool {
	return len(needle) > 0 && len(haystack) >= len(needle) && stringContains(haystack, needle)
}

// stringContains is split out so writer's import graph stays minimal.
func stringContains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}

	return false
}

func reportDrift(path string, drifted bool) {
	if !drifted {
		return
	}

	fmt.Fprintf(os.Stderr, "changelog: %s is out of date — run `make changelog`.\n", path)
}

// Ensure writer is referenced so future drift-by-file-rewrite variants
// can swap straight to writer.PrependBoth without import churn.
var _ = writer.PrependBoth
