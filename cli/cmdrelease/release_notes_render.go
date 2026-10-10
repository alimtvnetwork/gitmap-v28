// Package cmdrelease — release notes rendering helpers (flat, grouped, json).
package cmdrelease

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// renderReleaseNotes turns parsed log lines into the chosen output format.
func renderReleaseNotes(opts ReleaseNotesOpts, lines []string) string {
	header := releaseNotesHeader(opts)
	switch opts.Format {
	case releaseNotesFormatFlat:
		return header + renderFlat(lines)
	case releaseNotesFormatJSON:
		return renderJSON(opts, lines)
	case releaseNotesFormatGrouped, releaseNotesFormatMarkdown:
		return header + renderGrouped(lines)
	default:
		return header + renderGrouped(lines)
	}
}

func releaseNotesHeader(opts ReleaseNotesOpts) string {
	scope := opts.Range
	if scope == "" {
		scope = "--since=" + opts.Since
	}

	return fmt.Sprintf("## Changes (%s)\n\n", scope)
}

func renderFlat(lines []string) string {
	var b strings.Builder
	for _, ln := range lines {
		b.WriteString("- " + formatLine(ln) + "\n")
	}

	return b.String()
}

func sortedGroupKeys(groups map[string][]string) []string {
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func renderGroupSection(b *strings.Builder, header string, items []string) {
	b.WriteString("### " + header + "\n")
	for _, ln := range items {
		b.WriteString("- " + formatLine(ln) + "\n")
	}

	b.WriteString("\n")
}

func renderGrouped(lines []string) string {
	groups := groupCommits(lines)
	var b strings.Builder
	for _, k := range sortedGroupKeys(groups) {
		renderGroupSection(&b, k, groups[k])
	}

	return b.String()
}

type releaseNotesJSONEntry struct {
	Group   string `json:"group"`
	Subject string `json:"subject"`
	SHA     string `json:"sha"`
}

type releaseNotesJSONOutput struct {
	Range   string                  `json:"range,omitempty"`
	Since   string                  `json:"since,omitempty"`
	Entries []releaseNotesJSONEntry `json:"entries"`
}

func buildJSONEntries(lines []string) []releaseNotesJSONEntry {
	entries := make([]releaseNotesJSONEntry, 0, len(lines))
	for _, ln := range lines {
		subj, sha := splitLine(ln)
		entries = append(entries, releaseNotesJSONEntry{Group: classifyCommit(ln), Subject: subj, SHA: sha})
	}

	return entries
}

func renderJSON(opts ReleaseNotesOpts, lines []string) string {
	out := releaseNotesJSONOutput{Range: opts.Range, Since: opts.Since, Entries: buildJSONEntries(lines)}
	buf, _ := json.MarshalIndent(out, "", "  ")

	return string(buf) + "\n"
}

func splitLine(ln string) (string, string) {
	if idx := strings.LastIndex(ln, "|"); idx >= 0 {
		return ln[:idx], ln[idx+1:]
	}

	return ln, ""
}

func formatLine(ln string) string {
	subj, sha := splitLine(ln)
	if sha == "" {
		return subj
	}

	return fmt.Sprintf("%s (%s)", subj, sha)
}
