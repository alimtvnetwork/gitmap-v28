// Package cluster — types.go defines domain payload models and single reusable Result envelopes.
package cluster

import (
	"context"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// AliasEntry represents a cluster node path alias mapping.
	AliasEntry struct {
		Alias string
		Path  string
	}

	// AliasEntrySliceResult is the canonical single reusable result envelope for alias entry slices.
	AliasEntrySliceResult = result.ResultSlice[AliasEntry]

	// LifecycleExecParams encapsulates arguments for cluster machine lifecycle execution (restart, shutdown, logoff).
	LifecycleExecParams struct {
		Context          context.Context
		Node             ClusterNode
		IsForceLifecycle bool
		ProvidedPassword string
	}

	// PreflightParams encapsulates arguments for cluster preflight confirmation checks.
	PreflightParams struct {
		Selector      TargetSelectorType
		Effective     []ClusterNode
		Command       string
		RunRef        string
		IsAutoConfirm bool
	}

	// NodeExecutionReportParams encapsulates arguments for reporting execution outcomes in cluster pools.
	NodeExecutionReportParams struct {
		Spinner    *pterm.SpinnerPrinter
		NodeLabel  string
		DisplayCmd string
		DurationMs int
		ExitCode   int
		CtxErr     error
		IsAllOk    bool
	}

	// FinishClusterPoolParams encapsulates arguments for finishing cluster pool executions.
	FinishClusterPoolParams struct {
		Multi         *pterm.MultiPrinter
		IsMultiActive bool
		UpdateCounts  func(bool)
		RunId         int64
		TotalNodes    int
		Succeeded     int
		Failed        int
		Skipped       int
	}
)
