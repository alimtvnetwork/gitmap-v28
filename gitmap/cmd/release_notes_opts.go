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

// parseReleaseNotesArgs converts CLI args into ReleaseNotesOpts.
func parseReleaseNotesArgs(args []string) (ReleaseNotesOpts, error) {
	opts := ReleaseNotesOpts{Format: releaseNotesFormatMarkdown}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--since" && i+1 < len(args):
			opts.Since = args[i+1]
			i++
		case a == "--since-tag" && i+1 < len(args):
			opts.SinceTag = args[i+1]
			i++
		case a == "--format" && i+1 < len(args):
			opts.Format = args[i+1]
			i++
		case strings.Contains(a, ".."):
			opts.Range = a
		default:
			return opts, fmt.Errorf("unknown arg %q", a)
		}
	}
	if opts.SinceTag != "" && opts.Range == "" {
		opts.Range = opts.SinceTag + "..HEAD"
	}
	if opts.Range == "" && opts.Since == "" {
		return opts, fmt.Errorf("need <tagA>..<tagB>, --since, or --since-tag")
	}
	return opts, nil
}

// gitLogForOpts runs git log honoring range + --since.
func gitLogForOpts(opts ReleaseNotesOpts) ([]string, error) {
	args := []string{"log", "--pretty=format:%s|%h"}
	if opts.Since != "" {
		args = append(args, "--since="+opts.Since)
	}
	if opts.Range != "" {
		args = append(args, opts.Range)
	}
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git log: %v\n%s", err, out)
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

func classifyCommit(line string) string {
	lower := strings.ToLower(line)
	switch {
	case strings.HasPrefix(lower, "feat"):
		return "Features"
	case strings.HasPrefix(lower, "fix"):
		return "Fixes"
	case strings.HasPrefix(lower, "docs"):
		return "Docs"
	case strings.HasPrefix(lower, "refactor"), strings.HasPrefix(lower, "perf"):
		return "Refactor"
	case strings.HasPrefix(lower, "test"):
		return "Tests"
	case strings.HasPrefix(lower, "chore"), strings.HasPrefix(lower, "ci"), strings.HasPrefix(lower, "build"):
		return "Chore"
	default:
		return "Other"
	}
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

func renderGrouped(lines []string) string {
	groups := groupCommits(lines)
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString("### " + k + "\n")
		for _, ln := range groups[k] {
			b.WriteString("- " + formatLine(ln) + "\n")
		}
		b.WriteString("\n")
	}
	return b.String()
}

func renderJSON(opts ReleaseNotesOpts, lines []string) string {
	type entry struct {
		Group   string `json:"group"`
		Subject string `json:"subject"`
		SHA     string `json:"sha"`
	}
	out := struct {
		Range   string  `json:"range,omitempty"`
		Since   string  `json:"since,omitempty"`
		Entries []entry `json:"entries"`
	}{Range: opts.Range, Since: opts.Since}
	for _, ln := range lines {
		subj, sha := splitLine(ln)
		out.Entries = append(out.Entries, entry{Group: classifyCommit(ln), Subject: subj, SHA: sha})
	}
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

// runReleaseNotesV2 is the flag-aware entry point used by the dispatcher.
func runReleaseNotesV2(args []string) {
	opts, err := parseReleaseNotesArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "release-notes: ERROR %v\n", err)
		fmt.Fprintln(os.Stderr, "usage: gitmap release-notes [<tagA>..<tagB>] [--since <when>] [--since-tag <tag>] [--format flat|grouped|markdown|json]")
		os.Exit(2)
	}
	lines, err := gitLogForOpts(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "release-notes: ERROR %v\n", err)
		os.Exit(1)
	}
	if len(lines) == 0 {
		fmt.Fprintln(os.Stderr, "release-notes: no commits in selected range")
		return
	}
	fmt.Print(renderReleaseNotes(opts, lines))
}
