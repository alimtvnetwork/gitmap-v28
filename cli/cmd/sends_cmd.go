// Package cmd — sends_cmd.go handles semantic multi-repo commit dispatch with conventional prefixes.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

var defaultSendsGitExecutor gitCmdExecutor = func(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

var currentSendsGitExecutor gitCmdExecutor = defaultSendsGitExecutor

// runSends is the CLI dispatch entry point.
func runSends(args []string) error {
	return RunSends(args)
}

// RunSends executes the semantic commit dispatch across single or multiple repositories.
func RunSends(args []string) error {
	if isSendsHelpRequested(args) {
		printSendsHelp()
		return nil
	}

	opts, errParse := parseSendsArguments(args)
	if errParse != nil {
		return errParse
	}

	return executeSendsDispatch(opts)
}

func isSendsHelpRequested(args []string) bool {
	if len(args) == 0 {
		return true
	}
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "help" {
			return true
		}
	}
	return false
}

func parseSendsArguments(args []string) (SendsOptions, error) {
	var positional []string
	opts := SendsOptions{
		IsPushed: true,
	}

	isLocalOnlyRequested := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-n" || arg == "--dry-run":
			opts.IsDryRun = true
		case arg == "--no-push":
			isLocalOnlyRequested = true
		case arg == "-j" || arg == "--json":
			opts.IsJSON = true
		case arg == "-v" || arg == "--verbose":
			opts.IsVerbose = true
		case strings.HasPrefix(arg, "-"):
			// Ignore unknown flags or handle as needed
		default:
			positional = append(positional, arg)
		}
	}

	if len(positional) < 3 {
		return opts, apperror.NewSimple("Usage: gitmap sends <verb> <target> \"<message>\" [flags]", "E9000")
	}

	opts.Verb = strings.ToLower(positional[0])
	opts.Target = positional[1]
	opts.RawMessage = strings.TrimSpace(strings.Join(positional[2:], " "))

	if !isValidSendsVerb(opts.Verb) {
		return opts, apperror.NewSimple(
			fmt.Sprintf("unknown semantic verb '%s'. Supported verbs: cpf, cpb, cpr, commit-fix, cp, cm", opts.Verb),
			"E9000",
		)
	}

	// cm verb defaults to local commit only (no push)
	if opts.Verb == "cm" || isLocalOnlyRequested {
		opts.IsPushed = false
	}

	return opts, nil
}

func isValidSendsVerb(verb string) bool {
	switch verb {
	case "cpf", "cpb", "cpr", "commit-fix", "cp", "cm":
		return true
	default:
		return false
	}
}

func resolveSemanticCommitPrefix(verb string) string {
	switch verb {
	case "cpf":
		return "Feature: "
	case "cpb":
		return "Bug: "
	case "cpr":
		return "Release: "
	case "commit-fix":
		return "Fix: "
	default:
		return ""
	}
}

func executeSendsDispatch(opts SendsOptions) error {
	cwd, errCwd := os.Getwd()
	if errCwd != nil {
		cwd = "."
	}

	allRepos := resolveWorkspaceRepositories(cwd)
	targets, errTargets := filterTargetRepositories(allRepos, opts.Target)
	if errTargets != nil {
		return errTargets
	}

	prefix := resolveSemanticCommitPrefix(opts.Verb)
	finalCommitMsg := prefix + opts.RawMessage

	var results []RepoSendResultRecord
	totalCommitted := 0
	totalSkipped := 0

	for _, repo := range targets {
		rec := processSingleRepoSend(repo, opts, finalCommitMsg)
		results = append(results, rec)

		if rec.Status == "committed" || rec.Status == "pushed" || rec.Status == "dry-run-simulated" {
			totalCommitted++
		} else if rec.Status == "clean-skipped" || rec.Status == "already-clean" {
			totalSkipped++
		}
	}

	payload := SendsExecutionPayload{
		Timestamp:          time.Now().UTC(),
		Verb:               opts.Verb,
		PrefixApplied:      prefix,
		RawMessage:         opts.RawMessage,
		FinalCommitMessage: finalCommitMsg,
		TargetScope:        opts.Target,
		IsDryRun:           opts.IsDryRun,
		IsPushed:           opts.IsPushed,
		TotalProcessed:     len(targets),
		TotalCommitted:     totalCommitted,
		TotalSkipped:       totalSkipped,
		Results:            results,
	}

	if opts.IsJSON {
		return emitSendsJSON(payload)
	}

	renderSendsTerminalSummary(payload)
	return nil
}

