// Package cmdssh — ssh_deploy_json.go implements JSON telemetry rendering for deploy commands.
package cmdssh

import (
	"encoding/json"
	"os"
)

// DeployTargetNodeJSON holds remote node identification in JSON summaries.
type DeployTargetNodeJSON struct {
	ID      string `json:"id"`
	Alias   string `json:"alias"`
	IP      string `json:"ip"`
	OS      string `json:"os"`
	WorkDir string `json:"workDir"`
}

// DeployMetricsJSON contains performance and count metrics for a deployment.
type DeployMetricsJSON struct {
	FilesProcessed   int   `json:"filesProcessed"`
	FilesTransferred int   `json:"filesTransferred"`
	FilesSkipped     int   `json:"filesSkipped"`
	BytesTransferred int64 `json:"bytesTransferred"`
	DurationMs       int64 `json:"durationMs"`
	ParallelWorkers  int   `json:"parallelWorkers"`
}

// DeployFileRecord represents a single file transfer record.
type DeployFileRecord struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Bytes       int64  `json:"bytes,omitempty"`
	Status      string `json:"status"`
	Direction   string `json:"direction,omitempty"`
}

// DeployConflictRecord captures a detected file conflict.
type DeployConflictRecord struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	LocalMtime  int64  `json:"localMtime"`
	RemoteMtime int64  `json:"remoteMtime"`
	LocalSize   int64  `json:"localSize"`
	RemoteSize  int64  `json:"remoteSize"`
	Reason      string `json:"reason"`
}

// DeploySummaryJSON represents the root JSON telemetry summary.
type DeploySummaryJSON struct {
	Status     string                 `json:"status"`
	Command    string                 `json:"command"`
	Direction  string                 `json:"direction"`
	TargetNode DeployTargetNodeJSON   `json:"targetNode"`
	Metrics    DeployMetricsJSON      `json:"metrics"`
	Transfers  []DeployFileRecord     `json:"transfers"`
	Conflicts  []DeployConflictRecord `json:"conflicts,omitempty"`
	ExitCode   int                    `json:"exitCode"`
}

// DeployConflictPromptJSON represents the payload emitted on conflict in JSON mode.
type DeployConflictPromptJSON struct {
	Status    string                 `json:"status"`
	ExitCode  int                    `json:"exitCode"`
	Message   string                 `json:"message"`
	Conflicts []DeployConflictRecord `json:"conflicts"`
}

// RenderDeployJSONSummary renders the deploy summary to stdout as formatted JSON.
func RenderDeployJSONSummary(summary DeploySummaryJSON) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(summary)
}

// RenderConflictJSONPrompt renders non-interactive conflict notice to stdout.
func RenderConflictJSONPrompt(conflicts []DeployConflictRecord) error {
	payload := DeployConflictPromptJSON{
		Status:    "conflict",
		ExitCode:  3,
		Message:   "Conflicts detected: specify --overwrite, --skip, --sync, --sync-right, or --sync-left",
		Conflicts: conflicts,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
