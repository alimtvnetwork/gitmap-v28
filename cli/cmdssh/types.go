// Package cmdssh provides SSH, cluster management, and remote execution tooling.
package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// ClusterRunResultSlice is the canonical single reusable result envelope for cluster run slices.
	ClusterRunResultSlice = result.ResultSlice[ClusterRunResult]

	// BootstrapResultSlice is the canonical single reusable result envelope for bootstrap result slices.
	BootstrapResultSlice = result.ResultSlice[BootstrapResult]

	// SSHHealthResultSlice is the canonical single reusable result envelope for SSH health result slices.
	SSHHealthResultSlice = result.ResultSlice[SSHHealthResult]

	// ClusterNodeResult is a canonical single reusable result envelope for cluster node operations.
	ClusterNodeResult = result.ErrorWrapper

	// ClusterK8sResult is a canonical single reusable result envelope for Kubernetes cluster operations.
	ClusterK8sResult = result.ErrorWrapper
)

// SSHMacroSyncOptions defines parameters for synchronizing macros across SSH nodes.
type SSHMacroSyncOptions struct {
	Target   string
	Except   string
	Exclude  string
	Name     string
	IsForce  bool
	IsDryRun bool
	IsHelp   bool
}

// SSHMacroIOOptions defines parameters for importing or exporting macros over SSH.
type SSHMacroIOOptions struct {
	Target   string
	Except   string
	Exclude  string
	FilePath string
	Name     string
	RenameAs string
	IsForce  bool
	IsDryRun bool
	IsAll    bool
	IsHelp   bool
}
