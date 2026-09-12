package replay

import "time"

// commitEnv builds the GIT_AUTHOR_*/GIT_COMMITTER_* env vars that pin
// BOTH dates byte-for-byte. RFC3339 (with offset) is the format git
// stores in the object header, so re-reading the commit later yields
// the same wall-clock and offset we passed in.
func commitEnv(p Plan) []string {
	return []string{
		"GIT_AUTHOR_NAME=" + p.AuthorName,
		"GIT_AUTHOR_EMAIL=" + p.AuthorEmail,
		"GIT_AUTHOR_DATE=" + p.AuthorDate.Format(time.RFC3339),
		"GIT_COMMITTER_NAME=" + p.AuthorName,
		"GIT_COMMITTER_EMAIL=" + p.AuthorEmail,
		"GIT_COMMITTER_DATE=" + p.CommitterDate.Format(time.RFC3339),
	}
}
