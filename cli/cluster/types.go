// Package cluster — types.go defines domain payload models and single reusable Result envelopes.
package cluster

import "github.com/alimtvnetwork/gitmap-v28/cli/result"

type (
	// AliasEntry represents a cluster node path alias mapping.
	AliasEntry struct {
		Alias string
		Path  string
	}

	// AliasEntrySliceResult is the canonical single reusable result envelope for alias entry slices.
	AliasEntrySliceResult = result.ResultSlice[AliasEntry]
)
