// Package cmd — release-notes flag parsing & grouped formatting.
//
// Supports:
//
//	--since <date|ref>   git log --since= window (e.g. "2 weeks ago", "2025-01-01")
//	--since-tag <tag>    shorthand for <tag>..HEAD
//	--format <fmt>       flat | grouped | markdown | json
//
// A bare positional <tagA>..<tagB> is still accepted for back-compat.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// ReleaseNotesOpts holds parsed flags for release-notes.
type ReleaseNotesOpts struct {
	Range    string // "vA..vB" or "" when using Since
	Since    string // git --since= value
	SinceTag string // shorthand: <tag>..HEAD
	Format   string // flat | grouped | markdown | json
}

const (
	releaseNotesFormatFlat     = "flat"
	releaseNotesFormatGrouped  = "grouped"
	releaseNotesFormatMarkdown = "markdown"
	releaseNotesFormatJSON     = "json"
)

func applyReleaseNotesFlag(args []string, i int, opts *ReleaseNotesOpts) (int, bool) {
	if i+1 >= len(args) {
		return i, false
	}
	switch args[i] {
	case "--since":
		opts.Since = args[i+1]
		return i + 1, true
	case "--since-tag":
		opts.SinceTag = args[i+1]
		return i + 1, true
	case "--format":
		opts.Format = args[i+1]
		return i + 1, true
	}
	return i, false
}

func applyReleaseNotesArg(args []string, i int, opts *ReleaseNotesOpts) (int, error) {
	if nextI, matched := applyReleaseNotesFlag(args, i, opts); matched {
		return nextI, nil
	}
	if strings.Contains(args[i], "..") {
		opts.Range = args[i]
		return i, nil
	}
	return i, fmt.Errorf("unknown arg %q", args[i])
}

func validateReleaseNotesOpts(opts *ReleaseNotesOpts) error {
	if opts.SinceTag != "" && opts.Range == "" {
		opts.Range = opts.SinceTag + "..HEAD"
	}
	if opts.Range == "" && opts.Since == "" {
		return fmt.Errorf("need <tagA>..<tagB>, --since, or --since-tag")
	}
	return nil
}

// parseReleaseNotesArgs converts CLI args into ReleaseNotesOpts.
func parseReleaseNotesArgs(args []string) (ReleaseNotesOpts, error) {
	opts := ReleaseNotesOpts{Format: releaseNotesFormatMarkdown}
	for i := 0; i < len(args); i++ {
		nextI, err := applyReleaseNotesArg(args, i, &opts)
		if err != nil {
			return opts, err
		}
		i = nextI
	}
	return opts, validateReleaseNotesOpts(&opts)
}

func buildGitLogArgs(opts ReleaseNotesOpts) []string {
	args := []string{"log", "--pretty=format:%s|%h"}
	if opts.Since != "" {
		args = append(args, "--since="+opts.Since)
	}
	if opts.Range != "" {
		args = append(args, opts.Range)
	}
	return args
}

// gitLogForOpts runs git log honoring range + --since.
func gitLogForOpts(opts ReleaseNotesOpts) ([]string, error) {
	args := buildGitLogArgs(opts)
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git log: %w\n%s", err, out)
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return nil, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

// groupCommits buckets messages by conventional-commit prefix.
func groupCommits(lines []string) map[string][]string {
	groups := map[string][]string{}
	for _, ln := range lines {
		bucket := classifyCommit(ln)
		groups[bucket] = append(groups[bucket], ln)
	}
	return groups
}

var commitPrefixMap = []struct {
	prefixes []string
	category string
}{
	{[]string{"feat"}, "Features"},
	{[]string{"fix"}, "Fixes"},
	{[]string{"docs"}, "Docs"},
	{[]string{"refactor", "perf"}, "Refactor"},
	{[]string{"test"}, "Tests"},
	{[]string{"chore", "ci", "build"}, "Chore"},
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func classifyCommit(line string) string {
	lower := strings.ToLower(line)
	for _, entry := range commitPrefixMap {
		if hasAnyPrefix(lower, entry.prefixes) {
			return entry.category
		}
	}
	return "Other"
}

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

func handleReleaseNotesArgsError(err error) {
	fmt.Fprintf(os.Stderr, "release-notes: ERROR %v\n", err)
	fmt.Fprintln(os.Stderr, "usage: gitmap release-notes [<tagA>..<tagB>] [--since <when>] [--since-tag <tag>] [--format flat|grouped|markdown|json]")
	os.Exit(2)
}

// runReleaseNotesV2 is the flag-aware entry point used by the dispatcher.
func runReleaseNotesV2(args []string) error {
	opts, err := parseReleaseNotesArgs(args)
	if err != nil {
		handleReleaseNotesArgsError(err)
	}
	lines, err := gitLogForOpts(opts)
	if err != nil {
		return apperror.WrapSimple(err, "release-notes: ERROR")
	}
	if len(lines) == 0 {
		fmt.Fprintln(os.Stderr, "release-notes: no commits in selected range")
		return nil
	}
	fmt.Print(renderReleaseNotes(opts, lines))
	return nil
}
