package cmdpipeline

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func handlePipelineErrorLogs(args []string) error {
	if hasArgFlag(args, "--help") || hasArgFlag(args, "-h") || hasArgFlag(args, "help") {
		printPipelineErrorLogsHelp()

		return nil
	}
	if hasArgFlag(args, "last-failed-logs") {
		return HandlePipelineLastFailedLogs(args)
	}
	if handled, err := HandlePipelineHistoryAI(args); handled {
		return err
	}
	if handled, err := HandlePEFormatCommands(args); handled {
		return err
	}

	return handlePipelineHistoryOrExecute(args)
}

func handlePipelineHistoryOrExecute(args []string) error {
	flags := ParsePipelineErrorFlags(args)
	if flags.ResolvedPath != "" {
		_ = os.Chdir(flags.ResolvedPath)
	}
	if len(flags.CommitTarget) > 0 {
		return executePipelineErrorLogs(args)
	}
	if handled, err := HandlePipelineHistoryErrors(args); handled {
		return err
	}

	return executePipelineErrorLogs(args)
}

func executePipelineErrorLogs(args []string) error {
	flags := ParsePipelineErrorFlags(args)
	if flags.ResolvedPath != "" {
		_ = os.Chdir(flags.ResolvedPath)
	}

	if flags.IsAll {
		return executeAllPipelineErrorLogs(flags, args)
	}

	if err := validateTargetRepo(flags); err != nil {
		return err
	}
	repo := resolveFlagsRepoSlug(flags)

	if isInstant, decision := tryInstantCommitHashRetrieval(repo, flags); isInstant {
		return renderCachedErrorLogs(repo, decision, flags)
	}

	if flags.HasTimeline || flags.HasUntilDone {
		return executeTimelineOrUntilDoneErrorLogs(repo, flags, args)
	}

	return processAndRenderErrorLogs(repo, flags)
}

func validateTargetRepo(flags PipelineErrorFlags) error {
	if flags.IsAll {
		return nil
	}

	if flags.RawRepoTarget != "" && flags.RepoTarget == "" {
		PrintRepoTargetNotFoundDiagnostic(flags.RawRepoTarget, QueryRepoSuggestionsFromDB(flags.RawRepoTarget))
		return fmt.Errorf("target repository %q not found", flags.RawRepoTarget)
	}
	return nil
}

func resolveFlagsRepoSlug(flags PipelineErrorFlags) string {
	if len(flags.RepoTarget) > 0 {
		return flags.RepoTarget
	}
	return resolveCurrentRepoSlug()
}

func executeTimelineOrUntilDoneErrorLogs(repo string, flags PipelineErrorFlags, args []string) error {
	WaitForRunnerETAIfActive()
	if !flags.IsJSON && !flags.HasSuppressOutputLog {
		printTimelineWatchHeader(repo, flags)
	}

	return runTimelineOrUntilDoneLoop(repo, flags)
}

func printTimelineWatchHeader(repo string, flags PipelineErrorFlags) {
	mode := "dynamic timeline"
	if flags.HasUntilDone {
		mode = "until-done"
	}
	fmt.Printf("%s● Watching pipeline [%s] (%s, interval: 2m)...%s\n",
		constants.ColorCyan, repo, mode, constants.ColorReset)
}

func runTimelineOrUntilDoneLoop(repo string, flags PipelineErrorFlags) error {
	for {
		isDone, err := pollTimelineLoopStep(repo, flags)
		if isDone {
			return err
		}

		if isTimelineTestMode() {
			return nil
		}

		time.Sleep(120 * time.Second)
	}
}

func pollTimelineLoopStep(repo string, flags PipelineErrorFlags) (bool, error) {
	runs := queryWorkflowRunsForTarget(repo, flags.CommitTarget)
	RecordFetchedRunsToSplitDb(repo, runs)
	payload := buildErrorLogsPayloadWithFlags(repo, runs, flags)
	applyPayloadOptions(&payload, runs, flags)

	if hasEarlyDetectedErrors(payload) {
		renderTimelineErrorAndExit(payload, flags)
		return true, fmt.Errorf("pipeline errors detected in active run")
	}

	if !payload.IsRunning {
		return true, writeTimelineOutputPayload(payload, flags)
	}

	printTimelineProgress(payload.ActiveRunName, payload.EtaSeconds, 120)
	return false, nil
}

func renderTimelineErrorAndExit(payload PipelineErrorLogsPayload, flags PipelineErrorFlags) {
	_ = writeTimelineOutputPayload(payload, flags)
	cliexit.Exit(1)
}

