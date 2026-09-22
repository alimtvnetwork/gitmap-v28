// Package cmdagy — agy_clean_cache_types.go defines data types for cache cleanup.
package cmdagy

// AgyCacheTarget represents a cache directory target to be cleaned.
type AgyCacheTarget struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Exists      bool   `json:"exists"`
	SizeBytes   int64  `json:"sizeBytes"`
	FileCount   int    `json:"fileCount"`
}

// AgyProcessInfo represents a discovered Antigravity/browser process.
type AgyProcessInfo struct {
	PID  int    `json:"pid"`
	Name string `json:"name"`
}

// ConvPruneCandidate represents a conversation evaluated for retention.
type ConvPruneCandidate struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ProjectID    string `json:"projectId"`
	LastModified string `json:"lastModified"`
	IsPinned     bool   `json:"isPinned"`
}

// CleanCacheOptions holds CLI flags and execution parameters.
type CleanCacheOptions struct {
	DryRun      bool `json:"dryRun"`
	Preflight   bool `json:"preflight"`
	Force       bool `json:"force"`
	Yes         bool `json:"yes"`
	NoKill      bool `json:"noKill"`
	JSON        bool `json:"json"`
	IncludeTemp bool `json:"includeTemp"`
	Keep        int  `json:"keep"`
}

// CleanCacheReport summarizes the results of the clean-cache execution.
type CleanCacheReport struct {
	Timestamp           string           `json:"timestamp"`
	DryRun              bool             `json:"dryRun"`
	Preflight           bool             `json:"preflight"`
	KeepCount           int              `json:"keepCount"`
	Targets             []AgyCacheTarget `json:"targets"`
	Processes           []AgyProcessInfo `json:"processes"`
	ProcessesTerminated int              `json:"processesTerminated"`
	BytesFreed          int64            `json:"bytesFreed"`
	HumanFreed          string           `json:"humanFreed"`
	FilesDeleted        int              `json:"filesDeleted"`
	ConvsDeleted        int              `json:"convsDeleted"`
	Warnings            []string         `json:"warnings,omitempty"`
	DurationMs          int64            `json:"durationMs"`
	BackupDir           string           `json:"backupDir,omitempty"`
}
