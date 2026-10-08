package cmdpipeline

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func queryReleaseAndPRsConcurrently(repo string) (string, int) {
	var tag string
	var prs int
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		tag = queryLatestTagRelease(repo)
	}()
	go func() {
		defer wg.Done()
		prs = queryPendingPRs(repo)
	}()
	wg.Wait()

	return tag, prs
}

func resolveRepoWebURL(repo string) string {
	if len(repo) > 0 {
		return normalizeRepoURL(repo)
	}

	remoteURL, err := gitutil.RemoteURL(".")
	if err == nil && len(remoteURL) > 0 {
		return formatWebURL(remoteURL)
	}

	return ""
}

func normalizeRepoURL(repo string) string {
	clean := strings.TrimPrefix(repo, "github.com/")
	if strings.HasPrefix(clean, "http://") || strings.HasPrefix(clean, "https://") {
		return clean
	}
	if !strings.Contains(clean, "/") {
		return "https://github.com/alimtvnetwork/" + clean
	}
	return "https://github.com/" + clean
}

func formatWebURL(raw string) string {
	clean := strings.TrimSuffix(strings.TrimSpace(raw), ".git")
	if strings.HasPrefix(clean, "git@github.com:") {
		return "https://github.com/" + strings.TrimPrefix(clean, "git@github.com:")
	}

	return clean
}

func resolveLatestCommitHash(p *PipelineErrorLogsPayload, runs []ghRunItem) string {
	if len(p.Sha) > 0 {
		return gitutil.TruncSha(p.Sha)
	}
	if len(runs) > 0 && len(runs[0].HeadSha) > 0 {
		return gitutil.TruncSha(runs[0].HeadSha)
	}

	targetDir := resolveRepoLocalDir(p.Repo)
	if targetDir != "" {
		return resolveTargetDirCommitSHA(targetDir)
	}
	if isCurrentRepoMatching(p.Repo) {
		return resolveLocalCommitSHA()
	}

	return ""
}

func resolveTargetDirCommitSHA(targetDir string) string {
	sha := gitutil.GetLastCommitSHA(targetDir)
	if len(sha) == 0 || sha == "-" {
		return ""
	}
	return sha
}

func resolveLocalCommitSHA() string {
	sha := gitutil.GetLastCommitSHA(".")
	if len(sha) > 0 && sha != "-" {
		return sha
	}

	return ""
}

func initBaseErrorLogsPayload(repo string) PipelineErrorLogsPayload {
	return PipelineErrorLogsPayload{
		Repo:            repo,
		DbPath:          filepath.ToSlash(pipelinedb.PipelineDbPath(repo)),
		SavedReportFile: filepath.ToSlash(resolvePipelineErrorReportPathForRepo(repo)),
	}
}

func checkAndApplyRunningState(p *PipelineErrorLogsPayload, runs []ghRunItem, targetCommit ...string) {
	if len(runs) == 0 {
		return
	}

	targetSha := resolveTargetCommitSha(runs, targetCommit...)
	bestRun := selectPrimaryActiveRun(runs, targetSha)
	if bestRun == nil {
		return
	}

	eta := CalculateRunETA(*bestRun, runs)
	setPayloadRunningState(p, *bestRun, eta)
}

func selectPrimaryActiveRun(runs []ghRunItem, targetSha string) *ghRunItem {
	candidates := collectActiveRunsForSha(runs, targetSha)
	if len(candidates) == 0 {
		return nil
	}
	if len(candidates) == 1 {
		return candidates[0]
	}

	return pickBottleneckActiveRun(runs, candidates)
}

func collectActiveRunsForSha(runs []ghRunItem, targetSha string) []*ghRunItem {
	var candidates []*ghRunItem
	for i := range runs {
		if (targetSha == "" || runs[i].HeadSha == targetSha) && isRunActive(runs[i]) {
			candidates = append(candidates, &runs[i])
		}
	}

	return candidates
}

func pickBottleneckActiveRun(runs []ghRunItem, candidates []*ghRunItem) *ghRunItem {
	best := candidates[0]
	bestDur := calculateAverageDuration(runs, best.Name)
	for _, c := range candidates[1:] {
		dur := calculateAverageDuration(runs, c.Name)
		if dur > bestDur {
			best = c
			bestDur = dur
		}
	}

	return best
}

func isRunActive(r ghRunItem) bool {
	return r.Status == "in_progress" || r.Status == "queued"
}

func initLatestRunMeta(payload *PipelineErrorLogsPayload, latest ghRunItem) {
	payload.WorkflowName = latest.Name
	payload.RunId = latest.DatabaseId
	payload.Status = latest.Status
	payload.Conclusion = latest.Conclusion
	payload.Url = latest.Url
	payload.Branch = latest.HeadBranch
	payload.Sha = latest.HeadSha
	payload.CreatedAt = latest.CreatedAt
	payload.UpdatedAt = latest.UpdatedAt
	payload.DurationSeconds = calculateRunDuration(latest.CreatedAt, latest.UpdatedAt)
	payload.SavedLogFile = filepath.ToSlash(getCachedPipelineLogPathForRepo(payload.Repo, latest.DatabaseId))
}

func setPayloadRunningState(p *PipelineErrorLogsPayload, r ghRunItem, eta int) {
	p.IsRunning = true
	p.ActiveRunName = r.Name
	p.ActiveRunId = r.DatabaseId
	p.ActiveRunUrl = r.Url
	p.EtaSeconds = eta
	p.ErrorLogs = formatRunningErrorLogs(r, eta)
}

func formatRunningErrorLogs(r ghRunItem, eta int) string {
	if eta < 0 {
		return fmt.Sprintf("Pipeline [%s #%d] is currently running (%s).",
			r.Name, r.DatabaseId, formatEtaDisplay(eta))
	}

	return fmt.Sprintf("Pipeline [%s #%d] is currently running. Estimated completion in %d seconds.",
		r.Name, r.DatabaseId, eta)
}
