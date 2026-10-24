// Package cmd — visibilitydriftguard.go: pure decision helper for the
// `vu` / `vr` drift guard. Extracted as a seam so the policy
// (force-override vs drift-skip vs proceed) is unit-testable without
// a real provider client.
//
// Spec: 02-spec/01-app/116-bulk-visibility-mapub-mapri.md §undo-redo.
package cmdvisibility

import "github.com/alimtvnetwork/gitmap-v28/cli/cmdvisibility/enums"

// DriftActionType is the canonical drift action type, defined in enums.
type DriftActionType = enums.DriftActionType

type driftAction = DriftActionType

const (
	DriftActionTypeProceed = enums.DriftActionTypeProceed
	DriftActionTypeForce   = enums.DriftActionTypeForce
	DriftActionTypeSkip    = enums.DriftActionTypeSkip
)

const (
	driftActionProceed = DriftActionTypeProceed
	driftActionForce   = DriftActionTypeForce
	driftActionSkip    = DriftActionTypeSkip
)

// decideDriftAction is total (every input maps to exactly one action)
// and has no side effects — safe to table-test.
func decideDriftAction(current, expected string, force bool) DriftActionType {
	if force {
		return DriftActionTypeForce
	}

	if current != expected {
		return DriftActionTypeSkip
	}

	return DriftActionTypeProceed
}
