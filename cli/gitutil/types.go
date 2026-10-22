package gitutil

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// CaseCollisionSlice defines a collection of case collision records.
	CaseCollisionSlice = []CaseCollision

	// CaseCollisionSliceWrap is the canonical single reusable container for CaseCollision collections.
	CaseCollisionSliceWrap = result.Wrap[CaseCollisionSlice]

	// RemoteBranchSlice defines a collection of remote branch metadata.
	RemoteBranchSlice = []RemoteBranchInfo

	// RemoteBranchSliceWrap is the canonical single reusable container for remote branch collections.
	RemoteBranchSliceWrap = result.Wrap[RemoteBranchSlice]
)
