// Package cmdfixreleasetags defines data models, safety rules, audit detection,
// and deletion orchestration for release tags across GitHub and Git.
package cmdfixreleasetags

import (
	"os/exec"
	"time"
)

// AuditReason identifies the primary cause for flagging a release tag.
type AuditReason string

const (
	ReasonHealthy         AuditReason = "HEALTHY"
	ReasonOrphanTag       AuditReason = "ORPHAN_TAG"
	ReasonDraftRelease    AuditReason = "DRAFT_RELEASE"
	ReasonMissingAssets   AuditReason = "MISSING_ASSETS"
	ReasonCorruptedAssets AuditReason = "CORRUPTED_ASSETS"
	ReasonCICDFailed      AuditReason = "CICD_FAILED"
)

// ProtectionStatus indicates whether a tag is immune to deletion.
type ProtectionStatus string

const (
	StatusEligible        ProtectionStatus = "ELIGIBLE"
	StatusProtectedActive ProtectionStatus = "PROTECTED_ACTIVE_VERSION"
	StatusProtectedLatest ProtectionStatus = "PROTECTED_LATEST_HEALTHY"
	StatusProtectedGrace  ProtectionStatus = "PROTECTED_GRACE_WINDOW"
)

// StepStatus tracks execution outcomes for individual deletion steps.
type StepStatus string

const (
	StepPending StepStatus = "PENDING"
	StepSuccess StepStatus = "SUCCESS"
	StepSkipped StepStatus = "SKIPPED"
	StepFailed  StepStatus = "FAILED"
)

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
	Tag                   string              `json:"tag"`
	CommitSha             string              `json:"commitSha"`
	HasLocalTag           bool                `json:"hasLocalTag"`
	HasRemoteTag          bool                `json:"hasRemoteTag"`
	HasGitHubRelease      bool                `json:"hasGitHubRelease"`
	IsDraft               bool                `json:"isDraft"`
	HasAssets             bool                `json:"hasAssets"`
	AssetCount            int                 `json:"assetCount"`
	Assets                []ReleaseAssetInfo  `json:"assets,omitempty"`
	HasChecksumFile       bool                `json:"hasChecksumFile"`
	HasBinaryArchives     bool                `json:"hasBinaryArchives"`
	WorkflowRuns          []CIWorkflowRunInfo `json:"workflowRuns,omitempty"`
	AuditReason           AuditReason         `json:"auditReason"`
	ProtectionStatus      ProtectionStatus    `json:"protectionStatus"`
	IsProtected           bool                `json:"isProtected"`
	IsEligibleForDeletion bool                `json:"isEligibleForDeletion"`
}

// StepExecutionResult tracks the outcome of an individual deletion step.
type StepExecutionResult struct {
	StepName string     `json:"stepName"`
	Status   StepStatus `json:"status"`
	Detail   string     `json:"detail,omitempty"`
}

// ReleaseTagDeletionResult tracks the full 4-tier execution for a tag.
type ReleaseTagDeletionResult struct {
	Tag          string                `json:"tag"`
	AuditReason  AuditReason           `json:"auditReason"`
	Steps        []StepExecutionResult `json:"steps"`
	IsCompleted  bool                  `json:"isCompleted"`
	IsSuccess    bool                  `json:"isSuccess"`
	ErrorMessage string                `json:"errorMessage,omitempty"`
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
