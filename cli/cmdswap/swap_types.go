package cmdswap

import (
	"time"
)

// SwapOptions encapsulates all input parameters for the binary swap workflow.
type SwapOptions struct {
	FromPath       string `json:"fromPath"`
	ToPath         string `json:"toPath"`
	TargetPID      int    `json:"targetPid"`
	HasKill        bool   `json:"hasKill"`
	HasRestart     bool   `json:"hasRestart"`
	HasForce       bool   `json:"hasForce"`
	IsScriptOnly   bool   `json:"isScriptOnly"`
	IsJSONOutput   bool   `json:"isJsonOutput"`
	MaxWaitSeconds int    `json:"maxWaitSeconds"`
}

// ScriptPlan models the platform-specific detached execution plan.
type ScriptPlan struct {
	Platform      string        `json:"platform"`
	ScriptPath    string        `json:"scriptPath"`
	ScriptContent string        `json:"scriptContent"`
	TargetPIDList []int         `json:"targetPidList"`
	IsDetached    bool          `json:"isDetached"`
	ExecutionArgs []string      `json:"executionArgs"`
	EstimatedWait time.Duration `json:"estimatedWait"`
}

// SwapResult encapsulates the execution outcome returned to callers or emitted as JSON.
type SwapResult struct {
	IsSuccess      bool        `json:"isSuccess"`
	FromPath       string      `json:"fromPath"`
	ToPath         string      `json:"toPath"`
	TerminatedPIDs []int       `json:"terminatedPids"`
	ScriptPlan     *ScriptPlan `json:"scriptPlan,omitempty"`
	ErrorMessage   string      `json:"errorMessage,omitempty"`
	OutputMessage  string      `json:"outputMessage"`
	DurationMs     int64       `json:"durationMs"`
}

// ProcessLockInfo describes a detected process holding a file handle or executing the target binary.
type ProcessLockInfo struct {
	PID         int    `json:"pid"`
	ProcessName string `json:"processName"`
	Executable  string `json:"executable"`
	IsParent    bool   `json:"isParent"`
}
