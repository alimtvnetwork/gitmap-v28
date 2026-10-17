package cmdpe

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os/exec"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const peRerunHelp = `pe rerun — re-run the failed jobs of a GitHub Actions run.

Description:
  Re-runs only the failed jobs of a workflow run via
  ` + "`gh run rerun <run-id> --failed`" + `. Safe by default: prints the
  plan unless --yes is passed.

Usage:
  gitmap pe rerun <run-id> [--repo <owner/repo>] [--dry-run] [--yes]
  gitmap perr <run-id> [--yes]

Flags:
  --repo <slug>   target owner/repo (default: gh resolves from cwd)
  --dry-run       print what would be re-run without doing it (default: true)
  --yes           actually request the re-run (required for live execution)
  --help, -h      show this help

Example:
  gitmap pe rerun 11823445501 --dry-run
`

// RunPeRerun implements `gitmap pe rerun <run-id>` (alias: perr).
func RunPeRerun(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(peRerunHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("pe rerun", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "owner/repo target")
	dryRun := fs.Bool("dry-run", true, "print the plan without executing")
	yes := fs.Bool("yes", false, "actually request the re-run")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(peRerunHelp)
		return apperror.WrapSimple(parseErr, "pe.rerun.flags")
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Print(peRerunHelp)
		return appfault.NewValidationError("missing required argument: run-id")
	}
	runID := strings.TrimSpace(rest[0])
	if runID == "" {
		fmt.Print(peRerunHelp)
		return appfault.NewValidationError("missing required argument: run-id")
	}
	if !hasGh() {
		return appfault.NewValidationError("pe rerun requires the GitHub CLI (gh) on PATH")
	}

	ghArgs := []string{"run", "rerun", runID, "--failed"}
	if strings.TrimSpace(*repoFlag) != "" {
		ghArgs = append(ghArgs, "-R", strings.TrimSpace(*repoFlag))
	}
	planLine := "gh " + strings.Join(ghArgs, " ")

	isDryRun := *dryRun && !*yes
	if isDryRun {
		fmt.Println("pe rerun plan (dry-run, no changes made):")
		fmt.Printf("  would execute: %s\n", planLine)
		fmt.Printf("  scope: failed jobs only, run %s\n", runID)
		fmt.Println("  re-run with --yes to execute.")

		return nil
	}

	fmt.Printf("requesting re-run: %s\n", planLine)
	cmd := exec.Command("gh", ghArgs...)
	out, runErr := cmd.CombinedOutput()
	fmt.Print(string(out))
	if runErr != nil {
		return apperror.WrapSimple(runErr, "pe.rerun.exec")
	}
	fmt.Printf("re-run requested for run %s (failed jobs only)\n", runID)

	return nil
}
