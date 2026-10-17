// Package cmdrelease — release_train.go: `gitmap release train <repo...>` (alias `rtrain`).
//
// Ordered multi-repo release: for each repo — version bump (next tag from
// the latest release) → CI gate check → tag → release notes. DEFAULT IS
// --dry-run (prints the plan); real execution only with --yes.
//
// Wiring (lead): same pattern as RunReleaseVerify — add `train`/`rtrain` to
// RunRelease's sub routing with an explicit nil check.
package cmdrelease

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

// releaseTrainOpts carries the parsed flags for `release train`.
type releaseTrainOpts struct {
	order string
	isYes bool
}

// releaseTrainPlan is the resolved 4-step plan for one repo.
type releaseTrainPlan struct {
	repo          string
	latestTag     string
	nextTag       string
	defaultBranch string
	ciStatus      string
	ciConclusion  string
	commitsSince  int
	isGateGreen   bool
}

// RunReleaseTrain executes `gitmap release train <repo...>`.
// args is the tail after the `train` token.
func RunReleaseTrain(args []string) *appfault.AppError {
	if isReleaseNewCmdHelpRequested(args) {
		fmt.Print(releaseTrainHelpText)
		cliexit.Exit(0)
	}

	opts, appErr := parseReleaseTrainFlags(args)
	if appErr != nil {
		return appErr
	}
	repos := resolveTrainRepos(opts, args)
	if len(repos) == 0 {
		return appfault.NewValidationError("missing required argument: <repo...>\n" + releaseTrainUsageLine +
			"example: gitmap release train agm gitmap-v28")
	}

	if appErr := requireGhCLI(); appErr != nil {
		return appErr
	}

	plans := make([]releaseTrainPlan, 0, len(repos))
	for _, repo := range repos {
		plan, appErr := buildReleaseTrainPlan(repo)
		if appErr != nil {
			return appErr
		}
		plans = append(plans, plan)
	}

	printReleaseTrainPlan(plans, opts.isYes)
	if !opts.isYes {
		fmt.Println("dry-run: no changes made — pass --yes to execute")

		return nil
	}

	return executeReleaseTrain(plans)
}

// parseReleaseTrainFlags parses the `release train` flag set.
func parseReleaseTrainFlags(args []string) (releaseTrainOpts, *appfault.AppError) {
	opts := releaseTrainOpts{}
	fs := flag.NewFlagSet("release train", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.order, "order", "", "comma-separated repo order (overrides positional order)")
	fs.BoolVar(&opts.isYes, "yes", false, "execute the release train (default is dry-run)")
	fs.BoolVar(&opts.isYes, "y", false, "execute the release train (shorthand)")
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return opts, appfault.NewValidationError(releaseTrainUsageLine + "error: " + err.Error())
	}

	return opts, nil
}

// resolveTrainRepos normalizes the repo list; --order wins over positionals.
func resolveTrainRepos(opts releaseTrainOpts, args []string) []string {
	raw := args
	if strings.TrimSpace(opts.order) != "" {
		raw = strings.Split(opts.order, ",")
	}
	repos := make([]string, 0, len(raw))
	for _, item := range raw {
		if trimmed := strings.TrimSpace(item); trimmed != "" && !strings.HasPrefix(trimmed, "-") {
			repos = append(repos, normalizeReleaseRepo(trimmed))
		}
	}

	return repos
}

// buildReleaseTrainPlan resolves the 4-step plan for one repo.
func buildReleaseTrainPlan(repo string) (releaseTrainPlan, *appfault.AppError) {
	plan := releaseTrainPlan{repo: repo}

	latest, appErr := latestReleaseTag(repo)
	if appErr != nil {
		return plan, appErr
	}
	plan.latestTag = latest

	next, appErr := bumpPatchVersion(latest)
	if appErr != nil {
		return plan, appErr
	}
	plan.nextTag = next

	branch, appErr := runGhCLI("api", "repos/"+repo, "--jq", ".default_branch")
	if appErr != nil {
		return plan, appErr
	}
	plan.defaultBranch = branch

	plan.ciStatus, plan.ciConclusion = latestCIRun(repo, branch)
	plan.isGateGreen = plan.ciConclusion == "success"
	plan.commitsSince = commitsSinceTag(repo, latest, branch)

	return plan, nil
}

