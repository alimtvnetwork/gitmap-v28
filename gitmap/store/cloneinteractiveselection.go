package store

// cloneinteractiveselection.go -- persistence for `gitmap clone-pick`
// (spec 100). One row per successful invocation. See
// constants/constants_clonepick_store.go for the schema rationale.

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonepick"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// SaveClonePickSelection inserts one row representing plan and
// returns the new SelectionId. Implements clonepick.Persister so
// *store.DB can be passed straight to clonepick.Execute.
func (db *DB) SaveClonePickSelection(plan clonepick.Plan) (int64, error) {
	res, err := db.conn.Exec(constants.SQLInsertClonePickSelection,
		plan.Name,
		plan.RepoCanonicalId,
		plan.RepoUrl,
		plan.Mode,
		plan.Branch,
		plan.Depth,
		boolToInt(plan.Cone),
		boolToInt(plan.KeepGit),
		plan.DestDir,
		strings.Join(plan.Paths, ","),
		boolToInt(plan.UsedAsk),
	)
	if err != nil {
		return 0, fmt.Errorf(constants.ErrClonePickDBInsert, err)
	}

	return res.LastInsertId()
}

// boolToInt is shared with release.go (same package). Defined once
// there to avoid `redeclared in this block` (go vet / go build).
