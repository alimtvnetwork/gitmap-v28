// Package cmdagy — agy_wpr_types.go defines types for the watch-prompts-running subsystem.
package cmdagy

import "time"

// WPROptions encapsulates options and flags parsed from watch-prompts-running invocations.
type WPROptions struct {
	Subcommand     string        `json:"Subcommand"`
	Target         string        `json:"Target"`
	TargetProject  string        `json:"TargetProject"`
	Email          string        `json:"Email,omitempty"`
	IntervalStr    string        `json:"IntervalStr"`
	Interval       time.Duration `json:"Interval"`
	PrefixTemplate string        `json:"PrefixTemplate"`
	SuffixTemplate string        `json:"SuffixTemplate"`
	FilePath       string        `json:"FilePath,omitempty"`
	IsSSH          bool          `json:"IsSSH"`
	IsJSON         bool          `json:"IsJSON"`
	IsDryRun       bool          `json:"IsDryRun"`
	IsRestart      bool          `json:"IsRestart"`
	IsOnce         bool          `json:"IsOnce"`
}

// WPRMachineCacheEntry models cached machine identity and status for SSH fleet nodes.
type WPRMachineCacheEntry struct {
	NodeId      string `json:"NodeId"`
	Alias       string `json:"Alias"`
	MachineName string `json:"MachineName"`
	IPAddress   string `json:"IPAddress"`
	OS          string `json:"OS"`
	Status      string `json:"Status"`
	CachedAt    string `json:"CachedAt"`
}

// WPRProjectStatus captures live telemetry for a watched project workspace.
type WPRProjectStatus struct {
	RepoSlug       string   `json:"RepoSlug"`
	ProjectName    string   `json:"ProjectName"`
	ProjectPath    string   `json:"ProjectPath"`
	ProjectId      string   `json:"ProjectId"`
	SequenceId     string   `json:"SequenceId"`
	ConversationId string   `json:"ConversationId"`
	PromptSnippet  string   `json:"PromptSnippet"`
	PromptStatus   string   `json:"PromptStatus"`
	WordCount      int      `json:"WordCount"`
	MediaCount     int      `json:"MediaCount"`
	MediaPaths     []string `json:"MediaPaths"`
	IsRunning      bool     `json:"IsRunning"`
	IsWatching     bool     `json:"IsWatching"`
	NodeAlias      string   `json:"NodeAlias,omitempty"`
	MachineName    string   `json:"MachineName,omitempty"`
	UpdatedAt      string   `json:"UpdatedAt"`
}

// WPRRuntimeConfig persists watcher settings and registered targets list.
type WPRRuntimeConfig struct {
	IntervalSec   int      `json:"IntervalSec"`
	Prefix        string   `json:"Prefix"`
	Suffix        string   `json:"Suffix"`
	TargetSlugs   []string `json:"TargetSlugs"`
	IsEnabled     bool     `json:"IsEnabled"`
	IsFastForward bool     `json:"IsFastForward"`
	UpdatedAt     string   `json:"UpdatedAt"`
}

// WPRDeployResult records outcome of a WPR database deployment to a remote machine.
type WPRDeployResult struct {
	NodeAlias      string `json:"NodeAlias"`
	MachineName    string `json:"MachineName"`
	IPAddress      string `json:"IPAddress"`
	RepoSlug       string `json:"RepoSlug"`
	ProjectName    string `json:"ProjectName"`
	RemotePath     string `json:"RemotePath"`
	IsCopied       bool   `json:"IsCopied"`
	IsTaskEnqueued bool   `json:"IsTaskEnqueued"`
	IsRunning      bool   `json:"IsRunning"`
	ErrorMessage   string `json:"ErrorMessage,omitempty"`
	DurationMs     int64  `json:"DurationMs"`
}

// WPRDeploySummary summarizes multi-node deployment results.
type WPRDeploySummary struct {
	Status        string            `json:"Status"`
	TargetAlias   string            `json:"TargetAlias"`
	TargetProject string            `json:"TargetProject"`
	TotalNodes    int               `json:"TotalNodes"`
	SuccessCount  int               `json:"SuccessCount"`
	Results       []WPRDeployResult `json:"Results"`
	DurationMs    int64             `json:"DurationMs"`
}
