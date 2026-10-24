package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func emitSendsJSON(payload cmdpending.SendsExecutionPayload) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal sends json")
	}

	fmt.Println(string(b))

	return nil
}

func renderSendsTerminalSummary(payload cmdpending.SendsExecutionPayload) {
	border := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold

	prefixDisplay := payload.PrefixApplied
	if prefixDisplay == "" {
		prefixDisplay = "none"
	}

	fmt.Println()
	fmt.Printf("  %s┌──────────────────────────────────────────────────────────────────────────────┐%s\n", border, reset)
	fmt.Printf("  %s│%s%s                        GITMAP SENDS EXECUTION CARD                           %s%s│%s\n", border, reset, bold, reset, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s Verb: %-4s (%-8s) │ Target: %-6s │ Dry-Run: %-5t │ Push: %-5t          %s│%s\n",
		border, reset,
		payload.Verb, prefixDisplay, payload.TargetScope, payload.IsDryRun, payload.IsPushed,
		border, reset)
	fmt.Printf("  %s│%s Message: %-65s %s│%s\n",
		border, reset, truncateText(payload.FinalCommitMessage, 65), border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)

	failedCount := 0
	for _, r := range payload.Results {
		if r.Status == string(constants.SendResultStatusFailed) {
			failedCount++
		}
		line := formatSendResultLine(r)
		fmt.Printf("  %s│%s %-76s %s│%s\n", border, reset, line, border, reset)
	}

	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	summaryLine := fmt.Sprintf("Summary: %d committed & pushed | %d skipped clean | %d failed",
		payload.TotalCommitted, payload.TotalSkipped, failedCount)
	fmt.Printf("  %s│%s %-76s %s│%s\n", border, reset, summaryLine, border, reset)
	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}

func formatSendResultLine(r cmdpending.RepoSendResultRecord) string {
	reset := constants.ColorReset
	switch constants.SendResultStatusType(r.Status) {
	case constants.SendResultStatusPushed:
		return fmt.Sprintf(" %s✓%s %-14s -> Staged %d files -> Committed %s -> Pushed origin",
			constants.ColorGreen, reset, r.RepoName, r.FilesStaged, r.HeadSHA)
	case constants.SendResultStatusCommitted:
		return fmt.Sprintf(" %s✓%s %-14s -> Staged %d files -> Committed %s (local only)",
			constants.ColorGreen, reset, r.RepoName, r.FilesStaged, r.HeadSHA)
	case constants.SendResultStatusDryRunSimulated:
		return fmt.Sprintf(" %s●%s %-14s -> Simulation: %d files staged -> Push planned",
			constants.ColorYellow, reset, r.RepoName, r.FilesStaged)
	case constants.SendResultStatusCleanSkipped:
		return fmt.Sprintf(" %s-%s %-14s -> Clean working tree -> Skipped",
			constants.ColorDim, reset, r.RepoName)
	case constants.SendResultStatusAlreadyClean:
		return fmt.Sprintf(" %s-%s %-14s -> Already clean and up-to-date",
			constants.ColorDim, reset, r.RepoName)
	case constants.SendResultStatusFailed:
		return fmt.Sprintf(" %s✗%s %-14s -> FAILED: %s",
			constants.ColorRed, reset, r.RepoName, truncateText(r.ErrorMessage, 45))
	default:
		return fmt.Sprintf(" ? %-14s -> %s", r.RepoName, r.Status)
	}
}

func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return s[:maxLen]
	}

	return s[:maxLen-3] + "..."
}

func printSendsHelp() {
	cyan := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold
	yellow := constants.ColorYellow

	fmt.Printf(`
%s┌──────────────────────────────────────────────────────────────────────────────┐
│                           GITMAP SENDS HELP                                  │
└──────────────────────────────────────────────────────────────────────────────┘%s

%sDESCRIPTION:%s
  Dispatches conventional commits across targeted repositories with standardized
  semantic prefixes, non-destructive clean repository skipping, and push automation.

%sUSAGE:%s
  gitmap sends <verb> <target> "<message>" [flags]

%sSEMANTIC VERBS:%s
  cpf         Injects "Feature: " prefix, commits, and pushes
  cpb, cbf    Injects "Bug: " prefix, commits, and pushes
  cfr         Injects "Refactor: " prefix, commits, and pushes
  cpc         Injects "Chore: " prefix, commits, and pushes
  cpr         Injects "Release: " prefix, commits, and pushes
  commit-fix  Injects "Fix: " prefix, commits, and pushes
  cp          Flat commit without prefix, commits, and pushes
  cm          Local commit only without prefix, does not push

%sTARGETS:%s
  all         Iterates across all workspace repositories (skips clean repos)
  <repoName>  Targets a specific repository by name, slug, or relative path

%sFLAGS:%s
  -n, --dry-run     Preview commit and push actions without mutating Git refs
      --no-push     Commit locally but do not push to upstream remote
  -j, --json        Output structured JSON telemetry
  -v, --verbose     Print detailed git subprocess execution output
  -h, --help        Display this help menu

%sEXAMPLES:%s
  %sgitmap sends cpf all "add unified pending commits suite"%s
      Stages and commits all dirty repositories with "Feature: " prefix and pushes.

  %sgitmap sends cpb gitmap "fix null pointer in table renderer" -n%s
      Simulates a bug fix commit against 'gitmap' in dry-run mode.
`,
		cyan, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		bold, reset,
		yellow, reset,
		yellow, reset,
	)
}
