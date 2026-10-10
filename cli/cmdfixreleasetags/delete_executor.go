// Package cmdfixreleasetags provides 3-tier deletion execution across
// GitHub releases, remote Git tags, local Git tags, and sidecar files.
package cmdfixreleasetags

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ExecuteDeletion executes the deletion plan for target tags.
func ExecuteDeletion(repoPath string, plan *DeletionPlan) (*ReleaseTagDeletionResult, error) {
	if plan == nil || len(plan.Records) == 0 {
		return nil, fmt.Errorf("empty deletion plan")
	}

	exec := plan.CommandExecutor
	if exec == nil {
		exec = DefaultCommandExecutor{}
	}

	results, err := ExecuteBatchDeletions(repoPath, plan.Records, plan.IsDryRun, exec)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no deletion results generated")
	}

	return &results[0], nil
}

// ExecuteBatchDeletions runs sequential deletions for all eligible candidate records.
func ExecuteBatchDeletions(repoPath string, records []ReleaseTagAuditRecord, isDryRun bool, executor CommandExecutor) ([]ReleaseTagDeletionResult, error) {
	if executor == nil {
		executor = DefaultCommandExecutor{}
	}

	var results []ReleaseTagDeletionResult
	for _, rec := range records {
		if !rec.IsEligibleForDeletion && !isDryRun {
			continue
		}

		res, err := ExecuteTagDeletion(repoPath, rec, isDryRun, executor)
		if err != nil {
			return results, err
		}

		results = append(results, *res)
	}

	return results, nil
}

// ExecuteTagDeletion executes the 4-step sequential deletion flow for a single tag.
func ExecuteTagDeletion(repoPath string, record ReleaseTagAuditRecord, isDryRun bool, executor CommandExecutor) (*ReleaseTagDeletionResult, error) {
	if executor == nil {
		executor = DefaultCommandExecutor{}
	}

	result := &ReleaseTagDeletionResult{
		Tag:         record.Tag,
		AuditReason: record.AuditReason,
	}

	s1 := executeStep1GitHubRelease(repoPath, record, isDryRun, executor)
	result.Steps = append(result.Steps, s1)

	s2 := executeStep2RemoteTag(repoPath, record, isDryRun, executor)
	result.Steps = append(result.Steps, s2)

	s3 := executeStep3LocalTag(repoPath, record, isDryRun, executor)
	result.Steps = append(result.Steps, s3)

	s4 := executeStep4SidecarMetadata(repoPath, record.Tag, isDryRun)
	result.Steps = append(result.Steps, s4)

	finalizeDeletionResult(result)

	return result, nil
}

func executeStep1GitHubRelease(repoPath string, record ReleaseTagAuditRecord, isDryRun bool, exec CommandExecutor) StepExecutionResult {
	const stepName = "GitHub Release Deletion"
	if isDryRun {
		return StepExecutionResult{StepName: stepName, Status: StepSuccess, Detail: "[dry-run] Would delete GitHub release " + record.Tag}
	}

	if !record.HasGitHubRelease {
		return StepExecutionResult{StepName: stepName, Status: StepSkipped, Detail: "No GitHub release object exists"}
	}

	out, err := exec.Run(repoPath, "gh", "release", "delete", record.Tag, "-y")
	if err != nil {
		return StepExecutionResult{StepName: stepName, Status: StepFailed, Detail: strings.TrimSpace(string(out))}
	}

	return StepExecutionResult{StepName: stepName, Status: StepSuccess, Detail: "GitHub release deleted"}
}

func executeStep2RemoteTag(repoPath string, record ReleaseTagAuditRecord, isDryRun bool, exec CommandExecutor) StepExecutionResult {
	const stepName = "Remote Git Tag Deletion"
	if isDryRun {
		return StepExecutionResult{StepName: stepName, Status: StepSuccess, Detail: "[dry-run] Would delete remote git tag refs/tags/" + record.Tag}
	}

	if !record.HasRemoteTag {
		return StepExecutionResult{StepName: stepName, Status: StepSkipped, Detail: "No remote git tag exists"}
	}

	refSpec := ":refs/tags/" + record.Tag
	out, err := exec.Run(repoPath, "git", "push", "origin", refSpec)
	if err != nil {
		return StepExecutionResult{StepName: stepName, Status: StepFailed, Detail: strings.TrimSpace(string(out))}
	}

	return StepExecutionResult{StepName: stepName, Status: StepSuccess, Detail: "Remote git tag deleted"}
}

func executeStep3LocalTag(repoPath string, record ReleaseTagAuditRecord, isDryRun bool, exec CommandExecutor) StepExecutionResult {
	const stepName = "Local Git Tag Deletion"
	if isDryRun {
		return StepExecutionResult{StepName: stepName, Status: StepSuccess, Detail: "[dry-run] Would delete local git tag " + record.Tag}
	}

	if !record.HasLocalTag {
		return StepExecutionResult{StepName: stepName, Status: StepSkipped, Detail: "No local git tag exists"}
	}

	out, err := exec.Run(repoPath, "git", "tag", "-d", record.Tag)
	if err != nil {
		return StepExecutionResult{StepName: stepName, Status: StepFailed, Detail: strings.TrimSpace(string(out))}
	}

	return StepExecutionResult{StepName: stepName, Status: StepSuccess, Detail: "Local git tag deleted"}
}

func executeStep4SidecarMetadata(repoPath string, tag string, isDryRun bool) StepExecutionResult {
	const stepName = "Sidecar Metadata Cleanup"
	relFile := filepath.Join(repoPath, constants.GitMapDir, constants.ReleaseDirName, tag+".json")
	if isDryRun {
		return StepExecutionResult{StepName: stepName, Status: StepSuccess, Detail: "[dry-run] Would remove sidecar metadata " + relFile}
	}

	_ = os.Remove(relFile)
	cleanStagingTempFiles(repoPath, tag)

	return StepExecutionResult{StepName: stepName, Status: StepSuccess, Detail: "Sidecar metadata cleaned"}
}

func cleanStagingTempFiles(repoPath, tag string) {
	tempDir := filepath.Join(repoPath, constants.GitMapDir, constants.ReleaseDirName, "temp")
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if strings.Contains(entry.Name(), tag) {
			_ = os.Remove(filepath.Join(tempDir, entry.Name()))
		}
	}
}

func finalizeDeletionResult(res *ReleaseTagDeletionResult) {
	res.IsCompleted = true
	res.IsSuccess = true

	var errs []string
	for _, step := range res.Steps {
		if step.Status == StepFailed {
			res.IsSuccess = false
			errs = append(errs, fmt.Sprintf("%s: %s", step.StepName, step.Detail))
		}
	}

	if len(errs) > 0 {
		res.ErrorMessage = strings.Join(errs, "; ")
	}
}
