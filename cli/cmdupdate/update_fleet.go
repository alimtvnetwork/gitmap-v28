package cmdupdate

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
)

// FleetTarget represents an active cluster node for remote fleet updates.
type FleetTarget struct {
	ID       string
	Alias    string
	IP       string
	Username string
	Port     int
	Password string
	KeyPath  string
	OS       string
}

// FleetUpdateOptions contains parsed parameters for fleet update.
type FleetUpdateOptions struct {
	Target        string
	Pkg           string
	Except        string
	IsAll         bool
	IsZip         bool
	IncludeOthers bool
	IsDryRun      bool
	IsForce       bool
}

// FleetUpdateTelemetry is the structured JSON summary returned from a remote node.
type FleetUpdateTelemetry struct {
	NodeID          string   `json:"node_id,omitempty"`
	Alias           string   `json:"alias,omitempty"`
	IP              string   `json:"ip,omitempty"`
	Success         bool     `json:"success"`
	Updated         []string `json:"updated,omitempty"`
	Failed          []string `json:"failed,omitempty"`
	CurrentVersion  string   `json:"current_version,omitempty"`
	PreviousVersion string   `json:"previous_version,omitempty"`
	DurationMs      int64    `json:"duration_ms,omitempty"`
	Details         string   `json:"details,omitempty"`
}

// FleetUpdateNodeResult stores execution outcome on a single node.
type FleetUpdateNodeResult struct {
	Alias      string
	IP         string
	IsSuccess  bool
	IsOffline  bool
	Status     string
	Telemetry  FleetUpdateTelemetry
	DurationMs int64
	Details    string
	Error      error
}

// LoadFleetTargetsFn is a mockable target provider.
var LoadFleetTargetsFn = loadDefaultFleetTargets

// LoadClusterTargetsFn is a mockable additional cluster target provider for --include-others.
var LoadClusterTargetsFn = loadDefaultClusterTargets

// ExecuteRemoteUpdateFn is a mockable remote executor.
var ExecuteRemoteUpdateFn = executeDefaultRemoteUpdate

// ExecuteFleetZipUpdateFn is a mockable zip update executor.
var ExecuteFleetZipUpdateFn = executeSSHFleetZipUpdate

// CheckConnLivenessFn is a mockable connectivity liveness checker.
var CheckConnLivenessFn = cmdssh.CheckConnLiveness

// StreamFileToRemoteFn is a mockable SSH file streaming provider.
var StreamFileToRemoteFn = cmdssh.StreamFileToRemote

// StreamFileFromRemoteFn is a mockable reader for the remote update zip.
var StreamFileFromRemoteFn = cmdssh.StreamFileFromRemote

// CreateUpdateZipFn is a mockable zip archive packaging provider.
var CreateUpdateZipFn = createUpdatePackageZip
