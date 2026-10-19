package cmdnodes

import (
	"os"
)

// FleetNodesEnvelope models root JSON envelope for distributed commands.
type FleetNodesEnvelope struct {
	Attributes FleetNodesAttributes `json:"attributes"`
	Data       FleetNodesData       `json:"data"`
}

// FleetNodesAttributes models metadata attributes of the fleet nodes command run.
type FleetNodesAttributes struct {
	GeneratedAt     string `json:"generatedAt"`
	TotalNodes      int    `json:"totalNodes"`
	ReachableNodes  int    `json:"reachableNodes"`
	LocalReposCount int    `json:"localReposCount"`
	DelegatedCount  int    `json:"delegatedCount"`
}

// FleetNodesData models the data payload containing local and remote results.
type FleetNodesData struct {
	LocalNode   NodeSummaryResult   `json:"localNode"`
	RemoteNodes []NodeSummaryResult `json:"remoteNodes"`
}

// NodeSummaryResult aggregates execution result for one fleet node.
type NodeSummaryResult struct {
	NodeAlias       string      `json:"nodeAlias"`
	Host            string      `json:"host"`
	IsOnline        bool        `json:"isOnline"`
	IsLocalMaster   bool        `json:"isLocalMaster"`
	UniqueRepoCount int         `json:"uniqueRepoCount"`
	Account         string      `json:"account,omitempty"`
	Instances       int         `json:"instances,omitempty"`
	Payload         interface{} `json:"payload,omitempty"`
	RawOutput       string      `json:"rawOutput,omitempty"`
	ErrorMessage    string      `json:"errorMessage,omitempty"`
}

func resolveLocalMasterAccount() string {
	if u := os.Getenv("USERNAME"); u != "" {
		return u
	}
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return "local-master"
}

type nodesSummaryOptions struct {
	isJSON        bool
	withPE        bool
	forceAll      bool
	verbose       bool
	releasesCount int
}
