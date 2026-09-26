// Package cmdagy — agy_recreate_types.go defines types and options for recreating Antigravity projects.
package cmdagy

// AgyRecreateOptions holds execution options for the recreate-project command suite.
type AgyRecreateOptions struct {
	CustomPrompt string
	Model        string
	Profile      string
	IsDryRun     bool
}

// AgyRecreateResult holds the outcome of a single project recreate lifecycle.
type AgyRecreateResult struct {
	ProjectName  string
	ProjectPath  string
	OldProjectId string
	NewProjectId string
	PurgedConvs  int
	ConvId       string
	IsSuccess    bool
	Error        error
}
