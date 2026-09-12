// Package cmd — fixgit_types.go: data models and options for gitmap fix-git.
package cmd

// FixGitOptions holds execution flags for gitmap fix-git.
type FixGitOptions struct {
	IsDryRun          bool   `json:"isDryRun"`
	IsVerbose         bool   `json:"isVerbose"`
	IsJSON            bool   `json:"isJSON"`
	IsPermissionsOnly bool   `json:"isPermissionsOnly"`
	IsLocksOnly       bool   `json:"isLocksOnly"`
	IsIndexOnly       bool   `json:"isIndexOnly"`
	IsSafeDirOnly     bool   `json:"isSafeDirOnly"`
	TargetDir         string `json:"targetDir"`
}

// FixGitIssue records a diagnosed issue and its resolution status.
type FixGitIssue struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	IsFixed     bool   `json:"isFixed"`
	Remedy      string `json:"remedy"`
	ErrorDetail string `json:"errorDetail,omitempty"`
}

// FixGitResult represents the overall outcome of the fix-git operation.
type FixGitResult struct {
	TargetDir   string        `json:"targetDir"`
	IsClean     bool          `json:"isClean"`
	IssuesFound int           `json:"issuesFound"`
	IssuesFixed int           `json:"issuesFixed"`
	Issues      []FixGitIssue `json:"issues"`
}
