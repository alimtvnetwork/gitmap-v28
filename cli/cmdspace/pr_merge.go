// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunPRMerge implements `gitmap space PR-merge [target]` (spec 271):
// it merges the current branch into the target branch — the default branch
// when no target is given — with a regular merge commit (all commits kept).
// --squash / --single-commit / --flat select squash variants. The working
// tree must be clean; conflicts abort the merge instead of resolving it.
func RunPRMerge(args []string) error {
	if hasHelpFlag(args) {
		printPRMergeHelp()
		return nil
	}

	opts, err := parsePRMergeArgs(args)
	if err != nil {
		return err
	}

	source, err := runGitCapture("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return err
	}

	target := opts.target
	if target == "" {
		target, err = resolveDefaultBranch()
		if err != nil {
			return err
		}
	}

	if source == target {
		return apperror.NewValidationError(fmt.Sprintf("cannot merge %s into itself", source))
	}

	if cleanErr := requireCleanTree(); cleanErr != nil {
		return cleanErr
	}

	if fetchErr := fetchOrigin(); fetchErr != nil {
		return fetchErr
	}

	if coErr := runGitInherit("checkout", target); coErr != nil {
		return coErr
	}

	if hasOriginRemote() {
		if ffErr := runGitInherit("merge", "--ff-only", "origin/"+target); ffErr != nil {
			return apperror.NewExecutionError(fmt.Sprintf("local %s has diverged from origin/%s - sync it manually first", target, target))
		}
	}

	heatmap := buildChangeHeatmap(target, source)

	var mergeErr error
	if opts.squash || opts.singleCommit {
		mergeErr = runGitInherit("merge", "--squash", source)
	} else {
		mergeErr = runGitInherit("merge", "--no-ff", "--no-edit", source)
	}

	if mergeErr != nil {
		abortMerge()
		printConflictedFiles()
		return apperror.NewExecutionError("merge conflicts - aborted, resolve manually")
	}

	if opts.squash || opts.singleCommit {
		msg := squashMessage(source, target, opts.singleCommit, heatmap)
		if commitErr := runGitInherit("commit", "-m", msg); commitErr != nil {
			return commitErr
		}
	}

	sha, shaErr := runGitCapture("rev-parse", "--short", "HEAD")
	if shaErr != nil {
		return shaErr
	}

	fmt.Printf("merged %s into %s @ %s\n", source, target, sha)

	if !opts.noPush {
		if pushErr := runGitInherit("push", "origin", target); pushErr != nil {
			return pushErr
		}
	}

	return nil
}

// prMergeOptions holds the parsed flags for `gitmap space PR-merge`.
type prMergeOptions struct {
	target       string
	squash       bool
	singleCommit bool
	noPush       bool
}

// parsePRMergeArgs parses `gitmap space PR-merge [target] [flags]`.
// --force is rejected: merges never rewrite history, so the flag would lie.
func parsePRMergeArgs(args []string) (prMergeOptions, error) {
	var opts prMergeOptions

	for _, a := range args {
		switch a {
		case "--squash":
			opts.squash = true
		case "--single-commit", "--flat":
			opts.singleCommit = true
		case "--no-push":
			opts.noPush = true
		case "--force", "-f":
			return opts, apperror.NewValidationError("--force is not supported for PR-merge: merges never rewrite history")
		case "--help", "-h":
			// handled by the caller
		default:
			if strings.HasPrefix(a, "-") {
				return opts, apperror.NewValidationError(fmt.Sprintf("unknown flag: %s", a))
			}

			if opts.target != "" {
				return opts, apperror.NewValidationError(fmt.Sprintf("unexpected argument: %s", a))
			}

			opts.target = a
		}
	}

	return opts, nil
}

// resolveDefaultBranch returns the repo's default branch name: the
// origin/HEAD symbolic ref first, then main, then master.
func resolveDefaultBranch() (string, error) {
	ref, err := runGitCapture("symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil && ref != "" {
		name := strings.TrimPrefix(ref, "refs/remotes/origin/")
		if name != "" && name != ref {
			return name, nil
		}
	}

	for _, candidate := range []string{"main", "master"} {
		if err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+candidate).Run(); err == nil {
			return candidate, nil
		}
	}

	return "", apperror.NewExecutionError("could not determine the default branch")
}

// heatmapEntry is one row of the changed-files heatmap.
type heatmapEntry struct {
	path    string
	added   int
	deleted int
}

