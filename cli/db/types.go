// Package db — types.go defines domain payload models and single reusable Result envelopes.
package db

import "github.com/alimtvnetwork/gitmap-v28/cli/result"

type (
	// NodePathAlias represents a directory path alias mapping for a cluster node.
	NodePathAlias struct {
		NodePathAliasId int64
		NodeId          string
		Alias           string
		AbsolutePath    string
	}

	// NodePathAliasSliceResult is the canonical single reusable result envelope for node path alias slices.
	NodePathAliasSliceResult = result.ResultSlice[NodePathAlias]

	// ClusterNodeSliceResult is the canonical single reusable result envelope for cluster node slices.
	ClusterNodeSliceResult = result.ResultSlice[ClusterNode]

	// ClusterExecResultSliceResult is the canonical single reusable result envelope for cluster execution result slices.
	ClusterExecResultSliceResult = result.ResultSlice[ClusterExecResult]

	// ClusterRunSliceResult is the canonical single reusable result envelope for cluster run slices.
	ClusterRunSliceResult = result.ResultSlice[ClusterRun]

	// SSHConnectionSliceResult is the canonical single reusable result envelope for SSH connection slices.
	SSHConnectionSliceResult = result.ResultSlice[SSHConnection]
)
