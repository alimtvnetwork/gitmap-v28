// Package cmd — macro_add_history_seed.go: seed history and sanitize escape sequences.
package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
	"golang.org/x/term"
)

func seedMacroHistory(t *term.Terminal) {
	macros, err := macro.ListMacros()
	if err != nil {
		return
	}

	for _, m := range macros {
		for _, step := range m.Steps {
			if step.CommandLine != "" {
				t.History.Add(step.CommandLine)
			}
		}
	}
}

func sanitizeRawEscapeCodes(raw string) string {
	cleaned := strings.ReplaceAll(raw, "^[[A", "")
	cleaned = strings.ReplaceAll(cleaned, "^[[B", "")
	cleaned = strings.ReplaceAll(cleaned, "^[[C", "")
	cleaned = strings.ReplaceAll(cleaned, "^[[D", "")
	cleaned = strings.ReplaceAll(cleaned, "\x1b[A", "")
	cleaned = strings.ReplaceAll(cleaned, "\x1b[B", "")
	cleaned = strings.ReplaceAll(cleaned, "\x1b[C", "")
	cleaned = strings.ReplaceAll(cleaned, "\x1b[D", "")

	return cleaned
}
