// Package cmd — sends_cmd.go handles semantic multi-repo commit dispatch with conventional prefixes.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsends"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// gitCmdExecutor defines a function that executes a git command in a specific directory.
type gitCmdExecutor func(dir string, args ...string) (string, error)

var defaultSendsGitExecutor gitCmdExecutor = func(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

var currentSendsGitExecutor gitCmdExecutor = defaultSendsGitExecutor

var resolveAllIndexedRepositoriesFn func(cwd string) []model.ScanRecord

// runSends is the CLI dispatch entry point.
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
		if constants.IsHelpFlag(a) {
			return true
		}
	}

	return false
}

func parseSendsArguments(args []string) (cmdpending.SendsOptions, error) {
	var positional []string
	opts := cmdpending.SendsOptions{
		IsPushed: true,
	}

	isLocalOnlyRequested := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case constants.IsDryRunFlag(arg):
			opts.IsDryRun = true
		case arg == constants.FlagNoPush:
			isLocalOnlyRequested = true
		case constants.IsJSONFlag(arg):
			opts.IsJSON = true
		case constants.IsVerboseFlag(arg):
			opts.IsVerbose = true
		case strings.HasPrefix(arg, "-"):
			// Ignore unknown flags or handle as needed
		default:
			positional = append(positional, arg)
		}
	}

	if len(positional) < 3 {
		return opts, apperror.NewSimple("Usage: gitmap sends <verb> <target> \"<message>\" [flags]", string(constants.ErrCodeGenericExecution))
	}

	opts.Verb = strings.ToLower(positional[0])
	opts.Target = positional[1]
	opts.RawMessage = strings.TrimSpace(strings.Join(positional[2:], " "))

	if !isValidSendsVerb(opts.Verb) {
		return opts, apperror.NewSimple(
			fmt.Sprintf("unknown semantic verb '%s'. Supported verbs: cpf, cpb, cbf, cfr, cpc, cpr, commit-fix, cp, cm", opts.Verb),
			string(constants.ErrCodeGenericExecution),
		)
	}

	// cm verb defaults to local commit only (no push)
	if opts.Verb == string(constants.CommitVerbFlat) || isLocalOnlyRequested {
		opts.IsPushed = false
	}

	return opts, nil
}

func isValidSendsVerb(verb string) bool {
	switch constants.CommitVerbType(verb) {
	case constants.CommitVerbFeature, constants.CommitVerbBug, constants.CommitVerbBugAlt, constants.CommitVerbRefactor, constants.CommitVerbChore, constants.CommitVerbRelease, constants.CommitVerbFlat:
		return true
	default:
		return verb == "commit-fix" || verb == "cp"
	}
}

func resolveSemanticCommitPrefix(verb string) string {
	verbType := constants.CommitVerbType(verb)
	if verbType.IsValid() {
		return verbType.Prefix()
	}

	if verb == "commit-fix" {
		return "Fix: "
	}

	return ""
}

