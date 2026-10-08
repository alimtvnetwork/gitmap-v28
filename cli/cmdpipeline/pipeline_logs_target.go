package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func resolveExplicitTargetSha(runs []ghRunItem, targetCommit ...string) string {
	if len(targetCommit) == 0 || len(targetCommit[0]) == 0 {
		return ""
	}

	return resolveTargetShaFromGroups(runs, targetCommit[0])
}

func resolveTargetCommitSha(runs []ghRunItem, targetCommit ...string) string {
	if len(runs) == 0 {
		return ""
	}
	if sha := resolveExplicitTargetSha(runs, targetCommit...); len(sha) > 0 {
		return sha
	}
	if sha := findShaForActiveBranch(runs); len(sha) > 0 {
		return sha
	}

	return runs[0].HeadSha
}

func resolveTargetShaFromGroups(runs []ghRunItem, target string) string {
	groups := GroupRunsByCommit(runs)
	group, isFound := ResolveCommitGroupByTarget(groups, target)
	if isFound {
		return group.HeadSha
	}

	return ""
}

func findShaForActiveBranch(runs []ghRunItem) string {
	activeBranch := gitutil.GetActiveBranch(".")
	if len(activeBranch) == 0 || activeBranch == "-" {
		return ""
	}
	for _, r := range runs {
		if strings.EqualFold(r.HeadBranch, activeBranch) {
			return r.HeadSha
		}
	}

	return ""
}

func findPrimaryTargetRun(runs []ghRunItem, targetCommit ...string) ghRunItem {
	if len(runs) == 0 {
		return ghRunItem{}
	}

	targetSha := resolveTargetCommitSha(runs, targetCommit...)
	for _, r := range runs {
		if r.HeadSha == targetSha {
			return r
		}
	}

	return runs[0]
}

func collectFailedRuns(runs []ghRunItem) []ghRunItem {
	return collectFailedRunsForRepo("", runs)
}

func collectFailedRunsForRepo(repo string, runs []ghRunItem, targetCommit ...string) []ghRunItem {
	if len(runs) == 0 {
		return nil
	}
	targetSha := resolveTargetCommitSha(runs, targetCommit...)
	targetRuns := collectRunsMatchingSha(runs, targetSha)
	succeeded := make(map[string]bool)
	var activeFailed []ghRunItem
	for _, r := range targetRuns {
		checkAndCollectRunForRepo(repo, r, succeeded, &activeFailed)
	}

	return capFailedRuns(activeFailed, 5)
}

func resolveFailedRunsForPayloadWithTarget(repo string, runs []ghRunItem, target string) []ghRunItem {
	return collectFailedRunsForRepo(repo, runs, target)
}

func queryWorkflowRunsForTarget(repo string, target string) []ghRunItem {
	runs := queryWorkflowRuns(repo)
	if len(target) == 0 {
		return runs
	}

	groups := GroupRunsByCommit(runs)
	if _, isFound := ResolveCommitGroupByTarget(groups, target); isFound {
		return runs
	}

	targetSha := resolveTargetCommitCandidate(target)
	if !isCommitHexSha(targetSha) {
		return runs
	}

	commitRuns := queryWorkflowRunsByCommit(repo, targetSha)
	if len(commitRuns) > 0 {
		return append(commitRuns, runs...)
	}

	return runs
}

func resolveTargetCommitCandidate(target string) string {
	offset, isNegative := ParseNegativeIndex(target)
	if !isNegative {
		return target
	}

	resolved := resolveCommitShaByOffset(offset)
	if len(resolved) > 0 {
		return resolved
	}

	return target
}