// latestReleaseTag returns the newest release tag for the repo.
func latestReleaseTag(repo string) (string, *appfault.AppError) {
	out, appErr := runGhCLI("release", "list", "--repo", repo, "--limit", "1", "--json", "tagName")
	if appErr != nil {
		return "", appErr
	}

	var releases []struct {
		TagName string `json:"tagName"`
	}
	if err := json.Unmarshal([]byte(out), &releases); err != nil || len(releases) == 0 {
		return "", appfault.NewValidationError("no releases found for " + repo)
	}

	return releases[0].TagName, nil
}

// bumpPatchVersion computes the next patch version from a vX.Y.Z tag.
func bumpPatchVersion(tag string) (string, *appfault.AppError) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(tag), "v")
	parts := strings.Split(trimmed, ".")
	if len(parts) != 3 {
		return "", appfault.NewValidationError(
			"cannot bump version from non-semver tag " + strconv.Quote(tag))
	}

	nums := make([]int, 3)
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return "", appfault.NewValidationError(
				"cannot bump version from non-semver tag " + strconv.Quote(tag))
		}
		nums[i] = num
	}

	return fmt.Sprintf("v%d.%d.%d", nums[0], nums[1], nums[2]+1), nil
}

// latestCIRun returns the latest workflow run's status/conclusion on the branch.
func latestCIRun(repo, branch string) (string, string) {
	out, appErr := runGhCLI("run", "list", "--repo", repo, "--branch", branch,
		"--limit", "1", "--json", "status,conclusion")
	if appErr != nil {
		return "unknown", "unknown"
	}

	var runs []struct {
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
	}
	if err := json.Unmarshal([]byte(out), &runs); err != nil || len(runs) == 0 {
		return "unknown", "unknown"
	}

	return runs[0].Status, runs[0].Conclusion
}

// commitsSinceTag counts commits on the branch since the tag (best effort).
func commitsSinceTag(repo, tag, branch string) int {
	out, appErr := runGhCLI("api", "repos/"+repo+"/compare/"+tag+"..."+branch, "--jq", ".total_commits")
	if appErr != nil {
		return -1
	}

	count, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return -1
	}

	return count
}

// printReleaseTrainPlan renders the ordered dry-run (or pre-execution) plan.
func printReleaseTrainPlan(plans []releaseTrainPlan, isExecute bool) {
	mode := "DRY-RUN"
	if isExecute {
		mode = "EXECUTE"
	}
	fmt.Printf("RELEASE TRAIN [%s] — %d repos:\n", mode, len(plans))
	for i, plan := range plans {
		gate := "red ✗"
		if plan.isGateGreen {
			gate = "green ✓"
		}
		notes := "unknown commits"
		if plan.commitsSince >= 0 {
			notes = fmt.Sprintf("%d commits since %s", plan.commitsSince, plan.latestTag)
		}
		fmt.Printf("  [%d/%d] %s\n", i+1, len(plans), plan.repo)
		fmt.Printf("    bump: %s → %s · CI gate: %s (%s/%s) · tag: %s · notes: %s\n",
			plan.latestTag, plan.nextTag, gate, plan.ciStatus, plan.ciConclusion,
			plan.nextTag, notes)
	}
}

// executeReleaseTrain runs the train for real: re-check the CI gate per repo
// (abort on red), then create the tag + release with generated notes.
func executeReleaseTrain(plans []releaseTrainPlan) *appfault.AppError {
	for i, plan := range plans {
		status, conclusion := latestCIRun(plan.repo, plan.defaultBranch)
		if conclusion != "success" {
			return appfault.NewValidationError(
				fmt.Sprintf("release train aborted at [%d/%d] %s: CI gate is %s/%s (not success)",
					i+1, len(plans), plan.repo, status, conclusion))
		}

		fmt.Printf("[%d/%d] %s: creating %s ...\n", i+1, len(plans), plan.repo, plan.nextTag)
		if _, appErr := runGhCLI("release", "create", plan.nextTag,
			"--repo", plan.repo,
			"--target", plan.defaultBranch,
			"--title", plan.nextTag,
			"--generate-notes"); appErr != nil {
			return appErr
		}
		fmt.Printf("[%d/%d] %s: %s released ✓\n", i+1, len(plans), plan.repo, plan.nextTag)
	}

	return nil
}
