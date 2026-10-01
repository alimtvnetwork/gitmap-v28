package cmdignore

// IgnoreGroup represents a named collection of gitignore patterns.
type IgnoreGroup struct {
	Name      string   `json:"name"`
	Patterns  []string `json:"patterns"`
	IsDefault bool     `json:"isDefault"`
}

// IgnoreScanIssue represents an issue detected in a repository's gitignore configuration.
type IgnoreScanIssue struct {
	RepoPath         string   `json:"repoPath"`
	RepoName         string   `json:"repoName"`
	HasDuplicate     bool     `json:"hasDuplicate"`
	DuplicateCount   int      `json:"duplicateCount"`
	HasResumeTask    bool     `json:"hasResumeTask"`
	TrackedResume    []string `json:"trackedResume,omitempty"`
	MissingGitmapDir bool     `json:"missingGitmapDir"`
}
