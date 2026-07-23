package dedupe

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Verdict describes the outcome of a single ShaMap lookup. Maps to
// spec §3.1 stage 10 (`DedupeCheck`).
type Verdict struct {
	IsHit               bool  // true when SourceSha already in ShaMap
	PreviousRewrittenId int64 // populated only when IsHit
}

// Lookup queries ShaMap for sourceSha. A miss returns (Verdict{}, nil)
// — NOT an error — so the caller can treat it as a "proceed" signal.
// Any other DB error propagates so the caller can map it to
// constants.CommitInExitDbFailed.
func Lookup(db *sql.DB, sourceSha string) (Verdict, error) {
	if db == nil {
		return Verdict{}, fmt.Errorf("dedupe: nil db")
	}
	if sourceSha == "" {
		return Verdict{}, fmt.Errorf("dedupe: empty source sha")
	}
	row := db.QueryRow(sqlSelectShaMap, sourceSha)
	var rewrittenId int64
	if err := row.Scan(&rewrittenId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Verdict{}, nil
		}
		return Verdict{}, fmt.Errorf("dedupe: select %s: %w", constants.TableCommitInShaMap, err)
	}
	return Verdict{IsHit: true, PreviousRewrittenId: rewrittenId}, nil
}

// sqlSelectShaMap is the read used by Lookup. The unique index on
// SourceSha guarantees at most one row.
const sqlSelectShaMap = `SELECT RewrittenCommitId FROM ShaMap WHERE SourceSha = ?`
