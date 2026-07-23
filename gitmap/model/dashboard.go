// Package model defines data structures for gitmap features.
package model

// DashboardData is the top-level structure written to dashboard.json.
type DashboardData struct {
	Meta      DashboardMeta `json:"meta"`
	Branches  []BranchInfo  `json:"branches"`
	Tags      []TagInfo     `json:"tags"`
	Authors   []AuthorInfo  `json:"authors"`
	Commits   []CommitInfo  `json:"commits"`
	Frequency FrequencyData `json:"frequency"`
}

// DashboardMeta holds repository-level metadata and generation context.
type DashboardMeta struct {
	RepoName      string `json:"repoName"`
	GeneratedAt   string `json:"generatedAt"`
	Branch        string `json:"branch"`
	RemoteURL     string `json:"remoteURL"`
	TotalCommits  int    `json:"totalCommits"`
	TotalBranches int    `json:"totalBranches"`
	TotalTags     int    `json:"totalTags"`
	Limit         int    `json:"limit,omitempty"`
	Since         string `json:"since,omitempty"`
}

// BranchInfo describes a single local or remote branch.
type BranchInfo struct {
	Name           string `json:"name"`
	IsRemote       bool   `json:"isRemote"`
	LastCommitSHA  string `json:"lastCommitSHA"`
	LastCommitDate string `json:"lastCommitDate"`
	Ahead          int    `json:"ahead"`
	Behind         int    `json:"behind"`
}

// TagInfo describes a single tag with its distance from the previous tag.
type TagInfo struct {
	Name        string `json:"name"`
	SHA         string `json:"sha"`
	Date        string `json:"date"`
	CommitCount int    `json:"commitCount"`
}

// AuthorInfo aggregates contribution metrics for a single author.
type AuthorInfo struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	TotalCommits int    `json:"totalCommits"`
	FirstCommit  string `json:"firstCommit"`
	LastCommit   string `json:"lastCommit"`
	ActiveDays   int    `json:"activeDays"`
}

// CommitInfo holds metadata for a single commit.
type CommitInfo struct {
	SHA          string   `json:"sha"`
	ShortSHA     string   `json:"shortSHA"`
	Author       string   `json:"author"`
	Email        string   `json:"email"`
	Date         string   `json:"date"`
	Message      string   `json:"message"`
	IsMerge      bool     `json:"isMerge"`
	FilesChanged int      `json:"filesChanged"`
	Insertions   int      `json:"insertions"`
	Deletions    int      `json:"deletions"`
	Tags         []string `json:"tags,omitempty"`
}

// FrequencyData holds pre-aggregated commit counts by time period.
type FrequencyData struct {
	Daily   map[string]int `json:"daily"`
	Weekly  map[string]int `json:"weekly"`
	Monthly map[string]int `json:"monthly"`
}
