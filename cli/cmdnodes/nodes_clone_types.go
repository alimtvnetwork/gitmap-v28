// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdnodes/enums"
)

// NodesCloneKindType is the canonical clone operation mode type, defined in enums.
type NodesCloneKindType = enums.NodesCloneKindType

const (
	// CloneKindClone executes standard multi-repo or single-repo clone.
	CloneKindClone = enums.CloneKindClone
	// CloneKindCFR executes clone-fix-repo (remediation & auto-setup).
	CloneKindCFR = enums.CloneKindCFR
	// CloneKindCFRP executes clone-fix-repo with public visibility promotion.
	CloneKindCFRP = enums.CloneKindCFRP
)

// NodesCloneDestination describes resolved destination directory settings for clone operations.
type NodesCloneDestination struct {
	TargetDir          string `json:"targetDir,omitempty"`
	RelativeSubdir     string `json:"relativeSubdir,omitempty"`
	HasCustomTargetDir bool   `json:"hasCustomTargetDir"`
	IsInsideWorkDir    bool   `json:"isInsideWorkDir"`
}

// NodesCloneOptions stores configuration for fleet clone execution.
type NodesCloneOptions struct {
	Kind                NodesCloneKindType    `json:"kind"`
	TargetFilter        string                `json:"targetFilter,omitempty"`
	ExcludeFilter       string                `json:"excludeFilter,omitempty"`
	ExceptOS            string                `json:"exceptOS,omitempty"`
	TargetOS            string                `json:"targetOS,omitempty"`
	TargetDir           string                `json:"targetDir,omitempty"`
	TargetRepo          string                `json:"targetRepo,omitempty"`
	RelativeSubdir      string                `json:"relativeSubdir,omitempty"`
	HasCustomTargetDir  bool                  `json:"hasCustomTargetDir"`
	IsInsideWorkDir     bool                  `json:"isInsideWorkDir"`
	Destination         NodesCloneDestination `json:"destination"`
	IsDryRun            bool                  `json:"isDryRun"`
	IsJSON              bool                  `json:"isJSON"`
	IsSkipLocal         bool                  `json:"isSkipLocal"`
	IsIDESyncSuppressed bool                  `json:"isIdeSyncSuppressed"`
	RawArgs             []string              `json:"rawArgs"`
	PassArgs            []string              `json:"passArgs"`
	DetectedFile        string                `json:"detectedFile,omitempty"`
	HasFile             bool                  `json:"hasFile"`
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
	Details    string        `json:"details,omitempty"`
}

// SetFleetCloneActive toggles recursion suppression for inner suggestions.
func SetFleetCloneActive(isActive bool) {
	cmdclone.SetFleetCloneActive(isActive)
}

// IsFleetCloneActive checks if fleet clone execution is currently active.
func IsFleetCloneActive() bool {
	return cmdclone.IsFleetCloneActive()
}
