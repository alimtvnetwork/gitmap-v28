// Package constants — clone-recency tracking (schema v26).
//
// Adds a single TEXT column to the Repo table that records the
// CURRENT_TIMESTAMP of the last clone-style operation (clone, cn,
// cfr, cfrp). Empty string means "never cloned by gitmap".
//
// Used by `gitmap release` to fall back into the most recently
// cloned repo when the user runs `gitmap r vX.Y.Z` from a parent
// directory that is not itself a git repo.
package constants

// Idempotent additive ALTER (handled via addColumnIfNotExists).
const SQLAddRepoLastClonedAt = "ALTER TABLE Repo ADD COLUMN LastClonedAt TEXT DEFAULT ''"

// Stamp the LastClonedAt for the row matching AbsolutePath.
const SQLUpdateRepoLastClonedAt = `
UPDATE Repo SET LastClonedAt = CURRENT_TIMESTAMP, UpdatedAt = CURRENT_TIMESTAMP
WHERE AbsolutePath = ?`

// Find the most-recently-cloned repo. Excludes rows that were never
// stamped (empty LastClonedAt) so legacy DBs don't surface stale rows.
const SQLSelectMostRecentClone = `
SELECT AbsolutePath, RepoName, LastClonedAt
FROM Repo
WHERE LastClonedAt IS NOT NULL AND LastClonedAt != ''
ORDER BY LastClonedAt DESC
LIMIT 1`

// User-facing messages for the auto-cd-into-recent-clone fallback.
const (
	MsgReleaseAutoCdRecent = "  " + ColorCyan + "↳ release: not inside a git repo — switching to most recently cloned: " +
		ColorReset + "%s" + ColorDim + " (cloned %s)" + ColorReset + "\n"
	MsgReleaseAutoCdReturn = "  " + ColorDim + "↩ release: returned to %s" + ColorReset + "\n"
	WarnReleaseAutoCdNone  = "  " + ColorYellow + "⚠ release: not inside a git repo and no recently-cloned repo recorded" + ColorReset + "\n"
)
