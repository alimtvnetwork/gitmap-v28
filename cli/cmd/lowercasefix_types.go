// Package cmd — lowercasefix_types.go defines types for lowercase fix command.
package cmd

// LowerCaseFixOptions holds flags and configuration for lowercase fix.
type LowerCaseFixOptions struct {
	Patterns      []string `json:"patterns"`
	IsDryRun      bool     `json:"isDryRun"`
	IsNoCommit    bool     `json:"isNoCommit"`
	IsYes         bool     `json:"isYes"`
	IsReadmeOnly  bool     `json:"isReadmeOnly"`
	CommitMessage string   `json:"commitMessage"`
}

// RenamePair represents a file to be renamed to lowercase.
type RenamePair struct {
	OldPath      string `json:"oldPath"`
	NewPath      string `json:"newPath"`
	OldBase      string `json:"oldBase"`
	NewBase      string `json:"newBase"`
	RelPath      string `json:"relPath"`
	IsGitTracked bool   `json:"isGitTracked"`
	MatchedBy    string `json:"matchedBy"`
}

// RenameSummary captures the overall execution statistics.
type RenameSummary struct {
	TotalScanned int    `json:"totalScanned"`
	TotalMatched int    `json:"totalMatched"`
	TotalRenamed int    `json:"totalRenamed"`
	IsGitRepo    bool   `json:"isGitRepo"`
	CommitSHA    string `json:"commitSha"`
}
