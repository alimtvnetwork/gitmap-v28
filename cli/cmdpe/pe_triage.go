package cmdpe

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"

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

// RunPeTriage implements `gitmap pe triage` (alias: pet).
func RunPeTriage(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(peTriageHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("pe triage", flag.ContinueOnError)
	workdirFlag := fs.String("workdir", "", "repo workdir")
	repoFilter := fs.String("repo", "", "repo name filter")
	repoFilterShort := fs.String("t", "", "repo name filter (shorthand)")
	limit := fs.Int("limit", 50, "max repos to probe")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(peTriageHelp)
		return apperror.WrapSimple(parseErr, "pe.triage.flags")
	}
	if !hasGh() {
		return appfault.NewValidationError("pe triage requires the GitHub CLI (gh) on PATH")
	}
	filter := *repoFilter
	if filter == "" {
		filter = *repoFilterShort
	}
	workdir := resolveWorkdir(*workdirFlag)
	repos := listGitRepos(workdir)
	if len(repos) == 0 {
		fmt.Printf("pe triage: no git repos found in %s\n", workdir)
		return nil
	}

	type failure struct {
		repo triageRepo
		run  ghRun
	}
	var failures []failure
	probed := 0
	for _, r := range repos {
		if filter != "" && !strings.Contains(strings.ToLower(r.name), strings.ToLower(filter)) {
			continue
		}
		if r.slug == "" {
			continue
		}
		if probed >= *limit {
			break
		}
		probed++
		var runs []ghRun
		if listErr := ghJSON(&runs, "run", "list", "-R", r.slug, "--limit", "1",
			"--json", "databaseId,conclusion,name,headBranch,createdAt,url"); listErr != nil {
			continue
		}
		if len(runs) == 0 {
			continue
		}
		if strings.EqualFold(runs[0].Conclusion, "failure") {
			failures = append(failures, failure{repo: r, run: runs[0]})
		}
	}

	if len(failures) == 0 {
		fmt.Printf("pe triage: probed %d repos, no failing latest runs. all green.\n", probed)
		return nil
	}

	worst := failures[0]
	for _, f := range failures[1:] {
		if f.run.CreatedAt > worst.run.CreatedAt {
			worst = f
		}
	}

	jobName, stepName := worstFailureDetail(worst.repo.slug, worst.run.DatabaseID)
	branch := worst.run.HeadBranch
	if branch == "" {
		branch = "unknown"
	}
	age := gitLastCommitAge(worst.repo.path)

	fmt.Println("worst failure:")
	fmt.Printf("  repo:        %s (%s)\n", worst.repo.name, worst.repo.slug)
	fmt.Printf("  run:         %d %s (%s)\n", worst.run.DatabaseID, worst.run.Name, worst.run.Conclusion)
	fmt.Printf("  job:         %s\n", jobName)
	fmt.Printf("  step:        %s\n", stepName)
	fmt.Printf("  branch:      %s\n", branch)
	fmt.Printf("  last commit: %s ago\n", age)
	fmt.Printf("  run url:     %s\n", worst.run.URL)
	if len(failures) > 1 {
		fmt.Printf("  also red:    %d more repo(s)\n", len(failures)-1)
	}
	fmt.Println("suggested next:")
	fmt.Printf("  gitmap pe rerun %d --dry-run\n", worst.run.DatabaseID)
	fmt.Printf("  gitmap ci open %s\n", worst.repo.name)

	return nil
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
	for _, j := range view.Jobs {
		if !strings.EqualFold(j.Conclusion, "failure") {
			continue
		}
		for _, s := range j.Steps {
			if strings.EqualFold(s.Conclusion, "failure") {
				return j.Name, s.Name
			}
		}
		return j.Name, "unknown"
	}

	return "unknown", "unknown"
}