func filterTargetRepositories(allRepos []model.ScanRecord, target string) ([]model.ScanRecord, error) {
	if strings.EqualFold(target, "all") {
		return allRepos, nil
	}

	var matched []model.ScanRecord
	for _, r := range allRepos {
		if strings.EqualFold(r.RepoName, target) ||
			strings.EqualFold(r.Slug, target) ||
			strings.EqualFold(r.RelativePath, target) {
			matched = append(matched, r)
		}
	}

	if len(matched) == 0 {
		return nil, apperror.NewSimple(fmt.Sprintf("repository '%s' not found", target), "E9002")
	}

	return matched, nil
}

func processSingleRepoSend(repo model.ScanRecord, opts SendsOptions, finalMsg string) RepoSendResultRecord {
	dir := repo.AbsolutePath
	if dir == "" {
		dir = repo.RelativePath
	}

	porcelain, _ := currentSendsGitExecutor(dir, "status", "--porcelain")
	untracked, modified, staged, _ := parsePorcelainStatusLines(porcelain)
	isDirty := (untracked + modified + staged) > 0

	branch, _ := currentSendsGitExecutor(dir, "rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.TrimSpace(branch)
	if branch == "" {
		branch = repo.Branch
	}
	if branch == "" {
		branch = "main"
	}

	relPath := repo.RelativePath
	if relPath == "" {
		relPath = repo.RepoName
	}

	isTargetAll := strings.EqualFold(opts.Target, "all")

	if !isDirty {
		return handleCleanRepoSend(dir, repo.RepoName, relPath, branch, isTargetAll, opts)
	}

	filesStaged := untracked + modified + staged
	if opts.IsDryRun {
		return RepoSendResultRecord{
			RepoName:     repo.RepoName,
			RelativePath: relPath,
			Status:       "dry-run-simulated",
			Branch:       branch,
			FilesStaged:  filesStaged,
			IsSuccess:    true,
		}
	}

	return executeRealRepoCommitAndPush(dir, repo.RepoName, relPath, branch, filesStaged, finalMsg, opts.IsPushed)
}

func handleCleanRepoSend(dir, repoName, relPath, branch string, isTargetAll bool, opts SendsOptions) RepoSendResultRecord {
	if isTargetAll {
		return RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "clean-skipped",
			Branch:       branch,
			FilesStaged:  0,
			IsSuccess:    true,
		}
	}

	unpushedCount := checkRepoUnpushedCount(dir)
	rec, handled := tryHandleUnpushedCleanRepo(dir, repoName, relPath, branch, unpushedCount, opts.IsPushed, opts.IsDryRun)
	if handled {
		return rec
	}

	return RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       "already-clean",
		Branch:       branch,
		FilesStaged:  0,
		IsSuccess:    true,
	}
}

func tryHandleUnpushedCleanRepo(dir, repoName, relPath, branch string, unpushedCount int, isPushed, isDryRun bool) (RepoSendResultRecord, bool) {
	if unpushedCount <= 0 || !isPushed {
		return RepoSendResultRecord{}, false
	}
	return handleUnpushedCleanRepo(dir, repoName, relPath, branch, isDryRun)
}

func handleUnpushedCleanRepo(dir, repoName, relPath, branch string, isDryRun bool) (RepoSendResultRecord, bool) {
	if isDryRun {
		return RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "dry-run-simulated",
			Branch:       branch,
			FilesStaged:  0,
			IsSuccess:    true,
		}, true
	}

	_, errPush := currentSendsGitExecutor(dir, "push")
	if errPush != nil {
		return RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "failed",
			Branch:       branch,
			FilesStaged:  0,
			IsSuccess:    false,
			ErrorMessage: fmt.Sprintf("git push failed: %v", errPush),
		}, true
	}

	headSha, _ := currentSendsGitExecutor(dir, "rev-parse", "--short", "HEAD")
	return RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       "pushed",
		Branch:       branch,
		HeadSHA:      strings.TrimSpace(headSha),
		FilesStaged:  0,
		IsSuccess:    true,
	}, true
}

