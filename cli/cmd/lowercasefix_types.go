// Package cmd — lowercasefix_types.go defines types for lower-case-fix command.
package cmd

// LowerCaseFixOptions holds flags and configuration for lower-case-fix.
type LowerCaseFixOptions struct {
	Patterns      []string `json:"patterns"`
	IsDryRun      bool     `json:"isDryRun"`
	IsCommit      bool     `json:"isCommit"`
	IsYes         bool     `json:"isYes"`
	CommitMessage string   `json:"commitMessage"`
}

// RenamePair represents a file to be renamed to lowercase.
type RenamePair struct {
	OldPath string `json:"oldPath"`
	NewPath string `json:"newPath"`
	OldBase string `json:"oldBase"`
	NewBase string `json:"newBase"`
}