func resolveCommitShaByOffset(offset int) string {
	absOffset := NormalizeCommitOffset(offset)
	if absOffset <= 0 {
		return ""
	}

	rev := fmt.Sprintf("HEAD~%d", absOffset)
	cmd := exec.Command("git", "rev-parse", "--short=7", rev)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func queryWorkflowRunsByCommit(repo, sha string) []ghRunItem {
	out, err := runGHCommandWithTimeout("run", "list", "--repo", repo, "--commit", sha, "--limit", "10", "--json",
		"databaseId,name,status,conclusion,createdAt,updatedAt,headBranch,headSha,url,displayTitle,event")
	if err != nil || len(out) == 0 {
		return nil
	}

	var runs []ghRunItem
	if jsonErr := json.Unmarshal(out, &runs); jsonErr != nil {
		return nil
	}

	return runs
}

func normalizeWorkflowKey(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	lower = strings.TrimPrefix(lower, ".github/workflows/")
	lower = strings.TrimPrefix(lower, "workflows/")
	lower = strings.TrimSuffix(lower, ".yml")
	lower = strings.TrimSuffix(lower, ".yaml")

	return lower
}

func buildWorkflowScopeKey(r ghRunItem) string {
	nameKey := normalizeWorkflowKey(r.Name)
	branchKey := strings.ToLower(strings.TrimSpace(r.HeadBranch))
	if len(nameKey) == 0 {
		return ""
	}
	if len(branchKey) == 0 {
		return nameKey
	}

	return branchKey + ":" + nameKey
}

func checkAndCollectRunForRepo(repo string, r ghRunItem, succeeded map[string]bool, active *[]ghRunItem) {
	scopeKey := buildWorkflowScopeKey(r)
	if r.Conclusion == "success" {
		recordWorkflowSuccess(scopeKey, succeeded)

		return
	}
	if isWorkflowScopeSucceeded(scopeKey, succeeded) {
		return
	}
	if isFailingConclusion(r.Conclusion) {
		*active = append(*active, r)

		return
	}
	if isRunActive(r) && hasActiveRunFailedJobs(repo, r) {
		failingRun := r
		failingRun.Conclusion = "failure"
		*active = append(*active, failingRun)
	}
}

func hasActiveRunFailedJobs(repo string, r ghRunItem) bool {
	jobs := queryRunJobs(repo, r.DatabaseId)
	for _, j := range jobs {
		if isJobFailing(j) || hasAnyFailingStep(j.Steps) {
			return true
		}
	}

	return false
}

func recordWorkflowSuccess(scopeKey string, succeeded map[string]bool) {
	if len(scopeKey) > 0 {
		succeeded[scopeKey] = true
	}
}

func isWorkflowScopeSucceeded(scopeKey string, succeeded map[string]bool) bool {
	return len(scopeKey) > 0 && succeeded[scopeKey]
}

func collectRunsMatchingSha(runs []ghRunItem, targetSha string) []ghRunItem {
	var filtered []ghRunItem
	for _, r := range runs {
		if len(targetSha) == 0 || r.HeadSha == targetSha {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

func capFailedRuns(runs []ghRunItem, limit int) []ghRunItem {
	effectiveLimit := resolveFailedRunsLimit(limit)
	if len(runs) > effectiveLimit {
		return runs[:effectiveLimit]
	}

	return runs
}

func resolveFailedRunsLimit(limit int) int {
	if activeLogLineLimit > 0 {
		return activeLogLineLimit
	}

	return limit
}

func buildLocalOrEmptyErrorPayload(payload PipelineErrorLogsPayload) PipelineErrorLogsPayload {
	localErr := readLocalLastErrorLog()

	if len(localErr) > 0 {
		payload.WorkflowName = "Local Command Execution"
		payload.Conclusion = "failure"
		payload.ErrorLogs = localErr

		return payload
	}

	payload.ErrorLogs = "No failed pipeline steps or local error logs found."

	return payload
}

func readLocalLastErrorLog() string {
	content, err := os.ReadFile(".gitmap/last_error.log")

	if err == nil && len(content) > 0 {
		return strings.TrimSpace(string(content))
	}

	return ""
}
