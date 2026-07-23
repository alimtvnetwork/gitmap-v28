package clonepick

// persist.go: thin wrapper over the store layer so the cmd entry can
// auto-save selections without importing store directly.
//
// The Persister interface keeps Execute testable (a fake recorder is
// enough to exercise the success path) and lets the cmd layer plug
// in the real *store.DB at call time.

import "fmt"

// Persister is the surface needed to save + look up selections.
// Implemented by *store.DB at call time; mocked in tests.
type Persister interface {
	SaveClonePickSelection(plan Plan) (int64, error)
}

// SaveSelection records the plan and returns the new SelectionId.
// Skipped (returns 0, nil) when plan.DryRun is true so dry-runs
// never touch the database.
func SaveSelection(p Persister, plan Plan) (int64, error) {
	if plan.DryRun {
		return 0, nil
	}
	if p == nil {
		// Persister optional -- when the cmd layer can't open the DB
		// we still want the clone to succeed; just log "not saved".
		return 0, nil
	}
	id, err := p.SaveClonePickSelection(plan)
	if err != nil {
		return 0, fmt.Errorf("clone-pick: save selection: %w", err)
	}

	return id, nil
}
