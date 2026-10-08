package cmdupdate

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// UpdateOptions specifies parameters for self/fleet updates.
type UpdateOptions struct {
	TargetVersion   string `json:"targetVersion"`
	IsForce         bool   `json:"isForce"`
	IsDryRun        bool   `json:"isDryRun"`
	IsJSON          bool   `json:"isJSON"`
	IsQuiet         bool   `json:"isQuiet"`
	MaxFallbackTags int    `json:"maxFallbackTags"`
}

// ReleaseCacheRecord mirrors the Split SQLite ReleaseCache table.
type ReleaseCacheRecord = store.ReleaseCacheRecord

// UpdateResult encapsulates the full telemetry of an update operation.
type UpdateResult struct {
	Success          bool     `json:"success"`
	Status           string   `json:"status"`
	PreviousVersion  string   `json:"previous_version"`
	CurrentVersion   string   `json:"current_version"`
	TargetTag        string   `json:"target_tag"`
	InstallDir       string   `json:"install_dir"`
	BinaryPath       string   `json:"binary_path"`
	DownloaderEngine string   `json:"downloader_engine"`
	DownloadBytes    int64    `json:"download_size_bytes"`
	DurationMs       int64    `json:"duration_ms"`
	FallbackDepth    int      `json:"fallback_depth"`
	IsCached         bool     `json:"cached"`
	IsAlreadyUpdated bool     `json:"already_updated"`
	ScriptsUpdated   []string `json:"scripts_updated"`
	ErrorMessage     *string  `json:"error"`
}

// AlreadyUpdatedResult represents the JSON payload when update is skipped because already on target.
type AlreadyUpdatedResult struct {
	Status         string `json:"status"`
	CurrentVersion string `json:"current_version"`
	TargetVersion  string `json:"target_version"`
	Updated        bool   `json:"updated"`
	Message        string `json:"message"`
}

// ResultUpdate wraps UpdateResult in monadic Result.
type ResultUpdate = result.Result[UpdateResult]
