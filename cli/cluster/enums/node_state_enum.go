package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// NodeStateType indicates the current connection status of a cluster node.
type NodeStateType string

const (
	// StateConnected means the node has recently sent a heartbeat.
	StateConnected NodeStateType = "Connected"
	// StateDisconnected means the node's heartbeat has timed out.
	StateDisconnected NodeStateType = "Disconnected"
)

// IsValid reports whether the node state matches a known state.
func (s NodeStateType) IsValid() bool {
	switch s {
	case StateConnected, StateDisconnected:
		return true
	default:
		return false
	}
}

// String returns the string representation of the node state.
func (s NodeStateType) String() string {
	return string(s)
}

// ParseNodeStateType parses a raw string into a validated NodeStateType.
func ParseNodeStateType(s string) (NodeStateType, *apperror.AppError) {
	val := NodeStateType(strings.TrimSpace(s))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseNodeStateType", "E1000", fmt.Sprintf("invalid cluster node state: %q", s)).WithContext("value", s)
}

// AllNodeStateTypes returns all valid NodeStateType variants.
func AllNodeStateTypes() []NodeStateType {
	return []NodeStateType{StateConnected, StateDisconnected}
}

// IsConnected reports whether the node state is connected.
func (s NodeStateType) IsConnected() bool {
	return s == StateConnected
}

// IsDisconnected reports whether the node state is disconnected.
func (s NodeStateType) IsDisconnected() bool {
	return s == StateDisconnected
}
