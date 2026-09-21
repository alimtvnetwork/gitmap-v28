package cmdssh

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SSHHistoryTask represents a snapshot of deleted SSH nodes for undo/restore.
type SSHHistoryTask struct {
	TaskID     string          `json:"taskId"`
	Action     string          `json:"action"`
	Target     string          `json:"target"`
	Nodes      []store.SSHHost `json:"nodes"`
	CreatedAt  time.Time       `json:"createdAt"`
	RestoredAt string          `json:"restoredAt,omitempty"`
}

const (
	ActionRmNode     = "rm_node"
	ActionClearNodes = "clear_nodes"
	ActionResetNodes = "reset_nodes"
)
