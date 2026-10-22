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

type peRerunConfig struct {
	repoFlag string
	dryRun   bool
	yes      bool
	runID    string
}

// RunPeRerun implements `gitmap pe rerun <run-id>` (alias: perr).
func RunPeRerun(args []string) *appfault.AppError {
	if isPeRerunHelp(args) {
		fmt.Print(peRerunHelp)

		return nil
	}

	cfg, parseErr := parsePeRerunFlags(args)
	if parseErr != nil {
		return parseErr
	}

	return executePeRerun(cfg)
}

func isPeRerunHelp(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}

	return false
}

func parsePeRerunFlags(args []string) (peRerunConfig, *appfault.AppError) {
	fs := flag.NewFlagSet("pe rerun", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "owner/repo target")
	dryRun := fs.Bool("dry-run", true, "print the plan without executing")
	yes := fs.Bool("yes", false, "actually request the re-run")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(peRerunHelp)

		return peRerunConfig{}, apperror.WrapSimple(parseErr, "pe.rerun.flags")
	}

	return validatePeRerunArgs(fs.Args(), *repoFlag, *dryRun, *yes)
}

func validatePeRerunArgs(rest []string, repo string, dryRun, yes bool) (peRerunConfig, *appfault.AppError) {
	if len(rest) == 0 || strings.TrimSpace(rest[0]) == "" {
		fmt.Print(peRerunHelp)

		return peRerunConfig{}, appfault.NewValidationError("missing required argument: run-id")
	}

	return peRerunConfig{
		repoFlag: strings.TrimSpace(repo),
		dryRun:   dryRun,
		yes:      yes,
		runID:    strings.TrimSpace(rest[0]),
	}, nil
}

func executePeRerun(cfg peRerunConfig) *appfault.AppError {
	if !hasGh() {
		return appfault.NewValidationError("pe rerun requires the GitHub CLI (gh) on PATH")
	}

	ghArgs := buildPeRerunGhArgs(cfg)
	if cfg.dryRun && !cfg.yes {
		printPeRerunPlan(ghArgs, cfg.runID)

		return nil
	}

	return executePeRerunCmd(ghArgs, cfg.runID)
}

func buildPeRerunGhArgs(cfg peRerunConfig) []string {
	ghArgs := []string{"run", "rerun", cfg.runID, "--failed"}
	if cfg.repoFlag != "" {
		ghArgs = append(ghArgs, "-R", cfg.repoFlag)
	}

	return ghArgs
}

func printPeRerunPlan(ghArgs []string, runID string) {
	planLine := "gh " + strings.Join(ghArgs, " ")
	fmt.Println("pe rerun plan (dry-run, no changes made):")
	fmt.Printf("  would execute: %s\n", planLine)
	fmt.Printf("  scope: failed jobs only, run %s\n", runID)
	fmt.Println("  re-run with --yes to execute.")
}

func executePeRerunCmd(ghArgs []string, runID string) *appfault.AppError {
	planLine := "gh " + strings.Join(ghArgs, " ")
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