// buildChangeHeatmap returns the top 12 changed files between the merge base
// of target/source and source, sorted by churn (added+deleted), rendered for
// squash messages.
func buildChangeHeatmap(target, source string) []heatmapEntry {
	base, err := runGitCapture("merge-base", target, source)
	if err != nil || base == "" {
		return nil
	}

	out, err := runGitCapture("diff", "--numstat", base+".."+source)
	if err != nil || out == "" {
		return nil
	}

	var entries []heatmapEntry

	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		added, aErr := strconv.Atoi(fields[0])
		deleted, dErr := strconv.Atoi(fields[1])

		if aErr != nil || dErr != nil {
			continue // binary files show "-" instead of counts
		}

		entries = append(entries, heatmapEntry{path: fields[2], added: added, deleted: deleted})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].added+entries[i].deleted > entries[j].added+entries[j].deleted
	})

	if len(entries) > 12 {
		entries = entries[:12]
	}

	return entries
}

// prDetails holds a GitHub PR's title/body for squash messages.
type prDetails struct {
	title  string
	body   string
	number int
}

// lookupPRDetails returns the open PR for the branch, or nil when gh is
// missing or no PR exists. Never errors — callers fall back gracefully.
func lookupPRDetails(branch string) *prDetails {
	out, err := exec.Command("gh", "pr", "view", branch, "--json", "title,body,number").Output()
	if err != nil {
		return nil
	}

	text := string(out)
	title := extractJSONField(text, "title")
	body := extractJSONField(text, "body")
	number := extractJSONNumber(text, "number")

	if title == "" {
		return nil
	}

	return &prDetails{title: title, body: strings.TrimSpace(body), number: number}
}

// extractJSONField pulls a string field value out of flat gh JSON output.
func extractJSONField(json, field string) string {
	key := `"` + field + `":"`

	idx := strings.Index(json, key)
	if idx < 0 {
		return ""
	}

	rest := json[idx+len(key):]
	var sb strings.Builder

	escaped := false

	for _, r := range rest {
		if escaped {
			sb.WriteRune(r)
			escaped = false

			continue
		}

		if r == '\\' {
			escaped = true

			continue
		}

		if r == '"' {
			break
		}

		sb.WriteRune(r)
	}

	return sb.String()
}

// extractJSONNumber pulls an integer field value out of flat gh JSON output.
func extractJSONNumber(json, field string) int {
	key := `"` + field + `":`

	idx := strings.Index(json, key)
	if idx < 0 {
		return 0
	}

	rest := json[idx+len(key):]
	end := strings.IndexAny(rest, ",}")

	if end < 0 {
		return 0
	}

	n, err := strconv.Atoi(strings.TrimSpace(rest[:end]))
	if err != nil {
		return 0
	}

	return n
}

// squashMessage composes the squash commit message: PR title/body when a PR
// exists (single-commit mode uses a one-liner), plus the file heatmap.
func squashMessage(source, target string, singleCommit bool, heatmap []heatmapEntry) string {
	var sb strings.Builder

	if singleCommit {
		sb.WriteString(source + " -> " + target + "\n")
	} else if pr := lookupPRDetails(source); pr != nil {
		sb.WriteString(pr.title)
		if pr.number > 0 {
			sb.WriteString(fmt.Sprintf(" (#%d)", pr.number))
		}

		sb.WriteString("\n")

		if pr.body != "" {
			sb.WriteString("\n" + pr.body + "\n")
		}
	} else {
		sb.WriteString(fmt.Sprintf("Merge %s into %s (squash)\n", source, target))
	}

	if len(heatmap) > 0 {
		sb.WriteString("\n---\nChanged files:\n")

		for _, h := range heatmap {
			sb.WriteString(fmt.Sprintf("  %s  +%d -%d\n", h.path, h.added, h.deleted))
		}
	}

	return sb.String()
}

// hasOriginRemote reports whether the repo has an origin remote.
func hasOriginRemote() bool {
	out, err := runGitCapture("remote")
	if err != nil {
		return false
	}

	for _, r := range strings.Split(out, "\n") {
		if strings.TrimSpace(r) == "origin" {
			return true
		}
	}

	return false
}

// fetchOrigin fetches origin when the remote exists; without an origin
// remote there is nothing to fetch and the merge proceeds locally.
func fetchOrigin() error {
	if !hasOriginRemote() {
		fmt.Println("no origin remote - merging locally")
		return nil
	}

	return runGitInherit("fetch", "origin")
}

// abortMerge rolls back a conflicted merge.
func abortMerge() {
	_ = exec.Command("git", "merge", "--abort").Run()
}

// printConflictedFiles lists files with merge conflicts for the user.
func printConflictedFiles() {
	out, err := runGitCapture("diff", "--name-only", "--diff-filter=U")
	if err != nil || out == "" {
		return
	}

	fmt.Println("Conflicted files:")

	for _, f := range strings.Split(out, "\n") {
		if f != "" {
			fmt.Println("  " + f)
		}
	}
}
