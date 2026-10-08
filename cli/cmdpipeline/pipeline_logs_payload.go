package cmdpipeline

import (
	"fmt"
	"strings"
)

func applyPayloadOptions(p *PipelineErrorLogsPayload, runs []ghRunItem, flags PipelineErrorFlags) {
	if len(runs) > 0 {
		p.RerunEtaSeconds = calculateAverageDuration(runs, p.WorkflowName)
	}
	if !flags.IsDetailed {
		compactErrorPayload(p)
	}
	if flags.HasFix || flags.HasCheck {
		p.CICDChecks = runInternalCICDChecks(flags.HasFix)
	}
	if flags.FormatProfile != "" {
		applyFormatProfileToPayload(p, flags.FormatProfile)
	}
}

func applyFormatProfileToPayload(p *PipelineErrorLogsPayload, formatName string) {
	profile, err := LoadPEFormatProfile(formatName)
	if err != nil {
		fmt.Printf("⚠️  Warning: could not load format profile '%s': %v\n", formatName, err)

		return
	}
	applyProfileToFailedRuns(p.FailedRuns, profile)
	applyProfileToSectionFailures(p.SectionFailures, profile)
	if len(p.ErrorLogs) > 0 {
		p.ErrorLogs = FilterLogWithProfile(p.ErrorLogs, profile)
	}
	if len(p.CombinedErrors) > 0 {
		p.CombinedErrors = FilterLogWithProfile(p.CombinedErrors, profile)
	}
}

func applyProfileToFailedRuns(runs []FailedRunItem, profile *PEFormatProfile) {
	for i := range runs {
		applyProfileToFailedJobs(runs[i].FailedJobs, profile)
	}
}

func applyProfileToFailedJobs(jobs []FailedJobItem, profile *PEFormatProfile) {
	for i := range jobs {
		var filtered []string
		for _, l := range jobs[i].ErrorLines {
			if !IsLineStrippedByProfile(l, profile) {
				filtered = append(filtered, l)
			}
		}
		jobs[i].ErrorLines = filtered
	}
}

func applyProfileToSectionFailures(sections []SectionFailure, profile *PEFormatProfile) {
	for i := range sections {
		var filtered []string
		for _, l := range sections[i].ErrorLines {
			if !IsLineStrippedByProfile(l, profile) {
				filtered = append(filtered, l)
			}
		}
		sections[i].ErrorLines = filtered
	}
}

func buildErrorLogsPayload(repo string, runs []ghRunItem) PipelineErrorLogsPayload {
	return buildErrorLogsPayloadWithFlags(repo, runs, PipelineErrorFlags{})
}

func buildErrorLogsPayloadWithFlags(repo string, runs []ghRunItem, flags PipelineErrorFlags) PipelineErrorLogsPayload {
	payload := initBaseErrorLogsPayload(repo)
	payload.Runs = runs
	if len(runs) == 0 {
		return handleEmptyRunsPayload(repo, runs, payload)
	}

	populateRunsIntoPayloadWithFlags(repo, runs, &payload, flags)
	enrichErrorLogsMetadata(&payload, repo, runs)

	return payload
}

func handleEmptyRunsPayload(repo string, runs []ghRunItem, p PipelineErrorLogsPayload) PipelineErrorLogsPayload {
	if ApplyPreviousDbFallbackToPayload(&p, repo) {
		enrichErrorLogsMetadata(&p, repo, runs)

		return p
	}
	enrichErrorLogsMetadata(&p, repo, runs)

	return buildLocalOrEmptyErrorPayload(p)
}

func populateRunsIntoPayloadWithFlags(repo string, runs []ghRunItem, p *PipelineErrorLogsPayload, flags PipelineErrorFlags) {
	initTargetRunMeta(p, runs, flags.CommitTarget)
	failedRuns := resolveFailedRunsForPayloadWithTarget(repo, runs, flags.CommitTarget)
	if len(failedRuns) > 0 {
		populateFailedRunsPayload(repo, failedRuns, p)

		return
	}
	applyCleanOrFallbackState(p, repo, runs)
}

func initTargetRunMeta(p *PipelineErrorLogsPayload, runs []ghRunItem, targetCommit ...string) {
	initLatestRunMeta(p, findPrimaryTargetRun(runs, targetCommit...))
	targetSha := resolveTargetCommitSha(runs, targetCommit...)
	targetRuns := collectRunsMatchingSha(runs, targetSha)
	bestBranch := resolveBestBranchForRuns(targetRuns)
	if len(bestBranch) > 0 {
		p.Branch = bestBranch
	}
	checkAndApplyRunningState(p, runs, targetCommit...)
}

func applyCleanOrFallbackState(p *PipelineErrorLogsPayload, repo string, runs []ghRunItem) {
	if !p.IsRunning {
		p.Conclusion = "success"
	}
	if !isSkipFallback(p, runs) {
		_ = ApplyPreviousRunFallbackToPayload(p, repo, runs)
	}
}

func isSkipFallback(p *PipelineErrorLogsPayload, runs []ghRunItem) bool {
	if p.Conclusion == "success" || p.IsRunning {
		return true
	}

	if hasFailingRuns(runs) {
		return false
	}

	return true
}

func hasFailingRuns(runs []ghRunItem) bool {
	for _, r := range runs {
		if isFailingConclusion(r.Conclusion) {
			return true
		}
	}

	return false
}

func isFailingConclusion(conclusion string) bool {
	lower := strings.ToLower(strings.TrimSpace(conclusion))
	if strings.HasPrefix(lower, "cancel") {
		return true
	}

	switch lower {
	case "failure", "timed_out", "startup_failure", "action_required", "stale":
		return true
	default:
		return false
	}
}

func enrichErrorLogsMetadata(p *PipelineErrorLogsPayload, repo string, runs []ghRunItem) {
	p.RepoUrl = resolveRepoWebURL(repo)
	p.LastReleaseVersion, p.OpenPRsCount = queryReleaseAndPRsConcurrently(repo)
	p.LatestBranch = resolveLatestBranchName(p, runs)
	if isTagRef(p.Branch) || len(p.Branch) == 0 {
		p.Branch = p.LatestBranch
	}
	p.LastHash = resolveLatestCommitHash(p, runs)
}
