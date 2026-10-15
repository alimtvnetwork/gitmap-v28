// Package helpdoc — examples_derive.go: derived registry entries.
//
// Seeds the help-example registry for every documented topic without a
// curated entry, so `gitmap help example` covers the full help surface.
// WhatItDoes comes from the topic summaries (catalog.go) or the embedded
// help file's title + first paragraph; the invocation is the topic's own
// help screen (always runnable) with a normalized expected-output shape.
package helpdoc

import (
	"sort"
	"strings"
)

// deriveSummaryExamples seeds the registry for every topic summary that
// has no curated entry. The invocation is the topic's own help screen —
// always runnable — and the expected output is its normalized shape.
func deriveSummaryExamples(reg map[string]HelpExample) {
	keys := make([]string, 0, len(topicSummaries))
	for key := range topicSummaries {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		if _, exists := reg[key]; exists {
			continue
		}

		helpVerb := strings.ReplaceAll(key, "_", "-")
		reg[key] = HelpExample{
			WhatItDoes: topicSummaries[key],
			Examples: []Example{{
				Invocation: "gitmap help " + helpVerb,
				ExpectedOutput: "the " + helpVerb + " help screen: " +
					"usage synopsis, flags, and examples",
			}},
		}
	}
}

// deriveFileExamples seeds the registry from embedded help files that
// have neither a curated nor a summary entry. WhatItDoes is extracted
// from the file's title + first paragraph.
func deriveFileExamples(reg map[string]HelpExample) {
	entries, err := files.ReadDir(".")
	if err != nil {
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}

		key := strings.TrimSuffix(name, ".md")
		if _, exists := reg[key]; exists {
			continue
		}

		raw, err := files.ReadFile(name)
		if err != nil {
			continue
		}

		summary := extractFileSummary(string(raw))
		if len(summary) == 0 {
			continue
		}

		helpVerb := strings.ReplaceAll(key, "_", "-")
		reg[key] = HelpExample{
			WhatItDoes: summary,
			Examples: []Example{{
				Invocation:     "gitmap help " + helpVerb,
				ExpectedOutput: "the " + helpVerb + " help screen: usage synopsis, flags, and examples",
			}},
		}
	}
}

// extractFileSummary returns "<title>: <first paragraph>" from a help
// markdown file, capped for the WhatItDoes slot.
func extractFileSummary(md string) string {
	var title, firstPara string

	for _, line := range strings.Split(md, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(title) == 0 && strings.HasPrefix(trimmed, "# ") {
			title = strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
			continue
		}

		if len(title) > 0 && len(firstPara) == 0 && len(trimmed) > 0 &&
			!strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, ">") {
			firstPara = trimmed
		}

		if len(title) > 0 && len(firstPara) > 0 {
			break
		}
	}

	if len(title) == 0 {
		return ""
	}

	summary := title
	if len(firstPara) > 0 {
		summary += ": " + firstPara
	}

	if len(summary) > 220 {
		summary = summary[:217] + "..."
	}

	return summary
}
