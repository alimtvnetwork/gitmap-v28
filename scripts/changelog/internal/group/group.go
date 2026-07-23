// Package group classifies Commits into Conventional-Commit sections
// (Added / Fixed / Changed / Docs / etc.) in a deterministic order.
package group

import (
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/scripts/changelog/internal/gitlog"
)

// Section is a named bucket of commit subjects rendered together.
type Section struct {
	Name  string
	Items []string
}

// SectionOrder is the rendering order. Listed here once so the writer,
// renderer, and tests all agree.
var SectionOrder = []string{
	"Added",
	"Changed",
	"Fixed",
	"Performance",
	"Reverted",
	"Docs",
	"Build",
	"CI",
	"Refactor",
	"Tests",
	"Style",
	"Chore",
}

// prefixToSection maps Conventional-Commit prefixes to section names.
var prefixToSection = map[string]string{
	"feat":     "Added",
	"fix":      "Fixed",
	"docs":     "Docs",
	"chore":    "Chore",
	"refactor": "Refactor",
	"perf":     "Performance",
	"test":     "Tests",
	"build":    "Build",
	"ci":       "CI",
	"style":    "Style",
	"revert":   "Reverted",
}

// ByPrefix splits commits by Conventional-Commit prefix. Commits without
// a recognised prefix are returned in `skipped` so the caller can warn.
func ByPrefix(commits []gitlog.Commit) ([]Section, []gitlog.Commit) {
	buckets := map[string][]string{}

	var skipped []gitlog.Commit

	for _, c := range commits {
		section, item, ok := classify(c.Subject)
		if !ok {
			skipped = append(skipped, c)

			continue
		}

		buckets[section] = append(buckets[section], item)
	}

	return orderedSections(buckets), skipped
}

func classify(subject string) (string, string, bool) {
	idx := strings.Index(subject, ":")
	if idx <= 0 {
		return "", "", false
	}

	prefix := strings.TrimSpace(subject[:idx])
	prefix = strings.TrimSuffix(prefix, "!")
	if open := strings.Index(prefix, "("); open > 0 {
		prefix = prefix[:open]
	}

	section, ok := prefixToSection[strings.ToLower(prefix)]
	if !ok {
		return "", "", false
	}

	body := strings.TrimSpace(subject[idx+1:])

	return section, body, body != ""
}

// orderedSections renders sections in the fixed SectionOrder and sorts
// items inside each section lexicographically. Section order is the
// project's declared rendering priority (Added → Changed → Fixed → …);
// item order within a section is alphabetical so two regenerations on
// machines that received commits in different sequences still produce
// byte-identical output. Combined with gitlog's (author-date, hash)
// commit sort, this is the full deterministic ordering contract.
func orderedSections(buckets map[string][]string) []Section {
	out := make([]Section, 0, len(buckets))

	for _, name := range SectionOrder {
		items := buckets[name]
		if len(items) == 0 {
			continue
		}

		sorted := append([]string(nil), items...)
		sort.Strings(sorted)

		out = append(out, Section{Name: name, Items: sorted})
	}

	return out
}
