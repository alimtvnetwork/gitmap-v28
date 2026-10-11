package cmdsummary

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
)

// PipelineAllItem models pipeline health state for one repository in pe all.
type PipelineAllItem struct {
	RepoSlug     string             `json:"repoSlug"`
	RepoPath     string             `json:"repoPath"`
	IsActive     bool               `json:"isActive"`
	Status       string             `json:"status"` // "clean", "failed", "no-data"
	Pipeline     *RepoPipelineError `json:"pipeline,omitempty"`
	LastCommitAt string             `json:"lastCommitAt"`
}

// PipelineAllPayload models the JSON envelope for pe all.
type PipelineAllPayload struct {
	Attributes struct {
		GeneratedAt   string `json:"generatedAt"`
		TotalScanned  int    `json:"totalScanned"`
		ActiveScanned int    `json:"activeScanned"`
		ForceAll      bool   `json:"forceAll"`
	} `json:"attributes"`
	Data struct {
		CleanCount  int               `json:"cleanCount"`
		FailedCount int               `json:"failedCount"`
		Items       []PipelineAllItem `json:"items"`
	} `json:"data"`
}

// RunPipelineErrorsAll executes workspace-wide CI/CD pipeline error discovery with 48h filtering.
func RunPipelineErrorsAll(args []string) error {
	return cmdpipeline.RunPipelineErrors(append([]string{"all"}, args...))
}