func writeTimelineOutputPayload(payload PipelineErrorLogsPayload, flags PipelineErrorFlags) error {
	return writeOrRenderErrorLogs(ErrorLogOutputParams{
		Payload:              payload,
		IsJSON:               flags.IsJSON,
		HasSuppressOutputLog: flags.HasSuppressOutputLog,
		FilePath:             flags.FilePath,
		TempFile:             flags.TempFileName,
		WantFix:              flags.HasFix,
	})
}

func hasEarlyDetectedErrors(payload PipelineErrorLogsPayload) bool {
	if len(payload.FailedRuns) > 0 || len(payload.SectionFailures) > 0 {
		return true
	}
	if isFailingConclusion(payload.Conclusion) {
		return true
	}
	if len(payload.CombinedErrors) > 0 {
		return true
	}
	if isFailingErrorLogs(payload.ErrorLogs) {
		return true
	}

	return false
}

func isFailingErrorLogs(logs string) bool {
	if len(logs) == 0 {
		return false
	}
	if strings.Contains(logs, "No failed pipeline steps") {
		return false
	}

	return strings.Contains(logs, "FAIL") || strings.Contains(logs, "Error") || strings.Contains(logs, "error")
}

func tryInstantCommitHashRetrieval(repo string, flags PipelineErrorFlags) (bool, PipelineCacheDecision) {
	if flags.HasForce {
		return false, PipelineCacheDecision{}
	}

	db, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		return false, PipelineCacheDecision{}
	}
	defer db.Close()

	runRes := db.QueryRecentRuns(20)
	if runRes.IsFailure() || len(runRes.Data) == 0 {
		return false, PipelineCacheDecision{}
	}

	targetSha := resolveTargetOrHeadSha(repo, flags)
	if len(targetSha) == 0 {
		return false, PipelineCacheDecision{}
	}

	if checkCommitTargetCacheHit(runRes.Data, targetSha) {
		return true, buildCacheHitDecision(runRes.Data, targetSha, "instant_hash_retrieval_completed")
	}

	return false, PipelineCacheDecision{}
}

func resolveTargetOrHeadSha(repo string, flags PipelineErrorFlags) string {
	if len(flags.CommitTarget) > 0 && !flags.HasIndex {
		return flags.CommitTarget
	}

	return ResolveRepoHeadCommitSha(repo)
}

func processAndRenderErrorLogs(repo string, flags PipelineErrorFlags) error {
	activeLogLineLimit = flags.Limit
	if isInstant, decision := tryInstantCommitHashRetrieval(repo, flags); isInstant {
		return renderCachedErrorLogs(repo, decision, flags)
	}

	decision := EvaluatePipelineErrorsCache(repo, flags)
	if decision.IsFromCache {
		return renderCachedErrorLogs(repo, decision, flags)
	}

	return fetchAndRenderFreshErrorLogs(repo, flags)
}

func renderCachedErrorLogs(repo string, decision PipelineCacheDecision, flags PipelineErrorFlags) error {
	activeLogLineLimit = flags.Limit
	payload := buildErrorLogsPayloadWithFlags(repo, decision.CachedRuns, flags)
	payload.IsFromCache = true
	payload.CacheSource = decision.CacheSource
	applyPayloadOptions(&payload, decision.CachedRuns, flags)
	_ = RecordPipelineCheckInTask(repo, payload)

	return writeOrRenderErrorLogs(ErrorLogOutputParams{
		Payload:              payload,
		IsJSON:               flags.IsJSON,
		HasSuppressOutputLog: flags.HasSuppressOutputLog,
		FilePath:             flags.FilePath,
		TempFile:             flags.TempFileName,
		WantFix:              flags.HasFix,
	})
}

func fetchAndRenderFreshErrorLogs(repo string, flags PipelineErrorFlags) error {
	activeLogLineLimit = flags.Limit
	printReadingProgress(flags)
	runs := queryWorkflowRunsForTarget(repo, flags.CommitTarget)
	RecordFetchedRunsToSplitDb(repo, runs)
	payload := buildErrorLogsPayloadWithFlags(repo, runs, flags)
	applyPayloadOptions(&payload, runs, flags)
	_ = RecordPipelineCheckInTask(repo, payload)

	return writeOrRenderErrorLogs(ErrorLogOutputParams{
		Payload:              payload,
		IsJSON:               flags.IsJSON,
		HasSuppressOutputLog: flags.HasSuppressOutputLog,
		FilePath:             flags.FilePath,
		TempFile:             flags.TempFileName,
		WantFix:              flags.HasFix,
	})
}

func printReadingProgress(flags PipelineErrorFlags) {
	if flags.IsJSON || flags.HasSuppressOutputLog {
		return
	}

	fmt.Printf("\n  Reading pipeline logs...\n\n")
}
