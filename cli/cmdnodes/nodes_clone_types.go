// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
)

// NodesCloneKind defines the clone operation mode.
type NodesCloneKind string

const (
	// CloneKindClone executes standard multi-repo or single-repo clone.
	CloneKindClone NodesCloneKind = "clone"
	// CloneKindCFR executes clone-fix-repo (remediation & auto-setup).
	CloneKindCFR NodesCloneKind = "cfr"
	// CloneKindCFRP executes clone-fix-repo with public visibility promotion.
	CloneKindCFRP NodesCloneKind = "cfrp"
)

// NodesCloneOptions stores configuration for fleet clone execution.
type NodesCloneOptions struct {
	Kind          NodesCloneKind `json:"kind"`
	TargetFilter  string         `json:"targetFilter,omitempty"`
	ExcludeFilter string         `json:"excludeFilter,omitempty"`
	ExceptOS      string         `json:"exceptOS,omitempty"`
	TargetOS      string         `json:"targetOS,omitempty"`
	IsDryRun      bool           `json:"isDryRun"`
	IsJSON        bool           `json:"isJSON"`
	IsSkipLocal   bool           `json:"isSkipLocal"`
	RawArgs       []string       `json:"rawArgs"`
	PassArgs      []string       `json:"passArgs"`
	DetectedFile  string         `json:"detectedFile,omitempty"`
	HasFile       bool           `json:"hasFile"`
}

// RemoteCloneNodeResult captures the execution output from an individual node.
type RemoteCloneNodeResult struct {
	Alias      string        `json:"alias"`
	Host       string        `json:"host"`
	Role       string        `json:"role"`
	Status     string        `json:"status"`
	Duration   time.Duration `json:"duration"`
	DurationMs int64         `json:"durationMs"`
	Stdout     string        `json:"stdout,omitempty"`
	Stderr     string        `json:"stderr,omitempty"`
	Error      string        `json:"error,omitempty"`
}

// SetFleetCloneActive toggles recursion suppression for inner suggestions.
func SetFleetCloneActive(active bool) {
	cmdclone.SetFleetCloneActive(active)
}

// IsFleetCloneActive checks if fleet clone execution is currently active.
func IsFleetCloneActive() bool {
	return cmdclone.IsFleetCloneActive()
}
