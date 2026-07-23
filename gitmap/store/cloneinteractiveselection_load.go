package store

// cloneinteractiveselection_load.go -- read side of the
// CloneInteractiveSelection table. Powers `gitmap clone-pick --replay
// <id|name>` (spec 100). Kept in its own file so the hot insert path
// (cloneinteractiveselection.go) stays small and easy to audit.

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonepick"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// LoadClonePickByID resolves a SelectionId to a Plan + its id.
// Returns sql.ErrNoRows verbatim when the id is unknown so callers
// can distinguish "not found" from a transport failure.
func (db *DB) LoadClonePickByID(id int64) (clonepick.Plan, int64, error) {
	row := db.conn.QueryRow(constants.SQLSelectClonePickByID, id)

	return scanClonePickRow(row)
}

// LoadClonePickByName resolves a Name to the most recently created
// matching row (the SQL is `ORDER BY SelectionId DESC`). Empty names
// would match every "auto" row -- guarded explicitly so a typo can't
// silently pull the wrong selection.
func (db *DB) LoadClonePickByName(name string) (clonepick.Plan, int64, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return clonepick.Plan{}, 0, sql.ErrNoRows
	}
	row := db.conn.QueryRow(constants.SQLSelectClonePickByName, name)

	return scanClonePickRow(row)
}

// TouchClonePickCreatedAt bumps CreatedAt on the replayed row so the
// most-recently-used selections surface first in future listings.
func (db *DB) TouchClonePickCreatedAt(id int64) error {
	_, err := db.conn.Exec(constants.SQLTouchClonePickCreatedAt, id)
	if err != nil {
		return fmt.Errorf("clone-pick: touch CreatedAt: %w", err)
	}

	return nil
}

// scanClonePickRow centralizes the column->Plan mapping used by both
// lookup paths. Column order MUST mirror SQLSelectClonePickByID /
// SQLSelectClonePickByName -- both share the same prefix on purpose.
func scanClonePickRow(row *sql.Row) (clonepick.Plan, int64, error) {
	var (
		plan          clonepick.Plan
		id            int64
		coneInt       int
		keepGitInt    int
		usedAskInt    int
		pathsCsv      string
		createdAtStub sql.NullString
	)
	err := row.Scan(
		&id,
		&plan.Name,
		&plan.RepoCanonicalId,
		&plan.RepoUrl,
		&plan.Mode,
		&plan.Branch,
		&plan.Depth,
		&coneInt,
		&keepGitInt,
		&plan.DestDir,
		&pathsCsv,
		&usedAskInt,
		&createdAtStub,
	)
	if err != nil {
		return clonepick.Plan{}, 0, err
	}
	plan.Cone = coneInt != 0
	plan.KeepGit = keepGitInt != 0
	plan.UsedAsk = usedAskInt != 0
	plan.Paths = splitNonEmptyCsv(pathsCsv)
	_ = createdAtStub

	return plan, id, nil
}

// splitNonEmptyCsv guards against the "empty string -> []string{\"\"}"
// trap so a row written with zero paths round-trips as nil.
func splitNonEmptyCsv(csv string) []string {
	csv = strings.TrimSpace(csv)
	if len(csv) == 0 {
		return nil
	}

	return strings.Split(csv, ",")
}
