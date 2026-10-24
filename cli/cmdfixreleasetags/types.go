// Package cmdfixreleasetags defines data models, safety rules, audit detection,
// and deletion orchestration for release tags across GitHub and Git.
package cmdfixreleasetags

import (
	"os/exec"
	"time"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdfixreleasetags/enums"

// AuditReasonType identifies the primary cause for flagging a release tag.
// AuditReasonType aliases the strongly-typed enum from the enums package.
type AuditReasonType = enums.AuditReasonType

const (
	ReasonHealthy         = enums.ReasonHealthy
	ReasonOrphanTag       = enums.ReasonOrphanTag
	ReasonDraftRelease    = enums.ReasonDraftRelease
	ReasonMissingAssets   = enums.ReasonMissingAssets
	ReasonCorruptedAssets = enums.ReasonCorruptedAssets
	ReasonCICDFailed      = enums.ReasonCICDFailed
)

// AuditReason is an alias for AuditReasonType.
type AuditReason = enums.AuditReasonType

// ProtectionStatusType indicates whether a tag is immune to deletion.
// ProtectionStatusType aliases the strongly-typed enum from the enums package.
type ProtectionStatusType = enums.ProtectionStatusType

const (
	StatusEligible        = enums.StatusEligible
	StatusProtectedActive = enums.StatusProtectedActive
	StatusProtectedLatest = enums.StatusProtectedLatest
	StatusProtectedGrace  = enums.StatusProtectedGrace
)

// ProtectionStatus is an alias for ProtectionStatusType.
type ProtectionStatus = enums.ProtectionStatusType

// StepStatusType tracks execution outcomes for individual deletion steps.
// StepStatusType aliases the strongly-typed enum from the enums package.
type StepStatusType = enums.StepStatusType

const (
	StepPending = enums.StepPending
	StepSuccess = enums.StepSuccess
	StepSkipped = enums.StepSkipped
	StepFailed  = enums.StepFailed
)

// StepStatus is an alias for StepStatusType.
type StepStatus = enums.StepStatusType

// ReleaseAssetInfo details an uploaded asset file.
type ReleaseAssetInfo struct {
	Name              string `json:"name"`
	Size              int64  `json:"size"`
	State             string `json:"state"`
	IsValid           bool   `json:"isValid"`
	HasBinaryArchives bool   `json:"hasBinaryArchives"`
}

// CIWorkflowRunInfo summarizes a GitHub Actions workflow run for a tag/commit.
type CIWorkflowRunInfo struct {
	RunId        uint64        `json:"runId"`
	WorkflowName string        `json:"workflowName"`
	Status       string        `json:"status"`
	Conclusion   string        `json:"conclusion"`
	CreatedAt    time.Time     `json:"createdAt"`
	Duration     time.Duration `json:"duration"`
	IsSuccessful bool          `json:"isSuccessful"`
	IsInProgress bool          `json:"isInProgress"`
}

// ReleaseTagAuditRecord contains the complete audit diagnosis for a single tag.
type ReleaseTagAuditRecord struct {
	Tag                   string               `json:"tag"`
	CommitSha             string               `json:"commitSha"`
	HasLocalTag           bool                 `json:"hasLocalTag"`
	HasRemoteTag          bool                 `json:"hasRemoteTag"`
	HasGitHubRelease      bool                 `json:"hasGitHubRelease"`
	IsDraft               bool                 `json:"isDraft"`
	HasAssets             bool                 `json:"hasAssets"`
	AssetCount            int                  `json:"assetCount"`
	Assets                []ReleaseAssetInfo   `json:"assets,omitempty"`
	HasChecksumFile       bool                 `json:"hasChecksumFile"`
	HasBinaryArchives     bool                 `json:"hasBinaryArchives"`
	WorkflowRuns          []CIWorkflowRunInfo  `json:"workflowRuns,omitempty"`
	AuditReason           AuditReasonType      `json:"auditReason"`
	ProtectionStatus      ProtectionStatusType `json:"protectionStatus"`
	IsProtected           bool                 `json:"isProtected"`
	IsEligibleForDeletion bool                 `json:"isEligibleForDeletion"`
}

// StepExecutionResult tracks the outcome of an individual deletion step.
type StepExecutionResult struct {
	StepName string         `json:"stepName"`
	Status   StepStatusType `json:"status"`
	Detail   string         `json:"detail,omitempty"`
}

// ReleaseTagDeletionResult tracks the full 4-tier execution for a tag.
type ReleaseTagDeletionResult struct {
	Tag          string                `json:"tag"`
	AuditReason  AuditReasonType       `json:"auditReason"`
	Steps        []StepExecutionResult `json:"steps"`
	IsCompleted  bool                  `json:"isCompleted"`
	IsSuccess    bool                  `json:"isSuccess"`
	ErrorMessage string                `json:"errorMessage,omitempty"`
}

// IsFailed reports whether the deletion was incomplete or unsuccessful.
func (r ReleaseTagDeletionResult) IsFailed() bool {
	return !r.IsCompleted || !r.IsSuccess
}

// AuditSummary aggregates counts across the repository audit.
type AuditSummary struct {
	TotalTagsChecked       int `json:"totalTagsChecked"`
	HealthyReleasesCount   int `json:"healthyReleasesCount"`
	OrphanTagsCount        int `json:"orphanTagsCount"`
	DraftReleasesCount     int `json:"draftReleasesCount"`
	MissingAssetsCount     int `json:"missingAssetsCount"`
	CorruptedAssetsCount   int `json:"corruptedAssetsCount"`
	CICDFailedCount        int `json:"cicdFailedCount"`
	ProtectedTagsCount     int `json:"protectedTagsCount"`
	EligibleDeletionsCount int `json:"eligibleDeletionsCount"`
}

// AuditReport is the top-level payload for audit and deletion runs.
type AuditReport struct {
	RepoPath         string                     `json:"repoPath"`
	ActiveVersion    string                     `json:"activeVersion"`
	LatestHealthyTag string                     `json:"latestHealthyTag"`
	Summary          AuditSummary               `json:"summary"`
	Records          []ReleaseTagAuditRecord    `json:"records"`
	DeletionResults  []ReleaseTagDeletionResult `json:"deletionResults,omitempty"`
}

// CommandExecutor executes commands in a directory, injectable for mocking.
type CommandExecutor interface {
	Run(dir string, name string, args ...string) ([]byte, error)
}

// DefaultCommandExecutor runs commands directly on the host system.
type DefaultCommandExecutor struct{}

// Run executes the named command with args in dir.
func (e DefaultCommandExecutor) Run(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	return cmd.CombinedOutput()
}

// AuditOptions configures audit execution.
type AuditOptions struct {
	GracePeriod           time.Duration   `json:"gracePeriod"`
	ActiveVersionOverride string          `json:"activeVersionOverride"`
	MockExec              bool            `json:"mockExec"`
	CommandExecutor       CommandExecutor `json:"-"`
}

// AuditFilterOptions aliases AuditOptions for compatibility with both specifications.
type AuditFilterOptions = AuditOptions

// DeletionPlan details the targets and options for deletion execution.
type DeletionPlan struct {
	RepoPath        string                  `json:"repoPath"`
	IsDryRun        bool                    `json:"isDryRun"`
	Records         []ReleaseTagAuditRecord `json:"records"`
	CommandExecutor CommandExecutor         `json:"-"`
}
