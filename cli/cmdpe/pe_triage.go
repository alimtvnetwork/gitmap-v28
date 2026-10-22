package cmdpe

import (
	"flag"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// ghRun is the subset of `gh run list --json` fields triage needs.
type ghRun struct {
	DatabaseID int64  `json:"databaseId"`
	Conclusion string `json:"conclusion"`
	Name       string `json:"name"`
	HeadBranch string `json:"headBranch"`
	CreatedAt  string `json:"createdAt"`
	URL        string `json:"url"`
}

// ghJob is the subset of `gh run view --json jobs` fields triage needs.
type ghJob struct {
	Name       string   `json:"name"`
	Conclusion string   `json:"conclusion"`
	Steps      []ghStep `json:"steps"`
}

// ghStep is a single job step result.
type ghStep struct {
	Name       string `json:"name"`
	Conclusion string `json:"conclusion"`
}

type peTriageConfig struct {
	workdir string
	filter  string
	limit   int
}

type peFailure struct {
	repo triageRepo
	run  ghRun
}

const peTriageHelp = `pe triage — worst-failure bundle across workdir repos.

Description:
  Probes the latest GitHub Actions run for every repo in the workdir,
  picks the worst (most recent) failure, and prints a triage bundle:
  failing repo, run, job, step, branch, last commit age, and the
  suggested next command.

Usage:
  gitmap pe triage [--workdir <dir>] [--repo <name>] [--limit <n>]
  gitmap pet [--workdir <dir>] [--repo <name>]

Flags:
  --workdir <dir>   repo workdir (default: $GITMAP_WORKDIR or ~/workspace/repos)
  --repo, -t <name> only triage repos whose name contains <name>
  --limit <n>       probe at most n repos (default: 50)
  --help, -h        show this help

Example:
  gitmap pe triage
`

func isPeTriageHelp(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}

	return false
}

func setupPeTriageFlags(fs *flag.FlagSet) (*string, *string, *string, *int) {
	w := fs.String("workdir", "", "repo workdir")
	r := fs.String("repo", "", "repo name filter")
	s := fs.String("t", "", "repo name filter (shorthand)")
	l := fs.Int("limit", 50, "max repos to probe")

	return w, r, s, l
}

func resolveRepoFilter(long, short string) string {
	if long != "" {
		return long
	}

	return short
}

func parsePeTriageFlags(args []string) (peTriageConfig, *appfault.AppError) {
	fs := flag.NewFlagSet("pe triage", flag.ContinueOnError)
	w, r, s, l := setupPeTriageFlags(fs)
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(peTriageHelp)

		return peTriageConfig{}, apperror.WrapSimple(parseErr, "pe.triage.flags")
	}

	return peTriageConfig{
		workdir: resolveWorkdir(*w),
		filter:  resolveRepoFilter(*r, *s),
		limit:   *l,
	}, nil
}

func checkPeTriagePrereqs(args []string) (peTriageConfig, bool, *appfault.AppError) {
	if isPeTriageHelp(args) {
		fmt.Print(peTriageHelp)

		return peTriageConfig{}, true, nil
	}

	cfg, parseErr := parsePeTriageFlags(args)
	if parseErr != nil {
		return cfg, true, parseErr
	}

	if !hasGh() {
		return cfg, true, appfault.NewValidationError("pe triage requires the GitHub CLI (gh) on PATH")
	}

	return cfg, false, nil
}

// RunPeTriage implements gitmap pe triage (alias: pet).
func RunPeTriage(args []string) *appfault.AppError {
	cfg, isHandled, err := checkPeTriagePrereqs(args)
	if isHandled {
		return err
	}

	return executePeTriageFlow(cfg)
}

func executePeTriageFlow(cfg peTriageConfig) *appfault.AppError {
	repos, errRepos := listGitRepos(cfg.workdir)
	if errRepos != nil {
		return apperror.WrapSimple(errRepos, "pe.triage.listGitRepos")
	}

	if len(repos) == 0 {
		fmt.Printf("pe triage: no git repos found in %s\n", cfg.workdir)

		return nil
	}

	return processPeTriageFailures(repos, cfg)
}

func processPeTriageFailures(repos []triageRepo, cfg peTriageConfig) *appfault.AppError {
	failures, probed := collectPeFailures(repos, cfg)
	if len(failures) == 0 {
		fmt.Printf("pe triage: probed %d repos, no failing latest runs. all green.\n", probed)

		return nil
	}

	renderWorstFailureSummary(findMostRecentFailure(failures), len(failures))

	return nil
}