func executeRealRepoCommitAndPush(dir, repoName, relPath, branch string, filesStaged int, finalMsg string, isPushed bool) RepoSendResultRecord {
	_, errAdd := currentSendsGitExecutor(dir, "add", "-A")
	if errAdd != nil {
		return RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "failed",
			Branch:       branch,
			FilesStaged:  filesStaged,
			IsSuccess:    false,
			ErrorMessage: fmt.Sprintf("git add failed: %v", errAdd),
		}
	}

	_, errCommit := currentSendsGitExecutor(dir, "commit", "-m", finalMsg)
	if errCommit != nil {
		return RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "failed",
			Branch:       branch,
			FilesStaged:  filesStaged,
			IsSuccess:    false,
			ErrorMessage: fmt.Sprintf("git commit failed: %v", errCommit),
		}
	}

	headShaOut, _ := currentSendsGitExecutor(dir, "rev-parse", "--short", "HEAD")
	headSha := strings.TrimSpace(headShaOut)

	if isPushed {
		return pushCommittedRepo(dir, repoName, relPath, branch, headSha, filesStaged)
	}

	return RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       "committed",
		Branch:       branch,
		HeadSHA:      headSha,
		FilesStaged:  filesStaged,
		IsSuccess:    true,
	}
}

func pushCommittedRepo(dir, repoName, relPath, branch, headSha string, filesStaged int) RepoSendResultRecord {
	_, errPush := currentSendsGitExecutor(dir, "push", "origin", branch)
	if errPush != nil {
		return RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "failed",
			Branch:       branch,
			HeadSHA:      headSha,
			FilesStaged:  filesStaged,
			IsSuccess:    false,
			ErrorMessage: fmt.Sprintf("git push failed: %v", errPush),
		}
	}

	return RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       "pushed",
		Branch:       branch,
		HeadSHA:      headSha,
		FilesStaged:  filesStaged,
		IsSuccess:    true,
	}
}

func checkRepoUnpushedCount(dir string) int {
	revCount, errCount := currentSendsGitExecutor(dir, "rev-list", "--count", "@{u}..HEAD")
	if errCount != nil {
		return 0
	}
	count, _ := strconv.Atoi(strings.TrimSpace(revCount))
	return count
}

func emitSendsJSON(payload SendsExecutionPayload) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal sends json")
	}
	fmt.Println(string(b))
	return nil
}

func renderSendsTerminalSummary(payload SendsExecutionPayload) {
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
		if r.Status == "failed" {
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

func formatSendResultLine(r RepoSendResultRecord) string {
	reset := constants.ColorReset
	switch r.Status {
	case "pushed":
		return fmt.Sprintf(" %s✓%s %-14s -> Staged %d files -> Committed %s -> Pushed origin",
			constants.ColorGreen, reset, r.RepoName, r.FilesStaged, r.HeadSHA)
	case "committed":
		return fmt.Sprintf(" %s✓%s %-14s -> Staged %d files -> Committed %s (local only)",
			constants.ColorGreen, reset, r.RepoName, r.FilesStaged, r.HeadSHA)
	case "dry-run-simulated":
		return fmt.Sprintf(" %s●%s %-14s -> Simulation: %d files staged -> Push planned",
			constants.ColorYellow, reset, r.RepoName, r.FilesStaged)
	case "clean-skipped":
		return fmt.Sprintf(" %s-%s %-14s -> Clean working tree -> Skipped",
			constants.ColorDim, reset, r.RepoName)
	case "already-clean":
		return fmt.Sprintf(" %s-%s %-14s -> Already clean and up-to-date",
			constants.ColorDim, reset, r.RepoName)
	case "failed":
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
  cpb         Injects "Bug: " prefix, commits, and pushes
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
