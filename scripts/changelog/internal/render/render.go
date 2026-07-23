// Package render formats a release Entry into Markdown (CHANGELOG.md)
// and TypeScript (src/data/changelog.ts) fragments.
package render

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/group"
)

// Entry is the per-release payload rendered into both formats.
type Entry struct {
	Version string
	Date    string
	Groups  []group.Section
}

// Markdown renders the entry as a CHANGELOG.md block (no trailing
// newline beyond the usual blank-line separator).
func Markdown(e Entry) string {
	var b strings.Builder

	fmt.Fprintf(&b, "## %s — (%s)\n\n", e.Version, e.Date)
	for _, sec := range e.Groups {
		fmt.Fprintf(&b, "### %s\n\n", sec.Name)
		for _, item := range sec.Items {
			fmt.Fprintf(&b, "- %s\n", item)
		}

		b.WriteString("\n")
	}

	return b.String()
}

// TypeScript renders the entry as a src/data/changelog.ts object literal,
// including the leading two-space indent so it can be spliced into the
// existing array directly.
func TypeScript(e Entry) string {
	var b strings.Builder

	b.WriteString("  {\n")
	fmt.Fprintf(&b, "    version: %q,\n", e.Version)
	fmt.Fprintf(&b, "    date: %q,\n", e.Date)
	b.WriteString("    items: [\n")

	for _, sec := range e.Groups {
		for _, item := range sec.Items {
			fmt.Fprintf(&b, "      %q,\n", sec.Name+": "+item)
		}
	}

	b.WriteString("    ],\n")
	b.WriteString("  },\n")

	return b.String()
}