func matchRepoFilter(repoName, filter string) bool {
	hasFilter := len(filter) > 0
	if !hasFilter {
		return true
	}

	return strings.Contains(strings.ToLower(repoName), strings.ToLower(filter))
}

func probeLatestRepoRun(slug string) (ghRun, bool) {
	var runs []ghRun
	err := ghJSON(&runs, "run", "list", "-R", slug, "--limit", "1",
		"--json", "databaseId,conclusion,name,headBranch,createdAt,url")
	hasValidRun := err == nil && len(runs) > 0
	if !hasValidRun {
		return ghRun{}, false
	}

	isFailed := strings.EqualFold(runs[0].Conclusion, "failure")

	return runs[0], isFailed
}

func tryProbeRepo(r triageRepo, cfg peTriageConfig) (peFailure, bool, bool) {
	if !matchRepoFilter(r.name, cfg.filter) || r.slug == "" {
		return peFailure{}, false, false
	}

	run, isFailed := probeLatestRepoRun(r.slug)

	return peFailure{repo: r, run: run}, isFailed, true
}

func recordProbeResult(probed int, list []peFailure, item peFailure, isFailed, isProbed bool) (int, []peFailure) {
	if isProbed {
		probed++
	}
	if isFailed {
		list = append(list, item)
	}

	return probed, list
}

func collectPeFailures(repos []triageRepo, cfg peTriageConfig) ([]peFailure, int) {
	var failures []peFailure
	probed := 0
	for _, r := range repos {
		if probed >= cfg.limit {
			break
		}
		failItem, isFailed, isProbed := tryProbeRepo(r, cfg)
		probed, failures = recordProbeResult(probed, failures, failItem, isFailed, isProbed)
	}

	return failures, probed
}

func findMostRecentFailure(failures []peFailure) peFailure {
	worst := failures[0]
	for _, f := range failures[1:] {
		if f.run.CreatedAt > worst.run.CreatedAt {
			worst = f
		}
	}

	return worst
}

func resolveHeadBranch(branch string) string {
	if branch == "" {
		return "unknown"
	}

	return branch
}

func printFailureDetails(w peFailure, job, step, branch, age string) {
	fmt.Println("worst failure:")
	fmt.Printf("  repo:        %s (%s)\n", w.repo.name, w.repo.slug)
	fmt.Printf("  run:         %d %s (%s)\n", w.run.DatabaseID, w.run.Name, w.run.Conclusion)
	fmt.Printf("  job:         %s\n", job)
	fmt.Printf("  step:        %s\n", step)
	fmt.Printf("  branch:      %s\n", branch)
	fmt.Printf("  last commit: %s ago\n", age)
	fmt.Printf("  run url:     %s\n", w.run.URL)
}

func printAdditionalFailures(total int) {
	if total > 1 {
		fmt.Printf("  also red:    %d more repo(s)\n", total-1)
	}
}

func printSuggestedNext(runID int64, repoName string) {
	fmt.Println("suggested next:")
	fmt.Printf("  gitmap pe rerun %d --dry-run\n", runID)
	fmt.Printf("  gitmap ci open %s\n", repoName)
}

func renderWorstFailureSummary(worst peFailure, totalFailures int) {
	jobName, stepName := worstFailureDetail(worst.repo.slug, worst.run.DatabaseID)
	branch := resolveHeadBranch(worst.run.HeadBranch)
	age := gitLastCommitAge(worst.repo.path)

	printFailureDetails(worst, jobName, stepName, branch, age)
	printAdditionalFailures(totalFailures)
	printSuggestedNext(worst.run.DatabaseID, worst.repo.name)
}

// worstFailureDetail returns the first failed job and its first failed step.
func worstFailureDetail(slug string, runID int64) (string, string) {
	var view struct {
		Jobs []ghJob `json:"jobs"`
	}
	id := fmt.Sprintf("%d", runID)
	if viewErr := ghJSON(&view, "run", "view", id, "-R", slug, "--json", "jobs"); viewErr != nil {
		return "unknown", "unknown"
	}

	return findFailedJobAndStep(view.Jobs)
}

func findFailedJobAndStep(jobs []ghJob) (string, string) {
	for _, j := range jobs {
		if strings.EqualFold(j.Conclusion, "failure") {
			return j.Name, findFailingStepName(j.Steps)
		}
	}

	return "unknown", "unknown"
}

func findFailingStepName(steps []ghStep) string {
	for _, s := range steps {
		if strings.EqualFold(s.Conclusion, "failure") {
			return s.Name
		}
	}

	return "unknown"
}
