package clonepick

// replay.go: load a previously persisted Plan from the DB and refresh
// its CreatedAt timestamp so most-recently-replayed selections sort to
// the top of `gitmap clone-pick --list` (planned).
//
// Replay rules (spec/01-app/100-clone-pick.md §"Replay rules"):
//   - Numeric ref  -> SELECT by SelectionId
//   - Non-numeric  -> SELECT by Name (case-sensitive, newest match wins)
//   - Replay does NOT insert a duplicate row; it bumps CreatedAt instead.
//   - --dry-run never writes to the DB (Touch is skipped by the caller).

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Loader is the read surface needed by --replay. Split from Persister
// (the write surface) so a future read-only test fake can implement
// just the lookup half. The lookup methods return the SelectionId
// alongside the Plan so callers can bump CreatedAt without a second
// round-trip.
type Loader interface {
	LoadClonePickByID(id int64) (Plan, int64, error)
	LoadClonePickByName(name string) (Plan, int64, error)
	TouchClonePickCreatedAt(id int64) error
}

// LoadFromDB resolves ref to a saved Plan + its SelectionId. Numeric
// refs hit the ID index; everything else is treated as a Name lookup.
// Returns the user-facing "no saved selection" message when nothing
// matches so the cmd layer can print verbatim.
func LoadFromDB(loader Loader, ref string) (Plan, int64, error) {
	if loader == nil {
		return Plan{}, 0, fmt.Errorf("clone-pick: --replay requires database access")
	}
	trimmed := strings.TrimSpace(ref)
	if len(trimmed) == 0 {
		return Plan{}, 0, fmt.Errorf(constants.MsgClonePickReplayNotFound, ref)
	}
	if id, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		plan, sel, loadErr := loader.LoadClonePickByID(id)
		if loadErr != nil {
			return Plan{}, 0, fmt.Errorf(constants.MsgClonePickReplayNotFound, ref)
		}

		return plan, sel, nil
	}
	plan, sel, err := loader.LoadClonePickByName(trimmed)
	if err != nil {
		return Plan{}, 0, fmt.Errorf(constants.MsgClonePickReplayNotFound, ref)
	}

	return plan, sel, nil
}

// TouchAfterReplay bumps CreatedAt on the replayed row. Best-effort:
// a failure is logged by the caller but never fails the replay.
// Skipped when dryRun is true so dry-runs stay read-only.
func TouchAfterReplay(loader Loader, id int64, dryRun bool) error {
	if loader == nil || id <= 0 || dryRun {
		return nil
	}

	return loader.TouchClonePickCreatedAt(id)
}
