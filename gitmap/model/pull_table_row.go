package model

// PullTableRow represents a row in the rich pull results table.
type PullTableRow struct {
	RepoName     string `json:"repoName"`
	Branch       string `json:"branch"`
	LatestBranch string `json:"latestBranch"`
	LastSHA      string `json:"lastSha"`
	PRStatus     string `json:"prStatus"`
	PullStatus   string `json:"pullStatus"`
	Duration     string `json:"duration"`
	IsDirty      bool   `json:"isDirty"`
	Reason       string `json:"reason,omitempty"`
}
