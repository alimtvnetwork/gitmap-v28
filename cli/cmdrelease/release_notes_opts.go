// Package cmdrelease — release-notes flag parsing & git log execution.
//
// Supports:
//
//	--since <date|ref>   git log --since= window (e.g. "2 weeks ago", "2025-01-01")
//	--since-tag <tag>    shorthand for <tag>..HEAD
//	--format <fmt>       flat | grouped | markdown | json
//
// A bare positional <tagA>..<tagB> is still accepted for back-compat.
package cmdrelease

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
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

func handleReleaseNotesArgsError(err error) error {
	fmt.Fprintf(os.Stderr, "release-notes: ERROR %v\n", err)
	fmt.Fprintln(os.Stderr, "usage: gitmap release-notes [<tagA>..<tagB>] [--since <when>] [--since-tag <tag>] [--format flat|grouped|markdown|json]")
	cliexit.HandleUsageError(nil)

	return err
}

// runReleaseNotesV2 is the flag-aware entry point used by the dispatcher.
func runReleaseNotesV2(args []string) error {
	opts, err := parseReleaseNotesArgs(args)
	if err != nil {
		return handleReleaseNotesArgsError(err)
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
