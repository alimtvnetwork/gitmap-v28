package cmdci

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const ciOpenHelp = `ci open — print the deep link to the repo's latest Actions run.

Description:
  Resolves the latest GitHub Actions run for [repo] and prints its URL.
  When the run is red, the link deep-links to the failing job. The URL
  is printed for opening in Chrome; this command never launches a
  browser itself.

Usage:
  gitmap ci open [repo]

Flags:
  --help, -h    show this help

Example:
  gitmap ci open go-email-reader
`

// RunCiOpen implements `gitmap ci open [repo]`.
func RunCiOpen(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(ciOpenHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("ci open", flag.ContinueOnError)
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(ciOpenHelp)
		return apperror.WrapSimple(parseErr, "ci.open.flags")
	}
	if !hasCiGh() {
		return appfault.NewValidationError("ci open requires the GitHub CLI (gh) on PATH")
	}
	var repoArg string
	if rest := fs.Args(); len(rest) > 0 {
		repoArg = rest[0]
	}
	slug, slugErr := ciRepoSlug(repoArg)
	if slugErr != nil {
		fmt.Print(ciOpenHelp)
		return slugErr
	}

	var runs []ciRun
	if listErr := ciGhJSON(&runs, "run", "list", "-R", slug, "--limit", "1",
		"--json", "databaseId,conclusion,displayTitle,url"); listErr != nil {
		return listErr
	}
	if len(runs) == 0 {
		fmt.Printf("ci open: no workflow runs found for %s\n", slug)
		return nil
	}
	run := runs[0]
	link := run.URL
	failedJob := ""
	if strings.EqualFold(run.Conclusion, "failure") {
		var view struct {
			Jobs []ciJob `json:"jobs"`
		}
		id := fmt.Sprintf("%d", run.DatabaseID)
		if viewErr := ciGhJSON(&view, "run", "view", id, "-R", slug, "--json", "jobs"); viewErr == nil {
			for _, j := range view.Jobs {
				if strings.EqualFold(j.Conclusion, "failure") {
					failedJob = j.Name
					link = fmt.Sprintf("%s/job/%d", run.URL, j.DatabaseID)
					break
				}
			}
		}
	}

	fmt.Printf("latest run: %s (%s)\n", run.Title, run.Conclusion)
	if failedJob != "" {
		fmt.Printf("failed job: %s\n", failedJob)
	}
	fmt.Printf("open in Chrome: %s\n", link)

	return nil
}
