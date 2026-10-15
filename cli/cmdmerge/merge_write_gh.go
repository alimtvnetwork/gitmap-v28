package cmdmerge

// merge_write_gh.go — GitHub host call for `gitmap merge write` / `mr`
// (program 273 WS3, spec 05). v1 supports github.com remotes only; auth and
// the PR creation go through the `gh` CLI, the same pattern as the other
// gh-backed commands. Host errors surface verbatim — never swallowed.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// merge-write host messages.
const (
	mwErrGhMissing = "the `gh` CLI is required for `merge write` but was not found on PATH.\nFix: install gh (https://cli.github.com), then `gh auth login`, or set GITHUB_TOKEN / GH_TOKEN."
	mwErrNoOrigin  = "no git remote 'origin' found — run `merge write` inside a git repository with an origin remote"
	mwErrHostFmt   = "unsupported remote %q: `merge write` supports github.com remotes only in v1"
	mwErrBranch    = "cannot determine the current branch (detached HEAD?) — pass --source explicitly"
	mwErrBaseFmt   = "could not read the default branch for %s: %s"
	mwErrCreateFmt = "gh pr create failed: %s"
)

// mergeWritePRJSON decodes `gh pr create --json number,url`.
type mergeWritePRJSON struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
}

// requireMergeWriteGh ensures the `gh` CLI is available.
func requireMergeWriteGh() *appfault.AppError {
	if _, err := exec.LookPath("gh"); err != nil {
		return appfault.NewNotFoundError(mwErrGhMissing)
	}
	return nil
}

// resolveGitHubSlug returns "owner/repo" for the origin remote.
// GitHub only in v1 — any other host is a tracked validation error.
func resolveGitHubSlug() (string, *appfault.AppError) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	remote := strings.TrimSpace(string(out))
	if err != nil || remote == "" {
		return "", appfault.NewNotFoundError(mwErrNoOrigin)
	}
	owner, repo, isGitHub := parseGitHubSlug(remote)
	if !isGitHub {
		return "", appfault.NewValidationError(fmt.Sprintf(mwErrHostFmt, remote))
	}
	return owner + "/" + repo, nil
}

// parseGitHubSlug extracts owner/repo from common GitHub remote URL forms.
// isGitHub=false for any other host (never a silent mismatch).
func parseGitHubSlug(remote string) (owner, repo string, isGitHub bool) {
	rest := strings.TrimSuffix(remote, ".git")
	rest = strings.Replace(rest, "git@github.com:", "https://github.com/", 1)
	rest = strings.Replace(rest, "ssh://git@github.com/", "https://github.com/", 1)
	if !strings.HasPrefix(rest, "https://github.com/") && !strings.HasPrefix(rest, "http://github.com/") {
		return "", "", false
	}
	path := strings.TrimPrefix(strings.TrimPrefix(rest, "https://github.com/"), "http://github.com/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || strings.Contains(parts[1], "/") {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// resolveMergeWriteBranches fills --source/--target defaults.
func resolveMergeWriteBranches(opts mergeWriteOpts, slug string) (mergeWriteOpts, *appfault.AppError) {
	withSource, appErr := defaultMergeWriteSource(opts)
	if appErr != nil {
		return opts, appErr
	}
	return defaultMergeWriteTarget(withSource, slug)
}

// defaultMergeWriteSource fills --source from the current branch.
func defaultMergeWriteSource(opts mergeWriteOpts) (mergeWriteOpts, *appfault.AppError) {
	if opts.source != "" {
		return opts, nil
	}
	source, appErr := currentGitBranch()
	if appErr != nil {
		return opts, appErr
	}
	opts.source = source
	return opts, nil
}

// defaultMergeWriteTarget fills --target from the repo default branch.
func defaultMergeWriteTarget(opts mergeWriteOpts, slug string) (mergeWriteOpts, *appfault.AppError) {
	if opts.target != "" {
		return opts, nil
	}
	target, appErr := defaultGitHubBranch(slug)
	if appErr != nil {
		return opts, appErr
	}
	opts.target = target
	return opts, nil
}

// currentGitBranch returns the checked-out branch name.
func currentGitBranch() (string, *appfault.AppError) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	branch := strings.TrimSpace(string(out))
	if err != nil || branch == "" || branch == "HEAD" {
		return "", appfault.NewValidationError(mwErrBranch)
	}
	return branch, nil
}

// defaultGitHubBranch reads the repo's default branch via gh.
func defaultGitHubBranch(slug string) (string, *appfault.AppError) {
	out, err := exec.Command("gh", "repo", "view", slug, "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name").Output()
	branch := strings.TrimSpace(string(out))
	if err != nil || branch == "" {
		return "", apperror.NewExecutionError(fmt.Sprintf(mwErrBaseFmt, slug, strings.TrimSpace(string(out))))
	}
	return branch, nil
}

// openGitHubPullRequest shells out to `gh pr create` and parses the result.
func openGitHubPullRequest(slug string, opts mergeWriteOpts) (mergeWritePRJSON, *appfault.AppError) {
	result := mergeWritePRJSON{}
	cmd := exec.Command("gh", buildGhPrCreateArgs(slug, opts)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return result, apperror.NewExecutionError(fmt.Sprintf(mwErrCreateFmt, ghFailureDetail(stderr.String(), err)))
	}
	return parseMergeWriteResult(out)
}

// buildGhPrCreateArgs assembles the `gh pr create` argument list.
func buildGhPrCreateArgs(slug string, opts mergeWriteOpts) []string {
	args := []string{"pr", "create", "--repo", slug,
		"--head", opts.source, "--base", opts.target,
		"--title", opts.title, "--json", "number,url"}
	if opts.body != "" {
		args = append(args, "--body", opts.body)
	}
	if opts.isDraft {
		args = append(args, "--draft")
	}
	return args
}

// ghFailureDetail prefers gh's stderr; falls back to the exec error.
func ghFailureDetail(stderr string, err error) string {
	if detail := strings.TrimSpace(stderr); detail != "" {
		return detail
	}
	return err.Error()
}

// parseMergeWriteResult decodes `gh pr create --json number,url`
// (object form; tolerates the array form).
func parseMergeWriteResult(out []byte) (mergeWritePRJSON, *appfault.AppError) {
	result := mergeWritePRJSON{}
	trimmed := bytes.TrimSpace(out)
	if bytes.HasPrefix(trimmed, []byte("[")) {
		var arr []mergeWritePRJSON
		if err := json.Unmarshal(trimmed, &arr); err != nil || len(arr) == 0 {
			return result, apperror.NewExecutionError("could not parse `gh pr create` output: " + string(trimmed))
		}
		return arr[0], nil
	}
	if err := json.Unmarshal(trimmed, &result); err != nil || result.Number == 0 || result.URL == "" {
		return result, apperror.NewExecutionError("could not parse `gh pr create` output: " + string(trimmed))
	}
	return result, nil
}

// printMergeWriteResult prints the success lines (spec 05 §2).
func printMergeWriteResult(result mergeWritePRJSON, opts mergeWriteOpts) {
	draft := ""
	if opts.isDraft {
		draft = " (draft)"
	}
	fmt.Printf(mwResultFmt, result.Number, opts.source, opts.target, draft, result.URL)
}
