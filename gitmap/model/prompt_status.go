// Package model — prompt_status.go defines installation status models.
package model

type PromptInstallResult struct {
	RepoName  string `json:"repoName"`
	RepoPath  string `json:"repoPath"`
	IsSuccess bool   `json:"isSuccess"`
	Version   string `json:"version"`
	Duration  string `json:"duration"`
	Error     string `json:"error,omitempty"`
}

// PromptRunStatusRecord is an alias for PromptInstallResult.
type PromptRunStatusRecord = PromptInstallResult

// IsFailed reports whether the prompt installation failed.
func (r PromptInstallResult) IsFailed() bool {

	return !r.IsSuccess
}

// IsFail reports whether the prompt installation failed.
func (r PromptInstallResult) IsFail() bool {

	return !r.IsSuccess
}
