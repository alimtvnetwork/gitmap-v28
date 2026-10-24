package cluster

import (
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cluster/enums"
)

// NodeStateType is the canonical cluster node state type, defined in enums.
type NodeStateType = enums.NodeStateType

const (
	// StateConnected means the node has recently sent a heartbeat.
	StateConnected = enums.StateConnected
	// StateDisconnected means the node's heartbeat has timed out.
	StateDisconnected = enums.StateDisconnected
)

// Node represents a connected VM in the cluster.
type Node struct {
	ID       string
	LastSeen time.Time
	State    NodeStateType
}

// Registry maintains the server-side state of all connected nodes.
type Registry struct {
	nodes            map[string]*Node
	mu               sync.Mutex
	heartbeatTimeout time.Duration
}

// NewRegistry creates a new NodeRegistry with the specified heartbeat timeout.
func NewRegistry(timeout time.Duration) *Registry {
	return &Registry{
		nodes:            make(map[string]*Node),
		heartbeatTimeout: timeout,
	}
}

// Register adds a new node or updates an existing node's LastSeen time.
func (r *Registry) Register(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes[id] = &Node{
		ID:       id,
		LastSeen: time.Now(),
		State:    StateConnected,
	}
}

// Ping updates the LastSeen time for an existing node.
func (r *Registry) Ping(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if node, exists := r.nodes[id]; exists {
		node.LastSeen = time.Now()
		node.State = StateConnected
	}
}

// CheckHeartbeats detects dropped nodes by evaluating their LastSeen time.
func (r *Registry) CheckHeartbeats() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, node := range r.nodes {
		isDropped := now.Sub(node.LastSeen) > r.heartbeatTimeout
		if isDropped {
			node.State = StateDisconnected
		}
	}
}

// GetNodes returns a snapshot of all registered nodes.
func (r *Registry) GetNodes() []Node {
	r.mu.Lock()
	defer r.mu.Unlock()

	snapshot := make([]Node, 0, len(r.nodes))
	for _, node := range r.nodes {
		snapshot = append(snapshot, *node)
	}

	return snapshot
}

// Disconnect marks an existing node as disconnected gracefully.
func (r *Registry) Disconnect(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if node, exists := r.nodes[id]; exists {
		node.State = StateDisconnected
	}
}
