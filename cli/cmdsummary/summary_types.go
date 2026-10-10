// Package cmdsummary implements single repository and workspace-wide release summaries,
// heated file churn tracking, activity heatmaps, and CI/CD error fusion with Split-DB caching.
package cmdsummary

import (
	"time"
)

// HeatedFileMetric models churn and modification rationale for a frequently modified file.
type HeatedFileMetric struct {
	Path         string `json:"path"`
	ChangesCount int    `json:"changesCount"`
	Insertions   int    `json:"insertions"`
	Deletions    int    `json:"deletions"`
	Description  string `json:"description"`
}

// ReleaseSummaryRecord represents a summarized git release tag and its changelog gist.
type ReleaseSummaryRecord struct {
	TagName       string             `json:"tagName"`
	TagCommitHash string             `json:"tagCommitHash"`
	ReleaseDate   string             `json:"releaseDate"`
	SummaryGist   string             `json:"summaryGist"`
	WordCount     int                `json:"wordCount"`
	HeatedFiles   []HeatedFileMetric `json:"heatedFiles,omitempty"`
	IsCacheHit    bool               `json:"isCacheHit,omitempty"`
}

// RepoPipelineError models an extracted CI/CD failure trace for a repository.
type RepoPipelineError struct {
	RepoSlug     string `json:"repoSlug"`
	WorkflowName string `json:"workflowName"`
	JobName      string `json:"jobName"`
	StepName     string `json:"stepName"`
	ExitCode     int    `json:"exitCode"`
	ErrorSummary string `json:"errorSummary"`
	StackTrace   string `json:"stackTrace"`
	RunURL       string `json:"runUrl,omitempty"`
}

// RepoSummaryRecord models the summary state of an individual repository.
type RepoSummaryRecord struct {
	RepoName        string                 `json:"repoName"`
	CanonicalSlug   string                 `json:"canonicalSlug"`
	LocalPath       string                 `json:"localPath"`
	RemoteURL       string                 `json:"remoteUrl"`
	CurrentBranch   string                 `json:"currentBranch"`
	HeadCommitHash  string                 `json:"headCommitHash"`
	IsDirty         bool                   `json:"isDirty"`
	DirtyFilesCount int                    `json:"dirtyFilesCount"`
	PendingFiles    []string               `json:"pendingFiles,omitempty"`
	LastActivityAt  time.Time              `json:"lastActivityAt"`
	Releases        []ReleaseSummaryRecord `json:"releases,omitempty"`
	PipelineError   *RepoPipelineError     `json:"pipelineError,omitempty"`
	IsPipelineClean bool                   `json:"isPipelineClean,omitempty"`
	SuggestedCommit string                 `json:"suggestedCommit,omitempty"`
}

// FullSummaryOptions encapsulates runtime parameters for full summary commands.
type FullSummaryOptions struct {
	TargetRepo    string
	ReleasesCount int
	ActivityHours int
	IsJSON        bool
	WithPE        bool
	ForceAll      bool
	IsVerbose     bool
}

// FullSummaryPayload represents the root aggregated response envelope for full summary.
type FullSummaryPayload struct {
	Attributes struct {
		GeneratedAt   string `json:"generatedAt"`
		GitMapVersion string `json:"gitMapVersion"`
		Command       string `json:"command"`
		ActivityHours int    `json:"activityHours"`
	} `json:"attributes"`
	Data struct {
		TotalReposDiscovered int                 `json:"totalReposDiscovered"`
		ActiveReposCount     int                 `json:"activeReposCount"`
		DirtyReposCount      int                 `json:"dirtyReposCount"`
		MasterSanitizeCmd    string              `json:"masterSanitizeCmd,omitempty"`
		Repositories         []RepoSummaryRecord `json:"repositories"`
	} `json:"data"`
}