func executeSendsDispatch(opts cmdpending.SendsOptions) error {
	cwd, errCwd := os.Getwd()
	if errCwd != nil {
		cwd = "."
	}

	var allRepos []model.ScanRecord
	isTargetAll := strings.EqualFold(opts.Target, "all")
	if isTargetAll {
		allRepos = resolveAllIndexedRepositories(cwd)
	} else {
		allRepos = cmdpending.ResolveWorkspaceRepositories(cwd)
	}

	targets, errTargets := filterTargetRepositories(allRepos, opts.Target)
	if errTargets != nil {
		return errTargets
	}

	prefix := resolveSemanticCommitPrefix(opts.Verb)
	finalCommitMsg := prefix + opts.RawMessage

	var results []cmdpending.RepoSendResultRecord
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

	payload := cmdpending.SendsExecutionPayload{
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

func resolveAllIndexedRepositories(cwd string) []model.ScanRecord {
	if resolveAllIndexedRepositoriesFn != nil {
		return resolveAllIndexedRepositoriesFn(cwd)
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		return cmdpending.ResolveWorkspaceRepositories(cwd)
	}
	defer db.Close()

	records, errList := db.ListRepos()
	if errList != nil || len(records) == 0 {
		return cmdpending.ResolveWorkspaceRepositories(cwd)
	}

	deduplicated := deduplicateScanRecords(records)
	if len(deduplicated) > 0 {
		return deduplicated
	}

	return cmdpending.ResolveWorkspaceRepositories(cwd)
}

func deduplicateScanRecords(records []model.ScanRecord) []model.ScanRecord {
	seen := make(map[string]bool)
	var deduplicated []model.ScanRecord

	for _, r := range records {
		cleanPath := resolveRecordCleanPath(r)
		if cleanPath == "" || cleanPath == "." {
			continue
		}

		key := strings.ToLower(cleanPath)
		if seen[key] {
			continue
		}

		seen[key] = true
		r.AbsolutePath = cleanPath
		deduplicated = append(deduplicated, r)
	}

	return deduplicated
}

func resolveRecordCleanPath(r model.ScanRecord) string {
	cleanPath := filepath.Clean(r.AbsolutePath)
	if cleanPath == "" || cleanPath == "." {
		return filepath.Clean(r.RelativePath)
	}

	return cleanPath
}

func processSingleRepoSend(repo model.ScanRecord, opts cmdpending.SendsOptions, finalMsg string) cmdpending.RepoSendResultRecord {
	dir := repo.AbsolutePath
	if dir == "" {
		dir = repo.RelativePath
	}

	porcelain, _ := currentSendsGitExecutor(dir, "status", "--porcelain")
	untracked, modified, staged, _ := cmdpending.ParsePorcelainStatusLines(porcelain)
	filesStaged := untracked + modified + staged

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
	isClean := filesStaged == 0

	if isClean {
		return handleCleanRepoSend(dir, repo.RepoName, relPath, branch, isTargetAll, opts)
	}

	if opts.IsDryRun {
		return cmdpending.RepoSendResultRecord{
			RepoName:     repo.RepoName,
			RelativePath: relPath,
			Status:       "dry-run-simulated",
			Branch:       branch,
			FilesStaged:  filesStaged,
			IsSuccess:    true,
		}
	}

	repoCommitMsg := resolveRepoCommitMessage(opts, finalMsg, porcelain)
	return executeRealRepoCommitAndPush(dir, repo.RepoName, relPath, branch, filesStaged, repoCommitMsg, opts.IsPushed)
}

func resolveRepoCommitMessage(opts cmdpending.SendsOptions, finalMsg, porcelain string) string {
	if !isSemanticGroupingApplicable(opts) {
		return finalMsg
	}

	categories := cmdsends.CategorizeChangedFiles(porcelain)
	if len(categories) > 0 {
		return cmdsends.FormatMarkdownCommitBody(finalMsg, categories)
	}

	return finalMsg
}

func isSemanticGroupingApplicable(opts cmdpending.SendsOptions) bool {
	isTargetAll := strings.EqualFold(opts.Target, "all")
	if isTargetAll {
		return true
	}

	return isWIPTrigger(opts.RawMessage)
}

func isWIPTrigger(msg string) bool {
	clean := strings.TrimSpace(strings.ToLower(msg))
	isWIPExact := clean == "wip" || clean == "save changes"
	if isWIPExact {
		return true
	}
	hasWIPPrefix := strings.HasPrefix(clean, "wip:") || strings.HasPrefix(clean, "wip ") || strings.HasPrefix(clean, "wip-")
	if hasWIPPrefix {
		return true
	}

	return strings.Contains(clean, "wip") || strings.Contains(clean, "save changes")
}

func handleCleanRepoSend(dir, repoName, relPath, branch string, isTargetAll bool, opts cmdpending.SendsOptions) cmdpending.RepoSendResultRecord {
	if isTargetAll {
		return cmdpending.RepoSendResultRecord{
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

	return cmdpending.RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       "already-clean",
		Branch:       branch,
		FilesStaged:  0,
		IsSuccess:    true,
	}
}

func tryHandleUnpushedCleanRepo(dir, repoName, relPath, branch string, unpushedCount int, isPushed, isDryRun bool) (cmdpending.RepoSendResultRecord, bool) {
	if unpushedCount <= 0 || !isPushed {
		return cmdpending.RepoSendResultRecord{}, false
	}
	return handleUnpushedCleanRepo(dir, repoName, relPath, branch, isDryRun)
}

func handleUnpushedCleanRepo(dir, repoName, relPath, branch string, isDryRun bool) (cmdpending.RepoSendResultRecord, bool) {
	if isDryRun {
		return cmdpending.RepoSendResultRecord{
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
		return cmdpending.RepoSendResultRecord{
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
	return cmdpending.RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       "pushed",
		Branch:       branch,
		HeadSHA:      strings.TrimSpace(headSha),
		FilesStaged:  0,
		IsSuccess:    true,
	}, true
}
