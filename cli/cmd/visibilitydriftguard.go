// Package cmd — visibilitydriftguard.go: pure decision helper for the
// `vu` / `vr` drift guard. Extracted as a seam so the policy
// (force-override vs drift-skip vs proceed) is unit-testable without
// a real provider client.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §undo-redo.
package cmd

// DriftActionType answers: given the *current* visibility we just read
// from the provider, the *expected* visibility we persisted as
// NewVisibility on the original run, and the caller's --force flag,
// should reverseOneRepo proceed, force-override, or skip-with-drift?
//
//	action == "proceed"  → no drift, apply PrevVisibility normally.
//	action == "force"    → --force set, log override line then apply.
//	action == "skip"     → drift detected, leave repo untouched.
type DriftActionType string

type driftAction = DriftActionType

const (
	DriftActionTypeProceed DriftActionType = "proceed"
	DriftActionTypeForce   DriftActionType = "force"
	DriftActionTypeSkip    DriftActionType = "skip"
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
